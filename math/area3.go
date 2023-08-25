package math

type Area3 struct {
	MinX int32
	MinY int32
	MinZ int32
	MaxX int32
	MaxY int32
	MaxZ int32
}

func NewArea3(MinX int32, MinY int32, MinZ int32, MaxX int32, MaxY int32, MaxZ int32) *Area3 {
	return &Area3{
		MinX: MinX,
		MinY: MinY,
		MinZ: MinZ,
		MaxX: MaxX,
		MaxY: MaxY,
		MaxZ: MaxZ,
	}
}

func (a *Area3) IsXYZInside(x int32, y int32, z int32) bool {
	return x > a.MinX && x < a.MaxX &&
		y > a.MinY && y < a.MaxY &&
		z > a.MinZ && z < a.MaxZ
}
