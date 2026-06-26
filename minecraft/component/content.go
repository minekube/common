package component

import "go.minekube.com/common/minecraft/key"

type Score struct {
	Name      string
	Objective string
	Value     string
	S         Style
	Extra     []Component
}

func (s *Score) Children() []Component {
	return s.Extra
}

func (s *Score) Style() *Style {
	return &s.S
}

func (s *Score) SetChildren(children []Component) {
	s.Extra = children
}

type Selector struct {
	Pattern   string
	Separator Component
	S         Style
	Extra     []Component
}

func (s *Selector) Children() []Component {
	return s.Extra
}

func (s *Selector) Style() *Style {
	return &s.S
}

func (s *Selector) SetChildren(children []Component) {
	s.Extra = children
}

type Keybind struct {
	Key   string
	S     Style
	Extra []Component
}

func (k *Keybind) Children() []Component {
	return k.Extra
}

func (k *Keybind) Style() *Style {
	return &k.S
}

func (k *Keybind) SetChildren(children []Component) {
	k.Extra = children
}

type NBT struct {
	Path      string
	Interpret bool
	Plain     bool
	Separator Component
	S         Style
	Extra     []Component
}

type BlockNBT struct {
	NBT
	Block string
}

type EntityNBT struct {
	NBT
	Entity string
}

type StorageNBT struct {
	NBT
	Storage key.Key
}

func (n *NBT) Children() []Component {
	return n.Extra
}

func (n *NBT) Style() *Style {
	return &n.S
}

func (n *NBT) SetChildren(children []Component) {
	n.Extra = children
}
