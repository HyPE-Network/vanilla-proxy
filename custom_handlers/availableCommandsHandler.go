package custom_handlers

import (
	"log"
	"strings"

	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type AvailableCommandsHandler struct {
}

type CommandRequestHandler struct {
}

func (AvailableCommandsHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.AvailableCommands)

	// Remove `/me`
	for i, command := range dataPacket.Commands {
		if command.Name == "me" {
			dataPacket.Commands = append(dataPacket.Commands[:i], dataPacket.Commands[i+1:]...)
			break
		}
	}

	// Remove `/tell`, `/w`, and `/msg`
	for i, command := range dataPacket.Commands {
		if command.Name == "tell" {
			dataPacket.Commands = append(dataPacket.Commands[:i], dataPacket.Commands[i+1:]...)
			break
		}
	}

	return true, dataPacket, nil
}

func (CommandRequestHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.CommandRequest)
	playerData := player.GetData().GameData

	var command = strings.ToLower(strings.Split(dataPacket.CommandLine[1:], " ")[0])

	if command == "me" || command == "tell" || command == "w" || command == "msg" {
		player.SendMessage("§cThe command /" + command + " is disabled!")
		player.PlaySound("note.bass", playerData.PlayerPosition, 1, 1)
		return false, pk, nil
	}

	// Command should be a custom `-` command
	log.Println("Command: ", "-"+strings.TrimPrefix(dataPacket.CommandLine, "/"))
	textPk := &packet.Text{
		TextType:         packet.TextTypeChat,
		NeedsTranslation: false,
		SourceName:       player.GetName(),
		Message:          "-" + strings.TrimPrefix(dataPacket.CommandLine, "/"),
		Parameters:       []string{},
		XUID:             player.GetSession().IdentityData.XUID,
		PlatformChatID:   "",
	}
	player.DataPacketToServer(textPk)

	return true, dataPacket, nil
}
