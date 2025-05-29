package codec

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
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
	// Since Minecraft 1.21.5+ field names changed from camelCase to snake_case.
	// This setting decides whether to use the legacy camelCase field names (clickEvent, hoverEvent)
	// instead of the new snake_case field names (click_event, hover_event).
	//
	// This setting is false by default to use the new format (1.21.5+).
	// Set to true for compatibility with clients before 1.21.5.
	UseLegacyFieldNames bool
	// Since Minecraft 1.21.5+ click event field structure changed.
	// The "value" field was renamed to more specific names like "url", "path", "command", "page".
	// This setting decides whether to use the legacy "value" field structure.
	//
	// This setting is false by default to use the new format (1.21.5+).
	// Set to true for compatibility with clients before 1.21.5.
	UseLegacyClickEventStructure bool
	// Since Minecraft 1.21.5+ hover event field structure changed.
	// The "contents" field was inlined and some field names changed.
	// This setting decides whether to use the legacy "contents" field structure.
	//
	// This setting is false by default to use the new format (1.21.5+).
	// Set to true for compatibility with clients before 1.21.5.
	UseLegacyHoverEventStructure bool
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

	// New format (1.21.5+): snake_case field names
	clickEvent = "click_event"
	hoverEvent = "hover_event"

	// Legacy format (pre-1.21.5): camelCase field names
	clickEventLegacy = "clickEvent"
	hoverEventLegacy = "hoverEvent"

	// Common field names used in both formats
	clickEventAction = "action"
	hoverEventAction = "action"

	// Legacy click event structure (pre-1.21.5)
	clickEventValue = "value"

	// New click event structure (1.21.5+) - specific field names
	clickEventUrl     = "url"     // for open_url action
	clickEventPath    = "path"    // for open_file action
	clickEventCommand = "command" // for run_command and suggest_command actions
	clickEventPage    = "page"    // for change_page action

	// New actions in 1.21.6+
	clickEventDialog  = "dialog"  // for show_dialog action
	clickEventId      = "id"      // for custom action
	clickEventPayload = "payload" // for custom action (optional)

	// Legacy hover event structure (pre-1.21.5)
	hoverEventValue    = "value"
	hoverEventContents = "contents"

	// New hover event structure (1.21.5+) - inlined fields
	// For show_text action
	hoverEventText = "value" // Note: was "text" in 25w02a, changed back to "value" in 25w03a

	// For show_item action (inlined from contents)
	itemId    = "id"
	itemCount = "count"
	itemTag   = "tag"

	// For show_entity action (inlined from contents, with field renames)
	entityType = "id"   // renamed from "type" in 1.21.5+
	entityUuid = "uuid" // renamed from "id" in 1.21.5+
	entityName = "name"

	// Legacy show_entity field names (pre-1.21.5)
	entityTypeLegacy = "type"
	entityIdLegacy   = "id"
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
		clickEventKey := clickEvent
		if j.UseLegacyFieldNames {
			clickEventKey = clickEventLegacy
		}

		clickEventObj := obj{
			clickEventAction: s.ClickEvent.Action().Name(),
		}

		// Handle different field structures based on version
		if j.UseLegacyClickEventStructure {
			// Legacy structure: use "value" field for all actions
			clickEventObj[clickEventValue] = s.ClickEvent.Value()
		} else {
			// New structure: use specific field names based on action
			switch s.ClickEvent.Action().Name() {
			case "open_url":
				clickEventObj[clickEventUrl] = s.ClickEvent.Value()
			case "open_file":
				clickEventObj[clickEventPath] = s.ClickEvent.Value()
			case "run_command", "suggest_command":
				clickEventObj[clickEventCommand] = s.ClickEvent.Value()
			case "change_page":
				// Convert page number to int for new format
				if pageStr := s.ClickEvent.Value(); pageStr != "" {
					if pageInt, err := strconv.Atoi(pageStr); err == nil {
						clickEventObj[clickEventPage] = pageInt
					} else {
						// Fallback to string if conversion fails
						clickEventObj[clickEventPage] = pageStr
					}
				}
			case "copy_to_clipboard":
				clickEventObj[clickEventValue] = s.ClickEvent.Value() // Still uses "value" in new format
			case "show_dialog":
				clickEventObj[clickEventDialog] = s.ClickEvent.Value()
			case "custom":
				// For custom actions, we expect the value to contain the ID
				// and optionally a payload separated by a delimiter (e.g., "|")
				value := s.ClickEvent.Value()
				if parts := strings.SplitN(value, "|", 2); len(parts) >= 1 {
					clickEventObj[clickEventId] = parts[0]
					if len(parts) == 2 && parts[1] != "" {
						clickEventObj[clickEventPayload] = parts[1]
					}
				} else {
					clickEventObj[clickEventId] = value
				}
			default:
				// Unknown action, use value field as fallback
				clickEventObj[clickEventValue] = s.ClickEvent.Value()
			}
		}

		o[clickEventKey] = clickEventObj
	}
	if s.HoverEvent != nil {
		eventObj := obj{}
		if err := j.encodeHoverEvent(eventObj, s.HoverEvent); err != nil {
			return err
		}
		if len(eventObj) != 0 {
			hoverEventKey := hoverEvent
			if j.UseLegacyFieldNames {
				hoverEventKey = hoverEventLegacy
			}
			o[hoverEventKey] = eventObj
		}
	}
	return nil
}

