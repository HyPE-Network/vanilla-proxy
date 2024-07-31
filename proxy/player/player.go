package player

import (
	"math"

	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/data"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/HyPE-Network/vanilla-proxy/proxy/session"
	"github.com/HyPE-Network/vanilla-proxy/proxy/world"
	"github.com/HyPE-Network/vanilla-proxy/utils"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Player struct {
	Name       string
	Session    *session.Session
	PlayerData *data.PlayerData
}

// Creates a new player instance from a server conn
func NewPlayer(conn *minecraft.Conn, session *session.Session) *Player {
	parsedGameData := conn.GameData()
	// Remove Items & Blocks from `parsedGameData` to reduce the size of the struct
	parsedGameData.Items = nil
	parsedGameData.CustomBlocks = nil

	return &Player{
		Name:    conn.IdentityData().DisplayName,
		Session: session,
		PlayerData: &data.PlayerData{
			GameData:         parsedGameData,
			StartSessionTime: utils.GetTimestamp(),
			Authorized:       false,
		},
	}
}

// Gets a player from a connection
func GetPlayer(conn *minecraft.Conn, serverConn *minecraft.Conn) human.Human {
	ab := session.NewBridge(conn, serverConn)
	newSession := session.NewSession(conn, ab)
	var pl human.Human = NewPlayer(conn, newSession)
	return pl
}

func (player *Player) GetName() string {
	return player.Name
}

func (player *Player) GetData() *data.PlayerData {
	return player.PlayerData
}

func (player *Player) GetSession() *session.Session {
	return player.Session
}

func (player *Player) SendMessage(message string) {
	player.textPacket(message, packet.TextTypeRaw)
}

func (player *Player) SendPopup(message string) {
	player.textPacket(message, packet.TextTypePopup)
}

func (player *Player) SendTip(message string) {
	player.textPacket(message, packet.TextTypeTip)
}

func (player *Player) SendSound(sound string, volume float32, pitch float32) {
	pk := &packet.PlaySound{
		SoundName: sound,
		Position:  player.PlayerData.GameData.PlayerPosition,
		Volume:    volume,
		Pitch:     pitch,
	}

	player.DataPacket(pk)
}

func (player *Player) Transfer(address string, port uint16) {
	pk := &packet.Transfer{
		Address: address,
		Port:    port,
	}

	player.DataPacket(pk)
	log.Logger.Debugln("Player", player.Name, "transferred to", address, port)
}

func (player *Player) Kick() {
	player.Close("")
}

func (player *Player) Close(message string) {
	pk := &packet.Disconnect{
		HideDisconnectionScreen: false,
		Message:                 message,
	}
	player.DataPacket(pk)
}

func (player *Player) Distance(pos mgl32.Vec3) float64 {
	return math.Sqrt(player.DistanceXYZSquared(pos.X(), pos.Y(), pos.Z()))
}

func (player *Player) DistanceXYZ(x float32, y float32, z float32) float64 {
	return math.Sqrt(player.DistanceXYZSquared(x, y, z))
}

func (player *Player) DistanceSquared(target *Player) float64 {
	pos := target.PlayerData.GameData.PlayerPosition
	return player.DistanceXYZSquared(pos.X(), pos.Y(), pos.Z())
}

func (player *Player) DistanceXYZSquared(x float32, y float32, z float32) float64 {
	pos := player.PlayerData.GameData.PlayerPosition
	ex := pos.X() - x
	ey := pos.Y() - y
	ez := pos.Z() - z
	return float64(ex*ex + ey*ey + ez*ez)
}

func (player *Player) SendAirUpdate(pos protocol.BlockPos) {
	player.SendUpdateBlock(pos, world.AirRID)
}

func (player *Player) SendUpdateBlock(pos protocol.BlockPos, rid uint32) {
	pk := &packet.UpdateBlock{
		Position:          pos,
		NewBlockRuntimeID: rid,
		Flags:             0,
		Layer:             0,
	}

	player.DataPacket(pk)
}

func (player *Player) SetPlayerLocation(pos mgl32.Vec3) {
	player.PlayerData.LastUpdatedLocation = pos
}

func (player *Player) InOverworld() bool {
	return player.PlayerData.GameData.Dimension == packet.DimensionOverworld
}

func (player *Player) InNether() bool {
	return player.PlayerData.GameData.Dimension == packet.DimensionNether
}

func (player *Player) InEnd() bool {
	return player.PlayerData.GameData.Dimension == packet.DimensionEnd
}

func (player *Player) GetDimension() int32 {
	return player.PlayerData.GameData.Dimension
}

func (player *Player) GetWorldName() string {
	return player.PlayerData.GameData.WorldName
}

func (player *Player) GetPing() int64 {
	return player.Session.Connection.ClientConn.Latency().Milliseconds()
}

func (player *Player) GetSessionTime() int64 {
	return utils.GetTimestamp() - player.PlayerData.StartSessionTime
}

func (player *Player) textPacket(message string, textType byte) {
	pk := &packet.Text{
		TextType:         textType,
		NeedsTranslation: false,
		SourceName:       "",
		Message:          message,
		Parameters:       []string{},
		XUID:             "",
		PlatformChatID:   "",
	}

	player.DataPacket(pk)
}

func (player *Player) DataPacket(pk packet.Packet) {
	if err := player.Session.Connection.ClientConn.WritePacket(pk); err != nil {
		log.Logger.Errorln(err)
	}
}

func (player *Player) DataPacketToServer(pk packet.Packet) {
	if err := player.Session.Connection.ServerConn.WritePacket(pk); err != nil {
		log.Logger.Errorln(err)
	}
}

func (player *Player) PlaySound(soundName string, pos mgl32.Vec3, volume float32, pitch float32) {
	pk := &packet.PlaySound{
		SoundName: soundName,
		Position:  pos,
		Volume:    volume,
		Pitch:     pitch,
	}

	player.DataPacket(pk)
}

func (player *Player) SendXUIDToAddon() {
	playerXuid := player.GetSession().IdentityData.XUID
	playerXuidTextPacket := &packet.Text{
		TextType:         packet.TextTypeChat,
		NeedsTranslation: false,
		SourceName:       player.GetName(),
		Message:          "[PROXY_SYSTEM] XUID=" + playerXuid,
		Parameters:       nil,
		XUID:             playerXuid,
		PlatformChatID:   "",
	}
	player.DataPacketToServer(playerXuidTextPacket)
}

// IsOp checks if the player is an operator on the server.
func (player *Player) IsOP() bool {
	config := utils.ReadConfig()
	return utils.StringInSlice(player.GetName(), config.Server.Ops)
}

// PlayerPermissions is the permission level of the player as it shows up in the player list built up using the PlayerList packet.
func (player *Player) PlayerPermissions() byte {
	if player.IsOP() {
		return packet.PermissionLevelOperator
	} else {
		return packet.PermissionLevelMember
	}
}

// CommandPermissions is a set of permissions that specify what commands a player is allowed to execute.
func (player *Player) CommandPermissions() byte {
	if player.IsOP() {
		return packet.CommandPermissionLevelHost
	} else {
		return packet.CommandPermissionLevelNormal
	}
}

// GetItemEntry returns the item entry of an item with the specified network ID. If the item is not found, nil is returned.
func (player *Player) GetItemEntry(networkID int32) *protocol.ItemEntry {
	items := player.GetData().GameData.Items
	idx := slices.IndexFunc(items, func(item protocol.ItemEntry) bool {
		return item.RuntimeID == int16(networkID)
	})
	if idx == -1 {
		// Unknown item?
		return nil
	}
	item := items[idx]
	return &item
}

// SetOpenContainerWindowID sets the ID of the window that is currently open for the player.
func (player *Player) SetOpenContainerWindowID(windowId byte) {
	player.PlayerData.OpenContainerWindowId = windowId
}

func (player *Player) SetOpenContainerType(containerType byte) {
	player.PlayerData.OpenContainerType = containerType
}

// SetLastItemStackRequestID sets the last item stack request ID that was sent by the player.
func (player *Player) SetLastItemStackRequestID(id int32) {
	player.PlayerData.LastItemStackRequestID = id
}

// GetNextItemStackRequestID returns the next item stack request ID that can be used by the player.
func (player *Player) GetNextItemStackRequestID() int32 {
	if player.PlayerData.LastItemStackRequestID == math.MaxInt32 {
		player.PlayerData.LastItemStackRequestID = 0
	}
	player.PlayerData.LastItemStackRequestID -= 2
	return player.PlayerData.LastItemStackRequestID
}

// SetItemToContainerSlot sets the amount of items that are in the container slot that the player has put in.
func (player *Player) SetItemToContainerSlot(slotInfo protocol.StackRequestSlotInfo) {
	// Find if the slot is already in this list, if so update, else append
	for i, slot := range player.PlayerData.ItemsInContainers {
		if slot.ContainerID == slotInfo.ContainerID && slot.Slot == slotInfo.Slot {
			player.PlayerData.ItemsInContainers[i] = slotInfo
			return
		}
	}
	player.PlayerData.ItemsInContainers = append(player.PlayerData.ItemsInContainers, slotInfo)
}

func (player *Player) ClearItemsInContainers() {
	player.PlayerData.ItemsInContainers = nil
}

// GetItemsInContainerSlot returns the amount of items that are in the container slot that the player has put in.
func (player *Player) GetItemFromContainerSlot(containerID byte, slot byte) protocol.StackRequestSlotInfo {
	for _, slotInfo := range player.PlayerData.ItemsInContainers {
		if slotInfo.ContainerID == containerID && slotInfo.Slot == slot {
			return slotInfo
		}
	}
	return protocol.StackRequestSlotInfo{
		ContainerID:    containerID,
		Slot:           slot,
		StackNetworkID: 0, // Empty
	}
}

// GetCursorItem returns the item that is currently in the cursor of the player.
func (player *Player) GetCursorItem() protocol.StackRequestSlotInfo {
	return player.GetItemFromContainerSlot(protocol.ContainerCursor, 0)
}
