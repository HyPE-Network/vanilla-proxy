package human

import (
	"vanilla-proxy/proxy/inventory"
	"vanilla-proxy/proxy/player/data"
	"vanilla-proxy/proxy/player/scoreboard"
	"vanilla-proxy/proxy/session"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft"
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

	InOverworld() bool
	InNether() bool
	InEnd() bool
	GetDimension() int32
	GetWorldName() string

	GetPing() int64
	GetSessionTime() int64

	DataPacket(packet.Packet)
	DataPacketToServer(packet.Packet)

	SendInventory(inventory.Inventory)
}

type HumanManager interface {
	AddPlayer(*minecraft.Conn, *minecraft.Conn) Human
	DeletePlayer(Human)
	DeleteAll()
	GetPlayer(string) Human
	GetPlayerExact(string) Human
	PlayerList() map[string]Human
	PlayersCount() int
	IsOnline(string) bool
}
