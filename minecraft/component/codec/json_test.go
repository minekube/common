package codec

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	. "go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/common/minecraft/key"
	"go.minekube.com/common/minecraft/nbt"
)

// Test configurations for different Minecraft versions

// Pre-1.21.5 format: camelCase field names with legacy structures
var jPre1215 = JsonPre1_21_5

// 1.21.5+ format: snake_case field names with new structures
var j1215Plus = JsonModern

// Compatibility decoder: can decode all formats but encodes in new format
var jCompat = JsonUniversal

func BenchmarkJson_Marshal(b *testing.B) {
	insertion := "insert me"
	tx := &Text{Content: "Hello", S: Style{
		Obfuscated:    True,
		Bold:          False,
		Strikethrough: NotSet,
		Underlined:    True,
		Italic:        False,
		Font:          DefaultFont,
		Color:         Aqua,
		ClickEvent:    SuggestCommand("/help"),
		HoverEvent: ShowText(&Text{
			Content: " world",
			Extra: []Component{
				&Text{Content: "!"},
			},
		}),
		Insertion: &insertion,
	}}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := j1215Plus.Marshal(ioutil.Discard, tx)
		require.NoError(b, err)
	}
}

var (
	insert = "insert me"
	txt    = &Text{
		Content: "Hello",
		Extra: []Component{
			&Text{Content: " there!", S: Style{
				Color:      Red.RGB,
				Italic:     True,
				Obfuscated: False,
			}},
		},
		S: Style{
			Obfuscated:    True,
			Bold:          False,
			Strikethrough: NotSet,
			Underlined:    True,
			Italic:        False,
			Font:          DefaultFont,
			Color:         Aqua.RGB,
			ClickEvent:    SuggestCommand("/help"),
			HoverEvent: ShowText(&Text{
				Content: " world",
				Extra: []Component{
					&Text{Content: "!"},
				},
			}),
			Insertion: &insert,
		}}

	// New format JSON (1.21.5+) - uses snake_case field names and new structures
	jsonTxtNew = `{"bold":false,"click_event":{"action":"suggest_command","command":"/help"},"color":"#55ffff","extra":[{"color":"#ff5555","italic":true,"obfuscated":false,"text":" there!"}],"font":"minecraft:default","hover_event":{"action":"show_text","value":{"extra":["!"],"text":" world"}},"insertion":"insert me","italic":false,"obfuscated":true,"text":"Hello","underlined":true}`

	// Legacy format JSON (pre-1.21.5) - uses camelCase field names and legacy structures
	jsonTxtLegacy = `{"bold":false,"clickEvent":{"action":"suggest_command","value":"/help"},"color":"#55ffff","extra":[{"color":"#ff5555","italic":true,"obfuscated":false,"text":" there!"}],"font":"minecraft:default","hoverEvent":{"action":"show_text","contents":{"extra":["!"],"text":" world"},"value":{"extra":["!"],"text":" world"}},"insertion":"insert me","italic":false,"obfuscated":true,"text":"Hello","underlined":true}`
)

func TestJson_Marshal_text(t *testing.T) {
	b := new(strings.Builder)
	err := j1215Plus.Marshal(b, txt)
	require.NoError(t, err)
	require.Equal(t, jsonTxtNew, b.String())
}

func TestJson_Unmarshal_text(t *testing.T) {
	c, err := j1215Plus.Unmarshal([]byte(jsonTxtNew))
	require.NoError(t, err)
	require.Equal(t, txt, c)
}

func TestJson_Unmarshal_NBTStyleDecorationBooleans(t *testing.T) {
	c, err := j1215Plus.Unmarshal([]byte(`{"text":"Hello","obfuscated":"1B","italic":"0B","underlined":1,"bold":0}`))
	require.NoError(t, err)
	require.Equal(t, &Text{
		Content: "Hello",
		S: Style{
			Obfuscated: True,
			Italic:     False,
			Underlined: True,
			Bold:       False,
		},
	}, c)
}

func TestJson_translation(t *testing.T) {
	tr := &Translation{
		Key: "sample.key",
		S:   Style{Color: Red.RGB},
		With: []Component{
			&Text{
				Content: "Hello",
			},
			&Translation{
				Key: "another.key",
				S:   Style{},
			},
		},
	}
	s := new(strings.Builder)
	require.NoError(t, j1215Plus.Marshal(s, tr))
	const exp = `{"color":"#ff5555","translate":"sample.key","with":["Hello",{"translate":"another.key"}]}`
	require.Equal(t, exp, s.String())

	tr2, err := j1215Plus.Unmarshal([]byte(exp))
	require.NoError(t, err)
	require.Equal(t, tr, tr2)
}

func TestJson_TranslationFallback(t *testing.T) {
	tr := &Translation{
		Key:      "sample.missing",
		Fallback: "Fallback text",
	}
	s := new(strings.Builder)
	require.NoError(t, j1215Plus.Marshal(s, tr))
	require.Equal(t, `{"fallback":"Fallback text","translate":"sample.missing"}`, s.String())

	tr2, err := j1215Plus.Unmarshal([]byte(s.String()))
	require.NoError(t, err)
	require.Equal(t, tr, tr2)
}

