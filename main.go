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

	log.Logger.Debugln("Loading proxy data..")
	if err := load(); err != nil {
		log.Logger.Errorln("Error while loading proxy: ", err)
		panic(err)
	}

	log.Logger.Debugln("Proxy is starting...")
	if err := start(); err != nil {
		log.Logger.Errorln("Error while starting server: ", err)
		panic(err)
	}
}

func load() error {
	var err error = nil

	config := utils.ReadConfig()

	proxy.ProxyInstance = proxy.New(config, manager.NewPlayerManager())

	return err
}

func start() error {
	var err error = nil

	err = proxy.ProxyInstance.Start(handlers.New())

	return err
}
