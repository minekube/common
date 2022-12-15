// Package component provides a way to represent Minecraft text components.
// Use dot import alias to keep your code lean when using structs a lot from this package.
package component

import (
	"encoding/json"
)

type Component interface {
	Children() []Component
	SetChildren([]Component)
	Style() *Style
}

type Text struct {
	Content string
	S       Style
	Extra   []Component
}

type Translation struct {
	Key  string // Translation key
	S    Style
	With []Component
}

func (t *Text) Children() []Component {
	return t.Extra
}
func (t *Text) Style() *Style {
	return &t.S
}
func (t *Text) SetChildren(children []Component) {
	t.Extra = children
}

func (t *Translation) Children() []Component {
	return t.With
}
func (t *Translation) Style() *Style {
	return &t.S
}
func (t *Translation) SetChildren(children []Component) {
	t.With = children
}

var (
	_ json.Marshaler   = (*Text)(nil)
	_ json.Unmarshaler = (*Text)(nil)
	_ json.Marshaler   = (*Translation)(nil)
	_ json.Unmarshaler = (*Translation)(nil)
)

func (t *Text) MarshalJSON() ([]byte, error)        { panic("use codec.Json instead") }
func (t *Text) UnmarshalJSON(b []byte) error        { panic("use codec.Json instead") }
func (t *Translation) UnmarshalJSON(b []byte) error { panic("use codec.Json instead") }
func (t *Translation) MarshalJSON() ([]byte, error) { panic("use codec.Json instead") }
