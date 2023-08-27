package handlers

import (
	"strings"
	"vanilla-proxy/log"
	"vanilla-proxy/proxy"
	"vanilla-proxy/proxy/player/human"
	"vanilla-proxy/utils/color"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type AvailableCommandsHandler struct {
}

type CommandRequestHandler struct {
}

func (AvailableCommandsHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.AvailableCommands)

	for executor, command := range proxy.ProxyInstance.CommandManager.Commands {
		if executor.ForPlayer() {
			dataPacket.Commands = append(dataPacket.Commands, command)
		}
	}

	// todo: hack
	helpCommand := protocol.Command{
		Name:            "help",
		Description:     "Help info",
		PermissionLevel: protocol.CommandEnumConstraintCheatsEnabled,
		Overloads:       []protocol.CommandOverload{},
	}

	dataPacket.Commands = append(dataPacket.Commands, helpCommand)

	return true, pk, nil
}

func (CommandRequestHandler) Handle(pk packet.Packet, player human.Human) (bool, packet.Packet, error) {
	dataPacket := pk.(*packet.CommandRequest)

	var command = strings.ToLower(strings.Split(dataPacket.CommandLine[1:], " ")[0])

	for executor, cmd := range proxy.ProxyInstance.CommandManager.Commands {
		if cmd.Name == command {
			args := formatArgs(dataPacket.CommandLine[1:])

			err := executor.Execute(player, args)
			if err != nil {
				player.SendMessage(color.Red + "Error in command handler!")
				log.Logger.Errorln(err)
				return true, pk, nil
			}

			return false, pk, nil
		}
	}

	return true, pk, nil
}

func formatArgs(command string) []string {
	var args []string
	command = strings.TrimSpace(command)
	command += " "

	arg := ""
	bigArg := false
	for _, value := range command {
		if value == ' ' && !bigArg {
			if arg != "" && arg != " " {
				args = append(args, arg)
			}
			arg = ""
		} else if value == '"' && !bigArg {
			bigArg = true
		} else if value == '"' && bigArg {
			bigArg = false
			if arg != "" && arg != " " {
				args = append(args, arg)
			}
			arg = ""
		} else {
			arg += string(value)
		}
	}

	if len(args) > 1 {
		return args[1:]
	}

	return []string{}
}
