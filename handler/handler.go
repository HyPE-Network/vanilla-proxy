package handler

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type PacketHandler interface {
	Handle(packet.Packet, human.Human) (bool, packet.Packet, error) // send packet to destination, packet, error
}

type HandlerManager interface {
	HandlePacket(packet.Packet, human.Human, string) (bool, packet.Packet, error)
}
