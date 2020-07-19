package text

import (
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/common/minecraft/component/color"
)

type Text struct {
	Content    string
	Color      color.TextColor
	Decoration component.TextDecoration
	Style      *component.Style
	Children   []component.Component
	Insertion  *string
	ClickEvent component.ClickEvent
	HoverEvent component.HoverEvent
}
