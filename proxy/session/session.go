package session

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
)

type Session struct {
	IdentityData login.IdentityData
	ClientData   login.ClientData
	Connection   *Bridge
}

func NewSession(identityData login.IdentityData, clientData login.ClientData, connection *Bridge) *Session {
	return &Session{IdentityData: identityData, ClientData: clientData, Connection: connection}
}

type Bridge struct {
	ClientConn *minecraft.Conn
	ServerConn *minecraft.Conn
}

func NewBridge(clientConn *minecraft.Conn, serverConn *minecraft.Conn) *Bridge {
	return &Bridge{ClientConn: clientConn, ServerConn: serverConn}
}
