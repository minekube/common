package component

import (
	"go.minekube.com/common/minecraft/component/color"
)

type Component interface {
	Append(...Component) Component // Append components.

	Color() color.TextColor
	Decoration() TextDecoration
	HasDecoration() bool
	Style() *Style
	Insertion() *string
	SetInsertion(*string) Component

	ClickEvent() ClickEvent
	SetClickEvent(ClickEvent) Component

	HoverEvent() HoverEvent
	SetHoverEvent(HoverEvent) Component

	DetectCycle(Component) // Detect cycle between another component.

	Children() []Component
	SetChildren([]Component) Component // Set the children components.
}

type Text interface {
	Component
	Content() string             // Returns the plain text content
	SetContent(string) Component // Set plain text content
}

type Translatable interface {
}
