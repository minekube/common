package component

import (
	"go.minekube.com/common/minecraft/component/color"
	"go.minekube.com/common/minecraft/key"
)

var (
	DefaultFont = key.New(key.MinecraftNamespace, "default")
)

type Style struct {
	Obfuscated, Bold, Strikethrough, Underlined, Italic State

	// Nil-ables
	Font       key.Key
	Color      color.TextColor
	ClickEvent ClickEvent
	HoverEvent HoverEvent
	Insertion  *string // Gets the string to be inserted when this component is shift-clicked.
}

func (c *Style) IsEmpty() bool {
	return c.Obfuscated == NotSet &&
		c.Bold == NotSet &&
		c.Strikethrough == NotSet &&
		c.Underlined == NotSet &&
		c.Italic == NotSet &&
		c.Font == nil &&
		c.Color == nil &&
		c.ClickEvent == nil &&
		c.HoverEvent == nil &&
		c.Insertion == nil
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
