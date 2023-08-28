package world

import (
	"vanilla-proxy/math"
)

type Worlds struct {
	BoarderEnabled bool
	Border         math.Area2
}

func Init(enabled bool, border *math.Area2) *Worlds {
	return &Worlds{
		BoarderEnabled: enabled,
		Border:         *border,
	}
}
