package main

import (
	"github.com/HyPE-Network/vanilla-proxy/custom_handlers"
	"github.com/HyPE-Network/vanilla-proxy/handler"
	"github.com/HyPE-Network/vanilla-proxy/handler/handlers"
	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/HyPE-Network/vanilla-proxy/proxy"
	"github.com/HyPE-Network/vanilla-proxy/utils"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func main() {
	log.Logger = log.New()
	log.Logger.Debugln("Logger has been started!")

	config := utils.ReadConfig()

	proxy.ProxyInstance = proxy.New(config)

	err := proxy.ProxyInstance.Start(loadHandlers())
	if err != nil {
		log.Logger.Panicln("Error while starting server: ", err)
	}

}

func loadHandlers() handler.HandlerManager {
	utils.NewRepeatingTask(60, func() {
		custom_handlers.FetchClaims()
	})

	h := handlers.New()
	h.RegisterHandler(packet.IDAvailableCommands, custom_handlers.AvailableCommandsHandler{})
	h.RegisterHandler(packet.IDCommandRequest, custom_handlers.CommandRequestHandler{})
	h.RegisterHandler(packet.IDBlockActorData, custom_handlers.SignEditHandler{})
	h.RegisterHandler(packet.IDInventoryTransaction, custom_handlers.ClaimInventoryTransactionHandler{})
	h.RegisterHandler(packet.IDPlayerAuthInput, custom_handlers.ClaimPlayerAuthInputHandler{})
	h.RegisterHandler(packet.IDText, custom_handlers.CustomCommandRegisterHandler{})
	h.RegisterHandler(packet.IDItemComponent, custom_handlers.ItemComponentHandler{})
	h.RegisterHandler(packet.IDContainerOpen, custom_handlers.OpenContainerHandler{})
	h.RegisterHandler(packet.IDContainerClose, custom_handlers.ContainerCloseHandler{})
	h.RegisterHandler(packet.IDItemStackRequest, custom_handlers.ItemStackRequestHandler{})
	h.RegisterHandler(packet.IDPlayerList, custom_handlers.PlayerListHandler{})

	return h
}
