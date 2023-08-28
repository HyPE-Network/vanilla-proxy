package manager

import (
	"math"
	"strings"

	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/HyPE-Network/vanilla-proxy/proxy/session"

	"github.com/sandertv/gophertunnel/minecraft"
)

type PlayerManager struct {
	Players map[string]human.Human
}

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		Players: make(map[string]human.Human),
	}
}

func (pm *PlayerManager) AddPlayer(conn *minecraft.Conn, serverConn *minecraft.Conn) human.Human {
	ab := session.NewBridge(conn, serverConn)
	newSession := session.NewSession(conn.IdentityData(), conn.ClientData(), ab)
	var pl human.Human = player.NewPlayer(conn.IdentityData().DisplayName, newSession, conn.GameData())

	pm.Players[conn.IdentityData().DisplayName] = pl

	return pl
}

func (pm *PlayerManager) DeletePlayer(player human.Human) {
	if _, ok := pm.Players[player.GetName()]; ok {
		delete(pm.Players, player.GetName())
		log.Logger.Infoln(player.GetName(), "left the server")
	}

	player.GetSession().Connection.ServerConn.Close()
	player.GetSession().Connection.ClientConn.Close()
}

func (pm *PlayerManager) DeleteAll() {
	for _, pl := range pm.Players {
		pm.DeletePlayer(pl)
	}
}

func (pm *PlayerManager) GetPlayer(name string) human.Human {
	name = strings.ToLower(name)
	dt := math.MaxUint8
	var found human.Human
	for _, pl := range pm.Players {
		if strings.HasPrefix(strings.ToLower(pl.GetName()), name) {
			cdt := len(pl.GetName()) - len(name)
			if cdt < dt {
				found = pl
				dt = cdt
			}
			if cdt == 0 {
				break
			}
		}
	}

	return found
}

func (pm *PlayerManager) GetPlayerExact(name string) human.Human {
	name = strings.ToLower(name)
	for _, pl := range pm.Players {
		if strings.ToLower(pl.GetName()) == name {
			return pl
		}
	}

	return nil
}

func (pm *PlayerManager) IsOnline(name string) bool {
	h := pm.GetPlayerExact(name)

	return h != nil
}

func (pm *PlayerManager) PlayerList() map[string]human.Human {
	return pm.Players
}

func (pm *PlayerManager) PlayersCount() int {
	return len(pm.Players)
}
