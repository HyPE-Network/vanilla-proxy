package human

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/data"
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

	SetPlayerLocation(mgl32.Vec3)

	GetPing() int64
	GetSessionTime() int64

	DataPacket(packet.Packet)
	DataPacketToServer(packet.Packet)

	SendXUIDToAddon()

	IsOP() bool

	PlayerPermissions() byte
	CommandPermissions() byte

	SetOpenContainerWindowID(windowId byte)
	SetOpenContainerType(containerType byte)
	SetLastItemStackRequestID(id int32)
	GetNextItemStackRequestID() int32
	SetItemToContainerSlot(slotInfo protocol.StackRequestSlotInfo)
	ClearItemsInContainers()
	GetItemFromContainerSlot(containerID byte, slot byte) protocol.StackRequestSlotInfo
	GetCursorItem() protocol.StackRequestSlotInfo
}
