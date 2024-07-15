package custom_handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type DisconnectHandler struct {
}

func (DisconnectHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.Disconnect)
	log.Logger.Debugln("Player has been disconnected with reason: ", dataPacket.Message)

	player.Transfer("play.pokebedrock.com", 19132)

	return true, dataPacket, nil
}
