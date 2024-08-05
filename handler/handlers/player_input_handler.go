package handlers

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type PlayerInputHandler struct {
}

func (PlayerInputHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.PlayerAuthInput)
	playerData := player.GetData().GameData

	player.SetPlayerLocation(dataPacket.Position)

	// Verify new position
	if proxy.ProxyInstance.Worlds != nil && !proxy.ProxyInstance.Worlds.Border.IsXZInside(int32(dataPacket.Position.X()), int32(dataPacket.Position.Z())) {
		player.SendMessage("§cYou cannot move outside the world!")
		player.PlaySound("note.bass", playerData.PlayerPosition, 1, 1)
		movePlayerPk := &packet.MovePlayer{
			EntityRuntimeID: playerData.EntityRuntimeID,
			Position:        playerData.PlayerPosition,
			Pitch:           playerData.Pitch,
			Yaw:             playerData.Yaw,
			HeadYaw:         playerData.Yaw,
			OnGround:        true,
			Mode:            packet.MoveModeTeleport,
			TeleportCause:   packet.TeleportCauseCommand,
			Tick:            dataPacket.ClientTick,
		}
		player.DataPacket(movePlayerPk)

		return false, pk, nil
	}

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

			if proxy.ProxyInstance.Worlds != nil && !proxy.ProxyInstance.Worlds.Border.IsXZInside(bPos.X(), bPos.Z()) {
				player.SendMessage("§cActions outside the world are prohibited!")
				player.PlaySound("note.bass", playerData.PlayerPosition, 1, 1)
				return false, pk, nil
			}
		}
	}

	return true, pk, nil
}
