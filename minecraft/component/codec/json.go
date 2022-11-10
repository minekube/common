package codec

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/francoispqt/gojay"
	"github.com/google/uuid"

	col "go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/common/minecraft/key"
	"go.minekube.com/common/minecraft/nbt"
)

// Codec can marshal and unmarshal to/from components to/from a different format.
type Codec interface {
	Marshaler
	Unmarshaler
}

type Unmarshaler interface {
	// Unmarshal decodes a Component from data.
	Unmarshal(data []byte) (Component, error)
}

type Marshaler interface {
	// Marshal writes the encoding of c into wr.
	Marshal(wr io.Writer, c Component) error
}

// Json is a json serializer for Minecraft text components.
type Json struct {
	// Since Minecraft 1.16+ there can be hex colors (e.g. "#ff5555" instead of the named color "red").
	// This setting decides whether to use the nearest legacy named color of the hex color instead.
	//
	// This setting is false by default to support older client versions.
	NoDownsampleColor bool
	// Since Minecraft 1.16+ hoverEvent "value" is deprecated in favour of "contents".
	// This setting decides whether the "value" key in a hoverEvent object shall be included
	// next to "contents" or not.
	//
	// It is false by default to support older client versions.
	NoLegacyHover bool
	// Whether to use Go's standard json library for marshalling.
	// It can be set to true if features such as key sorting in objects is needed
	// (e.g. when testing to compare output).
	//
	// It is false by default to use a MUCH MORE efficient (faster, less B/op & allocs/op needed)
	// third-party json marshaller instead.
	StdJson bool
}

var _ Codec = (*Json)(nil)

// Marshal writes the json encoded Component to the Writer.
func (j *Json) Marshal(wr io.Writer, c Component) (err error) {
	o := obj{}
	if err = j.encode(o, c); err != nil {
		return err
	}

	if j.StdJson {
		var data []byte
		data, err = json.Marshal(o)
		if err != nil {
			return err
		}
		_, err = wr.Write(data)
	} else {
		// Gojay is WAY faster but has no object keys sorting,
		// which is generally not needed.
		enc := gojay.BorrowEncoder(wr)
		defer enc.Release()
		if err = enc.Encode(o); err != nil {
			return err
		}
		_, err = enc.Write()
	}
	return
}

func (j *Json) Unmarshal(data []byte) (Component, error) {
	var i interface{}
	if err := json.Unmarshal(data, &i); err != nil {
		return nil, err
	}
	return j.decodeFromInterface(i)
}

// encode
// encode
// encode
// encode
// encode
// encode
// encode
// encode
// encode

func (j *Json) encode(o obj, c Component) (err error) {
	switch t := c.(type) {
	case *Text:
		return j.encodeText(o, t)
	case *Translation:
		return j.encodeTranslation(o, t)
	default:
		return fmt.Errorf("codec.Json marshal: unsupported component type %T", c)
	}
}

// json object keys
const (
	text  = "text"
	extra = "extra"

	translate     = "translate"
	translateWith = "with"

	font      = "font"
	color     = "color"
	insertion = "insertion"

	clickEvent       = "clickEvent"
	clickEventAction = "action"
	clickEventValue  = "value"

	hoverEvent         = "hoverEvent"
	hoverEventAction   = "action"
	hoverEventValue    = "value"
	hoverEventContents = "contents"

	itemId    = "id"
	itemCount = "count"
	itemTag   = "tag"

	entityId   = "id"
	entityType = "type"
	entityName = "name"
)

func (j *Json) encodeText(o obj, t *Text) error {
	if t == nil {
		return nil
	}
	o[text] = t.Content
	return j.encodeComponent(o, t, extra)
}
func (j *Json) encodeTranslation(o obj, t *Translation) error {
	if t == nil {
		return nil
	}
	o[translate] = t.Key
	return j.encodeComponent(o, t, translateWith)
}

func (j *Json) encodeComponent(o obj, c Component, childrenKey string) (err error) {
	if c == nil {
		return nil
	}
	if err = j.encodeStyle(o, c.Style()); err != nil {
		return err
	}
	var children arr
	for _, child := range c.Children() {
		childObj := obj{}
		if err = j.encode(childObj, child); err != nil {
			return err
		}
		children = append(children, childObj)
	}
	if len(children) != 0 {
		o[childrenKey] = children
	}
	return nil
}

