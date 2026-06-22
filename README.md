# Minekube Commons Library

**Get it:**
`go get -u go.minekube.com/common`

## Minecraft Text Components (`minecraft/component`)

A comprehensive Minecraft text components library with **multi-version support** for Minecraft text component formats from legacy JSON through modern object components.

### 🎯 Key Features

- **Multi-Version Support**: Supports legacy JSON through modern 1.21.9+/26.1 component shapes
- **Backward & Forward Compatibility**: Can decode any format regardless of encoding settings
- **Configurable Encoding**: Choose between formats for different Minecraft versions
- **High Performance**: Uses optimized JSON marshalling (thanks to [Gojay](https://github.com/francoispqt/gojay))
- **Complete Component Support**: Text, translation, score, selector, keybind, NBT, object, click/hover events, styling, and more
- **Modern Object Components**: Supports 1.21.9+ atlas sprites/player heads and 26.1 object fallbacks

### 🚀 Quick Start

```go
import "go.minekube.com/common/minecraft/component/codec"

// Create a codec for the modern format (1.21.5+)
j := &codec.Json{
    UseLegacyFieldNames:          false, // Use snake_case field names
    UseLegacyClickEventStructure: false, // Use specific field names
    UseLegacyHoverEventStructure: false, // Use inlined structure
    StdJson:                      true,
}

// Create a text component with click and hover events
component := &component.Text{
    Content: "Click me!",
    S: component.Style{
        Color:      color.Green.RGB,
        ClickEvent: component.OpenUrl("https://example.com"),
        HoverEvent: component.ShowText(&component.Text{Content: "Tooltip"}),
    },
}

// Encode to JSON
var buf strings.Builder
err := j.Marshal(&buf, component)
// Output: {"text":"Click me!","color":"#55ff55","click_event":{"action":"open_url","url":"https://example.com"},"hover_event":{"action":"show_text","value":{"text":"Tooltip"}}}
```

### 📋 Supported Click Events

| Action              | Legacy Field | Modern Field(s) | Version | Description                       |
| ------------------- | ------------ | --------------- | ------- | --------------------------------- |
| `open_url`          | `value`      | `url`           | 1.21.5+ | Opens URL in browser              |
| `open_file`         | `value`      | `path`          | 1.21.5+ | Opens file (client-side only)     |
| `run_command`       | `value`      | `command`       | 1.21.5+ | Executes command                  |
| `suggest_command`   | `value`      | `command`       | 1.21.5+ | Suggests command in chat          |
| `change_page`       | `value`      | `page`          | 1.21.5+ | Changes book page (int in modern) |
| `copy_to_clipboard` | `value`      | `value`         | All     | Copies text (unchanged)           |
| `show_dialog`       | `value`      | `dialog`        | 1.21.6+ | Opens dialog                      |
| `custom`            | `value`      | `id`/`payload`  | 1.21.6+ | Custom server event               |

### 📦 Supported Content Types

| Component type | JSON key(s) | Notes |
| -------------- | ----------- | ----- |
| Text | `text` | Also decodes JSON string, boolean, and number shorthands |
| Translation | `translate`, `fallback`, `with` | Supports 1.19.4+ fallback |
| Score | `score` | Includes obsolete `value` for compatibility |
| Selector | `selector`, `separator` | Supports component separator |
| Keybind | `keybind` | Client-local keybind rendering |
| NBT | `nbt`, `block`/`entity`/`storage`, `interpret`, `plain`, `separator` | Supports 26.1 `plain` |
| Object | `atlas`/`sprite`, `player`, `hat`, `fallback` | Supports 1.21.9+ sprites/player heads and 26.1 fallback |

### 🎛️ Version Configuration Examples

**Using Preset Configurations (Recommended):**

```go
import "go.minekube.com/common/minecraft/component/codec"

// For clients before 1.16
j := codec.JsonPre1_16

// For clients 1.16+ but before 1.20.3
j := codec.JsonPre1_20_3

// For clients 1.20.3+ but before 1.21.5
j := codec.JsonPre1_21_5

// For clients 1.21.5+ (modern format)
j := codec.JsonModern

// Universal compatibility (encodes modern, decodes all)
j := codec.JsonUniversal
```

**Manual Configuration:**

```go
// Modern Format (1.21.5+ base, with newer component types supported)
j := &codec.Json{
    UseLegacyFieldNames:          false, // snake_case: click_event, hover_event
    UseLegacyClickEventStructure: false, // Specific fields: url, path, command, etc.
    UseLegacyHoverEventStructure: false, // Inlined hover structure
    StdJson:                      true,
}

// Legacy Format (Pre-1.21.5)
j := &codec.Json{
    UseLegacyFieldNames:          true,  // camelCase: clickEvent, hoverEvent
    UseLegacyClickEventStructure: true,  // Universal "value" field
    UseLegacyHoverEventStructure: true,  // "contents" wrapper structure
    NoLegacyHover:                false, // Include both contents and value
    StdJson:                      true,
}
```

**Preset Configurations:**

| Preset          | Target Version  | Field Names | Click Events      | Hover Events      | Features                    |
| --------------- | --------------- | ----------- | ----------------- | ----------------- | --------------------------- |
| `JsonPre1_16`   | Before 1.16     | camelCase   | Universal "value" | Legacy + value    | Color downsampling          |
| `JsonPre1_20_3` | 1.16 - 1.20.2   | camelCase   | Universal "value" | Legacy + contents | Hex colors                  |
| `JsonPre1_21_5` | 1.20.3 - 1.21.4 | camelCase   | Universal "value" | Legacy + contents | Compact text, int UUIDs     |
| `JsonModern`    | 1.21.5+         | snake_case  | Specific fields   | Inlined           | All modern features         |
| `JsonUniversal` | Any             | snake_case  | Specific fields   | Inlined           | Decodes all, encodes modern |

### ⚙️ Advanced JSON Configuration Options

For fine-grained control over JSON serialization behavior, additional options are available that correspond to Minecraft's internal JSON serialization flags:

```go
j := &codec.Json{
    // Basic structure control
    UseLegacyFieldNames:          false, // snake_case vs camelCase field names
    UseLegacyClickEventStructure: false, // Specific fields vs universal "value"
    UseLegacyHoverEventStructure: false, // Inlined vs "contents" wrapper

    // Version-specific behavior control
    EmitChangePageClickEventPageAsString:     false, // Page as int (1.21.6+) vs string (legacy)
    EmitCompactTextComponent:                 true,  // Plain text optimization (1.20.3+)
    EmitHoverShowEntityIdAsIntArray:          true,  // UUID as int array (1.20.3+) vs string
    EmitHoverShowEntityKeyAsTypeAndUuidAsId:  false, // Modern ("id"/"uuid") vs legacy ("type"/"id") field names
    ValidateStrictEvents:                     true,  // Strict event validation (1.20.3+)
    EmitDefaultItemHoverQuantity:             true,  // Always emit count=1 (1.20.5+)

    // Advanced formatting modes
    ShowItemHoverDataMode: codec.ShowItemHoverDataModeEither, // Legacy NBT vs modern data components
    ShadowColorMode:       codec.ShadowColorEmitModeInteger,          // Shadow color format (1.21.4+)

    StdJson: true,
}
```

**Configuration Options:**

| Option                                    | Description                               | Default     | Since Version |
| ----------------------------------------- | ----------------------------------------- | ----------- | ------------- |
| `EmitChangePageClickEventPageAsString`    | Emit page numbers as strings vs integers  | `false`     | 1.21.6+       |
| `EmitCompactTextComponent`                | Use plain text for simple components      | `false`     | 1.20.3+       |
| `EmitHoverShowEntityIdAsIntArray`         | UUID as `[int, int, int, int]` vs string  | `false`     | 1.20.3+       |
| `EmitHoverShowEntityKeyAsTypeAndUuidAsId` | Use legacy field names (`"type"`, `"id"`) | `true`      | Pre-1.21.5    |
| `ValidateStrictEvents`                    | Enable strict event validation            | `false`     | 1.20.3+       |
| `EmitDefaultItemHoverQuantity`            | Always emit `count: 1` for items          | `false`     | 1.20.5+       |
| `ShowItemHoverDataMode`                   | Item data format (NBT vs data components) | `LegacyNBT` | 1.20.5+       |
| `ShadowColorMode`                         | Shadow color emission format              | `None`      | 1.21.4+       |

**ShowItemHoverDataMode Values:**

- `ShowItemHoverDataModeLegacyNBT` - Use legacy NBT format
- `ShowItemHoverDataModeDataComponents` - Use modern data components
- `ShowItemHoverDataModeEither` - Use whichever the item has

**ShadowColorMode Values:**

- `ShadowColorEmitModeNone` - Don't emit shadow colors
- `ShadowColorEmitModeInteger` - Emit as packed ARGB integer
- `ShadowColorEmitModeArray` - Emit as `[r, g, b, a]` float array

### 🆕 Modern Minecraft Features

```go
// Dialog click event (1.21.6+)
dialogComponent := &component.Text{
    Content: "Open Settings",
    S: component.Style{
        ClickEvent: component.ShowDialog("settings_dialog"),
    },
}

// Custom server events (1.21.6+)
customComponent := &component.Text{
    Content: "Trigger Event",
    S: component.Style{
        ClickEvent: component.CustomEvent("my_event", "some_payload"),
    },
}

// Atlas sprite object component (1.21.9+), usable for item/block sprites in text
atlas, _ := key.Parse("minecraft:blocks")
sprite, _ := key.Parse("minecraft:item/diamond")
spriteComponent := component.AtlasSprite(atlas, sprite)

// Player head object component with a 26.1 fallback for contexts like server MOTDs
playerComponent := component.PlayerHead(&component.PlayerProfile{Name: "jeb_"}, true)
playerComponent.Fallback = &component.Text{Content: "jeb_"}

// Score, selector, keybind, and NBT content types
scoreComponent := &component.Score{Name: "@s", Objective: "kills"}
selectorComponent := &component.Selector{Pattern: "@a", Separator: &component.Text{Content: ", "}}
keybindComponent := &component.Keybind{Key: "key.jump"}
nbtComponent := &component.EntityNBT{
    NBT:    component.NBT{Path: "Health", Plain: true},
    Entity: "@s",
}
```

### 🔄 Format Examples

**Atlas Sprite Object (1.21.9+):**

```json
{
  "atlas": "minecraft:blocks",
  "sprite": "minecraft:item/diamond"
}
```

**Player Head Object With Fallback (26.1+):**

```json
{
  "fallback": {
    "text": "jeb_"
  },
  "hat": true,
  "player": "jeb_"
}
```

**Modern Hover Item Components (1.20.5+):**

```json
{
  "text": "Hover",
  "hover_event": {
    "action": "show_item",
    "id": "minecraft:diamond",
    "count": 1,
    "components": {
      "minecraft:custom_name": {
        "text": "Spark"
      }
    }
  }
}
```

**Modern Click/Hover Format (1.21.5+):**

```json
{
  "text": "Click me",
  "click_event": {
    "action": "open_url",
    "url": "https://example.com"
  },
  "hover_event": {
    "action": "show_item",
    "id": "minecraft:diamond",
    "count": 5
  }
}
```

**Legacy Format (Pre-1.21.5):**

```json
{
  "text": "Click me",
  "clickEvent": {
    "action": "open_url",
    "value": "https://example.com"
  },
  "hoverEvent": {
    "action": "show_item",
    "contents": {
      "id": "minecraft:diamond",
      "count": 5
    }
  }
}
```

### ✨ Additional Features

- **Legacy colors & formats**: Support for legacy color codes
- **Minecraft 1.16+ hex colors**: Full hex color support (`#ff5555`)
- **Minecraft 1.21.4+ shadow colors**: Packed integer and `[r, g, b, a]` float-array formats
- **Hover events**: `show_text`, `show_item`, `show_entity` with all format variations
- **Modern item hover data**: Supports 1.20.5+ `components` next to legacy `tag`
- **Object components**: 1.21.9+ atlas sprites and player heads, including 26.1 fallbacks
- **Translations**: Full translation component support with arguments and fallback
- **Cross-version compatibility**: Decode any format, encode in your preferred format
- **Performance optimized**: Much faster than Go's standard JSON encoding
- **Comprehensive testing**: Extensive test coverage for all versions and formats

### 🧪 Testing, formatting, linting

```bash
make
```

All format variations and cross-compatibility scenarios are thoroughly tested to ensure reliability across all Minecraft versions.
