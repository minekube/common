package component

import (
	"github.com/google/uuid"
	"go.minekube.com/common/minecraft/key"
	"go.minekube.com/common/minecraft/nbt"
	"reflect"
)

type HoverEvent interface {
	Action() HoverAction
	Value() interface{}
}

func NewHoverEvent(action HoverAction, value interface{}) HoverEvent {
	return &hoverEvent{action, value}
}

func ShowText(text Component) HoverEvent {
	return &hoverEvent{ShowTextAction, text}
}

func ShowItem(item *ShowItemHoverType) HoverEvent {
	return &hoverEvent{ShowItemAction, item}
}

func ShowEntity(entity *ShowEntityHoverType) HoverEvent {
	return &hoverEvent{ShowEntityAction, entity}
}

type hoverEvent struct {
	action HoverAction
	value  interface{}
}

func (h *hoverEvent) Action() HoverAction {
	return h.action
}

func (h *hoverEvent) Value() interface{} {
	return h.value
}

type HoverAction interface {
	Name() string
	Type() ActionType
	Readable() bool // When a HoverAction is not readable it will not be unmarshalled.
}

type ActionType reflect.Type

type hoverAction struct {
	name     string
	_type    ActionType
	readable bool
}

func (a *hoverAction) Name() string {
	return a.name
}

func (a *hoverAction) Type() ActionType {
	return a._type
}

func (a *hoverAction) Readable() bool {
	return a.readable
}

func (a *hoverAction) String() string {
	return a.name
}

type ShowItemHoverType struct {
	Item  key.Key
	Count int
	NBT   nbt.BinaryTagHolder // nil-able
}

type ShowEntityHoverType struct {
	Type key.Key
	Id   uuid.UUID // UUID
	Name Component
}

var (
	ShowTextAction   HoverAction = &hoverAction{"show_text", reflect.TypeOf(Text{}), true}
	ShowItemAction   HoverAction = &hoverAction{"show_item", reflect.TypeOf(ShowItemHoverType{}), true}
	ShowEntityAction HoverAction = &hoverAction{"show_entity", reflect.TypeOf(ShowEntityHoverType{}), true}

	HoverActions = func() map[string]HoverAction {
		m := map[string]HoverAction{}
		for _, a := range []HoverAction{
			ShowTextAction,
			ShowItemAction,
			ShowEntityAction,
		} {
			m[a.Name()] = a
		}
		return m
	}()
)