func TestJson_ModernContentTypes(t *testing.T) {
	storage, _ := key.Parse("example:storage")
	components := []Component{
		&Score{Name: "@s", Objective: "kills", Value: "10"},
		&Selector{Pattern: "@a", Separator: &Text{Content: ", "}},
		&Keybind{Key: "key.jump"},
		&BlockNBT{NBT: NBT{Path: "Items[0].id", Interpret: true, Separator: &Text{Content: " | "}}, Block: "~ ~ ~"},
		&EntityNBT{NBT: NBT{Path: "Health", Plain: true}, Entity: "@s"},
		&StorageNBT{NBT: NBT{Path: "value"}, Storage: storage},
	}

	for _, component := range components {
		encoded := new(strings.Builder)
		require.NoError(t, j1215Plus.Marshal(encoded, component))

		decoded, err := j1215Plus.Unmarshal([]byte(encoded.String()))
		require.NoError(t, err)
		require.Equal(t, component, decoded)
	}
}

func TestJson_DecodePrimitiveShorthand(t *testing.T) {
	c, err := j1215Plus.Unmarshal([]byte(`true`))
	require.NoError(t, err)
	require.Equal(t, &Text{Content: "true"}, c)

	c, err = j1215Plus.Unmarshal([]byte(`12`))
	require.NoError(t, err)
	require.Equal(t, &Text{Content: "12"}, c)
}

func TestJson_StyleShadowColor(t *testing.T) {
	component := &Text{
		Content: "Shadow",
		S: Style{
			ShadowColor: ShadowColorFromARGB(0x80445566),
		},
	}

	encoded := new(strings.Builder)
	require.NoError(t, j1215Plus.Marshal(encoded, component))
	require.Equal(t, `{"shadow_color":-2143005338,"text":"Shadow"}`, encoded.String())

	decoded, err := j1215Plus.Unmarshal([]byte(`{"shadow_color":[0.26666666666666666,0.3333333333333333,0.4,0.5019607843137255],"text":"Shadow"}`))
	require.NoError(t, err)
	require.Equal(t, component, decoded)
}

func TestJson_ShowItemModernComponents(t *testing.T) {
	diamond, _ := key.Parse("minecraft:diamond")
	component := &Text{
		Content: "Hover",
		S: Style{HoverEvent: ShowItem(&ShowItemHoverType{
			Item:  diamond,
			Count: 1,
			Components: map[string]interface{}{
				"minecraft:custom_name": map[string]interface{}{"text": "Spark"},
			},
		})},
	}

	encoded := new(strings.Builder)
	require.NoError(t, j1215Plus.Marshal(encoded, component))
	require.Equal(t, `{"hover_event":{"action":"show_item","components":{"minecraft:custom_name":{"text":"Spark"}},"count":1,"id":"minecraft:diamond"},"text":"Hover"}`, encoded.String())

	decoded, err := j1215Plus.Unmarshal([]byte(encoded.String()))
	require.NoError(t, err)
	require.Equal(t, component, decoded)
}

func TestJson_ShowItemModernComponentsWithGojay(t *testing.T) {
	diamond, _ := key.Parse("minecraft:diamond")
	component := &Text{
		Content: "Hover",
		S: Style{HoverEvent: ShowItem(&ShowItemHoverType{
			Item:  diamond,
			Count: 1,
			Components: map[string]interface{}{
				"minecraft:custom_name": map[string]interface{}{"text": "Spark"},
			},
		})},
	}

	encoded := new(strings.Builder)
	require.NoError(t, (&Json{
		UseLegacyFieldNames:          false,
		UseLegacyHoverEventStructure: false,
		EmitDefaultItemHoverQuantity: true,
		ShowItemHoverDataMode:        ShowItemHoverDataModeDataComponents,
	}).Marshal(encoded, component))
	require.JSONEq(t, `{"hover_event":{"action":"show_item","components":{"minecraft:custom_name":{"text":"Spark"}},"count":1,"id":"minecraft:diamond"},"text":"Hover"}`, encoded.String())
}

func TestJson_ShowItemHoverDataMode(t *testing.T) {
	diamond, _ := key.Parse("minecraft:diamond")
	component := &Text{
		Content: "Hover",
		S: Style{HoverEvent: ShowItem(&ShowItemHoverType{
			Item:  diamond,
			Count: 1,
			NBT:   nbt.NewBinaryTagHolder("{display:{Name:\"Legacy\"}}"),
			Components: map[string]interface{}{
				"minecraft:custom_name": map[string]interface{}{"text": "Modern"},
			},
		})},
	}

	t.Run("legacy_nbt", func(t *testing.T) {
		encoded := new(strings.Builder)
		require.NoError(t, (&Json{
			UseLegacyFieldNames:          true,
			UseLegacyHoverEventStructure: true,
			ShowItemHoverDataMode:        ShowItemHoverDataModeLegacyNBT,
			StdJson:                      true,
		}).Marshal(encoded, component))
		require.Contains(t, encoded.String(), `"tag":"{display:{Name:\"Legacy\"}}"`)
		require.NotContains(t, encoded.String(), `"components"`)
	})

	t.Run("data_components", func(t *testing.T) {
		encoded := new(strings.Builder)
		require.NoError(t, (&Json{
			UseLegacyFieldNames:          false,
			UseLegacyHoverEventStructure: false,
			EmitDefaultItemHoverQuantity: true,
			ShowItemHoverDataMode:        ShowItemHoverDataModeDataComponents,
			StdJson:                      true,
		}).Marshal(encoded, component))
		require.Contains(t, encoded.String(), `"components":{"minecraft:custom_name":{"text":"Modern"}}`)
		require.NotContains(t, encoded.String(), `"tag"`)
	})
}

