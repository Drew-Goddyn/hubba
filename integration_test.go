package main

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Test Complete Application Integration
func TestApplicationIntegration(t *testing.T) {
	model := initialModel()

	// Test initialization
	if model.entityManager == nil {
		t.Error("Entity manager should be initialized")
	}
	if model.physicsEngine == nil {
		t.Error("Physics engine should be initialized")
	}
	if model.animationEngine == nil {
		t.Error("Animation engine should be initialized")
	}
	if model.controlPanel == nil {
		t.Error("Control panel should be initialized")
	}

	// Test initial state
	if model.paused {
		t.Error("Simulation should start unpaused")
	}
	if model.ready {
		t.Error("Model should not be ready until window size is set")
	}
	if model.entityManager.Count() != 0 {
		t.Error("Should start with no entities")
	}
}

// Test Complete Entity Lifecycle Integration
func TestEntityLifecycleIntegration(t *testing.T) {
	model := initialModel()

	// Set up model with window size
	model.termWidth = 80
	model.termHeight = 24
	model.updatePaneDimensions()
	model.ready = true

	// Update physics engine bounds
	renderGridHeight := model.simHeight - 8
	model.physicsEngine.UpdateBounds(float64(model.simWidth), float64(renderGridHeight))

	// Test adding entities through keyboard commands
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(Model)
	if model.entityManager.Count() != 1 {
		t.Errorf("Expected 1 entity after 'a' key, got %d", model.entityManager.Count())
	}

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	model = updatedModel.(Model)
	if model.entityManager.Count() != 2 {
		t.Errorf("Expected 2 entities after 's' key, got %d", model.entityManager.Count())
	}

	// Test clearing entities
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	model = updatedModel.(Model)
	if model.entityManager.Count() != 0 {
		t.Error("Expected 0 entities after 'c' key")
	}

	// Test reset functionality
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	model = updatedModel.(Model)
	if model.entityManager.Count() != 0 {
		t.Error("Expected 0 entities after reset")
	}
	if model.paused {
		t.Error("Expected unpaused state after reset")
	}
}

// Test Physics and Animation Integration
func TestPhysicsAnimationIntegration(t *testing.T) {
	model := initialModel()
	model.termWidth = 80
	model.termHeight = 24
	model.updatePaneDimensions()
	model.ready = true

	// Update physics engine bounds
	renderGridHeight := model.simHeight - 8
	model.physicsEngine.UpdateBounds(float64(model.simWidth), float64(renderGridHeight))

	// Create test entity
	sphere := NewSphere(10.0, 5.0, 1, lipgloss.Color("32"))
	sphere.SetVelocity(2.0, -1.0)
	model.entityManager.AddEntity(sphere)

	// Store initial positions
	physicsX, physicsY := sphere.GetPosition()
	displayX, displayY := sphere.GetDisplayPosition()

	// Simulate tick message
	updatedModel, _ := model.Update(tickMsg(time.Now()))
	model = updatedModel.(Model)

	// Check that physics position changed
	newPhysicsX, newPhysicsY := sphere.GetPosition()
	if newPhysicsX == physicsX && newPhysicsY == physicsY {
		t.Error("Physics position should have changed after tick")
	}

	// Check that display position is being animated
	newDisplayX, newDisplayY := sphere.GetDisplayPosition()
	if newDisplayX == displayX && newDisplayY == displayY {
		t.Error("Display position should have changed after animation update")
	}
}

// Test Pause/Resume Integration
func TestPauseResumeIntegration(t *testing.T) {
	model := initialModel()
	model.termWidth = 80
	model.termHeight = 24
	model.updatePaneDimensions()
	model.ready = true

	// Create test entity
	sphere := NewSphere(10.0, 5.0, 1, lipgloss.Color("32"))
	sphere.SetVelocity(2.0, -1.0)
	model.entityManager.AddEntity(sphere)

	// Pause simulation
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	model = updatedModel.(Model)
	if !model.paused {
		t.Error("Model should be paused after 'p' key")
	}
	if model.physicsEngine.IsRunning() {
		t.Error("Physics engine should be paused")
	}

	// Store position while paused
	pausedX, pausedY := sphere.GetPosition()

	// Simulate tick while paused
	updatedModel, _ = model.Update(tickMsg(time.Now()))
	model = updatedModel.(Model)

	// Position should not change due to physics while paused
	currentX, currentY := sphere.GetPosition()
	if currentX != pausedX || currentY != pausedY {
		t.Error("Physics position should not change while paused")
	}

	// Resume simulation
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	model = updatedModel.(Model)
	if model.paused {
		t.Error("Model should be unpaused after second 'p' key")
	}
	if !model.physicsEngine.IsRunning() {
		t.Error("Physics engine should be running after resume")
	}
}

