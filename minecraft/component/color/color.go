package color

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

var InvalidFormatErr = errors.New("color.ParseHexString: invalid format")

// Color is a Minecraft text color.
type Color color.RGBA

func (c Color) RGBA() (r uint32, g uint32, b uint32, a uint32) {
	return color.RGBA(c).RGBA()
}

// HexString converts the Color to the '#' prefixed hex string.
func (c Color) HexString() string {
	r, g, b, a := c.RGBA()
	r, g, b, a = r/0x101, g/0x101, b/0x101, a/0x101
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	//return fmt.Sprintf("#%06x", int(c))
}

// ParseHexString parses a web color given by its hex RGB format.
// See https://en.wikipedia.org/wiki/Web_colors for input format.
//
// For details, see https://stackoverflow.com/a/54200713/1705598
func ParseHexString(hex string) (c Color, err error) {
	hex = strings.TrimPrefix(hex, "#")

	if true {
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
		case 7:
			c.R = hexToByte(hex[1])<<4 + hexToByte(hex[2])
			c.G = hexToByte(hex[3])<<4 + hexToByte(hex[4])
			c.B = hexToByte(hex[5])<<4 + hexToByte(hex[6])
		case 4:
			c.R = hexToByte(hex[1]) * 17
			c.G = hexToByte(hex[2]) * 17
			c.B = hexToByte(hex[3]) * 17
		default:
			err = InvalidFormatErr
		}
		c.A = 0xff
		return
	}

	values, err := strconv.ParseInt(hex, 16, 32)
	if err != nil {
		return c, InvalidFormatErr
	}
	return ParseHex(int(values)), nil
}

func ParseHex(hex int) (c Color) {
	c.A = 0xff
	c.R = uint8(hex >> 16)
	c.G = uint8((hex >> 8) & 0xff)
	c.B = uint8(hex & 0xff)
	return
}

func MustParseHexString(hex string) Color {
	c, err := ParseHexString(hex)
	if err != nil {
		panic(err)
	}
	return c
}
