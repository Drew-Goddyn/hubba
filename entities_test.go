package main

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// Test Entity Creation
func TestNewSphere(t *testing.T) {
	sphere := NewSphere(10.0, 5.0, 2, lipgloss.Color("32"))

	// Test position
	x, y := sphere.GetPosition()
	if x != 10.0 || y != 5.0 {
		t.Errorf("Expected position (10.0, 5.0), got (%.1f, %.1f)", x, y)
	}

	// Test properties
	if sphere.GetSize() != 2 {
		t.Errorf("Expected size 2, got %d", sphere.GetSize())
	}

	if sphere.GetType() != SphereType {
		t.Errorf("Expected type %s, got %s", SphereType, sphere.GetType())
	}

	if sphere.GetSymbol() != "●" {
		t.Errorf("Expected symbol '●', got '%s'", sphere.GetSymbol())
	}

	if sphere.GetRadius() != 0.5 {
		t.Errorf("Expected radius 0.5, got %.1f", sphere.GetRadius())
	}

	// Test ID is not empty
	if sphere.GetID() == "" {
		t.Error("Expected non-empty ID")
	}
}

func TestNewSprite(t *testing.T) {
	sprite := NewSprite(15.0, 8.0, 1, lipgloss.Color("34"), "★")

	// Test position
	x, y := sprite.GetPosition()
	if x != 15.0 || y != 8.0 {
		t.Errorf("Expected position (15.0, 8.0), got (%.1f, %.1f)", x, y)
	}

	// Test properties
	if sprite.GetSize() != 1 {
		t.Errorf("Expected size 1, got %d", sprite.GetSize())
	}

	if sprite.GetType() != SpriteType {
		t.Errorf("Expected type %s, got %s", SpriteType, sprite.GetType())
	}

	if sprite.GetSymbol() != "★" {
		t.Errorf("Expected symbol '★', got '%s'", sprite.GetSymbol())
	}

	// Test ID is not empty
	if sprite.GetID() == "" {
		t.Error("Expected non-empty ID")
	}
}

func TestSpriteWithRandomSymbol(t *testing.T) {
	sprite := NewSprite(0.0, 0.0, 1, lipgloss.Color("31"), "")

	// Should have a random symbol (not empty)
	if sprite.GetSymbol() == "" {
		t.Error("Expected non-empty random symbol")
	}

	// Should be one of the predefined symbols
	validSymbols := []string{"◆", "◇", "★", "☆", "▲", "△", "♦", "♢"}
	symbol := sprite.GetSymbol()
	found := false
	for _, validSymbol := range validSymbols {
		if symbol == validSymbol {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Symbol '%s' is not in the list of valid random symbols", symbol)
	}
}

// Test Entity Manager
func TestEntityManager(t *testing.T) {
	manager := NewEntityManager()

	// Test initial state
	if manager.Count() != 0 {
		t.Errorf("Expected count 0, got %d", manager.Count())
	}

	if len(manager.GetEntities()) != 0 {
		t.Errorf("Expected empty entities slice, got length %d", len(manager.GetEntities()))
	}
}

func TestEntityManagerAddEntity(t *testing.T) {
	manager := NewEntityManager()
	sphere := NewSphere(5.0, 5.0, 1, lipgloss.Color("32"))

	// Add entity
	manager.AddEntity(sphere)

	// Test count
	if manager.Count() != 1 {
		t.Errorf("Expected count 1, got %d", manager.Count())
	}

	// Test entities slice
	entities := manager.GetEntities()
	if len(entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(entities))
	}

	if entities[0].GetID() != sphere.GetID() {
		t.Errorf("Expected entity ID %s, got %s", sphere.GetID(), entities[0].GetID())
	}
}

