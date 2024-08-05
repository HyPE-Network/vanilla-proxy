package whitelist

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type WhitelistManager struct {
	Players map[string]string
}

func Init() *WhitelistManager {
	wm := &WhitelistManager{
		Players: make(map[string]string, 0),
	}

	if _, err := os.Stat("whitelist.json"); os.IsNotExist(err) {
		f, err := os.Create("whitelist.json")
		if err != nil {
			log.Fatalf("error creating whitelist: %v", err)
		}
		data, err := json.Marshal(wm)
		if err != nil {
			log.Fatalf("error encoding default whitelist: %v", err)
		}
		if _, err := f.Write(data); err != nil {
			log.Fatalf("error writing encoded default whitelist: %v", err)
		}
		_ = f.Close()
	}

	data, err := os.ReadFile("whitelist.json")
	if err != nil {
		log.Fatalf("error reading whitelist: %v", err)
	}

	if err := json.Unmarshal(data, wm); err != nil {
		log.Fatalf("error decoding whitelist: %v", err)
	}

	return wm
}

func (wm *WhitelistManager) HasPlayer(name string, xuid string) bool {
	for player_name, player_xuid := range wm.Players {
		if strings.EqualFold(player_xuid, xuid) {
			if !strings.EqualFold(player_name, name) {
				delete(wm.Players, player_name)
				wm.Players[name] = xuid
				wm.save()
			}

			return true
		}

		if strings.EqualFold(player_name, name) {
			if xuid != "" && player_xuid == "none" {
				wm.Players[player_name] = xuid
				wm.save()
			}

			return true
		}
	}

	return false
}

func (wm *WhitelistManager) HasPlayerName(name string) bool {
	for player_name := range wm.Players {
		if strings.EqualFold(player_name, name) {
			return true
		}
	}

	return false
}

func (wm *WhitelistManager) save() {
	file, err := os.Create("whitelist.json")
	if err != nil {
		log.Fatalf("error reading whitelist: %v", err)
	}

	p, err := json.MarshalIndent(wm, "", "\t")
	if err != nil {
		log.Fatalf("error marshal whitelist: %v", err)
	}

	_, err = file.Write(p)
	if err != nil {
		log.Fatalf("error write whitelist: %v", err)
	}
	file.Close()
}

func (wm *WhitelistManager) AddPlayer(name string) bool {
	if wm.HasPlayerName(name) {
		return false
	}

	wm.Players[strings.ToLower(name)] = "none"
	wm.save()
	return true
}

func (wm *WhitelistManager) RemovePlayer(name string) bool {
	if !wm.HasPlayerName(name) {
		return false
	}

	delete(wm.Players, strings.ToLower(name))
	wm.save()
	return true
}
