package world

import (
	"github.com/HyPE-Network/vanilla-proxy/math"
)

type Worlds struct {
	Border math.Area2
}

func Init(border *math.Area2) *Worlds {
	InitBlocks()

	return &Worlds{
		Border: *border,
	}
}
