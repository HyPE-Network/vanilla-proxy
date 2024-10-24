package data

import (
	"github.com/HyPE-Network/vanilla-proxy/proxy/player/form"

	"github.com/df-mc/atomic"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type PlayerData struct {
	GameData          minecraft.GameData
	Forms             map[uint32]form.Form
	Closed            bool
	CurrentScoreboard atomic.Value[string]
	CurrentLines      atomic.Value[[]string]
	StartSessionTime  int64
	Authorized        bool
	FakeChestOpen     bool
	FakeChestPos      protocol.BlockPos
	Windows           byte
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
