package color

import (
	"errors"
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"math"
	"strconv"
	"strings"
)

// Color is a Minecraft text color.
type Color colorful.Color

// Minecraft named color as RGB Color.
var (
	BlackColor       = HexInt(0x000000)
	DarkBlueColor    = HexInt(0x0000aa)
	DarkGreenColor   = HexInt(0x00aa00)
	DarkAquaColor    = HexInt(0x00aaaa)
	DarkRedColor     = HexInt(0xaa0000)
	DarkPurpleColor  = HexInt(0xaa00aa)
	GoldColor        = HexInt(0xffaa00)
	GrayColor        = HexInt(0xaaaaaa)
	DarkGrayColor    = HexInt(0x555555)
	BlueColor        = HexInt(0x5555ff)
	GreenColor       = HexInt(0x55ff55)
	AquaColor        = HexInt(0x55ffff)
	RedColor         = HexInt(0xff5555)
	LightPurpleColor = HexInt(0xff55ff)
	YellowColor      = HexInt(0xffff55)
	WhiteColor       = HexInt(0xffffff)
)

// Named is a color named by Minecraft.
type Named struct {
	Name string
	Color
}

// Minecraft named color.
var (
	Black       = Named{"black", BlackColor}
	DarkBlue    = Named{"dark_blue", DarkBlueColor}
	DarkGreen   = Named{"dark_green", DarkGreenColor}
	DarkAqua    = Named{"dark_aqua", DarkAquaColor}
	DarkRed     = Named{"dark_red", DarkRedColor}
	DarkPurple  = Named{"dark_purple", DarkPurpleColor}
	Gold        = Named{"gold", GoldColor}
	Gray        = Named{"gray", GrayColor}
	DarkGray    = Named{"dark_gray", DarkGrayColor}
	Blue        = Named{"blue", BlueColor}
	Green       = Named{"green", GreenColor}
	Aqua        = Named{"aqua", AquaColor}
	Red         = Named{"red", RedColor}
	LightPurple = Named{"light_purple", LightPurpleColor}
	Yellow      = Named{"yellow", YellowColor}
	White       = Named{"white", WhiteColor}

	NamesOrder = []Named{
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
	}

	Names = func() map[string]Named {
		m := map[string]Named{}
		for _, a := range NamesOrder {
			m[a.Name] = a
		}
		return m
	}()
)

// String implements component.Format.
func (n Named) String() string {
	return n.Name
}

// String implements component.Format.
func (c Color) String() string {
	return c.Hex()
}

// Hex returns the hex "html" representation of the color, as in #ff0080.
func (c Color) Hex() string {
	return colorful.Color(c).Hex()
}

// Distance computes the distance between two colors in RGB space.
func (c Color) Distance(c2 Color) float64 {
	return colorful.Color(c).DistanceRgb(colorful.Color(c2))
}

// RGBA makes Color implement the Go color.Color interface.
func (c Color) RGBA() (r uint32, g uint32, b uint32, a uint32) {
	return colorful.Color(c).RGBA()
}

// NearestNamed finds the nearest Named color for to this Color.
func (c *Color) NearestNamed() Named {
	matchedDistance := math.MaxFloat64
	match := Black
	for _, potential := range Names {
		if potential.Color == *c {
			return potential
		}
		distance := c.Distance(potential.Color)
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

// Make constructs a Color from Go's color.Color interface.
func Make(c color.Color) (Color, bool) {
	col, ok := colorful.MakeColor(c)
	return Color(col), ok
}

var InvalidFormatErr = errors.New("color.Hex: invalid format")

// Hex parses a web color given by its hex RGB format.
// See https://en.wikipedia.org/wiki/Web_colors for input format.
//
// Modified version of https://stackoverflow.com/a/54200713/1705598.
func Hex(hex string) (col Color, err error) {
	// This code is faster than colorful.Hex() since we don't use Scan and reflection.

	if !strings.HasPrefix(hex, "#") {
		return col, InvalidFormatErr
	}
	hex = strings.TrimPrefix(hex, "#")

	if true {
		var c color.RGBA
		c.A = 0xff

		hexToByte := func(b byte) byte {
			switch {
			case b >= '0' && b <= '9':
				return b - '0'
			case b >= 'a' && b <= 'f':
				return b - 'a' + 10
			case b >= 'A' && b <= 'F':
				return b - 'A' + 10
			}
			err = InvalidFormatErr
			return 0
		}

		switch len(hex) {
		case 6:
			c.R = hexToByte(hex[0])<<4 + hexToByte(hex[1])
			c.G = hexToByte(hex[2])<<4 + hexToByte(hex[3])
			c.B = hexToByte(hex[4])<<4 + hexToByte(hex[5])
		case 3:
			c.R = hexToByte(hex[0]) * 17
			c.G = hexToByte(hex[1]) * 17
			c.B = hexToByte(hex[2]) * 17
		default:
			err = InvalidFormatErr
		}
		if err != nil {
			return
		}
		col, _ := colorful.MakeColor(c)
		return Color(col), nil
	} else {
		values, err := strconv.ParseInt(hex, 16, 32)
		if err != nil {
			return col, InvalidFormatErr
		}
		return HexInt(int(values)), nil
	}

}

func HexInt(hex int) Color {
	var c color.RGBA
	c.A = 0xff
	c.R = uint8(hex >> 16)
	c.G = uint8((hex >> 8) & 0xff)
	c.B = uint8(hex & 0xff)

	col, _ := colorful.MakeColor(c)
	return Color(col)
}

// MustHex parses a "html" hex color-string, either in the 3 "#f0c" or 6 "#ff1034" digits form.
// It panics on error.
func MustHex(hex string) Color {
	c, err := Hex(hex)
	if err != nil {
		panic(err)
	}
	return c
}
