package codec

import (
	"bytes"
	"errors"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/common/minecraft/component/text"
)

var ()

type Codec interface {
	Marshaler
	Unmarshaler
}

type Unmarshaler interface {
	// Unmarshal parses the encoded data and stores the result
	// in the value pointed to by v. If v is nil or not a pointer,
	// Unmarshal returns an error.
	Unmarshal(data []byte, v interface{}) error
}

type Marshaler interface {
	// Marshal returns the encoding of v.
	Marshal(v component.Component) ([]byte, error)
}

// PlainComponent is a plain component serializer.
// Plain does not support more complex features such as, but not limited
// to, colours, decorations, ClickEvent and HoverEvent.
type PlainComponentCodec struct{}

var (
	PlainComponent = &PlainComponentCodec{}
)

func (p *PlainComponentCodec) Marshal(c component.Component) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := p.marshal(buf, c); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *PlainComponentCodec) marshal(buf *bytes.Buffer, c component.Component) (err error) {
	switch t := c.(type) {
	case text.Text:
		_, err = buf.WriteString(t.Content())
	}
	if err != nil {
		return err
	}

	for _, child := range c.Children() {
		if err = p.marshal(buf, child); err != nil {
			return err
		}
	}
	return nil
}

func (PlainComponentCodec) Unmarshal(data []byte, v interface{}) error {
	if v == nil {
		return errors.New("v cannot be nil")
	}
	ptr, ok := v.(*text.Text)
	if !ok {
		return errors.New("v is not a pointer")
	}
	*ptr = &text.Text{Content: string(data)}
	return nil
}