func TestJson_ShowEntityUUIDIntArray(t *testing.T) {
	entityID := uuid.MustParse("12345678-1234-1234-1234-123456789abc")
	entityKey, _ := key.Parse("minecraft:player")
	component := &Text{
		Content: "Hover for entity",
		S: Style{HoverEvent: ShowEntity(&ShowEntityHoverType{
			Type: entityKey,
			Id:   entityID,
			Name: &Text{Content: "TestPlayer"},
		})},
	}

	encoded := new(strings.Builder)
	require.NoError(t, (&Json{
		UseLegacyFieldNames:                     false,
		UseLegacyHoverEventStructure:            false,
		EmitCompactTextComponent:                true,
		EmitHoverShowEntityIdAsIntArray:         true,
		EmitHoverShowEntityKeyAsTypeAndUuidAsId: false,
		StdJson:                                 true,
	}).Marshal(encoded, component))
	require.JSONEq(t, `{"hover_event":{"action":"show_entity","id":"minecraft:player","name":"TestPlayer","uuid":[305419896,305402420,305402420,1450744508]},"text":"Hover for entity"}`, encoded.String())

	decoded, err := j1215Plus.Unmarshal([]byte(encoded.String()))
	require.NoError(t, err)
	require.Equal(t, component, decoded)
}

func TestJson_EmitCompactTextComponent(t *testing.T) {
	encoded := new(strings.Builder)
	require.NoError(t, (&Json{
		EmitCompactTextComponent: true,
		StdJson:                  true,
	}).Marshal(encoded, &Text{Content: "plain"}))
	require.Equal(t, `"plain"`, encoded.String())

	component := &Text{
		Content: "root",
		Extra: []Component{
			&Text{Content: " child"},
			&Text{Content: " styled", S: Style{Bold: True}},
		},
	}
	encoded.Reset()
	require.NoError(t, (&Json{
		EmitCompactTextComponent: true,
		StdJson:                  true,
	}).Marshal(encoded, component))
	require.JSONEq(t, `{"extra":[" child",{"bold":true,"text":" styled"}],"text":"root"}`, encoded.String())
}

type testClickAction struct {
	name     string
	readable bool
}

func (a testClickAction) Name() string   { return a.name }
func (a testClickAction) Readable() bool { return a.readable }

type testHoverAction struct {
	name     string
	readable bool
}

func (a testHoverAction) Name() string { return a.name }
func (a testHoverAction) Type() ActionType {
	return ActionType(reflect.TypeOf(Text{}))
}
func (a testHoverAction) Readable() bool { return a.readable }

func TestJson_ValidateStrictEvents(t *testing.T) {
	strict := &Json{ValidateStrictEvents: true, StdJson: true}

	err := strict.Marshal(new(strings.Builder), &Text{
		Content: "bad click",
		S:       Style{ClickEvent: NewClickEvent(testClickAction{name: "not_real", readable: true}, "value")},
	})
	require.ErrorContains(t, err, "unsupported click event action")

	err = strict.Marshal(new(strings.Builder), &Text{
		Content: "bad hover",
		S:       Style{HoverEvent: NewHoverEvent(testHoverAction{name: "not_real", readable: true}, &Text{Content: "hover"})},
	})
	require.ErrorContains(t, err, "unsupported hover event action")

	_, err = strict.Unmarshal([]byte(`{"text":"bad","click_event":{"action":"not_real","value":"x"}}`))
	require.ErrorContains(t, err, "unsupported click event action")

	_, err = strict.Unmarshal([]byte(`{"text":"bad","hover_event":{"action":"not_real","value":{"text":"x"}}}`))
	require.ErrorContains(t, err, "unsupported hover event action")

	_, err = strict.Unmarshal([]byte(`{"text":"bad","click_event":{"action":"run_command"}}`))
	require.ErrorContains(t, err, "missing value for click event action")

	_, err = strict.Unmarshal([]byte(`{"text":"bad","hover_event":{"action":"show_item"}}`))
	require.ErrorContains(t, err, "missing value for hover event action")
}

func TestJson_ShowTextHoverAnyComponent(t *testing.T) {
	component := &Text{
		Content: "Hover",
		S: Style{
			HoverEvent: ShowText(&Translation{Key: "item.minecraft.diamond", Fallback: "Diamond"}),
		},
	}

	encoded := new(strings.Builder)
	require.NoError(t, j1215Plus.Marshal(encoded, component))
	require.Equal(t, `{"hover_event":{"action":"show_text","value":{"fallback":"Diamond","translate":"item.minecraft.diamond"}},"text":"Hover"}`, encoded.String())

	decoded, err := j1215Plus.Unmarshal([]byte(encoded.String()))
	require.NoError(t, err)
	require.Equal(t, component, decoded)
}

