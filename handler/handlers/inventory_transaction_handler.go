package handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type InventoryTransactionHandler struct {
}

func (InventoryTransactionHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.InventoryTransaction)

	if len(dataPacket.Actions) > 0 && dataPacket.Actions[0].SourceType == 99999 { // cheats
		return false, pk, nil
	}

	return true, pk, nil
}
