package custom_handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type SignEditHandler struct {
}

func (SignEditHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.BlockActorData)

	//log.Logger.Debugln("Update Block: ", dataPacket.Position, dataPacket.NBTData)

	return true, dataPacket, nil
}
