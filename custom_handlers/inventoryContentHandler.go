package custom_handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ItemStackRequestHandler struct{}

func (ItemStackRequestHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.ItemStackRequest)

	for _, request := range dataPacket.Requests {
		for _, action := range request.Actions {
			//log.Println("Action:", action)
			switch td := action.(type) {
			case *protocol.PlaceStackRequestAction:
				//log.Println("PlaceStackRequestAction")
				destId := td.Destination.Container.ContainerID
				if destId == protocol.ContainerCraftingInput || destId == protocol.ContainerCombinedHotBarAndInventory {
					//log.Println("Item Being set to Crafting Table")
					// Most likely setting to a container, log the container ID
					copiedDestination := td.Destination
					copiedDestination.StackNetworkID = td.Source.StackNetworkID
					player.SetItemToContainerSlot(copiedDestination)
				}
			case *protocol.TakeStackRequestAction:
				// log.Println("TakeStackRequestAction")
				srcId := td.Source.Container.ContainerID
				if srcId == protocol.ContainerCraftingInput || srcId == protocol.ContainerCombinedHotBarAndInventory {
					// log.Println("Item Being taken from Crafting Table")
					copiedSource := td.Source
					copiedSource.StackNetworkID = 0 // Clear the item from the crafting table
					player.SetItemToContainerSlot(copiedSource)
				}
			}
		}
		player.SetLastItemStackRequestID(request.RequestID)
	}

	return true, dataPacket, nil
}
