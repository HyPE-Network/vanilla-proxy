package server

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy/command/sender"
)

type WhitelistCommandExecutor struct {
}

func (WhitelistCommandExecutor) Execute(commandSender sender.CommandSender, args []string) error {

	return nil
}

func (WhitelistCommandExecutor) ForPlayer() bool {
	return false
}
