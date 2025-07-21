package main

import (
	"math"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// Test Physics Engine Creation
func TestNewPhysicsEngine(t *testing.T) {
	pe := NewPhysicsEngine(100, 50)

	// Test default values
	if pe.Gravity != 25.0 {
		t.Errorf("Expected gravity 25.0, got %.1f", pe.Gravity)
	}

	if pe.AirResistance != 0.05 {
		t.Errorf("Expected air resistance 0.05, got %.2f", pe.AirResistance)
	}

	if pe.Restitution != 0.7 {
		t.Errorf("Expected restitution 0.7, got %.1f", pe.Restitution)
	}

	// Test bounds
	if pe.MaxX != 98.0 || pe.MaxY != 48.0 {
		t.Errorf("Expected bounds (98, 48), got (%.1f, %.1f)", pe.MaxX, pe.MaxY)
	}

	if pe.MinX != 1.0 || pe.MinY != 1.0 {
		t.Errorf("Expected min bounds (1, 1), got (%.1f, %.1f)", pe.MinX, pe.MinY)
	}

	if pe.DeltaTime != 0.1 {
		t.Errorf("Expected deltaTime 0.1, got %.1f", pe.DeltaTime)
	}
}

// Test Gravity Application
func TestApplyGravity(t *testing.T) {
	pe := NewPhysicsEngine(100, 50)
	sphere := NewSphere(10.0, 10.0, 1, lipgloss.Color("32"))

	// Initial velocity should be 0
	vx, vy := sphere.GetVelocity()
	if vx != 0.0 || vy != 0.0 {
		t.Errorf("Expected initial velocity (0, 0), got (%.1f, %.1f)", vx, vy)
	}

	// Apply physics (which includes gravity)
	pe.ApplyPhysics([]Entity{sphere})

	// Check that velocity has increased downward due to gravity
	vx, vy = sphere.GetVelocity()
	if vy <= 0 {
		t.Errorf("Expected positive Y velocity after gravity, got %.2f", vy)
	}

	if vx != 0 {
		t.Errorf("Expected X velocity to remain 0, got %.2f", vx)
	}
}

// Test Position Updates
func TestPositionUpdate(t *testing.T) {
	pe := NewPhysicsEngine(100, 50)
	sphere := NewSphere(10.0, 10.0, 1, lipgloss.Color("32"))

	// Set initial velocity
	sphere.SetVelocity(5.0, -3.0)

	// Apply physics
	pe.ApplyPhysics([]Entity{sphere})

	// Check position has updated based on velocity
	x, y := sphere.GetPosition()

	// Position should have changed from (10, 10)
	// X should increase by velocity * deltaTime: 10 + 5*0.1 = 10.5
	expectedX := 10.5
	if math.Abs(x-expectedX) > 0.1 {
		t.Errorf("Expected X position ~%.1f, got %.2f", expectedX, x)
	}

	// Y should change by velocity * deltaTime, but also affected by gravity
	// Initial VY=-3.0 (upward), but gravity should overcome this and pull down (increase Y)
	if y <= 10.0 {
		t.Errorf("Expected Y position to increase from 10.0 due to gravity, got %.2f", y)
	}
}

// Test Boundary Collisions
func TestBoundaryCollisions(t *testing.T) {
	pe := NewPhysicsEngine(20, 20) // Small bounds for testing

	// Test left boundary
	sphere := NewSphere(0.5, 10.0, 1, lipgloss.Color("32")) // Near left edge
	sphere.SetVelocity(-5.0, 0.0)                           // Moving left

	pe.ApplyPhysics([]Entity{sphere})

	x, _ := sphere.GetPosition()
	vx, _ := sphere.GetVelocity()

	// Should be pushed away from left boundary
	if x <= pe.MinX {
		t.Errorf("Entity should be moved away from left boundary, x=%.2f", x)
	}

	// Velocity should be reversed and reduced by restitution
	if vx <= 0 {
		t.Errorf("X velocity should be positive after left wall collision, got %.2f", vx)
	}

	// Test bottom boundary
	sphere2 := NewSphere(10.0, 18.5, 1, lipgloss.Color("33")) // Near bottom
	sphere2.SetVelocity(0.0, 5.0)                             // Moving down

	pe.ApplyPhysics([]Entity{sphere2})

	_, y2 := sphere2.GetPosition()
	_, vy2 := sphere2.GetVelocity()

	// Should be pushed away from bottom boundary
	if y2 >= pe.MaxY {
		t.Errorf("Entity should be moved away from bottom boundary, y=%.2f", y2)
	}

	// Velocity should be reversed
	if vy2 >= 0 {
		t.Errorf("Y velocity should be negative after bottom wall collision, got %.2f", vy2)
	}
}

// Test Physics Engine Entity Collision Detection
func TestPhysicsEntityCollisionDetection(t *testing.T) {
	pe := NewPhysicsEngine(100, 50)

	// Create two overlapping entities
	sphere1 := NewSphere(10.0, 10.0, 2, lipgloss.Color("32"))
	sphere2 := NewSphere(10.5, 10.5, 2, lipgloss.Color("33")) // Actually overlapping (distance ~0.7 < minDistance 0.9)
	sphere3 := NewSphere(20.0, 20.0, 2, lipgloss.Color("34")) // Not overlapping

	entities := []Entity{sphere1, sphere2, sphere3}

	// Check collision detection
	collisions := pe.findCollisions(entities)

	// Should detect one collision between sphere1 and sphere2
	if len(collisions) != 1 {
		t.Errorf("Expected 1 collision, got %d", len(collisions))
	}

	if len(collisions) > 0 {
		collision := collisions[0]

		// Verify it's the right collision
		id1 := collision.Entity1.GetID()
		id2 := collision.Entity2.GetID()

		if !((id1 == sphere1.GetID() && id2 == sphere2.GetID()) ||
			(id1 == sphere2.GetID() && id2 == sphere1.GetID())) {
			t.Error("Collision should be between sphere1 and sphere2")
		}
	}
}

// Test Entity Collision Resolution
func TestEntityCollisionResolution(t *testing.T) {
	pe := NewPhysicsEngine(100, 50)

	// Create two entities moving towards each other (close enough to collide)
	sphere1 := NewSphere(10.0, 10.0, 1, lipgloss.Color("32"))
	sphere2 := NewSphere(10.6, 10.0, 1, lipgloss.Color("33")) // Distance 0.6 < minDistance 0.7, so they collide

	sphere1.SetVelocity(5.0, 0.0)  // Moving right
	sphere2.SetVelocity(-5.0, 0.0) // Moving left

	// Store initial velocities
	vx1_initial, _ := sphere1.GetVelocity()
	vx2_initial, _ := sphere2.GetVelocity()

	// Apply physics including collision detection
	entities := []Entity{sphere1, sphere2}
	pe.HandleEntityCollisions(entities)

	// Get final velocities
	vx1_final, _ := sphere1.GetVelocity()
	vx2_final, _ := sphere2.GetVelocity()

	// After elastic collision, velocities should have changed
	// (exact values depend on collision physics, but they should be different)
	if vx1_final == vx1_initial {
		t.Error("Sphere1 velocity should change after collision")
	}

	if vx2_final == vx2_initial {
		t.Error("Sphere2 velocity should change after collision")
	}

	// Velocities should be opposite to initial directions (simplified test)
	if vx1_final > 0 {
		t.Error("Sphere1 should be moving left after collision")
	}

	if vx2_final < 0 {
		t.Error("Sphere2 should be moving right after collision")
	}
}

// Test Velocity Capping
func TestVelocityCapping(t *testing.T) {
	pe := NewPhysicsEngine(100, 50)
	sphere := NewSphere(10.0, 10.0, 1, lipgloss.Color("32"))

	// Set extreme velocity
	sphere.SetVelocity(100.0, -100.0)

	// Apply physics (which includes velocity capping)
	pe.ApplyPhysics([]Entity{sphere})

	vx, vy := sphere.GetVelocity()

	// Velocity should be capped at MaxVelocity (50.0)
	if math.Abs(vx) > pe.MaxVelocity {
		t.Errorf("X velocity should be capped at %.1f, got %.2f", pe.MaxVelocity, math.Abs(vx))
	}

	if math.Abs(vy) > pe.MaxVelocity {
		t.Errorf("Y velocity should be capped at %.1f, got %.2f", pe.MaxVelocity, math.Abs(vy))
	}
}

// Test Air Resistance
func TestAirResistance(t *testing.T) {
	pe := NewPhysicsEngine(100, 50)
	sphere := NewSphere(50.0, 25.0, 1, lipgloss.Color("32")) // Center position

	// Set initial velocity
	initialVX, initialVY := 10.0, 5.0
	sphere.SetVelocity(initialVX, initialVY)

	// Apply physics multiple times
	for i := 0; i < 10; i++ {
		pe.ApplyPhysics([]Entity{sphere})
	}

	vx, _ := sphere.GetVelocity()

	// Air resistance should reduce velocity magnitude over time
	// (though gravity affects Y velocity, so we mainly check X)
	if math.Abs(vx) >= math.Abs(initialVX) {
		t.Errorf("Air resistance should reduce X velocity from %.1f to less, got %.2f", initialVX, vx)
	}
}

// Test Physics Engine Settings
func TestPhysicsEngineSettings(t *testing.T) {
	pe := NewPhysicsEngine(100, 50)

	// Test gravity setting
	pe.SetGravity(50.0)
	if pe.GetGravity() != 50.0 {
		t.Errorf("Expected gravity 50.0, got %.1f", pe.GetGravity())
	}

	// Test restitution setting
	pe.SetRestitution(0.9)
	if pe.GetRestitution() != 0.9 {
		t.Errorf("Expected restitution 0.9, got %.1f", pe.GetRestitution())
	}

	// Test invalid restitution (should be ignored)
	pe.SetRestitution(1.5) // Invalid (> 1)
	if pe.GetRestitution() != 0.9 {
		t.Errorf("Invalid restitution should be ignored, still expect 0.9, got %.1f", pe.GetRestitution())
	}

	pe.SetRestitution(-0.1) // Invalid (< 0)
	if pe.GetRestitution() != 0.9 {
		t.Errorf("Invalid restitution should be ignored, still expect 0.9, got %.1f", pe.GetRestitution())
	}
}

// Test Pause/Resume Functionality
func TestPauseResume(t *testing.T) {
	pe := NewPhysicsEngine(100, 50)

	// Initially should be running
	if !pe.IsRunning() {
		t.Error("Physics engine should be running initially")
	}

	// Test pause
	pe.Pause()
	if pe.IsRunning() {
		t.Error("Physics engine should be paused after Pause()")
	}

	if pe.DeltaTime != 0 {
		t.Errorf("DeltaTime should be 0 when paused, got %.1f", pe.DeltaTime)
	}

	// Test resume
	pe.Resume()
	if !pe.IsRunning() {
		t.Error("Physics engine should be running after Resume()")
	}

	if pe.DeltaTime != 0.1 {
		t.Errorf("DeltaTime should be 0.1 when running, got %.1f", pe.DeltaTime)
	}
}

// Test Bounds Update
func TestUpdateBounds(t *testing.T) {
	pe := NewPhysicsEngine(100, 50)

	// Update bounds
	pe.UpdateBounds(200, 100)

	if pe.MaxX != 198.0 || pe.MaxY != 98.0 {
		t.Errorf("Expected updated bounds (198, 98), got (%.1f, %.1f)", pe.MaxX, pe.MaxY)
	}

	// Min bounds should remain the same
	if pe.MinX != 1.0 || pe.MinY != 1.0 {
		t.Errorf("Min bounds should remain (1, 1), got (%.1f, %.1f)", pe.MinX, pe.MinY)
	}
}

// Test Random Velocity Addition
func TestAddRandomVelocity(t *testing.T) {
	pe := NewPhysicsEngine(100, 50)
	sphere := NewSphere(10.0, 10.0, 1, lipgloss.Color("32"))

	// Initial velocity should be 0
	vx, vy := sphere.GetVelocity()
	if vx != 0.0 || vy != 0.0 {
		t.Errorf("Expected initial velocity (0, 0), got (%.1f, %.1f)", vx, vy)
	}

	// Add random velocity
	pe.AddRandomVelocity(sphere, 10.0)

	// Velocity should now be non-zero
	vx, vy = sphere.GetVelocity()
	if vx == 0.0 && vy == 0.0 {
		t.Error("Expected non-zero velocity after AddRandomVelocity")
	}

	// Velocity should be within the specified range
	if math.Abs(vx) > 5.0 || math.Abs(vy) > 5.0 {
		t.Errorf("Random velocity should be within ±5.0, got (%.2f, %.2f)", vx, vy)
	}
}

// Test Complete Physics Cycle
func TestCompletePhysicsCycle(t *testing.T) {
	pe := NewPhysicsEngine(50, 30)

	// Create entities at different positions
	sphere1 := NewSphere(25.0, 5.0, 1, lipgloss.Color("32"))  // Center top
	sphere2 := NewSphere(10.0, 10.0, 1, lipgloss.Color("33")) // Left middle

	// Add some initial velocity
	sphere1.SetVelocity(2.0, 0.0)
	sphere2.SetVelocity(-1.0, 3.0)

	entities := []Entity{sphere1, sphere2}

	// Run physics for several cycles
	for i := 0; i < 20; i++ {
		pe.ApplyPhysics(entities)
		pe.HandleEntityCollisions(entities)
	}

	// After running physics, entities should still be within bounds
	for _, entity := range entities {
		x, y := entity.GetPosition()
		size := float64(entity.GetSize())

		if x-size/2 < pe.MinX || x+size/2 > pe.MaxX {
			t.Errorf("Entity X position %.2f is outside bounds [%.1f, %.1f]", x, pe.MinX+size/2, pe.MaxX-size/2)
		}

		if y-size/2 < pe.MinY || y+size/2 > pe.MaxY {
			t.Errorf("Entity Y position %.2f is outside bounds [%.1f, %.1f]", y, pe.MinY+size/2, pe.MaxY-size/2)
		}
	}

	// Entities should have moved from their initial positions due to gravity and velocity
	x1, y1 := sphere1.GetPosition()
	x2, y2 := sphere2.GetPosition()

	if x1 == 25.0 && y1 == 5.0 {
		t.Error("Sphere1 should have moved from initial position")
	}

	if x2 == 10.0 && y2 == 10.0 {
		t.Error("Sphere2 should have moved from initial position")
	}
}
