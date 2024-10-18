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
	// Items holds a list of all custom items that are available in the server.
	Items []protocol.ItemEntry
	// CustomBlocks holds a list of all custom blocks that are available in the server.
	CustomBlocks []protocol.BlockEntry
}

func Init(border *math.Area2) *Worlds {
	InitBlocks()

	return &Worlds{
		Border: *border,
	}
}

func (worlds *Worlds) SetItems(items []protocol.ItemEntry) {
	worlds.Items = items
}

func (worlds *Worlds) GetItems() []protocol.ItemEntry {
	return worlds.Items
}

func (worlds *Worlds) SetCustomBlocks(blocks []protocol.BlockEntry) {
	worlds.CustomBlocks = blocks
}

func (worlds *Worlds) GetCustomBlocks() []protocol.BlockEntry {
	return worlds.CustomBlocks
}

// SetBDSAvailableCommands sets the AvailableCommands packet that is sent to the player when they join the server.
func (worlds *Worlds) SetBDSAvailableCommands(pk *packet.AvailableCommands) {
	worlds.BDSAvailableCommands = *pk
}

// GetItemEntry returns the item entry of an item with the specified network ID. If the item is not found, nil is returned.
func (worlds *Worlds) GetItemEntry(networkID int32) *protocol.ItemEntry {
	items := worlds.Items
	idx := slices.IndexFunc(items, func(item protocol.ItemEntry) bool {
		return item.RuntimeID == int16(networkID)
	})
	if idx == -1 {
		// Unknown item?
		return nil
	}
	item := items[idx]
	return &item
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
