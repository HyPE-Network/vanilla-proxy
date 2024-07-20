package custom_handlers

import (
	"strings"

	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type AvailableCommandsHandler struct {
}

type CommandRequestHandler struct {
}

// RemoveCommands removes a list of specified commands from the list of available commands.
func RemoveCommands(commands []protocol.Command, remove []string) []protocol.Command {
	removeMap := make(map[string]bool)
	for _, cmd := range remove {
		removeMap[cmd] = true
	}

	filteredCommands := make([]protocol.Command, 0, len(commands))
	for _, command := range commands {
		if !removeMap[command.Name] {
			filteredCommands = append(filteredCommands, command)
		}
	}

	return filteredCommands
}

func (AvailableCommandsHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.AvailableCommands)

	// Define an array of command names to remove
	commandsToRemove := []string{"me", "tell"}
	dataPacket.Commands = RemoveCommands(dataPacket.Commands, commandsToRemove)

	player.SetBDSAvailableCommands(dataPacket)

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

	minecraftCommands := player.GetData().BDSAvailableCommands.Commands
	// Check if {command} is a name inside {minecraftCommands}
	for _, cmd := range minecraftCommands {
		if cmd.Name == command {
			return true, pk, nil
		}
	}

	// Command should be a custom `-` command
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

	return false, dataPacket, nil
}
