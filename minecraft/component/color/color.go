package color

import (
	"fmt"
	"strconv"
	"strings"
)

type TextColor interface {
	Value() int // Returns the color, as an RGB value packed into an int.

	Red() int   // Get the red component of the text colour
	Green() int // Get the green component of the text colour
	Blue() int  // Get the blue component of the text colour

	HexString() string // Get the color as hex string
}

func FromHexString(s string) TextColor {
	if !strings.HasPrefix(s, "#") {
		return nil
	}
	hex, err := strconv.Atoi(strings.TrimPrefix(s, "#"))
	if err != nil {
		return nil
	}
	return FromValue(hex)
}

func FromValue(value int) TextColor {
	named := Exact(value)
	if named == nil {
		return RGB(value)
	}
	return named
}

// RGB is an rgb value packed into an int.
type RGB int

func (c RGB) Value() int {
	return int(c)
}

func (c RGB) Red() int {
	return (int(c) >> 16) & 0xff
}

func (c RGB) Green() int {
	return (int(c) >> 8) & 0xff
}

func (c RGB) Blue() int {
	return int(c) & 0xff
}

var _ TextColor = (*RGB)(nil)

func (c RGB) HexString() string {
	return fmt.Sprintf("#%06x", int(c))
}
