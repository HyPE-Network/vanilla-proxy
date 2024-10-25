package manager

import (
	"encoding/json"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player"
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/human"
	"github.com/HyPE-Network/vanilla-proxy/proxy/session"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
)

type PlayerManager struct {
	mu                  sync.Mutex
	Players             map[string]human.Human
	PlayersIdentityData map[string]IdentityData
}

type IdentityData struct {
	Identity, SelfSignedID string
}

func NewPlayerManager() *PlayerManager {
	playersIdentityData := make(map[string]IdentityData)
	if _, err := os.Stat("players.json"); os.IsNotExist(err) {
		f, err := os.Create("players.json")
		if err != nil {
			log.Logger.Fatalf("error creating players file: %v", err)
		}
		data, err := json.Marshal(playersIdentityData)
		if err != nil {
			log.Logger.Fatalf("error encoding default players: %v", err)
		}
		if _, err := f.Write(data); err != nil {
			log.Logger.Fatalf("error writing encoded default players: %v", err)
		}
		f.Close()
	} else {
		players, err := os.OpenFile("players.json", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Logger.Fatalf("error start open players file: %v", err)
		} else {
			info, err := players.Stat()
			if err != nil {
				log.Logger.Errorf("error stating playerlist file: %v", err)
			} else {
				data := make([]byte, info.Size())
				_, err = players.Read(data)
				if err != nil {
					log.Logger.Fatalf("error reading players file: %v", err)
				} else {
					err = json.Unmarshal(data, &playersIdentityData)
					if err != nil {
						log.Logger.Fatalf("error unmarshal players file: %v", err)
					}
				}
			}
		}
	}

	return &PlayerManager{
		Players:             make(map[string]human.Human),
		PlayersIdentityData: playersIdentityData,
	}
}

func (pm *PlayerManager) ProcessingPlayerData(cd *login.ClientData, idata *login.IdentityData) {
	if idata.XUID == "" {
		return
	}

	if pld, ok := pm.PlayersIdentityData[idata.XUID]; ok {
		cd.SelfSignedID = pld.SelfSignedID
		idata.Identity = pld.Identity
	} else {
		pld := IdentityData{
			SelfSignedID: cd.SelfSignedID,
			Identity:     idata.Identity,
		}

		pm.PlayersIdentityData[idata.XUID] = pld

		pm.mu.Lock()
		defer pm.mu.Unlock()

		players, err := os.OpenFile("players.json", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Logger.Fatalf("error open players file: %v", err)
		} else {
			js, err := json.MarshalIndent(pm.PlayersIdentityData, "", " ")
			if err != nil {
				log.Logger.Fatalf("error open players file: %v", err)
			} else {
				_, err := players.Write(js)
				if err != nil {
					log.Logger.Fatalf("error write players file: %v", err)
				}
			}
		}
	}
}

func (pm *PlayerManager) AddPlayer(conn *minecraft.Conn, serverConn *minecraft.Conn, cd login.ClientData) human.Human {
	ab := session.NewBridge(conn, serverConn)
	newSession := session.NewSession(conn.IdentityData(), cd, ab)
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
				found = pl
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
