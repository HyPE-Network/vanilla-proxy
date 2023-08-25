package sender

type CommandSender interface {
	SendMessage(message string)
}
