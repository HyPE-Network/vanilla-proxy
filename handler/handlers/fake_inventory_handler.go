package handlers

import (
	"vanilla-proxy/proxy/player/human"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type CloseInventoryHandler struct {
}

func (CloseInventoryHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.ContainerClose)

	if dataPacket.WindowID == player.GetData().Windows && player.GetData().FakeChestOpen {
		player.GetData().FakeChestOpen = false
		player.SendAirUpdate(player.GetData().FakeChestPos)
		player.GetData().Windows--
		player.GetData().FakeChestPos = protocol.BlockPos{}
		player.DataPacket(dataPacket)
		return false, pk, nil
	}

	return true, pk, nil
}

type OpenInventoryHandler struct {
}

func (OpenInventoryHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.ContainerOpen)

	player.GetData().Windows = dataPacket.WindowID

	return true, pk, nil
}