func (j *Json) encodeStyle(o obj, s *Style) error {
	if s == nil {
		return nil
	}
	if s.Font != nil {
		o[font] = s.Font.String()
	}
	if s.Color != nil {
		o[color] = j.encodeColor(s.Color)
	}
	for name := range Decorations {
		state := s.Decoration(name)
		if state != NotSet {
			o[string(name)] = state == True
		}
	}
	if s.Insertion != nil {
		o[insertion] = *s.Insertion
	}
	if s.ClickEvent != nil {
		o[clickEvent] = obj{
			clickEventAction: s.ClickEvent.Action().Name(),
			clickEventValue:  s.ClickEvent.Value(),
		}
	}
	if s.HoverEvent != nil {
		eventObj := obj{}
		if err := j.encodeHoverEvent(eventObj, s.HoverEvent); err != nil {
			return err
		}
		if len(eventObj) != 0 {
			o[hoverEvent] = eventObj
		}
	}
	return nil
}

func (j *Json) encodeHoverEvent(o obj, event HoverEvent) error {
	var value obj
	switch t := event.Value().(type) {
	case *Text:
		value = obj{}
		if err := j.encode(value, t); err != nil {
			return err
		}
	case *ShowItemHoverType:
		value = obj{
			itemTag:   t.Item.String(),
			itemCount: t.Count,
			itemId:    t.NBT.String(),
		}
	case *ShowEntityHoverType:
		name := obj{}
		if err := j.encode(name, t.Name); err != nil {
			return err
		}
		value = obj{
			entityType: t.Type.String(),
			entityId:   t.Id.String(),
			entityName: name,
		}
	}
	if value != nil {
		o[hoverEventAction] = event.Action().Name()
		o[hoverEventContents] = value
		if !j.NoLegacyHover {
			o[hoverEventValue] = value
		}
	}
	return nil
}

func (j *Json) encodeColor(c col.Color) (s string) {
	if c == nil {
		return
	}
	if !j.NoDownsampleColor {
		return c.Named().Name
	}
	return c.Hex()
}

// decode
// decode
// decode
// decode
// decode
// decode
// decode
// decode

func (j *Json) decodeFromInterface(i interface{}) (Component, error) {
	switch t := i.(type) {
	case map[string]interface{}:
		return j.decodeComponent(t)
	case string:
		return &Text{Content: t}, nil
	case []interface{}:
		return j.decodeFromInterfaceSlice(t)
	default:
		return nil, fmt.Errorf("codec.Json unmarshal: json input unmarshalled to unsupported type %T", i)
	}
}

func (j *Json) decodeFromInterfaceSlice(i []interface{}) (Component, error) {
	var parent Component
	for _, child := range i {
		c, err := j.decodeFromInterface(child)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			parent = c
		} else {
			parent.SetChildren(append(parent.Children(), c))
		}
	}
	if parent == nil {
		return nil, errors.New("component array must not be empty")
	}
	return parent, nil
}

func (j *Json) decodeComponent(o obj) (c Component, err error) {
	if o.Has(text) {
		c = &Text{Content: fmt.Sprint(o[text])}
	} else if o.Has(translate) {
		k := fmt.Sprint(o[translate])
		if o.Has(translateWith) {
			with, ok := o[translateWith].([]interface{})
			if !ok {
				return nil, fmt.Errorf(`found invalid translate component, value of key %q is not an array`, with)
			}
			args := make([]Component, 0, len(with))
			for _, arg := range with {
				a, err := j.decodeFromInterface(arg)
				if err != nil {
					return nil, err
				}
				args = append(args, a)
			}
			c = &Translation{
				Key:  k,
				With: args,
			}
		} else {
			c = &Translation{Key: k}
		}
	} else {
		c = &Text{}
	}

	if o.Has(extra) {
		ext, ok := o[extra].([]interface{})
		if !ok {
			return nil, fmt.Errorf(`value of key %q is not an array, but %T`, extra, o[extra])
		}
		for _, e := range ext {
			ex, err := j.decodeFromInterface(e)
			if err != nil {
				return nil, err
			}
			c.SetChildren(append(c.Children(), ex))
		}
	}

	style, err := j.decodeStyle(o)
	if err != nil {
		return nil, err
	}
	if !style.IsZero() {
		*c.Style() = *style
	}
	return c, nil
}

