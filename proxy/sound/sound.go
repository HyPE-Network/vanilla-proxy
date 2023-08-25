package sound

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func GetSound(soundId uint32, pos mgl32.Vec3) *packet.LevelSoundEvent {
	pk := &packet.LevelSoundEvent{
		SoundType:             soundId,
		Position:              pos,
		ExtraData:             -1,
		EntityType:            ":",
		BabyMob:               false,
		DisableRelativeVolume: false,
	}

	return pk
}
