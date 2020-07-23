package component

import (
	"go.minekube.com/common/minecraft/color"
	"go.minekube.com/common/minecraft/key"
)

var (
	DefaultFont = key.New(key.MinecraftNamespace, "default")
)

// Style is the style of a Component.
type Style struct {
	Obfuscated, Bold, Strikethrough, Underlined, Italic State

	Font       key.Key
	Color      color.Color
	ClickEvent ClickEvent
	HoverEvent HoverEvent
	Insertion  *string // Gets the string to be inserted when this component is shift-clicked.
}

// IsZero reports whether the Style is the zero value.
func (s *Style) IsZero() bool {
	return s == nil ||
		(s.Obfuscated == NotSet &&
			s.Bold == NotSet &&
			s.Strikethrough == NotSet &&
			s.Underlined == NotSet &&
			s.Italic == NotSet &&
			s.Font == nil &&
			s.Color == nil &&
			s.ClickEvent == nil &&
			s.HoverEvent == nil &&
			s.Insertion == nil)
}

func (s *Style) Decoration(decoration Decoration) State {
	switch decoration {
	case Obfuscated:
		return s.Obfuscated
	case Bold:
		return s.Bold
	case Strikethrough:
		return s.Strikethrough
	case Underlined:
		return s.Underlined
	case Italic:
		return s.Italic
	default:
		return NotSet // unknown decoration
	}
}

func (s *Style) SetDecoration(decoration Decoration, state State) {
	switch decoration {
	case Obfuscated:
		s.Obfuscated = state
	case Bold:
		s.Bold = state
	case Strikethrough:
		s.Strikethrough = state
	case Underlined:
		s.Underlined = state
	case Italic:
		s.Italic = state
	}
}

// State is a tri-state.
type State uint8

const (
	NotSet State = iota
	True
	False
)

func (s State) String() string {
	switch s {
	case True:
		return "true"
	case False:
		return "false"
	default:
		return "null"
	}
}

func StateByBool(b bool) State {
	if b {
		return True
	}
	return False
}
