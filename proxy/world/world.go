package world

import (
	"slices"

	"github.com/HyPE-Network/vanilla-proxy/math"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Worlds struct {
	Border math.Area2
	// Commands that were sent by BDS and are available to the player.
	BDSAvailableCommands packet.AvailableCommands
	// ItemComponentEntries holds a list of all custom items with their respective components set.
	ItemComponentEntries []protocol.ItemComponentEntry
}

func Init(border *math.Area2) *Worlds {
	InitBlocks()

	return &Worlds{
		Border: *border,
	}
}

// SetBDSAvailableCommands sets the AvailableCommands packet that is sent to the player when they join the server.
func (worlds *Worlds) SetBDSAvailableCommands(pk *packet.AvailableCommands) {
	worlds.BDSAvailableCommands = *pk
}

func (worlds *Worlds) GetItemComponentEntry(name string) *protocol.ItemComponentEntry {
	for _, entry := range worlds.ItemComponentEntries {
		if entry.Name == name {
			return &entry
		}
	}
	return nil
}

func (worlds *Worlds) GetItemComponentEntries() []protocol.ItemComponentEntry {
	return worlds.ItemComponentEntries
}

func (worlds *Worlds) AddItemComponentEntry(entry *protocol.ItemComponentEntry) {
	worlds.ItemComponentEntries = append(worlds.ItemComponentEntries, *entry)
}

func (worlds *Worlds) RemoveItemComponentEntry(entry *protocol.ItemComponentEntry) {
	idx := slices.IndexFunc(worlds.ItemComponentEntries, func(e protocol.ItemComponentEntry) bool {
		return e.Name == entry.Name
	})
	if idx == -1 {
		return
	}
	worlds.ItemComponentEntries = append(worlds.ItemComponentEntries[:idx], worlds.ItemComponentEntries[idx+1:]...)
}

func (worlds *Worlds) SetItemComponentEntries(entries []protocol.ItemComponentEntry) {
	worlds.ItemComponentEntries = entries
}
