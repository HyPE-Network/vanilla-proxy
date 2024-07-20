package data

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/form"

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
