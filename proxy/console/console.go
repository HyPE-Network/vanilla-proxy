package console

type CommandSender interface {
	SendCommand(string) error
	Close()
}
