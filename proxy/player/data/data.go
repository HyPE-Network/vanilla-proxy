package data

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type PlayerData struct {
	GameData         minecraft.GameData
	StartSessionTime int64
	Authorized       bool
	Windows          byte

	// Disconnected is true if the player is currently being disconnected from the server.
	Disconnected bool
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

func (pd *PlayerData) GetNextWindowId() byte {
	if pd.Windows == 0 {
		pd.Windows = 1
	} else if pd.Windows > 245 {
		pd.Windows = 1
	}

	pd.Windows++

	return pd.Windows
}
