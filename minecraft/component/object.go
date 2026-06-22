package component

import (
	"github.com/google/uuid"
	"go.minekube.com/common/minecraft/key"
)

var DefaultSpriteAtlas = key.New(key.MinecraftNamespace, "blocks")

type Object struct {
	Contents ObjectContents
	S        Style
	Extra    []Component
	Fallback Component
}

func (o *Object) Children() []Component {
	return o.Extra
}

func (o *Object) Style() *Style {
	return &o.S
}

func (o *Object) SetChildren(children []Component) {
	o.Extra = children
}

type ObjectContents interface {
	objectContents()
}

type SpriteObjectContents struct {
	Atlas  key.Key
	Sprite key.Key
}

func (*SpriteObjectContents) objectContents() {}

func AtlasSprite(atlas, sprite key.Key) *Object {
	if atlas == nil {
		atlas = DefaultSpriteAtlas
	}
	return &Object{Contents: &SpriteObjectContents{
		Atlas:  atlas,
		Sprite: sprite,
	}}
}

type PlayerHeadObjectContents struct {
	Profile *PlayerProfile
	Hat     bool
}

func (*PlayerHeadObjectContents) objectContents() {}

type PlayerProfile struct {
	Name       string
	Id         uuid.UUID
	Properties []ProfileProperty
	Texture    key.Key
}

type ProfileProperty struct {
	Name      string
	Value     string
	Signature string
}

func PlayerHead(profile *PlayerProfile, hat bool) *Object {
	return &Object{Contents: &PlayerHeadObjectContents{
		Profile: profile,
		Hat:     hat,
	}}
}
