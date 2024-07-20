package human

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/data"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/scoreboard"
	"github.com/HyPE-Network/vanilla-proxy/proxy/session"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type Human interface {
	GetName() string
	GetData() *data.PlayerData
	GetSession() *session.Session

	SendMessage(string)
	SendPopup(string)
	SendTip(string)

	HasScoreboard() bool
	SendScoreboard(*scoreboard.Scoreboard)
	RemoveScoreboard()

	Transfer(string, uint16)

	Kick()
	Close(string)

	Distance(mgl32.Vec3) float64
	DistanceXYZ(float32, float32, float32) float64

	SendUpdateBlock(protocol.BlockPos, uint32)
	SendAirUpdate(protocol.BlockPos)

	PlaySound(string, mgl32.Vec3, float32, float32)

	InOverworld() bool
	InNether() bool
	InEnd() bool
	GetDimension() int32
	GetWorldName() string

	GetPing() int64
	GetSessionTime() int64

	DataPacket(packet.Packet)
	DataPacketToServer(packet.Packet)

	SendXUIDToAddon()

	IsOP() bool

	PlayerPermissions() byte
	CommandPermissions() byte

	SetBDSAvailableCommands(*packet.AvailableCommands)
}
