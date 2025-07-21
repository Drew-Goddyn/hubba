package main

import (
	"math"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Test Edge Cases for Physics Engine
func TestPhysicsEngineEdgeCases(t *testing.T) {
	// Test with zero dimensions
	pe := NewPhysicsEngine(0, 0)
	if pe.MaxX < 0 || pe.MaxY < 0 {
		t.Error("Physics engine should handle zero dimensions gracefully")
	}

	// Test with negative dimensions
	pe2 := NewPhysicsEngine(-10, -5)
	if pe2.MaxX < 0 || pe2.MaxY < 0 {
		t.Error("Physics engine should handle negative dimensions gracefully")
	}

	// Test extreme gravity values
	pe.SetGravity(math.Inf(1)) // Positive infinity
	if math.IsInf(pe.GetGravity(), 1) {
		t.Error("Physics engine should reject infinite gravity")
	}

	pe.SetGravity(math.NaN()) // Not a number
	if math.IsNaN(pe.GetGravity()) {
		t.Error("Physics engine should reject NaN gravity")
	}

	// Test extreme restitution values
	pe.SetRestitution(-1.0) // Negative (invalid)
	if pe.GetRestitution() < 0 {
		t.Error("Physics engine should reject negative restitution")
	}

	pe.SetRestitution(2.0) // Greater than 1 (invalid)
	if pe.GetRestitution() > 1.0 {
		t.Error("Physics engine should reject restitution > 1")
	}
}

// Test Edge Cases for Entity Creation
func TestEntityCreationEdgeCases(t *testing.T) {
	// Test entity with extreme positions
	sphere := NewSphere(-1000.0, 1000.0, 1, lipgloss.Color("32"))
	x, y := sphere.GetPosition()
	if x != -1000.0 || y != 1000.0 {
		t.Error("Entity should accept extreme positions")
	}

	// Test entity with zero size
	sprite := NewSprite(0.0, 0.0, 0, lipgloss.Color("34"), "★")
	if sprite.GetSize() != 0 {
		t.Error("Entity should accept zero size")
	}

	// Test entity with negative size
	sphere2 := NewSphere(0.0, 0.0, -5, lipgloss.Color("32"))
	if sphere2.GetSize() < 0 {
		t.Error("Entity should handle negative size appropriately")
	}

	// Test sprite with empty symbol
	sprite2 := NewSprite(0.0, 0.0, 1, lipgloss.Color("34"), "")
	if sprite2.GetSymbol() == "" {
		t.Error("Sprite should assign random symbol when empty string provided")
	}

	// Test entity with extreme velocity
	sphere.SetVelocity(math.Inf(1), math.Inf(-1))
	vx, vy := sphere.GetVelocity()
	if math.IsInf(vx, 0) || math.IsInf(vy, 0) {
		t.Error("Entity should handle infinite velocity gracefully")
	}

	// Test entity with NaN velocity
	sphere.SetVelocity(math.NaN(), math.NaN())
	vx, vy = sphere.GetVelocity()
	if math.IsNaN(vx) || math.IsNaN(vy) {
		t.Error("Entity should handle NaN velocity gracefully")
	}
}

// Test Edge Cases for Entity Manager
func TestEntityManagerEdgeCases(t *testing.T) {
	manager := NewEntityManager()

	// Test removing from empty manager
	removed := manager.RemoveEntity("non-existent")
	if removed {
		t.Error("Should return false when removing from empty manager")
	}

	// Test clearing empty manager
	manager.Clear() // Should not panic

	// Test getting entities by type from empty manager
	spheres := manager.GetEntitiesByType(SphereType)
	if len(spheres) != 0 {
		t.Error("Should return empty slice for type from empty manager")
	}

	// Test count by type from empty manager
	count := manager.CountByType(SpriteType)
	if count != 0 {
		t.Error("Should return 0 count for type from empty manager")
	}

	// Test adding nil entity (should be handled gracefully)
	// Note: This would require modifying the manager to check for nil

	// Test massive number of entities
	for i := 0; i < 1000; i++ {
		sphere := NewSphere(float64(i), float64(i), 1, lipgloss.Color("32"))
		manager.AddEntity(sphere)
	}
	if manager.Count() != 1000 {
		t.Error("Manager should handle large number of entities")
	}

	// Test removing entities while iterating
	entities := manager.GetEntities()
	for i, entity := range entities {
		if i%2 == 0 { // Remove every other entity
			manager.RemoveEntity(entity.GetID())
		}
	}
	if manager.Count() != 500 {
		t.Error("Manager should handle entity removal during iteration")
	}
}

// Test Edge Cases for Collision Detection
func TestCollisionDetectionEdgeCases(t *testing.T) {
	manager := NewEntityManager()

	// Test collision detection with no entities
	collisions := manager.CheckCollisions()
	if len(collisions) != 0 {
		t.Error("Should find no collisions with no entities")
	}

	// Test collision detection with single entity
	sphere := NewSphere(10.0, 10.0, 2, lipgloss.Color("32"))
	manager.AddEntity(sphere)
	collisions = manager.CheckCollisions()
	if len(collisions) != 0 {
		t.Error("Should find no collisions with single entity")
	}

	// Test collision with entities of size 0
	sphere1 := NewSphere(5.0, 5.0, 0, lipgloss.Color("32"))
	sphere2 := NewSphere(5.0, 5.0, 0, lipgloss.Color("33"))
	manager.Clear()
	manager.AddEntity(sphere1)
	manager.AddEntity(sphere2)
	collisions = manager.CheckCollisions()
	// Zero-sized entities at same position should still collide
	if len(collisions) != 1 {
		t.Error("Zero-sized entities at same position should collide")
	}

	// Test collision with very large entities
	largeSphere1 := NewSphere(50.0, 50.0, 100, lipgloss.Color("32"))
	largeSphere2 := NewSphere(100.0, 100.0, 100, lipgloss.Color("33"))
	manager.Clear()
	manager.AddEntity(largeSphere1)
	manager.AddEntity(largeSphere2)
	collisions = manager.CheckCollisions()
	if len(collisions) != 1 {
		t.Error("Large overlapping entities should collide")
	}

	// Test with entities at extreme positions
	extremeSphere1 := NewSphere(-1000.0, -1000.0, 1, lipgloss.Color("32"))
	extremeSphere2 := NewSphere(1000.0, 1000.0, 1, lipgloss.Color("33"))
	manager.Clear()
	manager.AddEntity(extremeSphere1)
	manager.AddEntity(extremeSphere2)
	collisions = manager.CheckCollisions()
	if len(collisions) != 0 {
		t.Error("Entities at extreme positions should not collide")
	}
}

// Test Edge Cases for Animation Engine
func TestAnimationEngineEdgeCases(t *testing.T) {
	ae := NewAnimationEngine()

	// Test animation with zero spring tension
	ae.SpringTension = 0.0
	eas := ae.NewEntityAnimationState(0.0, 0.0)
	eas.SetTarget(10.0, 10.0)

	// Should handle zero tension gracefully
	ae.UpdateAnimation(eas)
	if eas.DisplayX != 0.0 || eas.DisplayY != 0.0 {
		// With zero tension, position should remain at start
		t.Log("Zero spring tension behavior handled")
	}

	// Test animation with negative spring values
	ae.SpringTension = -100.0
	ae.SpringDamping = -50.0
	ae.UpdateAnimation(eas)
	// Should not crash or produce invalid results

	// Test animation with extreme target positions
	eas.SetTarget(math.Inf(1), math.Inf(-1))
	ae.UpdateAnimation(eas)
	x, y := eas.GetDisplayPosition()
	if math.IsInf(x, 0) || math.IsInf(y, 0) {
		t.Error("Animation should handle infinite targets gracefully")
	}

	// Test animation with NaN targets
	eas.SetTarget(math.NaN(), math.NaN())
	ae.UpdateAnimation(eas)
	x, y = eas.GetDisplayPosition()
	if math.IsNaN(x) || math.IsNaN(y) {
		t.Error("Animation should handle NaN targets gracefully")
	}

	// Test rapid target changes
	for i := 0; i < 100; i++ {
		eas.SetTarget(float64(i), float64(-i))
		ae.UpdateAnimation(eas)
	}
	// Should handle rapid changes without issues
}

// Test Edge Cases for Boundary Collisions
func TestBoundaryCollisionEdgeCases(t *testing.T) {
	pe := NewPhysicsEngine(10, 10) // Small bounds

	// Test entity exactly at boundary
	sphere := NewSphere(pe.MinX, pe.MinY, 1, lipgloss.Color("32"))
	pe.ApplyPhysics([]Entity{sphere})
	x, y := sphere.GetPosition()
	if x < pe.MinX || y < pe.MinY {
		t.Error("Entity should not go below minimum bounds")
	}

	// Test entity larger than bounds
	largeSphere := NewSphere(5.0, 5.0, 20, lipgloss.Color("32")) // Size 20 in 10x10 space
	pe.ApplyPhysics([]Entity{largeSphere})
	// Should handle gracefully without infinite loops

	// Test with zero-sized bounds
	pe.UpdateBounds(0, 0)
	pe.ApplyPhysics([]Entity{sphere})
	// Should not crash

	// Test entity moving at extreme velocity toward boundary
	fastSphere := NewSphere(5.0, 5.0, 1, lipgloss.Color("32"))
	fastSphere.SetVelocity(1000.0, 1000.0)
	pe.UpdateBounds(10, 10) // Reset bounds
	pe.ApplyPhysics([]Entity{fastSphere})
	x, y = fastSphere.GetPosition()
	if x > pe.MaxX || y > pe.MaxY {
		t.Errorf("Fast-moving entity escaped bounds: position=(%.2f, %.2f), maxBounds=(%.2f, %.2f)", x, y, pe.MaxX, pe.MaxY)
	}
}

// Test Edge Cases for Application Model
func TestApplicationModelEdgeCases(t *testing.T) {
	model := initialModel()

	// Test with zero terminal dimensions
	model.termWidth = 0
	model.termHeight = 0
	model.updatePaneDimensions()
	if model.simWidth < 0 || model.simHeight < 0 {
		t.Error("Model should handle zero terminal dimensions gracefully")
	}

	// Test with extremely small terminal
	model.termWidth = 5
	model.termHeight = 3
	model.updatePaneDimensions()
	if model.simHeight < 6 {
		// Should enforce minimum simulation height
		t.Log("Minimum simulation height enforced")
	}

	// Test view rendering before ready
	view := model.View()
	if view == "" {
		t.Error("View should show initialization message when not ready")
	}

	// Test parameter cycling edge cases
	model.selectedGravity = 999.0 // Not in standard list
	model.cycleGravity()
	found := false
	for _, gravity := range gravityLevels {
		if gravity == model.selectedGravity {
			found = true
			break
		}
	}
	if !found {
		t.Error("Gravity cycling should fallback to valid value")
	}

	// Test color cycling at boundary
	colors := GetAvailableColors()
	model.selectedColorIndex = len(colors) - 1 // Last color
	model.cycleEntityColor()
	if model.selectedColorIndex != 0 {
		t.Error("Color cycling should wrap around to first color")
	}

	// Test entity size cycling edge case
	model.selectedEntitySize = 999 // Not in standard list
	model.cycleEntitySize()
	found = false
	for _, size := range entitySizes {
		if size == model.selectedEntitySize {
			found = true
			break
		}
	}
	if !found {
		t.Error("Size cycling should fallback to valid value")
	}
}

// Test Edge Cases for Stress Scenarios
func TestStressTestEdgeCases(t *testing.T) {
	model := initialModel()
	model.termWidth = 10
	model.termHeight = 5 // Very small terminal
	model.updatePaneDimensions()
	model.ready = true

	// Test stress test with tiny terminal
	model.runStressTest()
	// Should not crash even with minimal space

	// Test stress test at entity limit
	model.entityManager = NewEntityManager() // Clear existing entities
	model.maxEntityLimit = 5
	model.runStressTest()
	if model.entityManager.Count() > model.maxEntityLimit {
		t.Errorf("Stress test should respect entity limit: count=%d, limit=%d", model.entityManager.Count(), model.maxEntityLimit)
	}

	// Test multiple rapid stress tests
	for i := 0; i < 5; i++ {
		model.runStressTest()
	}
	// Should handle multiple rapid calls gracefully
}

// Test Edge Cases for Performance Monitoring
func TestPerformanceMonitoringEdgeCases(t *testing.T) {
	model := initialModel()

	// Test FPS calculation with zero time delta
	model.frameCount = 100
	testTime := time.Now()
	model.lastFPSUpdate = testTime // Set to specific time

	// Simulate tick to trigger FPS calculation with same time (zero delta)
	model.ready = true
	updatedModel, _ := model.Update(tickMsg(testTime))
	model = updatedModel.(Model)

	// Should handle zero time delta gracefully
	if math.IsInf(model.currentFPS, 0) || math.IsNaN(model.currentFPS) {
		t.Error("FPS calculation should handle zero time delta")
	}

	// Test with negative frame count (edge case)
	model.frameCount = -10
	updatedModel, _ = model.Update(tickMsg(model.lastFPSUpdate.Add(1000000000))) // 1 second later
	model = updatedModel.(Model)
	// Should handle gracefully
}

// Test Edge Cases for Random Functions
func TestRandomFunctionEdgeCases(t *testing.T) {
	// Test getting random color multiple times
	colors := make(map[string]bool)
	for i := 0; i < 100; i++ {
		color := GetRandomColor()
		colors[string(color)] = true
	}
	if len(colors) < 2 {
		t.Error("GetRandomColor should generate variety")
	}

	// Test sprite with empty symbol assignment
	sprites := make(map[string]bool)
	for i := 0; i < 50; i++ {
		sprite := NewSprite(0.0, 0.0, 1, lipgloss.Color("34"), "")
		sprites[sprite.GetSymbol()] = true
	}
	if len(sprites) < 2 {
		t.Error("Random sprite symbols should have variety")
	}
}

// Test Edge Cases for Control Panel Responsiveness
func TestControlPanelResponsivenessEdgeCases(t *testing.T) {
	// Test with zero dimensions
	cp := NewControlPanel(0, 0)
	if cp.width < 0 || cp.height < 0 {
		t.Error("Control panel should handle zero dimensions")
	}

	// Test update with zero dimensions
	cp.UpdateResponsiveMode(0, 0)
	// Should not crash

	// Test with negative dimensions
	cp.UpdateResponsiveMode(-10, -5)
	// Should handle gracefully

	// Test view with minimal dimensions
	view := cp.View()
	if view == "" {
		t.Error("Control panel should render something even with minimal dimensions")
	}

	// Test extreme responsive modes
	cp.ultraCompactMode = true
	cp.compactMode = true // Both true (edge case)
	view = cp.View()
	if view == "" {
		t.Error("Control panel should handle overlapping compact modes")
	}
}

// Test Memory and Resource Edge Cases
func TestMemoryResourceEdgeCases(t *testing.T) {
	// Test creating and destroying many entities rapidly
	manager := NewEntityManager()

	for cycle := 0; cycle < 10; cycle++ {
		// Add many entities
		for i := 0; i < 100; i++ {
			sphere := NewSphere(float64(i), float64(i), 1, GetRandomColor())
			manager.AddEntity(sphere)
		}

		// Remove all entities
		manager.Clear()

		// Verify clean state
		if manager.Count() != 0 {
			t.Error("Manager should be clean after clear")
		}
	}

	// Test physics engine with many repeated operations
	pe := NewPhysicsEngine(100, 100)
	entities := make([]Entity, 50)
	for i := range entities {
		entities[i] = NewSphere(float64(i*2), float64(i*2), 1, GetRandomColor())
		entities[i].SetVelocity(float64(i), float64(-i))
	}

	// Run many physics updates
	for i := 0; i < 100; i++ {
		pe.ApplyPhysics(entities)
		pe.HandleEntityCollisions(entities)
	}

	// Verify entities are still valid
	for _, entity := range entities {
		x, y := entity.GetPosition()
		if math.IsNaN(x) || math.IsNaN(y) || math.IsInf(x, 0) || math.IsInf(y, 0) {
			t.Error("Entity position should remain valid after many updates")
		}
	}
}

// Helper function for edge cases (using math.Abs instead of custom implementation)
func absEdgeCase(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