func (j *Json) decodeStyle(o obj) (s *Style, err error) {
	s = &Style{}
	if o.Has(font) {
		k, err := j.decodeKey(o[font])
		if err != nil {
			return nil, fmt.Errorf(`error decoding value of %q key: %v`, font, err)
		}
		s.Font = k
	}
	if o.Has(color) {
		c, dec, _, err := j.decodeColor(o[color])
		if err != nil {
			return nil, fmt.Errorf(`error decoding value of %q key: %v`, color, err)
		}
		if c != nil {
			s.Color = c
		} else if dec != nil {
			// Setting a decoration from a color is, unfortunately, something we need to support.
			s.SetDecoration(*dec, True)
		}
	}
	for dec := range Decorations {
		if o.Has(string(dec)) {
			if b, ok := o[string(dec)].(bool); ok {
				s.SetDecoration(dec, StateByBool(b))
			} else {
				return nil, fmt.Errorf(`value of key %q is not a bool, but %T`, dec, o[string(dec)])
			}
		}
	}
	if o.Has(insertion) {
		if i, ok := o[insertion].(string); ok {
			s.Insertion = &i
		} else {
			return nil, fmt.Errorf(`value of key %q is not a string, but %T`, insertion, o[insertion])
		}
	}

	if o.Has(clickEvent) {
		obj, ok := o[clickEvent].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf(`value of key %q is not a json object, but %T`, clickEvent, o[clickEvent])
		}
		s.ClickEvent = j.decodeClickEvent(obj)
	}
	if o.Has(hoverEvent) {
		obj, ok := o[hoverEvent].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf(`value of key %q is not a json object, but %T`, hoverEvent, o[hoverEvent])
		}
		s.HoverEvent, err = j.decodeHoverEvent(obj)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

