package playerlist

import (
	"encoding/json"
	"log"
	"os"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
)

type Player struct {
	// The player's name they connected with
	PlayerName string `json:"playerName"`
	// The player's identity they connected with, this was the first UUID they connected with
	Identity string `json:"identity"`
}

type PlayerlistManager struct {
	// A map of player XUIDs to their identity data
	Players map[string]Player `json:"players"`
}

// Init initializes the playerlist manager
func Init() *PlayerlistManager {
	plm := &PlayerlistManager{
		Players: make(map[string]Player),
	}

	if _, err := os.Stat("playerlist.json"); os.IsNotExist(err) {
		f, err := os.Create("playerlist.json")
		if err != nil {
			log.Fatalf("error creating playerlist: %v", err)
		}
		data, err := json.Marshal(plm.Players)
		if err != nil {
			log.Fatalf("error encoding default playerlist: %v", err)
		}
		if _, err := f.Write(data); err != nil {
			log.Fatalf("error writing encoded default playerlist: %v", err)
		}
		_ = f.Close()
	}

	data, err := os.ReadFile("playerlist.json")
	if err != nil {
		log.Fatalf("error reading playerlist: %v", err)
	}

	tempPlayers := make(map[string]Player)
	if err := json.Unmarshal(data, &tempPlayers); err != nil {
		log.Fatalf("error decoding playerlist: %v", err)
	}

	// Merge tempPlayers into the existing Players map
	for xuid, player := range tempPlayers {
		plm.Players[xuid] = player
	}

	return plm
}

// GetConnIdentity returns the identity data for a player's connection
func (plm *PlayerlistManager) GetConnIdentity(conn *minecraft.Conn) login.IdentityData {
	xuid := conn.IdentityData().XUID
	// If the player is not in the playerlist, return the identity data from the connection
	if player, ok := plm.Players[xuid]; ok {
		return login.IdentityData{
			XUID:        xuid,
			DisplayName: player.PlayerName,
			Identity:    player.Identity,
			TitleID:     conn.IdentityData().TitleID,
		}
	}

	// Set the player in the playerlist and return the identity data from the connection
	plm.SetPlayer(xuid, conn.IdentityData())

	return conn.IdentityData()
}

// GetXUIDFromName returns the XUID of a player by their name
func (plm *PlayerlistManager) GetXUIDFromName(playerName string) string {
	for xuid, player := range plm.Players {
		if player.PlayerName == playerName {
			return xuid
		}
	}
	return ""
}

// GetPlayer returns a player from the playerlist by their XUID
func (plm *PlayerlistManager) GetPlayer(xuid string) Player {
	return plm.Players[xuid]
}

func (plm *PlayerlistManager) SetPlayer(xuid string, identityData login.IdentityData) {
	player := Player{
		PlayerName: identityData.DisplayName,
		Identity:   identityData.Identity,
	}
	plm.Players[xuid] = player

	// Save the playerlist to disk
	data, err := json.MarshalIndent(plm.Players, "", "  ")
	if err != nil {
		log.Fatalf("error encoding playerlist: %v", err)
	}
	if err := os.WriteFile("playerlist.json", data, 0644); err != nil {
		log.Fatalf("error writing playerlist: %v", err)
	}
}
