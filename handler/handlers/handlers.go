package handlers

import (
	"vanilla-proxy/handler"
	"vanilla-proxy/log"
	"vanilla-proxy/proxy"
	"vanilla-proxy/proxy/player/human"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// TODO: DISABLE DEBUG BEFORE PRODUCTION RELEASE
var debug = disabled

var target = []uint32{
	packet.IDAddActor,
}

var ignored = []uint32{
	packet.IDAnimate,
	packet.IDSetActorData,
	packet.IDMoveActorDelta,
	packet.IDCreativeContent,
	packet.IDCraftingData,
	packet.IDBiomeDefinitionList,
	packet.IDPlayerList,
	packet.IDItemComponent,
	packet.IDLevelEvent,
	packet.IDSetActorMotion,
	packet.IDUpdateAttributes,
	packet.IDPlayerAuthInput,
	packet.IDLevelChunk,
	packet.IDSubChunk,
	packet.IDSubChunkRequest,
}

type handlerManager struct {
	PacketHandlers map[uint32][]handler.PacketHandler
}

func New() handlerManager {
	return handlerManager{PacketHandlers: registerHandlers()}
}

func registerHandlers() map[uint32][]handler.PacketHandler {
	var handlers = make(map[uint32][]handler.PacketHandler)

	handlers[packet.IDSubChunk] = []handler.PacketHandler{SubChunkHandler{}}

	if proxy.ProxyInstance.Worlds != nil {
		handlers[packet.IDSubChunkRequest] = []handler.PacketHandler{SubChunkRequestHandler{}}
		handlers[packet.IDSubChunk] = append(handlers[packet.IDSubChunk], SubChunkHandlerBoarder{})
		handlers[packet.IDLevelChunk] = []handler.PacketHandler{LevelChunkHandler{}}
		handlers[packet.IDContainerOpen] = []handler.PacketHandler{OpenInventoryHandlerBoarder{}}
	}

	handlers[packet.IDModalFormResponse] = []handler.PacketHandler{ModalFormResponseHandler{}}
	handlers[packet.IDPlayerAuthInput] = []handler.PacketHandler{PlayerInputHandler{}}

	handlers[packet.IDChunkRadiusUpdated] = []handler.PacketHandler{UpdateRadiusHandler{proxy.ProxyInstance.Config.Server.ViewDistance}}
	handlers[packet.IDRequestChunkRadius] = []handler.PacketHandler{RequestRadiusHandler{proxy.ProxyInstance.Config.Server.ViewDistance}}

	handlers[packet.IDInventoryTransaction] = []handler.PacketHandler{InventoryTransactionHandler{}}
	handlers[packet.IDContainerClose] = []handler.PacketHandler{CloseInventoryHandler{}}
	handlers[packet.IDContainerOpen] = []handler.PacketHandler{OpenInventoryHandler{}}

	handlers[packet.IDCommandRequest] = []handler.PacketHandler{CommandRequestHandler{}}
	handlers[packet.IDAvailableCommands] = []handler.PacketHandler{AvailableCommandsHandler{}}

	handlers[packet.IDPacketViolationWarning] = []handler.PacketHandler{MalformedHandler{}}

	return handlers
}

func (hm *handlerManager) RegisterHandler(id int, packetHandler handler.PacketHandler) {
	_, ok := hm.PacketHandlers[uint32(id)]
	if ok {
		hm.PacketHandlers[uint32(id)] = append(hm.PacketHandlers[uint32(id)], packetHandler)
	} else {
		hm.PacketHandlers[uint32(id)] = []handler.PacketHandler{packetHandler}
	}
}

func (hm handlerManager) HandlePacket(pk packet.Packet, player human.Human, sender string) (bool, packet.Packet, error) {
	var err error
	var packetHandlers []handler.PacketHandler
	var sendPacket = true // is packet will be sent to original (true by default, may be switched by handlers)

	if debug != disabled {
		sendDebug(pk, sender)
	}

	packetHandlers, hasHandler := hm.PacketHandlers[pk.ID()]
	if hasHandler {
		for _, packetHandler := range packetHandlers {
			if sendPacket {
				sendPacket, pk, err = packetHandler.Handle(pk, player)
			} else {
				_, pk, err = packetHandler.Handle(pk, player)
			}
		}
	}

	return sendPacket, pk, err
}

func sendDebug(pk packet.Packet, sender string) {
	switch debug {
	case debugLevelAll:
		log.Logger.Debugln(sender, ":", pk.ID(), ">", pk)

	case debugLevelNotIgnored:
		if !contains(ignored, pk.ID()) {
			log.Logger.Debugln(sender, ":", pk.ID(), ">", pk)
		}

	case debugLevelTarget:
		if contains(target, pk.ID()) {
			log.Logger.Debugln(sender, ":", pk.ID(), ">", pk)
		}
	}
}

func contains(a []uint32, x uint32) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

const (
	disabled = iota
	debugLevelAll
	debugLevelNotIgnored
	debugLevelTarget
)