func TestEntityManagerRemoveEntity(t *testing.T) {
	manager := NewEntityManager()
	sphere := NewSphere(5.0, 5.0, 1, lipgloss.Color("32"))
	sprite := NewSprite(10.0, 10.0, 1, lipgloss.Color("34"), "★")

	// Add entities
	manager.AddEntity(sphere)
	manager.AddEntity(sprite)

	if manager.Count() != 2 {
		t.Errorf("Expected count 2, got %d", manager.Count())
	}

	// Remove sphere
	removed := manager.RemoveEntity(sphere.GetID())
	if !removed {
		t.Error("Expected to successfully remove entity")
	}

	if manager.Count() != 1 {
		t.Errorf("Expected count 1 after removal, got %d", manager.Count())
	}

	// Verify sprite is still there
	entities := manager.GetEntities()
	if len(entities) != 1 || entities[0].GetID() != sprite.GetID() {
		t.Error("Expected sprite to remain after sphere removal")
	}

	// Try to remove non-existent entity
	removed = manager.RemoveEntity("non-existent-id")
	if removed {
		t.Error("Expected to fail when removing non-existent entity")
	}
}

func TestEntityManagerGetEntitiesByType(t *testing.T) {
	manager := NewEntityManager()

	// Add multiple entities of different types
	sphere1 := NewSphere(1.0, 1.0, 1, lipgloss.Color("32"))
	sphere2 := NewSphere(2.0, 2.0, 1, lipgloss.Color("33"))
	sprite1 := NewSprite(3.0, 3.0, 1, lipgloss.Color("34"), "★")
	sprite2 := NewSprite(4.0, 4.0, 1, lipgloss.Color("35"), "◆")

	manager.AddEntity(sphere1)
	manager.AddEntity(sprite1)
	manager.AddEntity(sphere2)
	manager.AddEntity(sprite2)

	// Test getting spheres
	spheres := manager.GetEntitiesByType(SphereType)
	if len(spheres) != 2 {
		t.Errorf("Expected 2 spheres, got %d", len(spheres))
	}

	// Test getting sprites
	sprites := manager.GetEntitiesByType(SpriteType)
	if len(sprites) != 2 {
		t.Errorf("Expected 2 sprites, got %d", len(sprites))
	}

	// Test count by type
	if manager.CountByType(SphereType) != 2 {
		t.Errorf("Expected 2 spheres by count, got %d", manager.CountByType(SphereType))
	}

	if manager.CountByType(SpriteType) != 2 {
		t.Errorf("Expected 2 sprites by count, got %d", manager.CountByType(SpriteType))
	}
}

func TestEntityManagerClear(t *testing.T) {
	manager := NewEntityManager()

	// Add entities
	manager.AddEntity(NewSphere(1.0, 1.0, 1, lipgloss.Color("32")))
	manager.AddEntity(NewSprite(2.0, 2.0, 1, lipgloss.Color("34"), "★"))

	if manager.Count() != 2 {
		t.Errorf("Expected count 2 before clear, got %d", manager.Count())
	}

	// Clear all
	manager.Clear()

	if manager.Count() != 0 {
		t.Errorf("Expected count 0 after clear, got %d", manager.Count())
	}

	if len(manager.GetEntities()) != 0 {
		t.Errorf("Expected empty entities slice after clear, got length %d", len(manager.GetEntities()))
	}
}

// Test Entity Behavior
func TestEntityVelocity(t *testing.T) {
	sphere := NewSphere(0.0, 0.0, 1, lipgloss.Color("32"))

	// Initial velocity should be 0
	vx, vy := sphere.GetVelocity()
	if vx != 0.0 || vy != 0.0 {
		t.Errorf("Expected initial velocity (0.0, 0.0), got (%.1f, %.1f)", vx, vy)
	}

	// Set velocity
	sphere.SetVelocity(5.0, -3.0)
	vx, vy = sphere.GetVelocity()
	if vx != 5.0 || vy != -3.0 {
		t.Errorf("Expected velocity (5.0, -3.0), got (%.1f, %.1f)", vx, vy)
	}
}

func TestEntityUpdate(t *testing.T) {
	sphere := NewSphere(10.0, 10.0, 1, lipgloss.Color("32"))
	sphere.SetVelocity(2.0, -1.0)

	// Update with deltaTime = 1.0
	sphere.Update(1.0)

	x, y := sphere.GetPosition()
	expectedX, expectedY := 12.0, 9.0
	if x != expectedX || y != expectedY {
		t.Errorf("Expected position (%.1f, %.1f), got (%.1f, %.1f)", expectedX, expectedY, x, y)
	}
}