func TestJson_ObjectComponent_AtlasSprite(t *testing.T) {
	blocks, _ := key.Parse("minecraft:blocks")
	diamond, _ := key.Parse("minecraft:item/diamond")
	component := AtlasSprite(blocks, diamond)

	encoded := new(strings.Builder)
	err := j1215Plus.Marshal(encoded, component)
	require.NoError(t, err)
	require.Equal(t, `{"atlas":"minecraft:blocks","sprite":"minecraft:item/diamond"}`, encoded.String())

	decoded, err := j1215Plus.Unmarshal([]byte(encoded.String()))
	require.NoError(t, err)
	require.Equal(t, component, decoded)
}

func TestJson_ObjectComponent_DecodeAdventureSpriteShape(t *testing.T) {
	decoded, err := j1215Plus.Unmarshal([]byte(`{"fallback":"diamond","sprite":"item/diamond"}`))
	require.NoError(t, err)

	diamond, _ := key.Parse("minecraft:item/diamond")
	expected := AtlasSprite(DefaultSpriteAtlas, diamond)
	expected.Fallback = &Text{Content: "diamond"}
	require.Equal(t, expected, decoded)
}

func TestJson_ObjectComponent_PlayerHeadWithFallback(t *testing.T) {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789abc")
	component := PlayerHead(&PlayerProfile{
		Name: "jeb_",
		Id:   id,
	}, false)
	component.Fallback = &Text{Content: "jeb_"}

	encoded := new(strings.Builder)
	err := j1215Plus.Marshal(encoded, component)
	require.NoError(t, err)
	require.Equal(t, `{"fallback":{"text":"jeb_"},"hat":false,"player":{"id":"12345678-1234-1234-1234-123456789abc","name":"jeb_"}}`, encoded.String())

	decoded, err := j1215Plus.Unmarshal([]byte(encoded.String()))
	require.NoError(t, err)
	require.Equal(t, component, decoded)
}

func TestJson_ObjectComponent_PlayerNameShortcut(t *testing.T) {
	component := PlayerHead(&PlayerProfile{Name: "jeb_"}, true)

	encoded := new(strings.Builder)
	err := j1215Plus.Marshal(encoded, component)
	require.NoError(t, err)
	require.Equal(t, `{"hat":true,"player":"jeb_"}`, encoded.String())

	decoded, err := j1215Plus.Unmarshal([]byte(encoded.String()))
	require.NoError(t, err)
	require.Equal(t, component, decoded)
}

func TestJson_ObjectComponent_PlayerPropertiesMap(t *testing.T) {
	decoded, err := j1215Plus.Unmarshal([]byte(`{"player":{"name":"jeb_","properties":{"textures":["value-a","value-b"]}}}`))
	require.NoError(t, err)

	require.Equal(t, PlayerHead(&PlayerProfile{
		Name: "jeb_",
		Properties: []ProfileProperty{
			{Name: "textures", Value: "value-a"},
			{Name: "textures", Value: "value-b"},
		},
	}, true), decoded)
}

// Test encoding with new format (1.21.5+)
func TestJson_Marshal_NewFormat(t *testing.T) {
	b := new(strings.Builder)
	err := j1215Plus.Marshal(b, txt)
	require.NoError(t, err)
	require.Equal(t, jsonTxtNew, b.String())
	require.Contains(t, b.String(), `"click_event":`)
	require.Contains(t, b.String(), `"hover_event":`)
}

// Test encoding with legacy format (pre-1.21.5)
func TestJson_Marshal_LegacyFormat(t *testing.T) {
	b := new(strings.Builder)
	err := jPre1215.Marshal(b, txt)
	require.NoError(t, err)
	require.Equal(t, jsonTxtLegacy, b.String())
	require.Contains(t, b.String(), `"clickEvent":`)
	require.Contains(t, b.String(), `"hoverEvent":`)
}

// Test decoding new format with new codec
func TestJson_Unmarshal_NewFormat(t *testing.T) {
	c, err := j1215Plus.Unmarshal([]byte(jsonTxtNew))
	require.NoError(t, err)
	require.Equal(t, txt, c)
}

// Test decoding legacy format with legacy codec
func TestJson_Unmarshal_LegacyFormat(t *testing.T) {
	c, err := jPre1215.Unmarshal([]byte(jsonTxtLegacy))
	require.NoError(t, err)
	require.Equal(t, txt, c)
}

// Test cross-compatibility: decode new format with legacy codec
func TestJson_CrossCompat_NewFormatWithLegacyCodec(t *testing.T) {
	c, err := jPre1215.Unmarshal([]byte(jsonTxtNew))
	require.NoError(t, err)
	require.Equal(t, txt, c)
}

// Test cross-compatibility: decode legacy format with new codec
func TestJson_CrossCompat_LegacyFormatWithNewCodec(t *testing.T) {
	c, err := j1215Plus.Unmarshal([]byte(jsonTxtLegacy))
	require.NoError(t, err)
	require.Equal(t, txt, c)
}

// Test cross-compatibility: decode both formats with compatibility codec
func TestJson_CrossCompat_BothFormatsWithCompatCodec(t *testing.T) {
	// Test new format
	c1, err := jCompat.Unmarshal([]byte(jsonTxtNew))
	require.NoError(t, err)
	require.Equal(t, txt, c1)

	// Test legacy format
	c2, err := jCompat.Unmarshal([]byte(jsonTxtLegacy))
	require.NoError(t, err)
	require.Equal(t, txt, c2)
}

