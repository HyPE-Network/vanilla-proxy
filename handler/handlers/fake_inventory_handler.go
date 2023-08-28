package handlers

import (
	"vanilla-proxy/proxy"
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

type OpenInventoryHandlerBoarder struct {
}

func (OpenInventoryHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.ContainerOpen)

	player.GetData().Windows = dataPacket.WindowID

	return true, pk, nil
}

func (OpenInventoryHandlerBoarder) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.ContainerOpen)

	isInside := proxy.ProxyInstance.Worlds.Border.IsPositionInside([]int32{dataPacket.ContainerPosition.X(), dataPacket.ContainerPosition.Y(), dataPacket.ContainerPosition.Z()})
	if !isInside {
		return false, pk, nil
	}

	return true, pk, nil
}
