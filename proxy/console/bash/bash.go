package bash

import "os/exec"

type Bash struct {
	Screen string
}

func NewBash(screen string) *Bash {
	return &Bash{
		Screen: screen,
	}
}

func (bash *Bash) SendCommand(command string) error {
	bashCommand := bash.CommandBuilder(command)
	cmd := exec.Command("bash", "-c", bashCommand)
	err := cmd.Start()

	return err
}

func (bash *Bash) CommandBuilder(command string) string {
	bashCommand := "screen -S " + bash.Screen + " -X eval 'stuff \"" + command + "\"\\015'"
	return bashCommand
}

func (bash *Bash) Close() {}
