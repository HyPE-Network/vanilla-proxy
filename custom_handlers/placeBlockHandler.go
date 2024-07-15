package custom_handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy/block/cube"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type PlaceBlockHandler struct {
}

func (PlaceBlockHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.InventoryTransaction)

	switch td := dataPacket.TransactionData.(type) {
	case *protocol.UseItemTransactionData:
		if td.ActionType == protocol.UseItemActionClickBlock {
			if td.HeldItem.Stack.BlockRuntimeID != 0 {
				if td.HeldItem.Stack.ItemType.NetworkID == 49 && player.InNether() { // obsidian block in nether
					pos := cube.Side(td.BlockPosition, td.BlockFace)

					player.SendMessage("§cYou can't place obsidian in this world!")
					player.SendAirUpdate(protocol.BlockPos{pos.X(), pos.Y(), pos.Z()})
					return false, pk, nil
				}
			}
		}
	}

	return true, pk, nil
}