func TestEntityApplyForce(t *testing.T) {
	sphere := NewSphere(0.0, 0.0, 2, lipgloss.Color("32")) // mass = 1.0 (effectiveSize for size 2)

	// Apply force
	sphere.ApplyForce(4.0, -2.0)

	// Check velocity (F = ma, so a = F/m, and velocity changes by acceleration)
	vx, vy := sphere.GetVelocity()
	expectedVX, expectedVY := 4.0, -2.0 // 4.0/1.0, -2.0/1.0 (mass is effectiveSize = 1.0)
	if vx != expectedVX || vy != expectedVY {
		t.Errorf("Expected velocity (%.1f, %.1f), got (%.1f, %.1f)", expectedVX, expectedVY, vx, vy)
	}
}

func TestEntityCollisionDetection(t *testing.T) {
	sphere1 := NewSphere(5.0, 5.0, 2, lipgloss.Color("32"))
	sphere2 := NewSphere(6.0, 6.0, 2, lipgloss.Color("33"))   // Overlapping
	sphere3 := NewSphere(10.0, 10.0, 2, lipgloss.Color("34")) // Not overlapping

	// Test collision detection
	if !sphere1.CheckCollision(sphere2) {
		t.Error("Expected collision between sphere1 and sphere2")
	}

	if sphere1.CheckCollision(sphere3) {
		t.Error("Expected no collision between sphere1 and sphere3")
	}

	if sphere2.CheckCollision(sphere3) {
		t.Error("Expected no collision between sphere2 and sphere3")
	}
}

func TestEntityManagerCollisions(t *testing.T) {
	manager := NewEntityManager()

	sphere1 := NewSphere(5.0, 5.0, 2, lipgloss.Color("32"))
	sphere2 := NewSphere(6.0, 6.0, 2, lipgloss.Color("33"))   // Overlapping with sphere1
	sphere3 := NewSphere(10.0, 10.0, 2, lipgloss.Color("34")) // Not overlapping

	manager.AddEntity(sphere1)
	manager.AddEntity(sphere2)
	manager.AddEntity(sphere3)

	collisions := manager.CheckCollisions()

	// Should have exactly one collision (sphere1 and sphere2)
	if len(collisions) != 1 {
		t.Errorf("Expected 1 collision, got %d", len(collisions))
	}

	if len(collisions) > 0 {
		collision := collisions[0]
		entity1ID := collision.Entity1.GetID()
		entity2ID := collision.Entity2.GetID()

		// Should be sphere1 and sphere2 (order may vary)
		if !((entity1ID == sphere1.GetID() && entity2ID == sphere2.GetID()) ||
			(entity1ID == sphere2.GetID() && entity2ID == sphere1.GetID())) {
			t.Error("Collision should be between sphere1 and sphere2")
		}
	}
}

func TestEntityRendering(t *testing.T) {
	sphere := NewSphere(0.0, 0.0, 1, lipgloss.Color("32"))
	sprite := NewSprite(0.0, 0.0, 1, lipgloss.Color("34"), "★")

	// Test sphere rendering
	sphereRender := sphere.Render()
	if sphereRender == "" {
		t.Error("Expected non-empty sphere render")
	}

	// Test sprite rendering
	spriteRender := sprite.Render()
	if spriteRender == "" {
		t.Error("Expected non-empty sprite render")
	}

	// Renders should be different (different symbols)
	if sphereRender == spriteRender {
		t.Error("Expected different renders for sphere and sprite")
	}
}

func TestGetRandomColor(t *testing.T) {
	color := GetRandomColor()

	// Should return a valid color (not empty)
	if string(color) == "" {
		t.Error("Expected non-empty color")
	}

	// Test that it generates different colors (probabilistic test)
	colors := make(map[string]bool)
	for i := 0; i < 20; i++ {
		color := GetRandomColor()
		colors[string(color)] = true
	}

	// Should have at least 2 different colors out of 20 attempts
	if len(colors) < 2 {
		t.Error("Expected at least 2 different colors from GetRandomColor")
	}
}