// Test round-trip: encode with one format, decode with another
func TestJson_RoundTrip_CrossFormat(t *testing.T) {
	// Encode with legacy, decode with new
	legacyEncoded := new(strings.Builder)
	err := jPre1215.Marshal(legacyEncoded, txt)
	require.NoError(t, err)

	decoded, err := j1215Plus.Unmarshal([]byte(legacyEncoded.String()))
	require.NoError(t, err)
	require.Equal(t, txt, decoded)

	// Encode with new, decode with legacy
	newEncoded := new(strings.Builder)
	err = j1215Plus.Marshal(newEncoded, txt)
	require.NoError(t, err)

	decoded2, err := jPre1215.Unmarshal([]byte(newEncoded.String()))
	require.NoError(t, err)
	require.Equal(t, txt, decoded2)
}

// Test that click events work in both formats
func TestJson_ClickEvent_BothFormats(t *testing.T) {
	component := &Text{
		Content: "Click me",
		S: Style{
			ClickEvent: CopyToClipboard("test"),
		},
	}

	// Test new format
	newEncoded := new(strings.Builder)
	err := j1215Plus.Marshal(newEncoded, component)
	require.NoError(t, err)
	require.Contains(t, newEncoded.String(), `"click_event":`)
	require.Contains(t, newEncoded.String(), `"copy_to_clipboard"`)

	decoded, err := j1215Plus.Unmarshal([]byte(newEncoded.String()))
	require.NoError(t, err)
	require.Equal(t, component, decoded)

	// Test legacy format
	legacyEncoded := new(strings.Builder)
	err = jPre1215.Marshal(legacyEncoded, component)
	require.NoError(t, err)
	require.Contains(t, legacyEncoded.String(), `"clickEvent":`)
	require.Contains(t, legacyEncoded.String(), `"copy_to_clipboard"`)

	decoded2, err := jPre1215.Unmarshal([]byte(legacyEncoded.String()))
	require.NoError(t, err)
	require.Equal(t, component, decoded2)
}

// Test that hover events work in both formats
func TestJson_HoverEvent_BothFormats(t *testing.T) {
	component := &Text{
		Content: "Hover me",
		S: Style{
			HoverEvent: ShowText(&Text{Content: "Tooltip"}),
		},
	}

	// Test new format
	newEncoded := new(strings.Builder)
	err := j1215Plus.Marshal(newEncoded, component)
	require.NoError(t, err)
	require.Contains(t, newEncoded.String(), `"hover_event":`)
	require.Contains(t, newEncoded.String(), `"show_text"`)

	decoded, err := j1215Plus.Unmarshal([]byte(newEncoded.String()))
	require.NoError(t, err)
	require.Equal(t, component, decoded)

	// Test legacy format
	legacyEncoded := new(strings.Builder)
	err = jPre1215.Marshal(legacyEncoded, component)
	require.NoError(t, err)
	require.Contains(t, legacyEncoded.String(), `"hoverEvent":`)
	require.Contains(t, legacyEncoded.String(), `"show_text"`)

	decoded2, err := jPre1215.Unmarshal([]byte(legacyEncoded.String()))
	require.NoError(t, err)
	require.Equal(t, component, decoded2)
}

// Test all click event actions in different formats
func TestJson_ClickEvent_AllActions(t *testing.T) {
	testCases := []struct {
		name      string
		component *Text
		action    string
		value     string
	}{
		{
			name: "open_url",
			component: &Text{
				Content: "Open URL",
				S:       Style{ClickEvent: OpenUrl("https://example.com")},
			},
			action: "open_url",
			value:  "https://example.com",
		},
		{
			name: "run_command",
			component: &Text{
				Content: "Run Command",
				S:       Style{ClickEvent: RunCommand("/say hello")},
			},
			action: "run_command",
			value:  "/say hello",
		},
		{
			name: "suggest_command",
			component: &Text{
				Content: "Suggest Command",
				S:       Style{ClickEvent: SuggestCommand("/help")},
			},
			action: "suggest_command",
			value:  "/help",
		},
		{
			name: "copy_to_clipboard",
			component: &Text{
				Content: "Copy Text",
				S:       Style{ClickEvent: CopyToClipboard("copied text")},
			},
			action: "copy_to_clipboard",
			value:  "copied text",
		},
		{
			name: "change_page",
			component: &Text{
				Content: "Change Page",
				S:       Style{ClickEvent: ChangePage("3")},
			},
			action: "change_page",
			value:  "3",
		},
		{
			name: "show_dialog",
			component: &Text{
				Content: "Show Dialog",
				S:       Style{ClickEvent: ShowDialog("my_dialog_id")},
			},
			action: "show_dialog",
			value:  "my_dialog_id",
		},
		{
			name: "custom",
			component: &Text{
				Content: "Custom Event",
				S:       Style{ClickEvent: CustomEvent("my_event", "some_payload")},
			},
			action: "custom",
			value:  "my_event|some_payload",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test pre-1.21.5 format
			t.Run("pre_1215_format", func(t *testing.T) {
				encoded := new(strings.Builder)
				err := jPre1215.Marshal(encoded, tc.component)
				require.NoError(t, err)

				// Should contain legacy field names and structure
				require.Contains(t, encoded.String(), `"clickEvent"`)
				require.Contains(t, encoded.String(), `"action":"`+tc.action+`"`)
				require.Contains(t, encoded.String(), `"value":"`+tc.value+`"`)

				// Should be decodable
				decoded, err := jPre1215.Unmarshal([]byte(encoded.String()))
				require.NoError(t, err)
				require.Equal(t, tc.component, decoded)
			})

			// Test 1.21.5+ format
			t.Run("1215_plus_format", func(t *testing.T) {
				encoded := new(strings.Builder)
				err := j1215Plus.Marshal(encoded, tc.component)
				require.NoError(t, err)

				// Should contain new field names
				require.Contains(t, encoded.String(), `"click_event"`)
				require.Contains(t, encoded.String(), `"action":"`+tc.action+`"`)

				// Check for specific field names based on action
				switch tc.action {
				case "open_url":
					require.Contains(t, encoded.String(), `"url":"`+tc.value+`"`)
				case "run_command", "suggest_command":
					require.Contains(t, encoded.String(), `"command":"`+tc.value+`"`)
				case "change_page":
					require.Contains(t, encoded.String(), `"page":3`) // Should be converted to int
				case "copy_to_clipboard":
					require.Contains(t, encoded.String(), `"value":"`+tc.value+`"`) // Still uses "value"
				}

				// Should be decodable
				decoded, err := j1215Plus.Unmarshal([]byte(encoded.String()))
				require.NoError(t, err)
				require.Equal(t, tc.component, decoded)
			})

			// Test cross-compatibility
			t.Run("cross_compatibility", func(t *testing.T) {
				// Encode with legacy, decode with new
				legacyEncoded := new(strings.Builder)
				err := jPre1215.Marshal(legacyEncoded, tc.component)
				require.NoError(t, err)

				decoded1, err := j1215Plus.Unmarshal([]byte(legacyEncoded.String()))
				require.NoError(t, err)
				require.Equal(t, tc.component, decoded1)

				// Encode with new, decode with legacy
				newEncoded := new(strings.Builder)
				err = j1215Plus.Marshal(newEncoded, tc.component)
				require.NoError(t, err)

				decoded2, err := jPre1215.Unmarshal([]byte(newEncoded.String()))
				require.NoError(t, err)
				require.Equal(t, tc.component, decoded2)
			})
		})
	}
}

