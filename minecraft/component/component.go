// Is is recommended to make this package a dot import alias to keep code clean.
package component

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
