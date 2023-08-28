package handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
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

type SubChunkHandler struct {
}

func (SubChunkHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.SubChunk)

	player.GetData().GameData.Dimension = dataPacket.Dimension

	return true, pk, nil
}

type SubChunkHandlerBoarder struct {
}

func (SubChunkHandlerBoarder) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.SubChunk)

	var newSubChunks = make([]protocol.SubChunkEntry, 0)
	for _, subChunk := range dataPacket.SubChunkEntries {
		isInside := proxy.ProxyInstance.Worlds.Border.IsPositionInside([]int32{dataPacket.Position.X() + int32(subChunk.Offset[0]), dataPacket.Position.Z() + int32(subChunk.Offset[2])})
		if isInside {
			newSubChunks = append(newSubChunks, subChunk)
		}
	}

	dataPacket.SubChunkEntries = newSubChunks

	player.GetData().GameData.Dimension = dataPacket.Dimension

	return true, pk, nil
}

type SubChunkRequestHandler struct {
}

func (SubChunkRequestHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.SubChunkRequest)

	var newPositions = make([]protocol.SubChunkOffset, 0)
	for _, pos := range dataPacket.Offsets {
		isInside := proxy.ProxyInstance.Worlds.Border.IsPositionInside([]int32{dataPacket.Position.X() + int32(pos[0]), dataPacket.Position.Z() + int32(pos[2])})
		if isInside {
			newPositions = append(newPositions, pos)
		}

		dataPacket.Offsets = newPositions
	}

	return true, pk, nil
}
