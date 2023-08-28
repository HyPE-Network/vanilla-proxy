package world

import (
	"vanilla-proxy/math"
)

type Worlds struct {
	Border math.Area2
}

func Init(border *math.Area2) *Worlds {
	return &Worlds{
		Border: *border,
	}
}
