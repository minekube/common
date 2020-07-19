package color

import (
	"math"
)

const (
	BlackValue       RGB = 0x000000
	DarkBlueValue    RGB = 0x0000aa
	DarkGreenValue   RGB = 0x00aa00
	DarkAquaValue    RGB = 0x00aaaa
	DarkRedValue     RGB = 0xaa0000
	DarkPurpleValue  RGB = 0xaa00aa
	GoldValue        RGB = 0xffaa00
	GrayValue        RGB = 0xaaaaaa
	DarkGrayValue    RGB = 0x555555
	BlueValue        RGB = 0x5555ff
	GreenValue       RGB = 0x55ff55
	AquaValue        RGB = 0x55ffff
	RedValue         RGB = 0xff5555
	LightPurpleValue RGB = 0xff55ff
	YellowValue      RGB = 0xffff55
	WhiteValue       RGB = 0xffffff
)

var (
	Black       = &Named{"black", BlackValue}
	DarkBlue    = &Named{"dark_blue", DarkBlueValue}
	DarkGreen   = &Named{"dark_green", DarkGreenValue}
	DarkAqua    = &Named{"dark_aqua", DarkAquaValue}
	DarkRed     = &Named{"dark_red", DarkRedValue}
	DarkPurple  = &Named{"dark_purple", DarkPurpleValue}
	Gold        = &Named{"gold", GoldValue}
	Gray        = &Named{"gray", GrayValue}
	DarkGray    = &Named{"dark_gray", DarkGrayValue}
	Blue        = &Named{"blue", BlueValue}
	Green       = &Named{"green", GreenValue}
	Aqua        = &Named{"aqua", AquaValue}
	Red         = &Named{"red", RedValue}
	LightPurple = &Named{"light_purple", LightPurpleValue}
	Yellow      = &Named{"yellow", YellowValue}
	White       = &Named{"white", WhiteValue}

	Names = func() map[string]*Named {
		m := map[string]*Named{}
		for _, a := range []*Named{
			Black,
			DarkBlue,
			DarkGreen,
			DarkAqua,
			DarkRed,
			DarkPurple,
			Gold,
			Gray,
			DarkGray,
			Blue,
			Green,
			Aqua,
			Red,
			LightPurple,
			Yellow,
			White,
		} {
			m[a.Name] = a
		}
		return m
	}()
)

// Exact gets the named color exactly matching the provided color value.
func Exact(value int) *Named {
	for _, n := range Names {
		if value == n.Value() {
			return n
		}
	}
	return nil
}

// Nearest finds the named color nearest to the provided colour.
// Always returns a value. Returns Black if any is nil.
func Nearest(any TextColor) *Named {
	if n, ok := any.(*Named); ok {
		return n
	}
	if any == nil {
		return Black
	}
	matchedDistance := math.MaxInt32
	match := Black
	for _, potential := range Names {
		distance := distanceSqrt(any, potential)
		if distance < matchedDistance {
			match = potential
			matchedDistance = distance
		}
		if distance == 0 {
			break // same color
		}
	}
	return match
}

func distanceSqrt(self, other TextColor) int {
	rAvg := (self.Red() + other.Red()) / 2
	dR := self.Red() - other.Red()
	dG := self.Green() - other.Green()
	dB := self.Blue() - other.Blue()
	return ((2 + (rAvg / 256)) * (dR * dR)) + (4 * (dG * dG)) + ((2 + ((255 - rAvg) / 256)) * (dB * dB))
}

type Named struct {
	Name string
	RGB
}

var _ TextColor = (*Named)(nil)
