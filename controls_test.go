package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewControlPanel(t *testing.T) {
	cp := NewControlPanel(80, 20)

	// Check basic initialization
	if cp == nil {
		t.Fatal("NewControlPanel returned nil")
	}

	if cp.width != 80 {
		t.Errorf("Expected width 80, got %d", cp.width)
	}

	if cp.height != 20 {
		t.Errorf("Expected height 20, got %d", cp.height)
	}

	if len(cp.buttons) != 5 {
		t.Errorf("Expected 5 buttons, got %d", len(cp.buttons))
	}

	// Check default focused button
	if cp.focused != 0 {
		t.Errorf("Expected focused button 0, got %d", cp.focused)
	}
}

func TestControlPanelNavigation(t *testing.T) {
	cp := NewControlPanel(80, 20)

	// Test tab navigation
	cp.Update(tea.KeyMsg{Type: tea.KeyTab})
	if cp.focused != 1 {
		t.Errorf("Expected focused button 1 after tab, got %d", cp.focused)
	}

	// Test multiple tabs wrapping around
	for i := 0; i < 6; i++ {
		cp.Update(tea.KeyMsg{Type: tea.KeyTab})
	}
	if cp.focused != 2 { // 7 total tabs: (0 + 7) % 5 = 2
		t.Errorf("Expected focused button 2 after 7 total tabs, got %d", cp.focused)
	}

	// Test shift+tab (reverse navigation) from position 0
	cp.focused = 0 // Reset to position 0
	cp.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if cp.focused != 4 { // Should wrap to last button (4)
		t.Errorf("Expected focused button 4 after shift+tab from 0, got %d", cp.focused)
	}
}

func TestControlPanelButtonActivation(t *testing.T) {
	cp := NewControlPanel(80, 20)

	// Test enter key activation
	cp.focused = 0 // Focus on first button (Add Sphere)
	_, cmd := cp.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Error("Expected command after enter key, got nil")
	}

	// Execute command to get the message
	msg := cmd()
	if buttonMsg, ok := msg.(ButtonMsg); ok {
		if buttonMsg.Action != AddSphereAction {
			t.Errorf("Expected AddSphereAction, got %v", buttonMsg.Action)
		}
	} else {
		t.Error("Expected ButtonMsg, got different type")
	}

	// Test space key activation
	cp.focused = 1 // Focus on second button (Add Sprite)
	_, cmd = cp.Update(tea.KeyMsg{Type: tea.KeySpace})

	if cmd == nil {
		t.Error("Expected command after space key, got nil")
	}

	msg = cmd()
	if buttonMsg, ok := msg.(ButtonMsg); ok {
		if buttonMsg.Action != AddSpriteAction {
			t.Errorf("Expected AddSpriteAction, got %v", buttonMsg.Action)
		}
	} else {
		t.Error("Expected ButtonMsg, got different type")
	}
}

func TestUpdatePauseButton(t *testing.T) {
	cp := NewControlPanel(80, 20)

	// Find pause button
	var pauseButtonIndex int = -1
	for i, button := range cp.buttons {
		if button.Action == PauseResumeAction {
			pauseButtonIndex = i
			break
		}
	}

	if pauseButtonIndex == -1 {
		t.Fatal("Pause button not found")
	}

	// Test updating to paused state
	cp.UpdatePauseButton(true)
	if cp.buttons[pauseButtonIndex].Label != "Resume" {
		t.Errorf("Expected 'Resume' label when paused, got '%s'", cp.buttons[pauseButtonIndex].Label)
	}

	// Test updating to running state
	cp.UpdatePauseButton(false)
	if cp.buttons[pauseButtonIndex].Label != "Pause" {
		t.Errorf("Expected 'Pause' label when running, got '%s'", cp.buttons[pauseButtonIndex].Label)
	}
}

func TestSetButtonActive(t *testing.T) {
	cp := NewControlPanel(80, 20)

	// Test setting button active
	cp.SetButtonActive(AddSphereAction, true)

	var found bool
	for _, button := range cp.buttons {
		if button.Action == AddSphereAction && button.Active {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected AddSphere button to be active")
	}

	// Test setting button inactive
	cp.SetButtonActive(AddSphereAction, false)

	found = false
	for _, button := range cp.buttons {
		if button.Action == AddSphereAction && button.Active {
			found = true
			break
		}
	}

	if found {
		t.Error("Expected AddSphere button to be inactive")
	}
}

func TestControlPanelView(t *testing.T) {
	cp := NewControlPanel(80, 20)

	view := cp.View()

	if view == "" {
		t.Error("View returned empty string")
	}

	// Check that view contains expected elements
	expectedStrings := []string{
		"🎮 PHYSICS CONTROLS",
		"Add Sphere",
		"Add Sprite",
		"Clear All",
		"Pause",
		"Reset",
		"Gravity",
		"Keys:",
	}

	for _, expected := range expectedStrings {
		if !contains(view, expected) {
			t.Errorf("View missing expected string: %s", expected)
		}
	}
}

func TestButtonActions(t *testing.T) {
	actions := []ButtonAction{
		AddSphereAction,
		AddSpriteAction,
		ClearAllAction,
		PauseResumeAction,
		ResetAction,
		GravityAction,
		BounceAction,
	}

	for _, action := range actions {
		if string(action) == "" {
			t.Errorf("Button action %v has empty string value", action)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(substr) == 0 ||
			s[:len(substr)] == substr ||
			(len(s) > len(substr) && (s[len(s)-len(substr):] == substr ||
				indexOf(s, substr) >= 0)))
}

// Simple indexOf implementation
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
