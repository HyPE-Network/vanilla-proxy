package proxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/HyPE-Network/vanilla-proxy/handler"
	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/HyPE-Network/vanilla-proxy/math"
	"github.com/HyPE-Network/vanilla-proxy/proxy/block"
	"github.com/HyPE-Network/vanilla-proxy/proxy/command"
	"github.com/HyPE-Network/vanilla-proxy/proxy/console"
	"github.com/HyPE-Network/vanilla-proxy/proxy/console/bash"
	"github.com/HyPE-Network/vanilla-proxy/proxy/console/bots"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/HyPE-Network/vanilla-proxy/proxy/whitelist"
	"github.com/HyPE-Network/vanilla-proxy/proxy/world"
	"github.com/HyPE-Network/vanilla-proxy/utils"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/resource"

	"github.com/sandertv/gophertunnel/minecraft"
)

var ProxyInstance *Proxy

type Proxy struct {
	Worlds           *world.Worlds
	Config           utils.Config
	PlayerManager    human.HumanManager
	Handlers         handler.HandlerManager
	CommandManager   *command.CommandManager
	Listener         *minecraft.Listener
	CommandSender    console.CommandSender
	WhitelistManager *whitelist.WhitelistManager
}

func New(config utils.Config, hm human.HumanManager) *Proxy {
	block.Init()

	commandManager := command.InitManager(config.Server.Ops)

	Proxy := &Proxy{
		Config:         config,
		PlayerManager:  hm,
		CommandManager: commandManager,
	}

	if config.WorldBorder.Enabled {
		Proxy.Worlds = world.Init(math.NewArea2(config.WorldBorder.MinX, config.WorldBorder.MinZ, config.WorldBorder.MaxX, config.WorldBorder.MaxZ))
	}

	if config.Server.Whitelist {
		Proxy.WhitelistManager = whitelist.Init(commandManager)
	}

	os := runtime.GOOS
	switch os {
	case "windows":
		if Proxy.Config.Bot.Enabled {
			log.Logger.Debugln("Creating new console bot instance..")
			Proxy.CommandSender = bots.NewBot(Proxy.Config)
		} else {
			log.Logger.Warnln("Console bot is disabled in config")
		}
	case "linux":
		log.Logger.Debugln("Creating new bash console instance..")
		Proxy.CommandSender = bash.NewBash(strings.Split(Proxy.Config.Connection.RemoteAddress, ":")[1])
	}

	if Proxy.CommandSender == nil {
		log.Logger.Warnln("CommandSender is not declared, functionality will be disabled")
	}

	return Proxy
}

func (arg *Proxy) Start(h handler.HandlerManager) error {
	arg.Handlers = h

	if arg.Config.Rcon.Enabled {
		go command.InitRCON(arg.CommandManager.Commands, arg.Config.Rcon.Port, arg.Config.Rcon.Password)
	}

	p, err := minecraft.NewForeignStatusProvider(arg.Config.Connection.RemoteAddress)
	if err != nil {
		return err
	}

	// Initialize an empty slice of *resource.Pack
	var resourcePacks []*resource.Pack

	// Loop through all the pack URLs and append each pack to the slice
	for _, url := range arg.Config.Resources.PackURLs {
		resourcePack := resource.MustReadURL(url)
		resourcePacks = append(resourcePacks, resourcePack)
	}

	arg.Listener, err = minecraft.ListenConfig{ // server settings
		AuthenticationDisabled: arg.Config.Server.DisableXboxAuth,
		StatusProvider:         p,
		ResourcePacks:          resourcePacks,
		TexturePacksRequired:   true,
	}.Listen("raknet", arg.Config.Connection.ProxyAddress)

	if err != nil {
		return err
	}

	defer arg.Stop()

	log.Logger.Debugln("Original server address:", arg.Config.Connection.RemoteAddress, "public address:", arg.Config.Connection.ProxyAddress)
	log.Logger.Println("Proxy has been started on Version", protocol.CurrentVersion, "protocol", protocol.CurrentProtocol)

	for {
		c, err := arg.Listener.Accept()

		if err != nil {
			log.Logger.Errorln(err)
			continue
		}

		go arg.handleConn(c.(*minecraft.Conn))
	}
}

func (arg *Proxy) Stop() {
	arg.CommandSender.Close()
	arg.PlayerManager.DeleteAll()
	arg.Listener.Close()
}