func (j *Json) encodeHoverEvent(o obj, event HoverEvent) error {
	o[hoverEventAction] = event.Action().Name()

	switch event.Action().Name() {
	case "show_text":
		switch t := event.Value().(type) {
		case *Text:
			if j.UseLegacyHoverEventStructure {
				// Legacy structure: use "contents" field
				textObj := obj{}
				if err := j.encode(textObj, t); err != nil {
					return err
				}
				o[hoverEventContents] = textObj
				if !j.NoLegacyHover {
					o[hoverEventValue] = textObj
				}
			} else {
				// New structure: inline the text component as "value"
				textObj := obj{}
				if err := j.encode(textObj, t); err != nil {
					return err
				}
				o[hoverEventText] = textObj
			}
		}

	case "show_item":
		switch t := event.Value().(type) {
		case *ShowItemHoverType:
			if j.UseLegacyHoverEventStructure {
				// Legacy structure: use "contents" field
				itemObj := obj{
					itemId:    t.Item.String(),
					itemCount: t.Count,
					itemTag:   t.NBT.String(),
				}
				o[hoverEventContents] = itemObj
				if !j.NoLegacyHover {
					o[hoverEventValue] = itemObj
				}
			} else {
				// New structure: inline the item data directly
				o[itemId] = t.Item.String()
				o[itemCount] = t.Count
				if t.NBT != nil {
					o[itemTag] = t.NBT.String()
				}
			}
		}

	case "show_entity":
		switch t := event.Value().(type) {
		case *ShowEntityHoverType:
			nameObj := obj{}
			if err := j.encode(nameObj, t.Name); err != nil {
				return err
			}

			if j.UseLegacyHoverEventStructure {
				// Legacy structure: use "contents" field with old field names
				entityObj := obj{
					entityTypeLegacy: t.Type.String(),
					entityIdLegacy:   t.Id.String(),
					entityName:       nameObj,
				}
				o[hoverEventContents] = entityObj
				if !j.NoLegacyHover {
					o[hoverEventValue] = entityObj
				}
			} else {
				// New structure: inline the entity data with new field names
				o[entityType] = t.Type.String()
				o[entityUuid] = t.Id.String()
				o[entityName] = nameObj
			}
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
			var b bool
			switch v := o[string(dec)].(type) {
			case string:
				b, err = strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`value of key %q is not a bool, but %T: %s`, dec, o[string(dec)], v)
				}
			case bool:
				b = v
			default:
				return nil, fmt.Errorf(`value of key %q is not a bool, but %T`, dec, o[string(dec)])
			}
			s.SetDecoration(dec, StateByBool(b))
		}
	}
	if o.Has(insertion) {
		if i, ok := o[insertion].(string); ok {
			s.Insertion = &i
		} else {
			return nil, fmt.Errorf(`value of key %q is not a string, but %T`, insertion, o[insertion])
		}
	}

	// Support both new (click_event) and legacy (clickEvent) field names for maximum compatibility
	if o.Has(clickEvent) || o.Has(clickEventLegacy) {
		var obj map[string]interface{}
		var ok bool
		var fieldName string

		// Try new format first, then fall back to legacy format
		if o.Has(clickEvent) {
			obj, ok = o[clickEvent].(map[string]interface{})
			fieldName = clickEvent
		} else {
			obj, ok = o[clickEventLegacy].(map[string]interface{})
			fieldName = clickEventLegacy
		}

		if !ok {
			return nil, fmt.Errorf(`value of key %q is not a json object, but %T`, fieldName, o[fieldName])
		}
		s.ClickEvent = j.decodeClickEvent(obj)
	}

	// Support both new (hover_event) and legacy (hoverEvent) field names for maximum compatibility
	if o.Has(hoverEvent) || o.Has(hoverEventLegacy) {
		var obj map[string]interface{}
		var ok bool
		var fieldName string

		// Try new format first, then fall back to legacy format
		if o.Has(hoverEvent) {
			obj, ok = o[hoverEvent].(map[string]interface{})
			fieldName = hoverEvent
		} else {
			obj, ok = o[hoverEventLegacy].(map[string]interface{})
			fieldName = hoverEventLegacy
		}

		if !ok {
			return nil, fmt.Errorf(`value of key %q is not a json object, but %T`, fieldName, o[fieldName])
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

	// Try different structures based on action type and available fields
	switch action {
	case "show_text":
		// For show_text, try different field names
		if o.Has(hoverEventText) {
			// New structure (1.21.5+): direct "value" field
			value, err = j.decodeFromInterface(o[hoverEventText])
		} else if o.Has(hoverEventContents) {
			// Legacy structure: "contents" field
			value, err = j.decodeFromInterface(o[hoverEventContents])
		} else if o.Has(hoverEventValue) {
			// Very old legacy structure: "value" field
			value, err = j.decodeHoverEventContents(o[hoverEventValue], hoverAction)
		}

	case "show_item":
		// Check if it's the new inlined structure (has direct id/count fields) or legacy structure
		if o.Has(itemId) && !o.Has(hoverEventContents) && !o.Has(hoverEventValue) {
			// New inlined structure (1.21.5+) - process directly here
			var h ShowItemHoverType
			h.Item, err = j.decodeKey(o[itemId])
			if err != nil {
				return nil, err
			}
			if o.Has(itemCount) {
				f, ok := o[itemCount].(float64)
				if !ok {
					return nil, fmt.Errorf(`show item hover event's value of key %q is not a number, but %T`,
						itemCount, o[itemCount])
				}
				h.Count = int(f)
			} else {
				h.Count = 1
			}
			if o.Has(itemTag) {
				s, ok := o[itemTag].(string)
				if !ok {
					return nil, fmt.Errorf(`show item hover event's value of key %q is not a string, but %T`,
						itemTag, o[itemTag])
				}
				h.NBT = nbt.NewBinaryTagHolder(s)
			}
			value = &h
		} else if o.Has(hoverEventContents) {
			// Legacy structure: "contents" field
			value, err = j.decodeHoverEventContents(o[hoverEventContents], hoverAction)
		} else if o.Has(hoverEventValue) {
			// Very old legacy structure: "value" field
			value, err = j.decodeHoverEventContents(o[hoverEventValue], hoverAction)
		}

	case "show_entity":
		// Check if it's the new inlined structure or legacy structure
		if (o.Has(entityType) || o.Has(entityTypeLegacy)) && !o.Has(hoverEventContents) && !o.Has(hoverEventValue) {
			// New inlined structure (1.21.5+) - process directly here
			var entityTypeField, entityIdField string
			if o.Has(entityType) && o.Has(entityUuid) {
				entityTypeField = entityType
				entityIdField = entityUuid
			} else if o.Has(entityTypeLegacy) && o.Has(entityIdLegacy) {
				entityTypeField = entityTypeLegacy
				entityIdField = entityIdLegacy
			} else {
				return nil, nil
			}

			var h ShowEntityHoverType
			h.Type, err = j.decodeKey(o[entityTypeField])
			if err != nil {
				return nil, err
			}
			h.Id, err = j.decodeUUID(o[entityIdField])
			if err != nil {
				return nil, err
			}
			if o.Has(entityName) {
				h.Name, err = j.decodeFromInterface(o[entityName])
				if err != nil {
					return nil, err
				}
			}
			value = &h
		} else if o.Has(hoverEventContents) {
			// Legacy structure: "contents" field
			value, err = j.decodeHoverEventContents(o[hoverEventContents], hoverAction)
		} else if o.Has(hoverEventValue) {
			// Very old legacy structure: "value" field
			value, err = j.decodeHoverEventContents(o[hoverEventValue], hoverAction)
		}

	default:
		// Unknown action, try legacy fields
		if o.Has(hoverEventContents) {
			value, err = j.decodeHoverEventContents(o[hoverEventContents], hoverAction)
		} else if o.Has(hoverEventValue) {
			value, err = j.decodeHoverEventContents(o[hoverEventValue], hoverAction)
		}
	}

	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
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
		// Try new field names first (1.21.5+)
		var entityTypeField, entityIdField string
		if o.Has(entityType) && o.Has(entityUuid) {
			entityTypeField = entityType
			entityIdField = entityUuid
		} else if o.Has(entityTypeLegacy) && o.Has(entityIdLegacy) {
			// Fallback to legacy field names (pre-1.21.5)
			entityTypeField = entityTypeLegacy
			entityIdField = entityIdLegacy
		} else {
			// Check what fields are available for error message
			availableFields := []string{}
			for k := range o {
				availableFields = append(availableFields, k)
			}
			return nil, fmt.Errorf(`show entity hover event misses required keys. Available fields: %v`, availableFields)
		}

		var h ShowEntityHoverType
		h.Type, err = j.decodeKey(o[entityTypeField])
		if err != nil {
			return nil, err
		}
		h.Id, err = j.decodeUUID(o[entityIdField])
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
	if !o.Has(clickEventAction) {
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

	// Try to extract value using different field names based on action and version
	var value string

	// First try the legacy "value" field (works for all versions)
	if o.Has(clickEventValue) {
		if v, ok := o[clickEventValue].(string); ok {
			value = v
		}
	}

	// Then try the new specific field names (1.21.5+)
	if value == "" {
		switch action {
		case "open_url":
			if o.Has(clickEventUrl) {
				if v, ok := o[clickEventUrl].(string); ok {
					value = v
				}
			}
		case "open_file":
			if o.Has(clickEventPath) {
				if v, ok := o[clickEventPath].(string); ok {
					value = v
				}
			}
		case "run_command", "suggest_command":
			if o.Has(clickEventCommand) {
				if v, ok := o[clickEventCommand].(string); ok {
					value = v
				}
			}
		case "change_page":
			if o.Has(clickEventPage) {
				// Handle both string and int formats
				switch v := o[clickEventPage].(type) {
				case string:
					value = v
				case int:
					value = strconv.Itoa(v)
				case float64:
					value = strconv.Itoa(int(v))
				}
			}
		case "copy_to_clipboard":
			// copy_to_clipboard still uses "value" field in new format
			if o.Has(clickEventValue) {
				if v, ok := o[clickEventValue].(string); ok {
					value = v
				}
			}
		case "show_dialog":
			if o.Has(clickEventDialog) {
				if v, ok := o[clickEventDialog].(string); ok {
					value = v
				}
			}
		case "custom":
			// For custom actions, reconstruct the value from id and optional payload
			if o.Has(clickEventId) {
				if id, ok := o[clickEventId].(string); ok {
					value = id
					if o.Has(clickEventPayload) {
						if payload, ok := o[clickEventPayload].(string); ok && payload != "" {
							value = id + "|" + payload
						}
					}
				}
			}
		}
	}

	if value == "" {
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
