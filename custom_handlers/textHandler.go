package custom_handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type TextHandler struct {
}

func (TextHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.Text)

	//log.Println("TextHandler: ", dataPacket.Message)

	return true, dataPacket, nil
}