// may return nil,nil in case object has missing/invalid keys to decode a HoverEvent or Readable() == false
func (j *Json) decodeHoverEvent(o obj) (h HoverEvent, err error) {
	if !o.Has(hoverEventAction) {
		return nil, nil
	}
	action, ok := o[hoverEventAction].(string)
	if !ok {
		return nil, nil
	}
	hoverAction, ok := HoverActions[action]
	if !ok || !hoverAction.Readable() {
		return nil, nil
	}
	var value interface{}
	if o.Has(hoverEventContents) {
		value, err = j.decodeHoverEventContents(o[hoverEventContents], hoverAction)
	} else if o.Has(hoverEventValue) {
		value, err = j.decodeHoverEventContents(o[hoverEventValue], hoverAction)
	} else {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return NewHoverEvent(hoverAction, value), nil
}

var errUnsupportedHoverEventAction = errors.New("unsupported hover event action")

func (j *Json) decodeHoverEventContents(v interface{}, action HoverAction) (value interface{}, err error) {
	var o obj

	switch t := v.(type) {
	case map[string]interface{}:
		o = t
	case []interface{}:
		return j.decodeFromInterfaceSlice(t)
	case string:
		// decode from legacy hover event value key which is json like "contents" but in a string
		switch {
		case equalFold(ShowTextAction, action):
			return j.Unmarshal([]byte(t))
		case equalFoldAny(action, ShowEntityAction, ShowItemAction):
			m := obj{}
			if err = json.Unmarshal([]byte(t), &m); err != nil {
				return nil, err
			}
			o = m
		default:
			return nil, fmt.Errorf("%w: %s of type %T", errUnsupportedHoverEventAction, action, v)
		}
	default:
		return nil, fmt.Errorf(
			`hover event's value of key %q is not a json object nor string, but %T`,
			hoverEventContents, v)
	}

	switch {
	case equalFold(ShowTextAction, action):
		return j.decodeComponent(o)
	case equalFold(ShowItemAction, action):
		if !o.Has(itemId) {
			return nil, fmt.Errorf(`show item hover event misses key %q`, itemId)
		}
		var h ShowItemHoverType
		h.Item, err = j.decodeKey(o[itemId])
		if err != nil {
			return nil, err
		}
		if o.Has(itemCount) {
			f, ok := o[itemCount].(float64)
			if !ok {
				return nil, fmt.Errorf(`show entity hover event's value of key %q is not a number, but %T`,
					itemCount, o[itemCount])
			}
			h.Count = int(f)
		} else {
			h.Count = 1
		}
		if o.Has(itemTag) {
			s, ok := o[itemTag].(string)
			if !ok {
				return nil, fmt.Errorf(`show entity hover event's value of key %q is not a string, but %T`,
					itemTag, o[itemTag])
			}
			h.NBT = nbt.NewBinaryTagHolder(s)
		}
		return &h, nil
	case equalFold(ShowEntityAction, action):
		if !o.Has(entityType) || !o.Has(entityId) {
			return nil, fmt.Errorf(`show entity hover event misses keys %q and/or %q`, entityType, entityId)
		}
		var h ShowEntityHoverType
		h.Type, err = j.decodeKey(o[entityType])
		if err != nil {
			return nil, err
		}
		h.Id, err = j.decodeUUID(o[entityId])
		if err != nil {
			return nil, err
		}
		if o.Has(entityName) {
			h.Name, err = j.decodeFromInterface(o[entityName])
			if err != nil {
				return nil, err
			}
		}
		return &h, nil
	}
	return nil, fmt.Errorf("%w: %s", errUnsupportedHoverEventAction, action)
}

// may return nil in case object has missing/invalid keys to decode a ClickEvent or Readable() == false
func (j *Json) decodeClickEvent(o obj) ClickEvent {
	if !o.Has(clickEventAction) || !o.Has(clickEventValue) {
		return nil
	}
	action, ok := o[clickEventAction].(string)
	if !ok {
		return nil
	}
	clickAction, ok := ClickActions[action]
	if !ok || !clickAction.Readable() {
		return nil
	}
	value, ok := o[clickEventValue].(string)
	if !ok {
		return nil
	}
	return NewClickEvent(clickAction, value)
}

func (j *Json) decodeColor(i interface{}) (c col.Color, dec *Decoration, reset bool, err error) {
	s, ok := i.(string)
	if !ok {
		err = errors.New("must be a string")
		return
	}
	if strings.HasPrefix(s, "#") {
		c, err = col.Hex(s)
		if err != nil {
			return
		}
	} else {
		c = col.Names[s]
	}
	_, ok = Decorations[Decoration(s)]
	if ok {
		dec = (*Decoration)(&s)
	}
	reset = dec == nil && strings.EqualFold(s, "reset")
	if c == nil && dec == nil && !reset {
		err = fmt.Errorf("don't know how to parse %s as color", s)
		return
	}
	return
}

func (j *Json) decodeKey(i interface{}) (key.Key, error) {
	if s, ok := i.(string); ok {
		return key.ParseValid(s)
	}
	return nil, errors.New("must be as string")
}

func (j *Json) decodeUUID(i interface{}) (uuid.UUID, error) {
	if s, ok := i.(string); ok {
		return uuid.Parse(s)
	}
	return [16]byte{}, errors.New("must be as string")
}

//
//
//
//
//
//
//

// used for json encoding/decoding
type (
	obj map[string]interface{}
	arr []interface{}
)

func (o obj) String(key string) string {
	if s, ok := o[key]; ok {
		if s, ok := s.(string); ok {
			return s
		}
	}
	return ""
}

func (o obj) Has(key string) bool {
	_, ok := o[key]
	return ok
}

func (o obj) MarshalJSONObject(enc *gojay.Encoder) {
	for k, v := range o {
		enc.AddInterfaceKey(k, v)
	}
}

func (o obj) IsNil() bool {
	return o == nil
}

func (a arr) MarshalJSONArray(enc *gojay.Encoder) {
	for _, v := range a {
		enc.AddInterface(v)
	}
}

func (a arr) IsNil() bool {
	return false
}

func equalFold(a, b HoverAction) bool {
	return a == b || strings.EqualFold(a.Name(), b.Name())
}

func equalFoldAny(s HoverAction, any ...HoverAction) bool {
	for _, a := range any {
		if equalFold(s, a) {
			return true
		}
	}
	return false
}