// Test all hover event actions in different formats
func TestJson_HoverEvent_AllActions(t *testing.T) {
	// Test show_text hover event
	t.Run("show_text", func(t *testing.T) {
		component := &Text{
			Content: "Hover for text",
			S: Style{
				HoverEvent: ShowText(&Text{Content: "Tooltip text", S: Style{Color: Red.RGB}}),
			},
		}

		// Test pre-1.21.5 format (uses "contents" field)
		t.Run("pre_1215_format", func(t *testing.T) {
			encoded := new(strings.Builder)
			err := jPre1215.Marshal(encoded, component)
			require.NoError(t, err)

			// Should contain legacy structure
			require.Contains(t, encoded.String(), `"hoverEvent"`)
			require.Contains(t, encoded.String(), `"action":"show_text"`)
			require.Contains(t, encoded.String(), `"contents"`)

			decoded, err := jPre1215.Unmarshal([]byte(encoded.String()))
			require.NoError(t, err)
			require.Equal(t, component, decoded)
		})

		// Test 1.21.5+ format (uses inlined "value" field)
		t.Run("1215_plus_format", func(t *testing.T) {
			encoded := new(strings.Builder)
			err := j1215Plus.Marshal(encoded, component)
			require.NoError(t, err)

			// Should contain new structure
			require.Contains(t, encoded.String(), `"hover_event"`)
			require.Contains(t, encoded.String(), `"action":"show_text"`)
			require.Contains(t, encoded.String(), `"value"`) // Direct value field, not contents
			require.NotContains(t, encoded.String(), `"contents"`)

			decoded, err := j1215Plus.Unmarshal([]byte(encoded.String()))
			require.NoError(t, err)
			require.Equal(t, component, decoded)
		})
	})

	// Test show_item hover event
	t.Run("show_item", func(t *testing.T) {
		itemKey, _ := key.Parse("minecraft:diamond")
		component := &Text{
			Content: "Hover for item",
			S: Style{
				HoverEvent: ShowItem(&ShowItemHoverType{
					Item:  itemKey,
					Count: 5,
					NBT:   nbt.NewBinaryTagHolder("{display:{Name:\"Special Diamond\"}}"),
				}),
			},
		}

		// Test pre-1.21.5 format (uses "contents" field)
		t.Run("pre_1215_format", func(t *testing.T) {
			encoded := new(strings.Builder)
			err := jPre1215.Marshal(encoded, component)
			require.NoError(t, err)

			// Should contain legacy structure
			require.Contains(t, encoded.String(), `"hoverEvent"`)
			require.Contains(t, encoded.String(), `"action":"show_item"`)
			require.Contains(t, encoded.String(), `"contents"`)
			require.Contains(t, encoded.String(), `"minecraft:diamond"`)

			decoded, err := jPre1215.Unmarshal([]byte(encoded.String()))
			require.NoError(t, err)
			require.Equal(t, component, decoded)
		})

		// Test 1.21.5+ format (inlined fields)
		t.Run("1215_plus_format", func(t *testing.T) {
			encoded := new(strings.Builder)
			err := j1215Plus.Marshal(encoded, component)
			require.NoError(t, err)

			// Should contain new inlined structure
			require.Contains(t, encoded.String(), `"hover_event"`)
			require.Contains(t, encoded.String(), `"action":"show_item"`)
			require.Contains(t, encoded.String(), `"id":"minecraft:diamond"`) // Direct field
			require.Contains(t, encoded.String(), `"count":5`)                // Direct field
			require.NotContains(t, encoded.String(), `"contents"`)

			decoded, err := j1215Plus.Unmarshal([]byte(encoded.String()))
			require.NoError(t, err)
			require.Equal(t, component, decoded)
		})
	})

	// Test show_entity hover event
	t.Run("show_entity", func(t *testing.T) {
		uuid := uuid.MustParse("12345678-1234-1234-1234-123456789abc")
		entityKey, _ := key.Parse("minecraft:player")
		component := &Text{
			Content: "Hover for entity",
			S: Style{
				HoverEvent: ShowEntity(&ShowEntityHoverType{
					Type: entityKey,
					Id:   uuid,
					Name: &Text{Content: "TestPlayer", S: Style{Color: Blue.RGB}},
				}),
			},
		}

		// Test pre-1.21.5 format (uses "contents" with old field names)
		t.Run("pre_1215_format", func(t *testing.T) {
			encoded := new(strings.Builder)
			err := jPre1215.Marshal(encoded, component)
			require.NoError(t, err)

			// Should contain legacy structure with old field names
			require.Contains(t, encoded.String(), `"hoverEvent"`)
			require.Contains(t, encoded.String(), `"action":"show_entity"`)
			require.Contains(t, encoded.String(), `"contents"`)
			require.Contains(t, encoded.String(), `"type":"minecraft:player"`) // Legacy field name
			require.Contains(t, encoded.String(), `"id":[305419896,305402420,305402420,1450744508]`)

			decoded, err := jPre1215.Unmarshal([]byte(encoded.String()))
			require.NoError(t, err)
			require.Equal(t, component, decoded)
		})

		// Test 1.21.5+ format (inlined with new field names)
		t.Run("1215_plus_format", func(t *testing.T) {
			encoded := new(strings.Builder)
			err := j1215Plus.Marshal(encoded, component)
			require.NoError(t, err)

			// Should contain new inlined structure with new field names
			require.Contains(t, encoded.String(), `"hover_event"`)
			require.Contains(t, encoded.String(), `"action":"show_entity"`)
			require.Contains(t, encoded.String(), `"id":"minecraft:player"`) // New field name (was "type")
			require.Contains(t, encoded.String(), `"uuid":[305419896,305402420,305402420,1450744508]`)
			require.NotContains(t, encoded.String(), `"contents"`)
			require.NotContains(t, encoded.String(), `"type":"minecraft:player"`) // Should not use legacy field name

			decoded, err := j1215Plus.Unmarshal([]byte(encoded.String()))
			require.NoError(t, err)
			require.Equal(t, component, decoded)
		})

		// Test cross-compatibility for field name changes
		t.Run("cross_compatibility_field_names", func(t *testing.T) {
			// Test that new decoder can read legacy field names
			legacyJSON := `{"text":"Hover for entity","hover_event":{"action":"show_entity","contents":{"type":"minecraft:player","id":"12345678-1234-1234-1234-123456789abc","name":{"text":"TestPlayer","color":"#5555ff"}}}}`
			decoded, err := j1215Plus.Unmarshal([]byte(legacyJSON))
			require.NoError(t, err)
			require.Equal(t, component, decoded)

			// Test that legacy decoder can read new field names
			newJSON := `{"text":"Hover for entity","hover_event":{"action":"show_entity","id":"minecraft:player","uuid":"12345678-1234-1234-1234-123456789abc","name":{"text":"TestPlayer","color":"#5555ff"}}}`
			decoded2, err := jPre1215.Unmarshal([]byte(newJSON))
			require.NoError(t, err)
			require.Equal(t, component, decoded2)
		})
	})
}

