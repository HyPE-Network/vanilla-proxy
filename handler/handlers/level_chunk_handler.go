package handlers

import (
	"vanilla-proxy/proxy"
	"vanilla-proxy/proxy/player/human"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type LevelChunkHandler struct {
}

func (LevelChunkHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.LevelChunk)

	isInside := proxy.ProxyInstance.Worlds.Border.IsPositionInside(dataPacket.Position[:])
	if !isInside {
		return false, pk, nil
	}

	return true, pk, nil
}
