package color

import "math"

var (
	BlackValue       = ParseHex(0x000000)
	DarkBlueValue    = ParseHex(0x0000aa)
	DarkGreenValue   = ParseHex(0x00aa00)
	DarkAquaValue    = ParseHex(0x00aaaa)
	DarkRedValue     = ParseHex(0xaa0000)
	DarkPurpleValue  = ParseHex(0xaa00aa)
	GoldValue        = ParseHex(0xffaa00)
	GrayValue        = ParseHex(0xaaaaaa)
	DarkGrayValue    = ParseHex(0x555555)
	BlueValue        = ParseHex(0x5555ff)
	GreenValue       = ParseHex(0x55ff55)
	AquaValue        = ParseHex(0x55ffff)
	RedValue         = ParseHex(0xff5555)
	LightPurpleValue = ParseHex(0xff55ff)
	YellowValue      = ParseHex(0xffff55)
	WhiteValue       = ParseHex(0xffffff)
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
func Nearest(any Color) *Named {
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

func distanceSqrt(self, other Color) int {
	rAvg := (self.R + other.R) / 2
	dR := self.R - other.R
	dG := self.G - other.G
	dB := self.B - other.B
	return math.Sqrt(Sq(c1.R-c2.R) + sq(c1.G-c2.G) + sq(c1.B-c2.B))
	return ((2 + (rAvg / 256)) * (dR * dR)) + (4 * (dG * dG)) + ((2 + ((255 - rAvg) / 256)) * (dB * dB))
}
*/
type Named struct {
	Name string
	Color
}
