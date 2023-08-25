package broadcaster

import (
	"vanilla-proxy/proxy/scheduler"
	"vanilla-proxy/server"
)

var messages []string
var timer = 0

func Init() {
	scheduler.NewRepeatingTask(60*5, Broadcast)
	messages = []string{}
}

func Broadcast() {
	if timer >= len(messages) {
		timer = 0
	}

	text := messages[timer]
	timer += 1
	server.BroadcastMessage(text)
}
