package color

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

// Color is the interface for a Minecraft text color,
// either Named or RGB (for Hex).
type Color interface {
	fmt.Stringer
	color.Color
	Hex() string   // Returns the hex "html" representation of the color, as in #ff0080.
	Named() *Named // Returns the exact or nearest Named color.
}

var (
	_ Color = (*RGB)(nil)
	_ Color = (*Named)(nil)
)

// RGB is a Minecraft text color.
type RGB colorful.Color

// Minecraft named color as RGB Color.
var (
	blackRGB       = HexInt(0x000000)
	darkBlueRGB    = HexInt(0x0000aa)
	darkGreenRGB   = HexInt(0x00aa00)
	darkAquaRGB    = HexInt(0x00aaaa)
	darkRedRGB     = HexInt(0xaa0000)
	darkPurpleRGB  = HexInt(0xaa00aa)
	goldRGB        = HexInt(0xffaa00)
	grayRGB        = HexInt(0xaaaaaa)
	darkGrayRGB    = HexInt(0x555555)
	blueRGB        = HexInt(0x5555ff)
	greenRGB       = HexInt(0x55ff55)
	aquaRGB        = HexInt(0x55ffff)
	redRGB         = HexInt(0xff5555)
	lightPurpleRGB = HexInt(0xff55ff)
	yellowRGB      = HexInt(0xffff55)
	whiteRGB       = HexInt(0xffffff)
)

// Named is a Minecraft named text color.
type Named struct {
	Name string
	*RGB
}

// Minecraft named color.
var (
	Black       = &Named{"black", blackRGB}
	DarkBlue    = &Named{"dark_blue", darkBlueRGB}
	DarkGreen   = &Named{"dark_green", darkGreenRGB}
	DarkAqua    = &Named{"dark_aqua", darkAquaRGB}
	DarkRed     = &Named{"dark_red", darkRedRGB}
	DarkPurple  = &Named{"dark_purple", darkPurpleRGB}
	Gold        = &Named{"gold", goldRGB}
	Gray        = &Named{"gray", grayRGB}
	DarkGray    = &Named{"dark_gray", darkGrayRGB}
	Blue        = &Named{"blue", blueRGB}
	Green       = &Named{"green", greenRGB}
	Aqua        = &Named{"aqua", aquaRGB}
	Red         = &Named{"red", redRGB}
	LightPurple = &Named{"light_purple", lightPurpleRGB}
	Yellow      = &Named{"yellow", yellowRGB}
	White       = &Named{"white", whiteRGB}

	// NamesOrder is the order of the named colors.
	NamesOrder = []*Named{
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

	// Names is a map of the named colors.
	Names = func() map[string]*Named {
		m := map[string]*Named{}
		for _, a := range NamesOrder {
			m[a.Name] = a
		}
		return m
	}()
)

// String implements Color interface.
func (n *Named) String() string {
	return n.Name
}

// String implements Color interface.
func (c *RGB) String() string {
	return c.Hex()
}
func (c *RGB) Named() *Named {
	return c.NearestNamed()
}
func (n *Named) Named() *Named {
	return n
}

// Hex returns the hex "html" representation of the color, as in #ff0080.
func (c *RGB) Hex() string {
	return colorful.Color(*c).Hex()
}

// Distance computes the distance between two colors in RGB space.
func (c *RGB) Distance(c2 *RGB) float64 {
	return colorful.Color(*c).DistanceRgb(colorful.Color(*c2))
}

// RGBA makes RGB implement the Go color.RGB interface.
func (c *RGB) RGBA() (r uint32, g uint32, b uint32, a uint32) {
	return colorful.Color(*c).RGBA()
}

// NearestNamed finds the nearest Named color for to this RGB.
func (c *RGB) NearestNamed() *Named {
	matchedDistance := math.MaxFloat64
	match := Black
	for _, potential := range Names {
		if potential.RGB == c {
			return potential
		}
		distance := c.Distance(potential.RGB)
		if distance == 0 {
			return potential // same color
		}
		if distance < matchedDistance {
			match = potential
			matchedDistance = distance
		}
	}
	return match
}

// Make constructs an RGB from Go's color.Color interface.
func Make(c color.Color) (*RGB, bool) {
	col, ok := colorful.MakeColor(c)
	return (*RGB)(&col), ok
}

var InvalidFormatErr = errors.New("color.Hex: invalid format")

// Hex parses a web color given by its hex RGB format.
// See https://en.wikipedia.org/wiki/Web_colors for input format.
//
// Modified version of https://stackoverflow.com/a/54200713/1705598.
func Hex(hex string) (col *RGB, err error) {
	// This code is faster than colorful.Hex() since we don't use Scan and reflection.

	if !strings.HasPrefix(hex, "#") {
		return col, InvalidFormatErr
	}
	hex = strings.TrimPrefix(hex, "#")

	//if true {
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
	co, _ := colorful.MakeColor(c)
	return (*RGB)(&co), nil
	//} else {
	//	values, err := strconv.ParseInt(hex, 16, 32)
	//	if err != nil {
	//		return col, InvalidFormatErr
	//	}
	//	return HexInt(int(values)), nil
	//}

}

// HexInt parses a web color given by its hex RGB format.
func HexInt(hex int) *RGB {
	var c color.RGBA
	c.A = 0xff
	c.R = uint8(hex >> 16)
	c.G = uint8((hex >> 8) & 0xff)
	c.B = uint8(hex & 0xff)

	col, _ := colorful.MakeColor(c)
	return (*RGB)(&col)
}

// MustHex parses a "html" hex color-string, either in the 3 "#f0c" or 6 "#ff1034" digits form.
// It panics on error.
func MustHex(hex string) *RGB {
	c, err := Hex(hex)
	if err != nil {
		panic(err)
	}
	return c
}
