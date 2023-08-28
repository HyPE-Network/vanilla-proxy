package command

import (
	"strings"

	"github.com/HyPE-Network/vanilla-proxy/proxy/command/sender"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type CommandManager struct {
	Ops      []string
	Commands map[CommandExecutor]protocol.Command
}

type CommandExecutor interface {
	Execute(sender.CommandSender, []string) error
	ForPlayer() bool
}

func InitManager(ops []string) *CommandManager {
	return &CommandManager{
		Ops:      ops,
		Commands: make(map[CommandExecutor]protocol.Command),
	}
}

func (cm *CommandManager) RegisterCommand(command protocol.Command, executor CommandExecutor) {
	cm.Commands[executor] = command
}

func (cm *CommandManager) IsOp(name string) bool {
	for _, player_name := range cm.Ops {
		if strings.EqualFold(player_name, name) {
			return true
		}
	}

	return false
}