// Test Parameter Changes Integration
func TestParameterChangesIntegration(t *testing.T) {
	model := initialModel()

	// Test gravity cycling
	initialGravity := model.selectedGravity
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	model = updatedModel.(Model)
	if model.selectedGravity == initialGravity {
		t.Error("Gravity should have changed after 'g' key")
	}
	if model.physicsEngine.GetGravity() != model.selectedGravity {
		t.Error("Physics engine gravity should match selected gravity")
	}

	// Test bounce cycling
	initialBounce := model.physicsEngine.GetRestitution()
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	model = updatedModel.(Model)
	if model.physicsEngine.GetRestitution() == initialBounce {
		t.Error("Bounce should have changed after 'b' key")
	}

	// Test size cycling
	initialSize := model.selectedEntitySize
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})
	model = updatedModel.(Model)
	if model.selectedEntitySize == initialSize {
		t.Error("Entity size should have changed after 'z' key")
	}

	// Test color cycling
	initialColor := model.selectedColorIndex
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	model = updatedModel.(Model)
	if model.selectedColorIndex == initialColor {
		t.Error("Color index should have changed after 'x' key")
	}
}

// Test Button Message Integration
func TestButtonMessageIntegration(t *testing.T) {
	model := initialModel()
	model.termWidth = 80
	model.termHeight = 24
	model.updatePaneDimensions()
	model.ready = true

	// Test AddSphere button
	buttonMsg := ButtonMsg{Action: AddSphereAction}
	updatedModel, _ := model.Update(buttonMsg)
	model = updatedModel.(Model)
	if model.entityManager.Count() != 1 {
		t.Error("Should have 1 entity after AddSphere button")
	}

	// Test AddSprite button
	buttonMsg = ButtonMsg{Action: AddSpriteAction}
	updatedModel, _ = model.Update(buttonMsg)
	model = updatedModel.(Model)
	if model.entityManager.Count() != 2 {
		t.Error("Should have 2 entities after AddSprite button")
	}

	// Test ClearAll button
	buttonMsg = ButtonMsg{Action: ClearAllAction}
	updatedModel, _ = model.Update(buttonMsg)
	model = updatedModel.(Model)
	if model.entityManager.Count() != 0 {
		t.Error("Should have 0 entities after ClearAll button")
	}

	// Test PauseResume button
	buttonMsg = ButtonMsg{Action: PauseResumeAction}
	updatedModel, _ = model.Update(buttonMsg)
	model = updatedModel.(Model)
	if !model.paused {
		t.Error("Should be paused after PauseResume button")
	}

	// Test Reset button
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}) // Add entity
	model = updatedModel.(Model)
	buttonMsg = ButtonMsg{Action: ResetAction}
	updatedModel, _ = model.Update(buttonMsg)
	model = updatedModel.(Model)
	if model.entityManager.Count() != 0 {
		t.Error("Should have 0 entities after Reset button")
	}
	if model.paused {
		t.Error("Should not be paused after Reset button")
	}
}

