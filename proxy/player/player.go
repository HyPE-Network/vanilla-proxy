package player

import (
	"math"
	"slices"

	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/data"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/form"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/scoreboard"
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
	return &Player{
		Name:    conn.IdentityData().DisplayName,
		Session: session,
		PlayerData: &data.PlayerData{
			GameData:         conn.GameData(),
			Forms:            make(map[uint32]form.Form),
			BrokenBlocks:     make(map[protocol.BlockPos]uint32),
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

func (player *Player) HasScoreboard() bool {
	return player.PlayerData.CurrentScoreboard.Load() != ""
}

func (player *Player) SendScoreboard(sb *scoreboard.Scoreboard) {
	currentName, currentLines := player.PlayerData.CurrentScoreboard.Load(), player.PlayerData.CurrentLines.Load()

	if currentName != sb.Name() {
		player.RemoveScoreboard()
		player.DataPacket(&packet.SetDisplayObjective{
			DisplaySlot:   "sidebar",
			ObjectiveName: sb.Name(),
			DisplayName:   sb.Name(),
			CriteriaName:  "dummy",
		})
		player.PlayerData.CurrentScoreboard.Store(sb.Name())
		player.PlayerData.CurrentLines.Store(append([]string(nil), sb.Lines()...))
	} else {
		// Remove all current lines from the scoreboard. We can't replace them without removing them.
		pk := &packet.SetScore{ActionType: packet.ScoreboardActionRemove}
		for i := range currentLines {
			pk.Entries = append(pk.Entries, protocol.ScoreboardEntry{
				EntryID:       int64(i),
				ObjectiveName: currentName,
				Score:         int32(i),
			})
		}
		if len(pk.Entries) > 0 {
			player.DataPacket(pk)
		}
	}
	pk := &packet.SetScore{ActionType: packet.ScoreboardActionModify}
	for k, line := range sb.Lines() {
		pk.Entries = append(pk.Entries, protocol.ScoreboardEntry{
			EntryID:       int64(k),
			ObjectiveName: sb.Name(),
			Score:         int32(k),
			IdentityType:  protocol.ScoreboardIdentityFakePlayer,
			DisplayName:   line,
		})
	}
	if len(pk.Entries) > 0 {
		player.DataPacket(pk)
	}
}

func (player *Player) RemoveScoreboard() {
	player.DataPacket(&packet.RemoveObjective{ObjectiveName: player.PlayerData.CurrentScoreboard.Load()})
	player.PlayerData.CurrentScoreboard.Store("")
	player.PlayerData.CurrentLines.Store([]string{})
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
	if player.GetData().Closed {
		return
	}

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

// SetBDSAvailableCommands sets the AvailableCommands packet that is sent to the player when they join the server.
func (player *Player) SetBDSAvailableCommands(pk *packet.AvailableCommands) {
	player.PlayerData.BDSAvailableCommands = *pk
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

func (player *Player) GetItemComponentEntry(name string) *protocol.ItemComponentEntry {
	for _, entry := range player.PlayerData.ItemComponentEntries {
		if entry.Name == name {
			return &entry
		}
	}
	return nil
}

func (player *Player) GetItemComponentEntries() []protocol.ItemComponentEntry {
	return player.PlayerData.ItemComponentEntries
}

func (player *Player) AddItemComponentEntry(entry *protocol.ItemComponentEntry) {
	player.PlayerData.ItemComponentEntries = append(player.PlayerData.ItemComponentEntries, *entry)
}

func (player *Player) RemoveItemComponentEntry(entry *protocol.ItemComponentEntry) {
	idx := slices.IndexFunc(player.PlayerData.ItemComponentEntries, func(e protocol.ItemComponentEntry) bool {
		return e.Name == entry.Name
	})
	if idx == -1 {
		return
	}
	player.PlayerData.ItemComponentEntries = append(player.PlayerData.ItemComponentEntries[:idx], player.PlayerData.ItemComponentEntries[idx+1:]...)
}

func (player *Player) SetItemComponentEntries(entries []protocol.ItemComponentEntry) {
	player.PlayerData.ItemComponentEntries = entries
}
