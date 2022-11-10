package component

import (
	"fmt"

	"go.minekube.com/common/minecraft/color"
)

// Format is text color or decoration.
type Format interface {
	fmt.Stringer
}

var (
	_ Format = (*Decoration)(nil)
	_ Format = (color.Color)(nil)
)

// Decoration styles
const (
	// Obfuscated is a decoration which makes text obfuscated/unreadable.
	Obfuscated Decoration = "obfuscated"
	// Bold is a decoration which makes text appear bold.
	Bold Decoration = "bold"
	// Strikethrough is a decoration which makes text have a strike through it.
	Strikethrough Decoration = "strikethrough"
	// Underlined is a decoration which makes text have an underline.
	Underlined Decoration = "underlined"
	// Italic is a decoration which makes text appear in italics.
	Italic Decoration = "italic"
)

var (
	// DecorationsOrder is the order of applied decorations.
	DecorationsOrder = []Decoration{
		Obfuscated,
		Bold,
		Strikethrough,
		Underlined,
		Italic,
	}

	// Decorations is a map of all decorations.
	Decorations = func() map[Decoration]struct{} {
		m := map[Decoration]struct{}{}
		for _, a := range DecorationsOrder {
			m[a] = struct{}{}
		}
		return m
	}()
)

// Decoration is a text decoration such as "underlined".
// Use the provided decoration constants.
type Decoration string

// String implements component.Format.
func (d Decoration) String() string {
	return string(d)
}
