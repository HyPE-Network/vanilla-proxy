package main

import (
	"github.com/HyPE-Network/vanilla-proxy/handler/handlers"
	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/HyPE-Network/vanilla-proxy/proxy"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/manager"
	"github.com/HyPE-Network/vanilla-proxy/utils"
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
