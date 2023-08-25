package math

type Area2 struct {
	MinX int32
	MinZ int32
	MaxX int32
	MaxZ int32
}

func NewArea2(MinX int32, MinZ int32, MaxX int32, MaxZ int32) *Area2 {
	return &Area2{
		MinX: MinX,
		MinZ: MinZ,
		MaxX: MaxX,
		MaxZ: MaxZ,
	}
}

func (a *Area2) IsXZInside(x int32, z int32) bool {
	return x > a.MinX && x < a.MaxX && z > a.MinZ && z < a.MaxZ
}

func (a *Area2) IsPositionInside(pos []int32) bool {
	x := pos[0] << 4
	z := pos[1] << 4

	return x >= a.MinX-16 && x <= a.MaxX &&
		z >= a.MinZ-16 && z <= a.MaxZ
}
