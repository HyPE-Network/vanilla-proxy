package custom_handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ItemComponentHandler struct {
}

func (ItemComponentHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.ItemComponent)

	proxy.ProxyInstance.Worlds.SetItemComponentEntries(dataPacket.Items)

	return true, dataPacket, nil
}
