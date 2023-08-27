package command

import (
	"vanilla-proxy/proxy/command/sender"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type CommandManager struct {
	Commands map[CommandExecutor]protocol.Command
}

type CommandExecutor interface {
	Execute(sender.CommandSender, []string) error
	ForPlayer() bool
}

func InitManager() *CommandManager {
	return &CommandManager{
		Commands: make(map[CommandExecutor]protocol.Command),
	}
}

func (cm CommandManager) RegisterCommand(command protocol.Command, executor CommandExecutor) {
	cm.Commands[executor] = command
}
