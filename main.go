package main

import (
	"vanilla-proxy/handler/handlers"
	"vanilla-proxy/log"
	"vanilla-proxy/proxy"
	"vanilla-proxy/proxy/player/manager"
	"vanilla-proxy/utils"
)

func main() {
	log.Logger = log.New()
	log.Logger.Debugln("Logger has been started")

	config := utils.ReadConfig()

	proxy.ProxyInstance = proxy.New(config, manager.NewPlayerManager())

	err := proxy.ProxyInstance.Start(handlers.New())
	if err != nil {
		log.Logger.Errorln("Error while starting server: ", err)
		panic(err)
	}
}
