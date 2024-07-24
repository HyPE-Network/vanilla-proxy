package proxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/HyPE-Network/vanilla-proxy/handler"
	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/HyPE-Network/vanilla-proxy/math"
	"github.com/HyPE-Network/vanilla-proxy/proxy/command"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/HyPE-Network/vanilla-proxy/proxy/whitelist"
	"github.com/HyPE-Network/vanilla-proxy/proxy/world"
	"github.com/HyPE-Network/vanilla-proxy/utils"
	"github.com/google/uuid"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/resource"

	"github.com/sandertv/gophertunnel/minecraft"
)

var ProxyInstance *Proxy

type Proxy struct {
	Worlds           *world.Worlds
	Config           utils.Config
	Handlers         handler.HandlerManager
	CommandManager   *command.CommandManager
	Listener         *minecraft.Listener
	WhitelistManager *whitelist.WhitelistManager
}

func New(config utils.Config) *Proxy {
	commandManager := command.InitManager(config.Server.Ops)

	Proxy := &Proxy{
		Config:         config,
		CommandManager: commandManager,
	}

	if config.WorldBorder.Enabled {
		Proxy.Worlds = world.Init(math.NewArea2(config.WorldBorder.MinX, config.WorldBorder.MinZ, config.WorldBorder.MaxX, config.WorldBorder.MaxZ))
	}

	if config.Server.Whitelist {
		Proxy.WhitelistManager = whitelist.Init(commandManager)
	}

	return Proxy
}

// The following program implements a proxy that forwards players from one local address to a remote address.
func (arg *Proxy) Start(h handler.HandlerManager) error {
	arg.Handlers = h

	p, err := minecraft.NewForeignStatusProvider(arg.Config.Connection.RemoteAddress)
	if err != nil {
		return err
	}

	// Initialize an empty slice of *resource.Pack
	var resourcePacks []*resource.Pack

	// Loop through all the pack URLs and append each pack to the slice
	for _, url := range arg.Config.Resources.PackURLs {
		resourcePack, err := resource.ReadURL(url)
		if err != nil {
			return err
		}
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

	log.Logger.Debugln("Original server address:", arg.Config.Connection.RemoteAddress, "public address:", arg.Config.Connection.ProxyAddress)
	log.Logger.Println("Proxy has been started on Version", protocol.CurrentVersion, "protocol", protocol.CurrentProtocol)

	defer arg.Listener.Close()
	for {
		c, err := arg.Listener.Accept()
		if err != nil {
			log.Logger.Errorln(err)
			continue
		}
		log.Logger.Debugln("New connection from", c.(*minecraft.Conn).RemoteAddr())
		go arg.handleConn(c.(*minecraft.Conn))
	}
}

// handleConn handles a new incoming minecraft.Conn from the minecraft.Listener passed.
func (arg *Proxy) handleConn(conn *minecraft.Conn) {
	if arg.Config.Server.Whitelist {
		if !arg.WhitelistManager.HasPlayer(conn.IdentityData().DisplayName, conn.IdentityData().XUID) {
			arg.Listener.Disconnect(conn, "You are not whitelisted on this server!")
			return
		}
	}

	serverConn, err := minecraft.Dialer{
		KeepXBLIdentityData: true,
		ClientData:          conn.ClientData(),
		IdentityData:        conn.IdentityData(),
		DownloadResourcePack: func(id uuid.UUID, version string, current int, total int) bool {
			return false
		},
	}.DialTimeout("raknet", arg.Config.Connection.RemoteAddress, time.Second*120)

	if err != nil {
		log.Logger.Errorln("Error in establishing serverConn: ", err)
		arg.Listener.Disconnect(conn, err.Error())
		return
	}

	log.Logger.Debugln("Server connection established for", serverConn.IdentityData().DisplayName)

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
		arg.Listener.Disconnect(conn, "Failed to establish a connection, please try again!")
		serverConn.Close()
		return
	}

	player := player.GetPlayer(conn, serverConn)
	log.Logger.Infoln(player.GetName(), "joined the server")
	player.SendXUIDToAddon()
	arg.UpdatePlayerDetails(player)

	go func() { // client->proxy
		defer arg.Listener.Disconnect(conn, "connection lost")
		defer serverConn.Close()
		for {
			pk, err := conn.ReadPacket()
			if err != nil {
				return
			}

			ok, pk, err := arg.Handlers.HandlePacket(pk, player, "Client")
			if err != nil {
				log.Logger.Errorln(err)
			}

			if ok {
				if err := serverConn.WritePacket(pk); err != nil {
					var disc minecraft.DisconnectError
					if ok := errors.As(err, &disc); ok {
						_ = arg.Listener.Disconnect(conn, disc.Error())
					}
					return
				}
			}
		}
	}()
	go func() { // proxy->server
		defer serverConn.Close()
		defer arg.Listener.Disconnect(conn, "connection lost")
		for {
			pk, err := serverConn.ReadPacket()
			if err != nil {
				var disc minecraft.DisconnectError
				if ok := errors.As(err, &disc); ok {
					_ = arg.Listener.Disconnect(conn, disc.Error())
				}
				return
			}

			ok, pk, err := arg.Handlers.HandlePacket(pk, player, "Server")
			if err != nil {
				log.Logger.Errorln(err)
			}

			if ok {
				if err := conn.WritePacket(pk); err != nil {
					return
				}
			}
		}
	}()
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
		IP:   strings.Split(player.GetSession().Connection.ClientConn.RemoteAddr().String(), ":")[0],
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
