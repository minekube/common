package text

import (
	"encoding/json"
	"fmt"
)

type Color string

const Reset = "reset"

// Text color
const (
	Black       Color = "black"
	DarkBlue    Color = "dark_blue"
	DarkGreen   Color = "dark_green"
	DarkAqua    Color = "dark_aqua"
	DarkRed     Color = "dark_red"
	DarkPurple  Color = "dark_purple"
	Gold        Color = "gold"
	Gray        Color = "gray"
	DarkGray    Color = "dark_gray"
	Blue        Color = "blue"
	Green       Color = "green"
	Aqua        Color = "aqua"
	Red         Color = "red"
	LightPurple Color = "light_purple"
	Yellow      Color = "yellow"
	White       Color = "white"
)

// Decoration is a text decoration.
type Decoration string

// Text decoration
const (
	Obfuscated    Decoration = "obfuscated"
	Bold          Decoration = "bold"
	Strikethrough Decoration = "strikethrough"
	Underlined    Decoration = "underlined"
	Italic        Decoration = "italic"
)

// Component is to be implemented by components.
type Component interface {
	fmt.Stringer
}

// Text is the base component which is also inherited by other components.
// See https://wiki.vg/Chat#Shared_between_all_components for details.
type Text struct {
	Text  string `json:"text"`
	Color Color  `json:"color,omitempty"`

	Bold          *bool `json:"bold,omitempty"`
	Italic        *bool `json:"italic,omitempty"`
	Underlined    *bool `json:"underlined,omitempty"`
	Strikethrough *bool `json:"strikethrough,omitempty"`
	Obfuscated    *bool `json:"obfuscated,omitempty"`

	Insertion  *string     `json:"insertion,omitempty"`
	ClickEvent *ClickEvent `json:"clickEvent,omitempty"`
	HoverEvent *HoverEvent `json:"hoverEvent,omitempty"`

	Extra []*Text `json:"extra,omitempty"`
}

var _ Component = (*Text)(nil)

// String returns the json format of the Text.
func (t *Text) String() string {
	obj, _ := json.Marshal(t) // error should never happen
	return string(obj)
}

// Clone clone a Text.
func (t *Text) Clone() *Text {
	return &Text{
		Text:          t.Text,
		Color:         t.Color,
		Bold:          t.Bold,
		Italic:        t.Italic,
		Underlined:    t.Underlined,
		Strikethrough: t.Strikethrough,
		Obfuscated:    t.Obfuscated,
		Insertion:     t.Insertion,
		ClickEvent:    t.ClickEvent.Clone(),
		HoverEvent:    t.HoverEvent.Clone(),
		Extra: func() []*Text {
			a := make([]*Text, 0, len(t.Extra))
			for _, e := range t.Extra {
				a = append(a, e.Clone())
			}
			return a
		}(),
	}
}

// Reset resets the formatting.
func (t *Text) Reset() {
	t.Color = ""
	t.Bold = nil
	t.Italic = nil
	t.Underlined = nil
	t.Strikethrough = nil
	t.Obfuscated = nil
}

// Translation is a text translation component.
// See https://wiki.vg/Chat#Translation_component for details.
type Translation struct {
	Translate string      `json:"translate"`
	With      []Component `json:"with,omitempty"`
}

var _ Component = (*Translation)(nil)

// String returns the json format of the Translation.
func (t *Translation) String() string {
	obj, _ := json.Marshal(t) // error should never happen
	return string(obj)
}

type ClickEvent struct {
	Action ClickEventAction `json:"action,omitempty"`
	Value  string           `json:"value,omitempty"`
}

func (e *ClickEvent) Clone() *ClickEvent {
	if e == nil {
		return nil
	}
	return &ClickEvent{Action: e.Action, Value: e.Value}
}

type ClickEventAction string

// Click clickevent actions
const (
	OpenUrl        ClickEventAction = "open_url"
	RunCommand     ClickEventAction = "run_command"
	SuggestCommand ClickEventAction = "suggest_command"
	ChangePage     ClickEventAction = "change_page"
)

func RunCommandClickEvent(command string) *ClickEvent {
	return &ClickEvent{Action: RunCommand, Value: command}
}
func OpenUrlClickEvent(url string) *ClickEvent {
	return &ClickEvent{Action: OpenUrl, Value: url}
}
func SuggestCommandClickEvent(command string) *ClickEvent {
	return &ClickEvent{Action: SuggestCommand, Value: command}
}
func ChangePageClickEvent(page string) *ClickEvent {
	return &ClickEvent{Action: ChangePage, Value: page}
}

type HoverEvent struct {
	Action string `json:"action,omitempty"`
	Value  *Text  `json:"value,omitempty"`
}

func (e *HoverEvent) Clone() *HoverEvent {
	if e == nil {
		return nil
	}
	return &HoverEvent{
		Action: e.Action,
		Value:  e.Value.Clone(),
	}
}

// Hover clickevent
const (
	HoverEventShowText   = "show_text"
	HoverEventShowItem   = "show_item"
	HoverEventShowEntity = "show_entity"
)

func HoverEventOfShowText(text *Text) *HoverEvent {
	return &HoverEvent{Action: HoverEventShowText, Value: text}
}

// Convenient bool pointer
var (
	True  = boolPtr(true)
	False = boolPtr(false)
)

func boolPtr(b bool) *bool { return &b }
