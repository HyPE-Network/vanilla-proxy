package custom_handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type OpenContainerHandler struct{}

func (OpenContainerHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.ContainerOpen)

	//log.Println("Player has opened a container with ID: ", dataPacket.WindowID, dataPacket.ContainerType)
	player.SetOpenContainerWindowID(dataPacket.WindowID)
	player.SetOpenContainerType(dataPacket.ContainerType)

	return true, dataPacket, nil
}

type ContainerCloseHandler struct{}

func (ContainerCloseHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.ContainerClose)

	//log.Println("Player has closed a container", dataPacket)
	player.SetOpenContainerWindowID(0)
	player.SetOpenContainerType(0)
	player.ClearItemsInContainers() // Clear bc container is closed

	return true, dataPacket, nil
}

type AddItemActorHandler struct{}

func (AddItemActorHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.AddItemActor)

	//log.Println("Player has added an item actor", dataPacket)

	return true, dataPacket, nil
}
