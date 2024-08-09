package proxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/HyPE-Network/vanilla-proxy/handler"
	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/HyPE-Network/vanilla-proxy/math"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/HyPE-Network/vanilla-proxy/proxy/playerlist"
	"github.com/HyPE-Network/vanilla-proxy/proxy/whitelist"
	"github.com/HyPE-Network/vanilla-proxy/proxy/world"
	"github.com/HyPE-Network/vanilla-proxy/utils"
	"github.com/google/uuid"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/resource"

	"github.com/sandertv/gophertunnel/minecraft"
)

var ProxyInstance *Proxy

type Proxy struct {
	Worlds            *world.Worlds
	Config            utils.Config
	Handlers          handler.HandlerManager
	Listener          *minecraft.Listener
	WhitelistManager  *whitelist.WhitelistManager
	PlayerListManager *playerlist.PlayerlistManager
}

func New(config utils.Config) *Proxy {
	playerListManager := playerlist.Init()

	Proxy := &Proxy{
		Config:            config,
		PlayerListManager: playerListManager,
	}

	if config.WorldBorder.Enabled {
		Proxy.Worlds = world.Init(math.NewArea2(config.WorldBorder.MinX, config.WorldBorder.MinZ, config.WorldBorder.MaxX, config.WorldBorder.MaxZ))
	}

	if config.Server.Whitelist {
		Proxy.WhitelistManager = whitelist.Init()
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

	// Loop through all the pack paths and append each pack to the slice
	for _, path := range arg.Config.Resources.PackPaths {
		resourcePack, err := resource.ReadPath(path)
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
			// The listener closed, so we should restart it.
			log.Logger.Errorln(err)
			utils.SendStaffAlertToDiscord("Proxy Listener Closed", err.Error(), 16711680, []map[string]interface{}{
				{
					"name":   "Connection From",
					"value":  c.RemoteAddr().String(),
					"inline": true,
				},
			})
			c.Close()
			arg.Start(h)
			return nil // Should return error, but we want to restart listener
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

	clientData := arg.PlayerListManager.GetConnClientData(conn)
	identityData := arg.PlayerListManager.GetConnIdentityData(conn)

	serverConn, err := minecraft.Dialer{
		KeepXBLIdentityData: true,
		ClientData:          clientData,
		IdentityData:        identityData,
		DownloadResourcePack: func(id uuid.UUID, version string, current int, total int) bool {
			return false
		},
	}.DialTimeout("raknet", arg.Config.Connection.RemoteAddress, time.Second*120)

	if err != nil {
		log.Logger.Errorln("Error in establishing serverConn: ", err)
		arg.Listener.Disconnect(conn, strings.Split(err.Error(), ": ")[1])
		return
	}

	log.Logger.Debugln("Server connection established for", serverConn.IdentityData().DisplayName)

	gameData := serverConn.GameData()
	gameData.WorldSeed = 0
	gameData.ClientSideGeneration = false
	arg.Worlds.SetItems(gameData.Items)
	arg.Worlds.SetCustomBlocks(gameData.CustomBlocks)

	var success = true
	var g sync.WaitGroup
	g.Add(2)
	go func() {
		if err := conn.StartGame(gameData); err != nil {
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
		defer arg.DisconnectPlayer(player, "Connection closed")
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
						arg.DisconnectPlayer(player, disc.Error())
					}
					log.Logger.Errorln(err)
					return
				}
			}
		}
	}()
	go func() { // proxy->server
		defer arg.DisconnectPlayer(player, "Connection closed")
		for {
			pk, err := serverConn.ReadPacket()
			if err != nil {
				var disc minecraft.DisconnectError
				if ok := errors.As(err, &disc); ok {
					arg.DisconnectPlayer(player, disc.Error())
				}
				log.Logger.Errorln("Failed to read Packet from Server", err)
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

// DisconnectPlayer disconnects a player from the proxy.
func (arg *Proxy) DisconnectPlayer(player human.Human, message string) {
	// Send close container packet
	if player.IsBeingDisconnected() {
		return // Player is already being disconnected, ignore this call
	}
	player.SetDisconnected(true)

	openContainerId := player.GetData().OpenContainerWindowId
	itemInContainers := player.GetData().ItemsInContainers

	if openContainerId != 0 {
		log.Logger.Println("Player has open container while disconnecting, *prob trying to dupe*")

		utils.SendStaffAlertToDiscord("Disconnecting With Open Container",
			"A Player Has disconnected with an open container, please investigate!",
			16711680,
			[]map[string]interface{}{
				{
					"name":   "Player Name",
					"value":  player.GetName(),
					"inline": true,
				},
				{
					"name":   "Container Type",
					"value":  openContainerId,
					"inline": true,
				},
				{
					"name":   "Player Location",
					"value":  player.GetData().LastUpdatedLocation,
					"inline": true,
				},
			})

		// Send Item Stack Requests to clear the container
		// Send Item Request to clear container id 13 (crafting table)
		// By sending from slot 32->40 (9 crafting slots) to `false` (throw on ground)
		request := protocol.ItemStackRequest{
			RequestID: player.GetNextItemStackRequestID(),
			Actions:   []protocol.StackRequestAction{},
		}
		// Loop through players container slots
		for _, slotInfo := range itemInContainers {
			action := &protocol.DropStackRequestAction{}
			action.Source = slotInfo
			action.Count = 64
			action.Randomly = false
			request.Actions = append(request.Actions, action)
		}
		pk := &packet.ItemStackRequest{
			Requests: []protocol.ItemStackRequest{request},
		}
		log.Logger.Debugln("Sending ItemStackRequest to clear container:")
		player.DataPacketToServer(pk)

		player.SetOpenContainerWindowID(0)
		player.SetOpenContainerType(0)

		// Sleep for 2 seconds to allow the packets to be sent
		time.Sleep(time.Second * 4)
	}

	cursorItem := player.GetItemFromContainerSlot(protocol.ContainerCombinedHotBarAndInventory, 0)
	if cursorItem.StackNetworkID != 0 {
		// Player left with a item in ContainerCombinedHotBarAndInventory
		utils.SendStaffAlertToDiscord("Disconnecting With Item",
			"A Player Has disconnected with a item in ContainerCombinedHotBarAndInventory, please investigate!",
			16711680,
			[]map[string]interface{}{
				{
					"name":   "Player Name",
					"value":  player.GetName(),
					"inline": true,
				},
				{
					"name":   "Stack Network ID",
					"value":  cursorItem.StackNetworkID,
					"inline": true,
				},
				{
					"name":   "Player Location",
					"value":  player.GetData().LastUpdatedLocation,
					"inline": true,
				},
			})

	}
	log.Logger.Debugln("Disconnecting player:", player.GetName(), "with reason:", message)

	// Disconnect
	player.GetSession().Connection.ServerConn.Close()
	arg.Listener.Disconnect(player.GetSession().Connection.ClientConn, message)
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
