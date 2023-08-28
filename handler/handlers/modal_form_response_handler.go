package handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/form"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ModalFormResponseHandler struct {
}

func (ModalFormResponseHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.ModalFormResponse)

	window, ok := player.GetData().Forms[dataPacket.FormID]
	if ok {
		closed, formData, err := form.GetResponseData(dataPacket, window.GetType())
		if err != nil {
			log.Logger.Errorln(err)
			return false, pk, nil
		}

		window.Do(closed, formData)
		delete(player.GetData().Forms, dataPacket.FormID)
	}

	return false, pk, nil
}
