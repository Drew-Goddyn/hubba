package main

import (
	"fmt"
	"strings"
	"testing"
)

// TestUserExperienceSpecification ensures the app behaves as users expect
func TestUserExperienceSpecification(t *testing.T) {
	testCases := []struct {
		name           string
		width          int
		height         int
		shouldHaveColors bool
		shouldHaveBorders bool
		description    string
	}{
		{
			name:           "very_small_terminal",
			width:          25,
			height:         10,
			shouldHaveColors: false,
			shouldHaveBorders: false,
			description:    "Ultra-compact terminals (≤30 chars) should use minimal mode",
		},
		{
			name:           "small_terminal", 
			width:          50,
			height:         15,
			shouldHaveColors: true,
			shouldHaveBorders: true,
			description:    "Small terminals (50 chars) should have full colorful UI",
		},
		{
			name:           "medium_terminal",
			width:          80,
			height:         24,
			shouldHaveColors: true,
			shouldHaveBorders: true,
			description:    "Standard terminals (80 chars) must have full colorful UI",
		},
		{
			name:           "large_terminal",
			width:          120,
			height:         30,
			shouldHaveColors: true,
			shouldHaveBorders: true,
			description:    "Large terminals (120 chars) must have full colorful UI",
		},
		{
			name:           "ultra_wide_terminal",
			width:          200,
			height:         50,
			shouldHaveColors: true,
			shouldHaveBorders: true,
			description:    "Ultra-wide terminals must have full colorful UI",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create model and simulate window size
			model := initialModel()
			model.termWidth = tc.width
			model.termHeight = tc.height
			model.updatePaneDimensions()
			model.ready = true

			// Update physics engine bounds
			renderGridHeight := model.simHeight - 8
			model.physicsEngine.UpdateBounds(float64(model.simWidth), float64(renderGridHeight))

			// Generate view
			view := model.View()

			// Test color requirement
			hasColors := strings.Contains(view, "\x1b[") // ANSI escape sequences
			if tc.shouldHaveColors && !hasColors {
				t.Errorf("%s: Expected colorful UI but got minimal mode. %s", tc.name, tc.description)
			}
			if !tc.shouldHaveColors && hasColors {
				t.Errorf("%s: Expected minimal mode but got colorful UI. %s", tc.name, tc.description)
			}

			// Test border requirement
			hasBorders := strings.Contains(view, "─") || strings.Contains(view, "│") || strings.Contains(view, "┌")
			if tc.shouldHaveBorders && !hasBorders {
				t.Errorf("%s: Expected bordered UI but got minimal mode. %s", tc.name, tc.description)
			}
			if !tc.shouldHaveBorders && hasBorders {
				t.Errorf("%s: Expected minimal mode but got bordered UI. %s", tc.name, tc.description)
			}

			// Ensure the app doesn't crash or return empty content
			if view == "" {
				t.Errorf("%s: App returned empty view", tc.name)
			}

			// For standard terminals and above, ensure rich content
			if tc.width >= 80 {
				if !strings.Contains(view, "PHYSICS") {
					t.Errorf("%s: Expected rich physics simulation UI for standard terminal", tc.name)
				}
				if !strings.Contains(view, "Add Sphere") {
					t.Errorf("%s: Expected interactive controls for standard terminal", tc.name)
				}
			}
		})
	}
}

// TestMinimalModeThreshold ensures the threshold is set correctly
func TestMinimalModeThreshold(t *testing.T) {
	// Test that 30 chars triggers minimal mode
	model := initialModel()
	model.termWidth = 30
	model.termHeight = 10
	model.updatePaneDimensions()
	model.ready = true

	view := model.View()
	hasColors := strings.Contains(view, "\x1b[")
	if hasColors {
		t.Error("30-character terminal should use minimal mode (no colors)")
	}

	// Test that 31 chars uses colorful mode
	model.termWidth = 31
	model.updatePaneDimensions()
	view = model.View()
	hasColors = strings.Contains(view, "\x1b[")
	if !hasColors {
		t.Error("31-character terminal should use colorful mode")
	}
}

// TestNoBlackAndWhiteRegression ensures the original bug doesn't return
func TestNoBlackAndWhiteRegression(t *testing.T) {
	// This test ensures we don't regress to the old 50-character threshold
	commonTerminalWidths := []int{60, 70, 80, 90, 100, 120}
	
	for _, width := range commonTerminalWidths {
		t.Run(fmt.Sprintf("width_%d", width), func(t *testing.T) {
			model := initialModel()
			model.termWidth = width
			model.termHeight = 24
			model.updatePaneDimensions()
			model.ready = true

			view := model.View()
			hasColors := strings.Contains(view, "\x1b[")
			
			if !hasColors {
				t.Errorf("Terminal width %d should have colorful UI, not black and white", width)
			}
		})
	}
} 