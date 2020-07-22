package component

import "fmt"

// Format is text color or decoration.
type Format interface {
	fmt.Stringer
}

var (
	// A decoration which makes text obfuscated/unreadable.
	Obfuscated Decoration = "obfuscated"
	// A decoration which makes text appear bold.
	Bold Decoration = "bold"
	// A decoration which makes text have a strike through it.
	Strikethrough Decoration = "strikethrough"
	// A decoration which makes text have an underline.
	Underlined Decoration = "underlined"
	// A decoration which makes text appear in italics.
	Italic Decoration = "italic"

	DecorationsOrder = []Decoration{
		Obfuscated,
		Bold,
		Strikethrough,
		Underlined,
		Italic,
	}

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
