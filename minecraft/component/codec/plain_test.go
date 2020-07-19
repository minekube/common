package codec

import (
	"fmt"
	"github.com/stretchr/testify/require"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/common/minecraft/component/text"
	"testing"
)

func TestPlainComponentCodec_Marshal(t *testing.T) {
	txt := &text.Text{
		Content: "Hello",
		Children_: []Component{
			&text.Text{
				Content: " world!",
			},
			ShowText(&text.Text{Content: ""}),
		},
	}

	p, err := PlainComponent.Marshal(txt)
	require.NoError(t, err)
	fmt.Println(string(p))
}
