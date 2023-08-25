package handlers

import (
	"vanilla-proxy/proxy"
	"vanilla-proxy/proxy/player/human"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type SubChunkHandler struct {
}

func (SubChunkHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.SubChunk)

	var newSubChunks = make([]protocol.SubChunkEntry, 0)
	for _, subChunk := range dataPacket.SubChunkEntries {
		isInside := proxy.ProxyInstance.Worlds.Border.IsPositionInside([]int32{dataPacket.Position.X() + int32(subChunk.Offset[0]), dataPacket.Position.Z() + int32(subChunk.Offset[2])})
		if isInside {
			newSubChunks = append(newSubChunks, subChunk)
		}
	}

	player.GetData().GameData.Dimension = dataPacket.Dimension

	dataPacket.SubChunkEntries = newSubChunks

	return true, pk, nil
}