func (arg *Proxy) handleConn(conn *minecraft.Conn) {
	if human, ok := arg.PlayerManager.PlayerList()[conn.IdentityData().DisplayName]; ok { // if the user is already in the system
		err := conn.Close()
		if err != nil {
			log.Logger.Errorln(err)
		}

		err = arg.Listener.Disconnect(conn, "connection lost")
		if err != nil {
			log.Logger.Errorln(err)
		}

		arg.deletePlayer(human)
		return
	}

	serverConn, err := minecraft.Dialer{
		KeepXBLIdentityData: true,
		ClientData:          conn.ClientData(),
		IdentityData:        conn.IdentityData(),
	}.DialTimeout("raknet", arg.Config.Connection.RemoteAddress, time.Second*60)

	if err != nil {
		log.Logger.Errorln("Error in establishing serverConn: ", err)
		arg.CloseConnections(conn, serverConn)
		return
	}

	var success = true
	var g sync.WaitGroup
	g.Add(2)
	go func() {
		if err := conn.StartGame(serverConn.GameData()); err != nil {
			log.Logger.Errorln(err)
			success = false
		}
		g.Done()
	}()
	go func() {
		if err := serverConn.DoSpawn(); err != nil {
			log.Logger.Errorln(err)
			success = false
		}
		g.Done()
	}()
	g.Wait()

	if !success {
		arg.CloseConnections(conn, serverConn)
		return
	}

	if arg.Config.Server.Whitelist {
		if !arg.WhitelistManager.HasPlayer(conn.IdentityData().DisplayName, conn.IdentityData().XUID) {
			arg.CloseConnections(conn, serverConn)
			return
		}
	}

	pl := arg.PlayerManager.AddPlayer(conn, serverConn)
	log.Logger.Infoln(pl.GetName(), "joined the server")
	pl.SendXUIDToAddon()
	arg.UpdatePlayerDetails(pl)

	go func() { // client-proxy
		for {
			pk, err := conn.ReadPacket()
			if err != nil {
				break
			}

			ok, pk, err := arg.Handlers.HandlePacket(pk, pl, "Client")
			if err != nil {
				log.Logger.Errorln(err)
			}

			if ok {
				if err := serverConn.WritePacket(pk); err != nil {
					if disconnect, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
						_ = arg.Listener.Disconnect(conn, disconnect.Error())
					}
					break
				}
			}
		}
		arg.deletePlayer(pl)
	}()
	go func() { // proxy-server
		for {
			pk, err := serverConn.ReadPacket()
			if err != nil {
				if disconnect, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
					_ = arg.Listener.Disconnect(conn, disconnect.Error())
				}
				break
			}

			ok, pk, err := arg.Handlers.HandlePacket(pk, pl, "Server")
			if err != nil {
				log.Logger.Errorln(err)
			}

			if ok {
				if err := conn.WritePacket(pk); err != nil {
					break
				}
			}
		}
		arg.deletePlayer(pl)
	}()
}

func (arg *Proxy) deletePlayer(human human.Human) {
	arg.PlayerManager.DeletePlayer(human)
	arg.Listener.Disconnect(human.GetSession().Connection.ClientConn, "connection lost")
}

func (arg *Proxy) SendConsoleCommand(cmd string) {
	if arg.CommandSender == nil {
		log.Logger.Warnln("CommandSender is not declared. Ignored:", cmd)
	} else {
		if err := arg.CommandSender.SendCommand(cmd); err != nil {
			log.Logger.Errorln("CommandSender:", err)
		}
	}
}

func (arg *Proxy) CloseConnections(conn *minecraft.Conn, serverConn *minecraft.Conn) {
	err := conn.Close()
	if err != nil {
		log.Logger.Errorln(err)
	}
	err = serverConn.Close()
	if err != nil {
		log.Logger.Errorln(err)
	}
	err = arg.Listener.Disconnect(conn, "connection lost")
	if err != nil {
		log.Logger.Errorln(err)
	}
}

type PlayerDetails struct {
	Xuid string `json:"xuid"`
	Name string `json:"name"`
	IP   string `json:"ip"`
}

func (arg *Proxy) UpdatePlayerDetails(player human.Human) {
	xuid := player.GetSession().IdentityData.XUID

	// Build the URI for the API request
	uri := arg.Config.Api.ApiHost + "/api/moderation/playerDetails"
	log.Logger.Printf("Sending \"%s\" playerDetails to: \"%s\"\n", player.GetName(), uri)

	// Create the player details payload
	playerDetails := PlayerDetails{
		Xuid: xuid,
		Name: player.GetName(),
		IP:   strings.Split(player.GetSession().ClientData.ServerAddress, ":")[0],
	}

	// Convert player details to JSON
	jsonData, err := json.Marshal(playerDetails)
	if err != nil {
		log.Logger.Errorln("Failed to marshal player details:", err)
		return
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Logger.Errorln("Failed to create new request:", err)
		return
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", arg.Config.Api.ApiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Logger.Errorln("Failed to send request:", err)
		return
	}
	defer resp.Body.Close()

	// Log the response status
	log.Logger.Printf("Sent playerDetails to: \"%s\", status: %d\n", uri, resp.StatusCode)
}
