package main

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Benchmark Physics Engine with Various Entity Counts
func BenchmarkPhysicsEngine10Entities(b *testing.B) {
	benchmarkPhysicsEngine(b, 10)
}

func BenchmarkPhysicsEngine50Entities(b *testing.B) {
	benchmarkPhysicsEngine(b, 50)
}

func BenchmarkPhysicsEngine100Entities(b *testing.B) {
	benchmarkPhysicsEngine(b, 100)
}

func BenchmarkPhysicsEngine500Entities(b *testing.B) {
	benchmarkPhysicsEngine(b, 500)
}

func benchmarkPhysicsEngine(b *testing.B, entityCount int) {
	pe := NewPhysicsEngine(100, 100)
	entities := make([]Entity, entityCount)

	// Create entities with random positions and velocities
	for i := range entities {
		x := float64(i%90) + 5.0
		y := float64((i/10)%90) + 5.0
		entities[i] = NewSphere(x, y, 1, GetRandomColor())
		entities[i].SetVelocity(float64(i%10)-5.0, float64((i/2)%10)-5.0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pe.ApplyPhysics(entities)
	}
}

// Benchmark Collision Detection with Various Entity Counts
func BenchmarkCollisionDetection10Entities(b *testing.B) {
	benchmarkCollisionDetection(b, 10)
}

func BenchmarkCollisionDetection50Entities(b *testing.B) {
	benchmarkCollisionDetection(b, 50)
}

func BenchmarkCollisionDetection100Entities(b *testing.B) {
	benchmarkCollisionDetection(b, 100)
}

func benchmarkCollisionDetection(b *testing.B, entityCount int) {
	pe := NewPhysicsEngine(100, 100)
	entities := make([]Entity, entityCount)

	// Create entities in a clustered pattern to ensure collisions
	for i := range entities {
		x := 50.0 + float64(i%10) // Cluster entities
		y := 50.0 + float64((i/10)%10)
		entities[i] = NewSphere(x, y, 2, GetRandomColor())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pe.HandleEntityCollisions(entities)
	}
}

// Benchmark Entity Manager Operations
func BenchmarkEntityManagerAdd(b *testing.B) {
	manager := NewEntityManager()
	entities := make([]Entity, b.N)

	for i := range entities {
		entities[i] = NewSphere(float64(i), float64(i), 1, GetRandomColor())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.AddEntity(entities[i])
	}
}

func BenchmarkEntityManagerRemove(b *testing.B) {
	manager := NewEntityManager()
	entities := make([]Entity, b.N)

	// Pre-populate manager
	for i := range entities {
		entities[i] = NewSphere(float64(i), float64(i), 1, GetRandomColor())
		manager.AddEntity(entities[i])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.RemoveEntity(entities[i].GetID())
	}
}

// Benchmark Animation Engine
func BenchmarkAnimationEngine(b *testing.B) {
	ae := NewAnimationEngine()
	states := make([]*EntityAnimationState, 100)

	for i := range states {
		states[i] = ae.NewEntityAnimationState(float64(i), float64(i))
		states[i].SetTarget(float64(i+50), float64(i+50))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, state := range states {
			ae.UpdateAnimation(state)
		}
	}
}

// Test Performance Under Load
func TestPhysicsPerformanceUnderLoad(t *testing.T) {
	pe := NewPhysicsEngine(200, 200)
	manager := NewEntityManager()

	// Add 1000 entities
	for i := 0; i < 1000; i++ {
		x := float64(i%190) + 5.0
		y := float64((i/10)%190) + 5.0
		entity := NewSphere(x, y, 1, GetRandomColor())
		entity.SetVelocity(float64(i%20)-10.0, float64((i/3)%20)-10.0)
		manager.AddEntity(entity)
	}

	entities := manager.GetEntities()

	// Measure time for 100 physics updates
	start := time.Now()
	for i := 0; i < 100; i++ {
		pe.ApplyPhysics(entities)
		pe.HandleEntityCollisions(entities)
	}
	duration := time.Since(start)

	// Should complete within reasonable time (adjust threshold as needed)
	if duration > time.Second*5 {
		t.Errorf("Physics updates took too long: %v", duration)
	}

	t.Logf("1000 entities, 100 updates: %v", duration)
}

// Test Memory Usage with Large Entity Count
func TestMemoryUsageWithLargeEntityCount(t *testing.T) {
	manager := NewEntityManager()

	// Add and remove entities in batches to test memory management
	for batch := 0; batch < 10; batch++ {
		// Add 500 entities
		for i := 0; i < 500; i++ {
			entity := NewSphere(float64(i), float64(i), 1, GetRandomColor())
			manager.AddEntity(entity)
		}

		// Remove half
		entities := manager.GetEntities()
		for i := 0; i < len(entities)/2; i++ {
			manager.RemoveEntity(entities[i].GetID())
		}
	}

	// Final count should be reasonable
	if manager.Count() > 2500 {
		t.Errorf("Unexpected entity count after batched operations: %d", manager.Count())
	}

	t.Logf("Final entity count after memory test: %d", manager.Count())
}

// Test Animation Performance
func TestAnimationPerformance(t *testing.T) {
	ae := NewAnimationEngine()
	states := make([]*EntityAnimationState, 500)

	// Create animation states
	for i := range states {
		states[i] = ae.NewEntityAnimationState(0.0, 0.0)
		states[i].SetTarget(100.0, 100.0)
	}

	// Measure time for 1000 animation updates
	start := time.Now()
	for i := 0; i < 1000; i++ {
		for _, state := range states {
			ae.UpdateAnimation(state)
		}
	}
	duration := time.Since(start)

	// Should complete within reasonable time
	if duration > time.Second*2 {
		t.Errorf("Animation updates took too long: %v", duration)
	}

	t.Logf("500 animations, 1000 updates: %v", duration)
}

// Test Collision Detection Performance
func TestCollisionDetectionPerformance(t *testing.T) {
	manager := NewEntityManager()
	pe := NewPhysicsEngine(100, 100)

	// Create entities in overlapping positions to maximize collisions
	for i := 0; i < 100; i++ {
		for j := 0; j < 5; j++ {
			x := 50.0 + float64(i%10)
			y := 50.0 + float64(j)
			entity := NewSphere(x, y, 2, GetRandomColor())
			manager.AddEntity(entity)
		}
	}

	entities := manager.GetEntities()

	// Measure collision detection time
	start := time.Now()
	for i := 0; i < 100; i++ {
		pe.HandleEntityCollisions(entities)
	}
	duration := time.Since(start)

	// Should complete within reasonable time
	if duration > time.Second*3 {
		t.Errorf("Collision detection took too long: %v", duration)
	}

	t.Logf("500 entities, 100 collision checks: %v", duration)
}

// Test Model Update Performance
func TestModelUpdatePerformance(t *testing.T) {
	model := initialModel()
	model.termWidth = 100
	model.termHeight = 50
	model.updatePaneDimensions()
	model.ready = true

	// Add many entities
	for i := 0; i < 200; i++ {
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
		model = updatedModel.(Model)
	}

	// Measure time for model updates
	start := time.Now()
	for i := 0; i < 100; i++ {
		updatedModel, _ := model.Update(tickMsg(time.Now()))
		model = updatedModel.(Model)
	}
	duration := time.Since(start)

	// Should complete within reasonable time
	if duration > time.Second*2 {
		t.Errorf("Model updates took too long: %v", duration)
	}

	t.Logf("200 entities, 100 model updates: %v", duration)
}

// Test View Rendering Performance
func TestViewRenderingPerformance(t *testing.T) {
	model := initialModel()
	model.termWidth = 120
	model.termHeight = 40
	model.updatePaneDimensions()
	model.ready = true

	// Add entities
	for i := 0; i < 100; i++ {
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
		model = updatedModel.(Model)
	}

	// Measure view rendering time
	start := time.Now()
	for i := 0; i < 100; i++ {
		view := model.View()
		if view == "" {
			t.Error("View should not be empty")
		}
	}
	duration := time.Since(start)

	// Should complete within reasonable time
	if duration > time.Second*1 {
		t.Errorf("View rendering took too long: %v", duration)
	}

	t.Logf("100 entities, 100 view renders: %v", duration)
}

// Test Stress Test Performance
func TestStressTestPerformance(t *testing.T) {
	model := initialModel()
	model.termWidth = 100
	model.termHeight = 50
	model.updatePaneDimensions()
	model.ready = true
	model.maxEntityLimit = 1000

	// Measure stress test execution time
	start := time.Now()
	for i := 0; i < 10; i++ {
		model.runStressTest()
	}
	duration := time.Since(start)

	// Should complete within reasonable time
	if duration > time.Second*5 {
		t.Errorf("Stress test took too long: %v", duration)
	}

	t.Logf("10 stress tests: %v, final entity count: %d", duration, model.entityManager.Count())
}

// Test Responsive Layout Performance
func TestResponsiveLayoutPerformance(t *testing.T) {
	model := initialModel()

	// Test multiple window size changes
	sizes := []struct{ w, h int }{
		{80, 24}, {120, 30}, {60, 20}, {200, 50}, {40, 15},
	}

	start := time.Now()
	for i := 0; i < 100; i++ {
		size := sizes[i%len(sizes)]
		windowMsg := tea.WindowSizeMsg{Width: size.w, Height: size.h}
		updatedModel, _ := model.Update(windowMsg)
		model = updatedModel.(Model)

		// Also test view rendering with each size
		view := model.View()
		if view == "" {
			t.Error("View should not be empty")
		}
	}
	duration := time.Since(start)

	// Should complete within reasonable time
	if duration > time.Second*2 {
		t.Errorf("Responsive layout updates took too long: %v", duration)
	}

	t.Logf("100 layout updates: %v", duration)
}

// Test Control Panel Performance
func TestControlPanelPerformance(t *testing.T) {
	cp := NewControlPanel(80, 20)

	start := time.Now()
	for i := 0; i < 1000; i++ {
		// Test navigation
		cp.Update(tea.KeyMsg{Type: tea.KeyTab})

		// Test view rendering
		view := cp.View()
		if view == "" {
			t.Error("Control panel view should not be empty")
		}

		// Test responsive mode updates
		cp.UpdateResponsiveMode(80+i%40, 20+i%10)
	}
	duration := time.Since(start)

	// Should complete within reasonable time
	if duration > time.Second*1 {
		t.Errorf("Control panel operations took too long: %v", duration)
	}

	t.Logf("1000 control panel operations: %v", duration)
}

// Test Concurrent Operations (if applicable)
func TestConcurrentSafety(t *testing.T) {
	manager := NewEntityManager()

	// Test that basic operations don't race
	// (Note: The current implementation may not be thread-safe,
	// but this tests for basic robustness)

	done := make(chan bool, 2)

	// Goroutine 1: Add entities
	go func() {
		for i := 0; i < 100; i++ {
			entity := NewSphere(float64(i), float64(i), 1, GetRandomColor())
			manager.AddEntity(entity)
		}
		done <- true
	}()

	// Goroutine 2: Count entities
	go func() {
		for i := 0; i < 100; i++ {
			_ = manager.Count()
			time.Sleep(time.Microsecond) // Small delay
		}
		done <- true
	}()

	// Wait for both to complete
	<-done
	<-done

	// Should have approximately 100 entities (exact count may vary due to timing)
	count := manager.Count()
	if count < 50 || count > 150 {
		t.Errorf("Unexpected entity count after concurrent operations: %d", count)
	}

	t.Logf("Entity count after concurrent test: %d", count)
}

// Benchmark Entity Creation
func BenchmarkEntityCreation(b *testing.B) {
	colors := GetAvailableColors()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			_ = NewSphere(float64(i), float64(i), i%4+1, colors[i%len(colors)])
		} else {
			_ = NewSprite(float64(i), float64(i), i%4+1, colors[i%len(colors)], "")
		}
	}
}

// Benchmark Random Functions
func BenchmarkGetRandomColor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetRandomColor()
	}
}

// Test Memory Stability Under Extended Load
func TestMemoryStabilityExtendedLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping extended load test in short mode")
	}

	model := initialModel()
	model.termWidth = 100
	model.termHeight = 50
	model.updatePaneDimensions()
	model.ready = true
	model.maxEntityLimit = 500

	// Run simulation for extended period
	for cycle := 0; cycle < 100; cycle++ {
		// Add some entities
		for i := 0; i < 10; i++ {
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			model = updatedModel.(Model)
		}

		// Run physics for a while
		for i := 0; i < 20; i++ {
			updatedModel, _ := model.Update(tickMsg(time.Now()))
			model = updatedModel.(Model)
		}

		// Occasionally clear and reset
		if cycle%20 == 0 {
			updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			model = updatedModel.(Model)
		}

		// Generate view
		view := model.View()
		if view == "" {
			t.Error("View should not be empty during extended load test")
		}
	}

	t.Logf("Extended load test completed successfully")
}