// TestJSONOptions tests the various JSON options and configurations
func TestJSONOptions(t *testing.T) {
	// Test EMIT_DEFAULT_ITEM_HOVER_QUANTITY option
	t.Run("emit_default_item_hover_quantity", func(t *testing.T) {
		diamond, _ := key.Parse("minecraft:diamond")
		component := &Text{
			Content: "Hover for item",
			S: Style{
				HoverEvent: ShowItem(&ShowItemHoverType{
					Item:  diamond,
					Count: 1, // Default count
					NBT:   nil,
				}),
			},
		}

		// Test with EmitDefaultItemHoverQuantity = false (default)
		t.Run("default_false", func(t *testing.T) {
			codec := &Json{
				NoLegacyHover:                false,
				NoDownsampleColor:            true,
				UseLegacyFieldNames:          true,
				UseLegacyClickEventStructure: true,
				UseLegacyHoverEventStructure: true,
				EmitDefaultItemHoverQuantity: false, // Don't emit count=1
				StdJson:                      true,
			}

			encoded := new(strings.Builder)
			err := codec.Marshal(encoded, component)
			require.NoError(t, err)

			// Should NOT contain count field when count=1 and EmitDefaultItemHoverQuantity=false
			require.NotContains(t, encoded.String(), `"count":1`)
			require.Contains(t, encoded.String(), `"id":"minecraft:diamond"`)
		})

		// Test with EmitDefaultItemHoverQuantity = true
		t.Run("emit_true", func(t *testing.T) {
			codec := &Json{
				NoLegacyHover:                false,
				NoDownsampleColor:            true,
				UseLegacyFieldNames:          true,
				UseLegacyClickEventStructure: true,
				UseLegacyHoverEventStructure: true,
				EmitDefaultItemHoverQuantity: true, // Emit count=1
				StdJson:                      true,
			}

			encoded := new(strings.Builder)
			err := codec.Marshal(encoded, component)
			require.NoError(t, err)

			// Should contain count field even when count=1
			require.Contains(t, encoded.String(), `"count":1`)
			require.Contains(t, encoded.String(), `"id":"minecraft:diamond"`)
		})
	})

	// Test EMIT_HOVER_SHOW_ENTITY_KEY_AS_TYPE_AND_UUID_AS_ID option
	t.Run("emit_hover_show_entity_field_names", func(t *testing.T) {
		uuid := uuid.MustParse("12345678-1234-1234-1234-123456789abc")
		entityKey, _ := key.Parse("minecraft:player")
		component := &Text{
			Content: "Hover for entity",
			S: Style{
				HoverEvent: ShowEntity(&ShowEntityHoverType{
					Type: entityKey,
					Id:   uuid,
					Name: &Text{Content: "TestPlayer", S: Style{Color: Blue.RGB}},
				}),
			},
		}

		// Test with EmitHoverShowEntityKeyAsTypeAndUuidAsId = true (legacy field names)
		t.Run("legacy_field_names", func(t *testing.T) {
			codec := &Json{
				NoLegacyHover:                           true,
				NoDownsampleColor:                       true,
				UseLegacyFieldNames:                     false, // Use snake_case
				UseLegacyClickEventStructure:            false,
				UseLegacyHoverEventStructure:            false, // Use inlined structure
				EmitHoverShowEntityKeyAsTypeAndUuidAsId: true,  // Use legacy field names
				StdJson:                                 true,
			}

			encoded := new(strings.Builder)
			err := codec.Marshal(encoded, component)
			require.NoError(t, err)

			// Should use legacy field names even in modern structure
			require.Contains(t, encoded.String(), `"type":"minecraft:player"`)                   // Legacy field name
			require.Contains(t, encoded.String(), `"id":"12345678-1234-1234-1234-123456789abc"`) // Legacy field name
			require.NotContains(t, encoded.String(), `"uuid":"12345678-1234-1234-1234-123456789abc"`)
		})

		// Test with EmitHoverShowEntityKeyAsTypeAndUuidAsId = false (modern field names)
		t.Run("modern_field_names", func(t *testing.T) {
			codec := &Json{
				NoLegacyHover:                           true,
				NoDownsampleColor:                       true,
				UseLegacyFieldNames:                     false, // Use snake_case
				UseLegacyClickEventStructure:            false,
				UseLegacyHoverEventStructure:            false, // Use inlined structure
				EmitHoverShowEntityKeyAsTypeAndUuidAsId: false, // Use modern field names
				StdJson:                                 true,
			}

			encoded := new(strings.Builder)
			err := codec.Marshal(encoded, component)
			require.NoError(t, err)

			// Should use modern field names
			require.Contains(t, encoded.String(), `"id":"minecraft:player"`)                       // Modern field name (was "type")
			require.Contains(t, encoded.String(), `"uuid":"12345678-1234-1234-1234-123456789abc"`) // Modern field name (was "id")
			require.NotContains(t, encoded.String(), `"type":"minecraft:player"`)
		})
	})

	// Test EMIT_CHANGE_PAGE_CLICK_EVENT_PAGE_AS_STRING option
	t.Run("emit_change_page_as_string", func(t *testing.T) {
		component := &Text{
			Content: "Click to change page",
			S: Style{
				ClickEvent: ChangePage("3"),
			},
		}

		// Test with EmitChangePageClickEventPageAsString = true (legacy string format)
		t.Run("page_as_string", func(t *testing.T) {
			codec := &Json{
				NoLegacyHover:                        false,
				NoDownsampleColor:                    true,
				UseLegacyFieldNames:                  false, // Use snake_case
				UseLegacyClickEventStructure:         false, // Use specific fields
				EmitChangePageClickEventPageAsString: true,  // Emit as string
				StdJson:                              true,
			}

			encoded := new(strings.Builder)
			err := codec.Marshal(encoded, component)
			require.NoError(t, err)

			// Should contain page as string
			require.Contains(t, encoded.String(), `"page":"3"`)
			require.NotContains(t, encoded.String(), `"page":3`)
		})

		// Test with EmitChangePageClickEventPageAsString = false (modern integer format)
		t.Run("page_as_integer", func(t *testing.T) {
			codec := &Json{
				NoLegacyHover:                        false,
				NoDownsampleColor:                    true,
				UseLegacyFieldNames:                  false, // Use snake_case
				UseLegacyClickEventStructure:         false, // Use specific fields
				EmitChangePageClickEventPageAsString: false, // Emit as integer
				StdJson:                              true,
			}

			encoded := new(strings.Builder)
			err := codec.Marshal(encoded, component)
			require.NoError(t, err)

			// Should contain page as integer
			require.Contains(t, encoded.String(), `"page":3`)
			require.NotContains(t, encoded.String(), `"page":"3"`)
		})
	})
}
