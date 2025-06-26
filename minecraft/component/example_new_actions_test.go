package component_test

import (
	"fmt"
	"strings"
	"testing"

	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/common/minecraft/component/codec"
)

// Example demonstrating the new Minecraft 1.21.6 click event actions
func Example_newClickEventActions() {
	// Modern codec for 1.21.6+ features
	j := &codec.Json{
		NoDownsampleColor:            true,
		UseLegacyFieldNames:          false, // Use snake_case field names
		UseLegacyClickEventStructure: false, // Use specific field names
		StdJson:                      true,
	}

	// show_dialog action (1.21.6+)
	dialogComponent := &Text{
		Content: "Open Dialog",
		S: Style{
			ClickEvent: ShowDialog("my_custom_dialog"),
		},
	}

	// custom action (1.21.6+) - with payload
	customComponentWithPayload := &Text{
		Content: "Custom Event with Payload",
		S: Style{
			ClickEvent: CustomEvent("my_event", "some_data"),
		},
	}

	// custom action (1.21.6+) - without payload
	customComponentSimple := &Text{
		Content: "Simple Custom Event",
		S: Style{
			ClickEvent: CustomEvent("simple_event"),
		},
	}

	// Encode dialog action
	dialogJSON := new(strings.Builder)
	err := j.Marshal(dialogJSON, dialogComponent)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Encode custom action with payload
	customWithPayloadJSON := new(strings.Builder)
	err = j.Marshal(customWithPayloadJSON, customComponentWithPayload)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Encode simple custom action
	simpleCustomJSON := new(strings.Builder)
	err = j.Marshal(simpleCustomJSON, customComponentSimple)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Dialog Action (1.21.6+):\n%s\n\n", dialogJSON.String())
	fmt.Printf("Custom Action with Payload (1.21.6+):\n%s\n\n", customWithPayloadJSON.String())
	fmt.Printf("Simple Custom Action (1.21.6+):\n%s\n\n", simpleCustomJSON.String())

	// Output:
	// Dialog Action (1.21.6+):
	// {"click_event":{"action":"show_dialog","dialog":"my_custom_dialog"},"text":"Open Dialog"}
	//
	// Custom Action with Payload (1.21.6+):
	// {"click_event":{"action":"custom","id":"my_event","payload":"some_data"},"text":"Custom Event with Payload"}
	//
	// Simple Custom Action (1.21.6+):
	// {"click_event":{"action":"custom","id":"simple_event"},"text":"Simple Custom Event"}
}

func TestNewClickEventActions_1216(t *testing.T) {
	// Test the new 1.21.6 actions in both formats
	j1216 := &codec.Json{
		UseLegacyFieldNames:          false, // Modern format
		UseLegacyClickEventStructure: false,
		StdJson:                      true,
	}

	jLegacy := &codec.Json{
		UseLegacyFieldNames:          true, // Legacy format
		UseLegacyClickEventStructure: true,
		StdJson:                      true,
	}

	// Test show_dialog action
	t.Run("show_dialog", func(t *testing.T) {
		component := &Text{
			Content: "Show Dialog",
			S:       Style{ClickEvent: ShowDialog("test_dialog")},
		}

		// Test modern format
		modern := new(strings.Builder)
		err := j1216.Marshal(modern, component)
		if err != nil {
			t.Fatalf("Failed to marshal show_dialog in modern format: %v", err)
		}

		if !strings.Contains(modern.String(), `"dialog":"test_dialog"`) {
			t.Errorf("Modern format should contain dialog field, got: %s", modern.String())
		}

		// Test legacy format (should fallback to value field)
		legacy := new(strings.Builder)
		err = jLegacy.Marshal(legacy, component)
		if err != nil {
			t.Fatalf("Failed to marshal show_dialog in legacy format: %v", err)
		}

		if !strings.Contains(legacy.String(), `"value":"test_dialog"`) {
			t.Errorf("Legacy format should contain value field, got: %s", legacy.String())
		}

		// Test round-trip compatibility
		decoded, err := j1216.Unmarshal([]byte(modern.String()))
		if err != nil {
			t.Fatalf("Failed to unmarshal modern format: %v", err)
		}

		if decoded.(*Text).S.ClickEvent.Value() != "test_dialog" {
			t.Errorf("Round-trip failed, expected 'test_dialog', got '%s'", decoded.(*Text).S.ClickEvent.Value())
		}
	})

	// Test custom action
	t.Run("custom", func(t *testing.T) {
		component := &Text{
			Content: "Custom Event",
			S:       Style{ClickEvent: CustomEvent("my_id", "my_payload")},
		}

		// Test modern format
		modern := new(strings.Builder)
		err := j1216.Marshal(modern, component)
		if err != nil {
			t.Fatalf("Failed to marshal custom in modern format: %v", err)
		}

		if !strings.Contains(modern.String(), `"id":"my_id"`) || !strings.Contains(modern.String(), `"payload":"my_payload"`) {
			t.Errorf("Modern format should contain id and payload fields, got: %s", modern.String())
		}

		// Test round-trip compatibility
		decoded, err := j1216.Unmarshal([]byte(modern.String()))
		if err != nil {
			t.Fatalf("Failed to unmarshal modern format: %v", err)
		}

		if decoded.(*Text).S.ClickEvent.Value() != "my_id|my_payload" {
			t.Errorf("Round-trip failed, expected 'my_id|my_payload', got '%s'", decoded.(*Text).S.ClickEvent.Value())
		}
	})
}
