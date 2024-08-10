package playerlist

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/HyPE-Network/vanilla-proxy/log"
	"github.com/gofrs/flock"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
)

type Player struct {
	PlayerName         string `json:"playerName"`
	Identity           string `json:"identity"`
	ClientSelfSignedID string `json:"clientSelfSignedID"`
}

type PlayerlistManager struct {
	mu      sync.Mutex
	Players map[string]Player `json:"players"`
}

// Init initializes the playerlist manager
func Init() (*PlayerlistManager, error) {
	plm := &PlayerlistManager{
		Players: make(map[string]Player),
	}

	log.Logger.Debugln("Attempting to acquire lock for Init")
	plm.mu.Lock()
	defer func() {
		log.Logger.Debugln("Releasing lock from Init")
		plm.mu.Unlock()
	}()

	// Create a file lock
	lock := flock.New("playerlist.json.lock")
	if err := lock.Lock(); err != nil {
		log.Logger.Errorf("error locking playerlist file: %v", err)
		return nil, err
	}
	defer lock.Unlock()

	// Open or create the playerlist.json file
	file, err := os.OpenFile("playerlist.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Logger.Errorf("error opening/creating playerlist: %v", err)
		return nil, err
	}
	defer file.Close()

	// Check if the file is empty (newly created)
	info, err := file.Stat()
	if err != nil {
		log.Logger.Errorf("error stating playerlist file: %v", err)
		return nil, err
	}
	if info.Size() == 0 {
		data, err := json.Marshal(plm.Players)
		if err != nil {
			log.Logger.Errorf("error encoding default playerlist: %v", err)
			return nil, err
		}
		if _, err := file.Write(data); err != nil {
			log.Logger.Errorf("error writing encoded default playerlist: %v", err)
			return nil, err
		}
	} else {
		// Read the existing data from the file
		data := make([]byte, info.Size())
		if _, err := file.Read(data); err != nil {
			log.Logger.Errorf("error reading playerlist: %v", err)
			return nil, err
		}

		// Unmarshal the data into the player map
		if err := json.Unmarshal(data, &plm.Players); err != nil {
			log.Logger.Errorf("error decoding playerlist: %v", err)
			return nil, err
		}
	}

	return plm, nil
}

// GetConnIdentityData returns the identity data for a player's connection
func (plm *PlayerlistManager) GetConnIdentityData(conn *minecraft.Conn) (login.IdentityData, error) {
	log.Logger.Debugln("Attempting to acquire lock in GetConnIdentityData for", conn.IdentityData().XUID)
	plm.mu.Lock()
	defer func() {
		log.Logger.Debugln("Releasing lock in GetConnIdentityData for", conn.IdentityData().XUID)
		plm.mu.Unlock()
	}()

	identityData := conn.IdentityData()
	xuid := identityData.XUID

	if player, ok := plm.Players[xuid]; ok {
		return login.IdentityData{
			XUID:        xuid,
			DisplayName: player.PlayerName,
			Identity:    player.Identity,
			TitleID:     identityData.TitleID,
		}, nil
	}

	if err := plm.SetPlayer(xuid, conn); err != nil {
		return login.IdentityData{}, err
	}
	return identityData, nil
}

// GetConnClientData returns the client data for a player's connection
func (plm *PlayerlistManager) GetConnClientData(conn *minecraft.Conn) (login.ClientData, error) {
	log.Logger.Debugln("Attempting to acquire lock in GetConnClientData for", conn.IdentityData().XUID)
	plm.mu.Lock()
	defer func() {
		log.Logger.Debugln("Releasing lock in GetConnClientData for", conn.IdentityData().XUID)
		plm.mu.Unlock()
	}()

	xuid := conn.IdentityData().XUID
	clientData := conn.ClientData()

	if player, ok := plm.Players[xuid]; ok {
		clientData.SelfSignedID = player.ClientSelfSignedID
		return clientData, nil
	}

	if err := plm.SetPlayer(xuid, conn); err != nil {
		return login.ClientData{}, err
	}
	return clientData, nil
}

// GetXUIDFromName returns the XUID of a player by their name
func (plm *PlayerlistManager) GetXUIDFromName(playerName string) (string, error) {
	log.Logger.Debugln("Attempting to acquire lock in GetXUIDFromName for", playerName)
	plm.mu.Lock()
	defer func() {
		log.Logger.Debugln("Releasing lock in GetXUIDFromName for", playerName)
		plm.mu.Unlock()
	}()

	for xuid, player := range plm.Players {
		if player.PlayerName == playerName {
			return xuid, nil
		}
	}

	return "", errors.New("player not found")
}

// GetPlayer returns a player from the playerlist by their XUID
func (plm *PlayerlistManager) GetPlayer(xuid string) (Player, error) {
	log.Logger.Debugln("Attempting to acquire lock in GetPlayer for", xuid)
	plm.mu.Lock()
	defer func() {
		log.Logger.Debugln("Releasing lock in GetPlayer for", xuid)
		plm.mu.Unlock()
	}()

	player, ok := plm.Players[xuid]
	if !ok {
		return Player{}, errors.New("player not found")
	}

	return player, nil
}

func (plm *PlayerlistManager) SetPlayer(xuid string, conn *minecraft.Conn) error {
	log.Logger.Debugln("Attempting to acquire lock in SetPlayer for", xuid)
	plm.mu.Lock()
	defer func() {
		log.Logger.Debugln("Releasing lock in SetPlayer for", xuid)
		plm.mu.Unlock()
	}()

	player := Player{
		PlayerName:         conn.IdentityData().DisplayName,
		Identity:           conn.IdentityData().Identity,
		ClientSelfSignedID: conn.ClientData().SelfSignedID,
	}
	plm.Players[xuid] = player

	// Create a file lock
	lock := flock.New("playerlist.json.lock")
	if err := lock.Lock(); err != nil {
		log.Logger.Errorf("error locking playerlist file: %v", err)
		return err
	}

	defer func() {
		if err := lock.Unlock(); err != nil {
			log.Logger.Errorf("error unlocking playerlist file: %v", err)
		}
	}()

	// Save the playerlist to disk
	data, err := json.MarshalIndent(plm.Players, "", "  ")
	if err != nil {
		log.Logger.Errorf("error encoding playerlist: %v", err)
		return err
	}

	// Open the file for writing
	file, err := os.OpenFile("playerlist.json", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Logger.Errorf("error opening playerlist for writing: %v", err)
		return err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		log.Logger.Errorf("error writing playerlist: %v", err)
		return err
	}

	return nil
}
