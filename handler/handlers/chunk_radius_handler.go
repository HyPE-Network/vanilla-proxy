package handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type UpdateRadiusHandler struct {
	radius int32
}

type RequestRadiusHandler struct {
	radius int32
}

func (uh UpdateRadiusHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.ChunkRadiusUpdated)
	dataPacket.ChunkRadius = uh.radius

	player.GetData().GameData.ChunkRadius = uh.radius

	return true, pk, nil
}

func (rh RequestRadiusHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.RequestChunkRadius)
	dataPacket.ChunkRadius = rh.radius

	return true, pk, nil
}
