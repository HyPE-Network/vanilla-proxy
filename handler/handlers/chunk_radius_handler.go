package handlers

import (
	"vanilla-proxy/proxy/player/human"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

var radius int32 = 8

type UpdateRadiusHandler struct {
}

type RequestRadiusHandler struct {
}

func (UpdateRadiusHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.ChunkRadiusUpdated)
	dataPacket.ChunkRadius = radius

	player.GetData().GameData.ChunkRadius = radius

	return true, pk, nil
}

func (RequestRadiusHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.RequestChunkRadius)
	dataPacket.ChunkRadius = radius

	return true, pk, nil
}
