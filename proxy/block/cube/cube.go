package cube

import "github.com/sandertv/gophertunnel/minecraft/protocol"

const (
	// FaceDown represents the bottom face of a block.
	FaceDown = iota
	// FaceUp represents the top face of a block.
	FaceUp
	// FaceNorth represents the north face of a block.
	FaceNorth
	// FaceSouth represents the south face of a block.
	FaceSouth
	// FaceWest represents the west face of the block.
	FaceWest
	// FaceEast represents the east face of the block.
	FaceEast
)

func Side(p protocol.BlockPos, face int32) protocol.BlockPos {
	switch face {
	case FaceUp:
		p[1]++
	case FaceDown:
		p[1]--
	case FaceNorth:
		p[2]--
	case FaceSouth:
		p[2]++
	case FaceWest:
		p[0]--
	case FaceEast:
		p[0]++
	}
	return p
}
