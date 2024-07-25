package custom_handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type InventoryContentHandler struct{}

func (InventoryContentHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.InventoryContent)

	// // Map the dataPacket.Content to a slice of strings
	// var content []string
	// for i, item := range dataPacket.Content {
	// 	itemEntry := player.GetItemEntry(item.Stack.ItemType.NetworkID)
	// 	if itemEntry == nil {
	// 		content = append(content, "empty")
	// 		continue
	// 	}
	// 	content = append(content, `Slot: `+strconv.Itoa(i)+`=>`+itemEntry.Name)
	// }

	// //log.Println(player.GetName(), "'s Inventory content:", dataPacket.WindowID, " has been set too:", content)

	return true, dataPacket, nil
}

type InventorySlotHandler struct{}

func (InventorySlotHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.InventorySlot)

	// itemEntry := player.GetItemEntry(dataPacket.NewItem.Stack.ItemType.NetworkID)
	// if itemEntry == nil {
	// 	//log.Println("Inventory slot:", dataPacket.Slot, " has been updated to empty")
	// 	return true, dataPacket, nil
	// } else {
	// 	//log.Println("Inventory slot", dataPacket.Slot, " has been updated to", itemEntry.Name)
	// }

	return true, dataPacket, nil
}

type ItemStackRequestHandler struct{}

func (ItemStackRequestHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.ItemStackRequest)

	for _, request := range dataPacket.Requests {
		for _, action := range request.Actions {
			//log.Println("Action:", action)
			switch td := action.(type) {
			case *protocol.PlaceStackRequestAction:
				//log.Println("PlaceStackRequestAction")
				destId := td.Destination.ContainerID
				if destId == protocol.ContainerCraftingInput || destId == protocol.ContainerCombinedHotBarAndInventory {
					//log.Println("Item Being set to Crafting Table")
					// Most likely setting to a container, log the container ID
					copiedDestination := td.Destination
					copiedDestination.StackNetworkID = td.Source.StackNetworkID
					player.SetItemToContainerSlot(copiedDestination)
				}
			case *protocol.TakeStackRequestAction:
				// log.Println("TakeStackRequestAction")
				srcId := td.Source.ContainerID
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

type ItemStackResponseHandler struct{}

func (ItemStackResponseHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.ItemStackResponse)

	//log.Println(player.GetName(), "has sent item stack responses:", dataPacket.Responses)

	return true, dataPacket, nil
}
