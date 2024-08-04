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
	// The player's self-signed ID, this needs to be set every time the player joins
	// Because it is what BDS uses to generate a ID from.
	// So if a player switches devices, we want them to use the same ID.
	ClientSelfSignedID string `json:"clientSelfSignedID"`
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

// GetConnIdentityData returns the identity data for a player's connection
func (plm *PlayerlistManager) GetConnIdentityData(conn *minecraft.Conn) login.IdentityData {
	identityData := conn.IdentityData()
	xuid := identityData.XUID

	// If the player is in the playerlist, return the identity data that is stored
	if player, ok := plm.Players[xuid]; ok {
		return login.IdentityData{
			XUID:        xuid,
			DisplayName: player.PlayerName,
			Identity:    player.Identity,
			TitleID:     identityData.TitleID,
		}
	}

	// Set the player in the playerlist and return the identity data from the connection
	plm.SetPlayer(xuid, conn)

	// If the player is not in the playerlist, return the identity data from the connection
	return identityData
}

// GetConnClientData returns the client data for a player's connection
func (plm *PlayerlistManager) GetConnClientData(conn *minecraft.Conn) login.ClientData {
	xuid := conn.IdentityData().XUID
	clientData := conn.ClientData()

	// If the player is in the playerlist, return the client data that is stored
	if player, ok := plm.Players[xuid]; ok {
		clientData.SelfSignedID = player.ClientSelfSignedID
		return clientData
	}

	// Set the player in the playerlist and return the identity data from the connection
	plm.SetPlayer(xuid, conn)

	// If the player is not in the playerlist, return the identity data from the connection
	return clientData
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

func (plm *PlayerlistManager) SetPlayer(xuid string, conn *minecraft.Conn) {
	player := Player{
		PlayerName:         conn.IdentityData().DisplayName,
		Identity:           conn.IdentityData().Identity,
		ClientSelfSignedID: conn.ClientData().SelfSignedID,
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
