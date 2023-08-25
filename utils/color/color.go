package color

import (
	"strings"
)

const (
	Black      = "§0"
	DarkBlue   = "§1"
	DarkGreen  = "§2"
	DarkAqua   = "§3"
	DarkRed    = "§4"
	DarkPurple = "§5"
	Gold       = "§6"
	Grey       = "§7"
	DarkGrey   = "§8"
	Blue       = "§9"
	Green      = "§a"
	Aqua       = "§b"
	Red        = "§c"
	Purple     = "§d"
	Yellow     = "§e"
	White      = "§f"
	DarkYellow = "§g"

	Obfuscated = "§k"
	Bold       = "§l"
	Italic     = "§o"
	Reset      = "§r"
)

var StrMap = map[string]string{
	"black":       Black,
	"dark-blue":   DarkBlue,
	"dark-green":  DarkGreen,
	"dark-aqua":   DarkAqua,
	"dark-red":    DarkRed,
	"dark-purple": DarkPurple,
	"gold":        Gold,
	"grey":        Grey,
	"dark-grey":   DarkGrey,
	"blue":        Blue,
	"green":       Green,
	"aqua":        Aqua,
	"red":         Red,
	"purple":      Purple,
	"yellow":      Yellow,
	"white":       White,
	"dark-yellow": DarkYellow,
	"obfuscated":  Obfuscated,
	"bold":        Bold,
	"b":           Bold,
	"italic":      Italic,
	"i":           Italic,
	"reset":       Reset,
}

func Colorize(str string) string { //%red%Hi %aqua%Egor %reset%
	for key, value := range StrMap {
		str = strings.ReplaceAll(str, "%"+key+"%", value)
	}

	return str
}

func Colorizef(str string, colors ...string) string { //%red%Hi %aqua%Egor %reset%
	for _, value := range colors {
		str = strings.Replace(str, "%color%", value, 1)
	}

	return str
}
