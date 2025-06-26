# Minekube Commons Library

**Get it:**
`go get -u go.minekube.com/common`

## Minecraft Text Components (`minecraft/component`)

A comprehensive Minecraft text components library with **full multi-version support** for all Minecraft versions from legacy to the latest (1.21.6+).

### 🎯 Key Features

- **Complete Multi-Version Support**: Supports all Minecraft versions from legacy to 1.21.6+
- **Backward & Forward Compatibility**: Can decode any format regardless of encoding settings
- **Configurable Encoding**: Choose between formats for different Minecraft versions
- **High Performance**: Uses optimized JSON marshalling (thanks to [Gojay](https://github.com/francoispqt/gojay))
- **Complete Component Support**: Text, translation, click events, hover events, styling, and more
- **Future-Ready**: Includes upcoming 1.21.6+ features like dialog and custom click events

### 🚀 Quick Start

```go
import "go.minekube.com/common/minecraft/component/codec"

// Create a codec for the latest format (1.21.5+)
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

### 🎛️ Version Configuration Examples

**Latest Format (1.21.6+):**

```go
j := &codec.Json{
    UseLegacyFieldNames:          false, // snake_case: click_event, hover_event
    UseLegacyClickEventStructure: false, // Specific fields: url, path, command, etc.
    UseLegacyHoverEventStructure: false, // Inlined hover structure
    StdJson:                      true,
}
```

**Legacy Format (Pre-1.21.5):**

```go
j := &codec.Json{
    UseLegacyFieldNames:          true,  // camelCase: clickEvent, hoverEvent
    UseLegacyClickEventStructure: true,  // Universal "value" field
    UseLegacyHoverEventStructure: true,  // "contents" wrapper structure
    NoLegacyHover:                false, // Include both contents and value
    StdJson:                      true,
}
```

**Universal Compatibility (Recommended):**

```go
j := &codec.Json{
    // Encodes in modern format but can decode ANY format
    UseLegacyFieldNames:          false,
    UseLegacyClickEventStructure: false,
    UseLegacyHoverEventStructure: false,
    StdJson:                      true,
}
```

### 🆕 New Minecraft 1.21.6 Features

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
```

### 🔄 Format Examples

**Modern Format (1.21.5+):**

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
- **Hover events**: `show_text`, `show_item`, `show_entity` with all format variations
- **Translations**: Full translation component support with arguments
- **Cross-version compatibility**: Decode any format, encode in your preferred format
- **Performance optimized**: Much faster than Go's standard JSON encoding
- **Comprehensive testing**: Extensive test coverage for all versions and formats

### 🧪 Testing, formatting, linting

```bash
make
```

All format variations and cross-compatibility scenarios are thoroughly tested to ensure reliability across all Minecraft versions.