// Test Window Resize Integration
func TestWindowResizeIntegration(t *testing.T) {
	model := initialModel()

	// Test initial state
	if model.ready {
		t.Error("Model should not be ready initially")
	}

	// Simulate window resize
	windowMsg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updatedModel, _ := model.Update(windowMsg)
	model = updatedModel.(Model)

	// Check that dimensions are updated
	if model.termWidth != 100 || model.termHeight != 30 {
		t.Error("Terminal dimensions should be updated after window resize")
	}
	if !model.ready {
		t.Error("Model should be ready after window resize")
	}

	// Check that pane dimensions are calculated
	if model.simWidth <= 0 || model.simHeight <= 0 {
		t.Error("Simulation pane dimensions should be positive")
	}
	if model.ctrlWidth <= 0 || model.ctrlHeight <= 0 {
		t.Error("Control pane dimensions should be positive")
	}

	// Simulate very small window
	smallWindowMsg := tea.WindowSizeMsg{Width: 20, Height: 10}
	updatedModel, _ = model.Update(smallWindowMsg)
	model = updatedModel.(Model)

	// Should handle minimum dimensions gracefully
	if model.simHeight < 6 {
		t.Error("Simulation height should have minimum value")
	}
}

// Test Performance Mode Integration
func TestPerformanceModeIntegration(t *testing.T) {
	model := initialModel()
	model.termWidth = 80
	model.termHeight = 24
	model.updatePaneDimensions()
	model.ready = true

	// Test performance mode toggle
	if model.performanceMode {
		t.Error("Performance mode should be off initially")
	}

	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	model = updatedModel.(Model)
	if !model.performanceMode {
		t.Error("Performance mode should be on after 'f' key")
	}

	// Test entity limit increase in performance mode
	if model.maxEntityLimit < 1000 {
		t.Error("Entity limit should increase in performance mode")
	}

	// Test stress test
	initialCount := model.entityManager.Count()
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	model = updatedModel.(Model)
	if model.entityManager.Count() <= initialCount {
		t.Error("Stress test should add entities")
	}
}

// Test Entity Limit Integration
func TestEntityLimitIntegration(t *testing.T) {
	model := initialModel()
	model.termWidth = 80
	model.termHeight = 24
	model.updatePaneDimensions()
	model.ready = true

	// Set low entity limit for testing
	model.maxEntityLimit = 2

	// Add entities up to limit
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(Model)
	if model.entityManager.Count() != 2 {
		t.Errorf("Should have 2 entities at limit, got %d", model.entityManager.Count())
	}

	// Try to add beyond limit
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(Model)
	if model.entityManager.Count() > 2 {
		t.Error("Should not exceed entity limit")
	}

	// Test entity limit cycling
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	model = updatedModel.(Model)
	if model.maxEntityLimit == 2 {
		t.Error("Entity limit should have changed after 'l' key")
	}
}

// Test Complete Simulation Workflow
func TestCompleteSimulationWorkflow(t *testing.T) {
	model := initialModel()

	// Initialize properly
	windowMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updatedModel, _ := model.Update(windowMsg)
	model = updatedModel.(Model)

	// Add some entities
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	model = updatedModel.(Model)

	// Change parameters
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	model = updatedModel.(Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})
	model = updatedModel.(Model)

	// Run simulation for several ticks
	for i := 0; i < 10; i++ {
		updatedModel, _ = model.Update(tickMsg(time.Now()))
		model = updatedModel.(Model)
	}

	// Pause and resume
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	model = updatedModel.(Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	model = updatedModel.(Model)

	// Test rendering
	view := model.View()
	if view == "" {
		t.Error("View should not be empty")
	}

	// Reset everything
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	model = updatedModel.(Model)

	// Verify clean state
	if model.entityManager.Count() != 0 {
		t.Error("Should have no entities after reset")
	}
	if model.paused {
		t.Error("Should not be paused after reset")
	}
}

// Test Control Panel Integration
func TestControlPanelIntegration(t *testing.T) {
	model := initialModel()
	model.termWidth = 80
	model.termHeight = 24
	model.updatePaneDimensions()
	model.ready = true

	// Test tab navigation
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updatedModel.(Model)

	// Test enter activation
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(Model)
	if cmd != nil {
		// Execute the command
		msg := cmd()
		if buttonMsg, ok := msg.(ButtonMsg); ok {
			// Process the button message
			updatedModel, _ = model.Update(buttonMsg)
			model = updatedModel.(Model)
		}
	}

	// Test view rendering with control panel
	view := model.View()
	if view == "" {
		t.Error("View should not be empty")
	}

	// Should contain expected control elements
	if !contains(view, "PHYSICS CONTROLS") {
		t.Error("View should contain control panel title")
	}
}

// Helper functions are defined in controls_test.go - we'll use those
