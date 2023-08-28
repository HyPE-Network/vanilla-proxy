package handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type MalformedHandler struct {
}

func (MalformedHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.PacketViolationWarning)

	log.Logger.Errorln(player.GetName(), "> Malformed", dataPacket)

	return true, pk, nil
}
