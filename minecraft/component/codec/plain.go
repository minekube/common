package codec

import (
	"fmt"
	"go.minekube.com/common/minecraft/component"
	"io"
	"strings"
)

// Plain is a plain text component serializer.
// Plain does not support more complex features such as, but not limited
// to, colours, decorations, ClickEvent and HoverEvent.
type Plain struct{}

var _ Codec = (*Plain)(nil)

func (p Plain) Marshal(wr io.Writer, c component.Component) error {
	b := new(strings.Builder)
	if err := p.encode(b, c); err != nil {
		return err
	}
	_, err := wr.Write([]byte(b.String()))
	return err
}

func (p Plain) encode(b *strings.Builder, c component.Component) (err error) {
	switch t := c.(type) {
	case *component.Text:
		_, err = b.WriteString(t.Content)
	default:
		err = fmt.Errorf("unsupported component type %T", c)
	}
	if err != nil {
		return err
	}

	for _, child := range c.Children() {
		if err = p.encode(b, child); err != nil {
			return err
		}
	}
	return nil
}

func (Plain) Unmarshal(str []byte) (component.Component, error) {
	return &component.Text{Content: string(str)}, nil
}
