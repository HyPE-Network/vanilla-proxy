package custom_handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type PlayerListHandler struct {
}

func (PlayerListHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.PlayerList)

	// Map the Entries to add the XUID to the playerlist
	for i, entry := range dataPacket.Entries {
		xuid, err := proxy.ProxyInstance.PlayerListManager.GetXUIDFromName(entry.Username)
		if err != nil {
			continue
		}
		entry.XUID = xuid
		dataPacket.Entries[i] = entry
	}

	return true, dataPacket, nil
}
