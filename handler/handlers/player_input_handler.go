package handlers

import (
	"vanilla-proxy/proxy"
	"vanilla-proxy/proxy/player/human"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type PlayerInputHandler struct {
}

func (PlayerInputHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.PlayerAuthInput)
	player.GetData().GameData.Pitch, player.GetData().GameData.Yaw = dataPacket.Pitch, dataPacket.Yaw
	player.GetData().GameData.PlayerPosition = dataPacket.Position

	if len(dataPacket.BlockActions) > 0 {
		for i, ba := range dataPacket.BlockActions {
			if ba.Action == protocol.PlayerActionCrackBreak { // continue break
				continue
			}

			bPos := ba.BlockPos
			if ba.Action == protocol.PlayerActionStopBreak && (i == 0 || i == 1) && len(dataPacket.BlockActions) > 1 { // break block action contains [0 0 0] position
				if i == 0 {
					bPos = dataPacket.BlockActions[i+1].BlockPos
				} else {
					bPos = dataPacket.BlockActions[i-1].BlockPos
				}
			}

			if proxy.ProxyInstance.Worlds.BoarderEnabled && !proxy.ProxyInstance.Worlds.Border.IsXZInside(bPos.X(), bPos.Z()) {
				player.SendMessage("§cActions outside the world are prohibited!")
				return false, pk, nil
			}
		}
	}

	return true, pk, nil
}
