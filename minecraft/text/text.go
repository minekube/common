package text

import (
	"encoding/json"
	"fmt"
)

type Color string

const Reset = "reset"

// Text color
const (
	ColorBlack       Color = "black"
	ColorDarkBlue    Color = "dark_blue"
	ColorDarkGreen   Color = "dark_green"
	ColorDarkAqua    Color = "dark_aqua"
	ColorDarkRed     Color = "dark_red"
	ColorDarkPurple  Color = "dark_purple"
	ColorGold        Color = "gold"
	ColorGray        Color = "gray"
	ColorDarkGray    Color = "dark_gray"
	ColorBlue        Color = "blue"
	ColorGreen       Color = "green"
	ColorAqua        Color = "aqua"
	ColorRed         Color = "red"
	ColorLightPurple Color = "light_purple"
	ColorYellow      Color = "yellow"
	ColorWhite       Color = "white"
)

// Decoration is a text decoration.
type Decoration string

// Text decoration
const (
	DecorationObfuscated    Decoration = "obfuscated"
	DecorationBold          Decoration = "bold"
	DecorationStrikethrough Decoration = "strikethrough"
	DecorationUnderlined    Decoration = "underlined"
	DecorationItalic        Decoration = "italic"
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
	Action string `json:"action,omitempty"`
	Value  string `json:"value,omitempty"`
}

func (e *ClickEvent) Clone() *ClickEvent {
	if e == nil {
		return nil
	}
	return &ClickEvent{Action: e.Action, Value: e.Value}
}

// Click events
const (
	ClickEventOpenUrl        = "open_url"
	ClickEventRunCommand     = "run_command"
	ClickEventSuggestCommand = "suggest_command"
	ClickEventChangePage     = "change_page"
)

func ClickEventOfRunCommand(command string) *ClickEvent {
	return &ClickEvent{Action: ClickEventRunCommand, Value: command}
}
func ClickEventOfOpenUrl(url string) *ClickEvent {
	return &ClickEvent{Action: ClickEventOpenUrl, Value: url}
}
func ClickEventOfSuggestCommand(command string) *ClickEvent {
	return &ClickEvent{Action: ClickEventSuggestCommand, Value: command}
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

// Hover event
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
