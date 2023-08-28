package whitelist

import (
	"fmt"
	"strings"
	"vanilla-proxy/proxy/command/sender"
)

type WhitelistCommandExecutor struct {
	WhitelistManager *WhitelistManager
}

func (wce WhitelistCommandExecutor) Execute(commandSender sender.CommandSender, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments")
	}

	do := args[0]
	player := strings.Join(args[1:], " ")

	switch do {
	case "add":
		ok := wce.WhitelistManager.AddPlayer(player)
		if ok {
			commandSender.SendMessage("Player added to whitelist!")
		} else {
			return fmt.Errorf("the player is already on the whitelist")
		}
	case "remove":
		ok := wce.WhitelistManager.AddPlayer(player)
		if ok {
			commandSender.SendMessage("Player removed from whitelist!")
		} else {
			return fmt.Errorf("the player is not on the whitelist")
		}
	}

	return nil
}

func (WhitelistCommandExecutor) ForPlayer() bool {
	return false
}
