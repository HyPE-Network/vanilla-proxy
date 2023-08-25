package handlers

import (
	"vanilla-proxy/proxy"
	"vanilla-proxy/proxy/player/human"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

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
	}

	dataPacket.Offsets = newPositions

	return true, pk, nil
}
