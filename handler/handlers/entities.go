package handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy"
	"github.com/HyPE-Network/vanilla-proxy/proxy/entity"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type AddActorHandler struct{}

func (ah AddActorHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.AddActor)

	proxy.ProxyInstance.Entities.SetEntity(dataPacket.TargetActorID, entity.EntityData{
		TypeID:    dataPacket.ActorType,
		RuntimeID: dataPacket.TargetRuntimeID,
	})

	return true, pk, nil
}

type RemoveActorHandler struct{}

func (ah RemoveActorHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.RemoveActor)

	proxy.ProxyInstance.Entities.RemoveEntity(dataPacket.EntityUniqueID)

	return true, pk, nil
}
