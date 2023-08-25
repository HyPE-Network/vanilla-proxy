package server

import (
	"vanilla-proxy/proxy"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func BroadcastMessage(message string) {
	BroadcastPacket(CreateTextPacket(message, packet.TextTypeChat))
}

func BroadcastPopup(message string) {
	BroadcastPacket(CreateTextPacket(message, packet.TextTypePopup))
}

func BroadcastTip(message string) {
	BroadcastPacket(CreateTextPacket(message, packet.TextTypeTip))
}

func BroadcastTransfer(address string, port uint16) {
	pk := &packet.Transfer{
		Address: address,
		Port:    port,
	}

	BroadcastPacket(pk)
}

func BroadcastPacket(pk packet.Packet) {
	for _, pl := range proxy.ProxyInstance.PlayerManager.PlayerList() {
		pl.DataPacket(pk)
	}
}

func CreateTextPacket(message string, textType byte) *packet.Text {
	return &packet.Text{
		TextType:         textType,
		NeedsTranslation: false,
		SourceName:       "",
		Message:          message,
		Parameters:       []string{},
		XUID:             "",
		PlatformChatID:   "",
	}
}
