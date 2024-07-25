package data

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/form"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/df-mc/atomic"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type PlayerData struct {
	GameData          minecraft.GameData
	Forms             map[uint32]form.Form
	Closed            bool
	BrokenBlocks      map[protocol.BlockPos]uint32
	CurrentScoreboard atomic.Value[string]
	CurrentLines      atomic.Value[[]string]
	StartSessionTime  int64
	Authorized        bool
	FakeChestOpen     bool
	FakeChestPos      protocol.BlockPos
	Windows           byte
	// Commands that were sent by BDS and are available to the player.
	BDSAvailableCommands packet.AvailableCommands
	// ItemComponentEntries holds a list of all custom items with their respective components set.
	ItemComponentEntries []protocol.ItemComponentEntry
	// OpenContainerWindowId is the ID of the window that is currently open for the player.
	OpenContainerWindowId byte
	// OpenContainerType is the type of container that is currently open for the player.
	OpenContainerType byte
	// LastItemStackRequestID is the last ID of an item stack request that was sent by the player.
	LastItemStackRequestID int32
	// ItemsInContainers holds a list of all items that are currently in containers the player has put in.
	ItemsInContainers []protocol.StackRequestSlotInfo
	// LastUpdatedLocation is the last location that was updated for the player (updated by auth-input).
	LastUpdatedLocation mgl32.Vec3
}

func (pd *PlayerData) SetClosed() {
	pd.Closed = true
}

func (pd *PlayerData) GetNextWindowId() byte {
	if pd.Windows == 0 {
		pd.Windows = 1
	} else if pd.Windows > 245 {
		pd.Windows = 1
	}

	pd.Windows++

	return pd.Windows
}
