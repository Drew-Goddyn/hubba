package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

// EntityType represents the type of entity
type EntityType string

const (
	SphereType EntityType = "sphere"
	SpriteType EntityType = "sprite"
)

// Entity interface defines the common behavior for all simulation entities
type Entity interface {
	// Position and movement
	GetPosition() (float64, float64)        // Physics position (target)
	GetDisplayPosition() (float64, float64) // Visual position (animated)
	SetPosition(x, y float64)               // Set physics target (animated)
	SetImmediatePosition(x, y float64)      // Set both physics and display (instant)
	GetVelocity() (float64, float64)
	SetVelocity(vx, vy float64)

	// Visual properties
	GetSymbol() string
	GetColor() lipgloss.Color
	GetSize() int

	// Entity properties
	GetType() EntityType
	GetID() string

	// Physics
	ApplyForce(fx, fy float64)
	Update(deltaTime float64)

	// Animation
	GetAnimationState() *EntityAnimationState
	UpdateAnimation(ae *AnimationEngine)

	// Collision
	GetBounds() (x, y, width, height float64)
	CheckCollision(other Entity) bool

	// Rendering
	Render() string
}

// BaseEntity provides common functionality for all entities
type BaseEntity struct {
	ID     string
	X, Y   float64 // Physics position (target)
	VX, VY float64
	Size   int
	Color  lipgloss.Color
	Symbol string
	Type   EntityType
	Mass   float64

	// Animation state
	AnimationState *EntityAnimationState
}

// Position methods
func (e *BaseEntity) GetPosition() (float64, float64) {
	return e.X, e.Y
}

func (e *BaseEntity) SetPosition(x, y float64) {
	e.X, e.Y = x, y
}

func (e *BaseEntity) SetImmediatePosition(x, y float64) {
	e.X, e.Y = x, y
	// Also immediately update animation state to prevent smooth interpolation
	if e.AnimationState != nil {
		e.AnimationState.SetInitialPosition(x, y)
	}
}

func (e *BaseEntity) GetVelocity() (float64, float64) {
	return e.VX, e.VY
}

func (e *BaseEntity) SetVelocity(vx, vy float64) {
	// Validate and sanitize velocity inputs
	if math.IsInf(vx, 0) || math.IsNaN(vx) {
		vx = 0 // Reset invalid X velocity
	}
	if math.IsInf(vy, 0) || math.IsNaN(vy) {
		vy = 0 // Reset invalid Y velocity
	}
	e.VX, e.VY = vx, vy
}

// Visual properties
func (e *BaseEntity) GetSymbol() string {
	return e.Symbol
}

func (e *BaseEntity) GetColor() lipgloss.Color {
	return e.Color
}

func (e *BaseEntity) GetSize() int {
	return e.Size
}

// Entity properties
func (e *BaseEntity) GetType() EntityType {
	return e.Type
}

func (e *BaseEntity) GetID() string {
	return e.ID
}

// Physics
func (e *BaseEntity) ApplyForce(fx, fy float64) {
	// F = ma, so a = F/m
	if e.Mass > 0 {
		e.VX += fx / e.Mass
		e.VY += fy / e.Mass
	}
}

func (e *BaseEntity) Update(deltaTime float64) {
	// Update physics position based on velocity
	e.X += e.VX * deltaTime
	e.Y += e.VY * deltaTime

	// Update animation target to physics position
	if e.AnimationState != nil {
		e.AnimationState.SetTarget(e.X, e.Y)
	}
}

// Collision detection
func (e *BaseEntity) GetBounds() (x, y, width, height float64) {
	// Adjust collision size to better match visual representation
	// Smaller collision boxes for single-character entities
	var effectiveSize float64
	switch e.Size {
	case 1:
		effectiveSize = 0.8 // Tiny
	case 2:
		effectiveSize = 1.0 // Small
	case 3:
		effectiveSize = 1.3 // Medium
	case 4:
		effectiveSize = 1.6 // Large
	default:
		effectiveSize = float64(e.Size) * 0.8
	}

	return e.X - effectiveSize/2, e.Y - effectiveSize/2, effectiveSize, effectiveSize
}

func (e *BaseEntity) CheckCollision(other Entity) bool {
	x1, y1, w1, h1 := e.GetBounds()
	x2, y2, w2, h2 := other.GetBounds()

	return !(x1+w1 < x2 || x2+w2 < x1 || y1+h1 < y2 || y2+h2 < y1)
}

// Animation methods
func (e *BaseEntity) GetDisplayPosition() (float64, float64) {
	if e.AnimationState != nil {
		return e.AnimationState.GetDisplayPosition()
	}
	// Fallback to physics position if no animation
	return e.X, e.Y
}

func (e *BaseEntity) GetAnimationState() *EntityAnimationState {
	return e.AnimationState
}

func (e *BaseEntity) UpdateAnimation(ae *AnimationEngine) {
	if e.AnimationState != nil {
		ae.UpdateAnimation(e.AnimationState)
	}
}

// Rendering with enhanced visual polish and effects
func (e *BaseEntity) Render() string {
	// Create enhanced style with better visibility
	style := lipgloss.NewStyle().
		Foreground(e.Color).
		Bold(true)

	// Create visual representation that matches collision size
	switch e.Size {
	case 1:
		// Tiny entities - single character
		if e.Type == SphereType {
			return style.Render("●") // Small filled circle
		} else {
			return style.Render("◆") // Small diamond
		}
	case 2:
		// Small entities - single character but different symbol
		if e.Type == SphereType {
			return style.Render("⬤") // Medium filled circle
		} else {
			return style.Render("◉") // Medium diamond with dot
		}
	case 3:
		// Medium entities - larger visual symbols
		if e.Type == SphereType {
			return style.Render("⭘") // Large circle with ring
		} else {
			return style.Render("⬢") // Large hexagon
		}
	case 4:
		// Large entities - biggest symbols
		if e.Type == SphereType {
			return style.Render("⬢") // Extra large hexagon
		} else {
			return style.Render("⬛") // Large black square
		}
	default:
		return style.Render(e.Symbol)
	}
}

// Sphere represents a circular entity
type Sphere struct {
	BaseEntity
	Radius float64
}

// NewSphere creates a new sphere entity
func NewSphere(x, y float64, size int, color lipgloss.Color) *Sphere {
	// Validate and sanitize size input
	if size < 0 {
		size = 1 // Default to minimum valid size
	}

	// Create animation engine for this entity
	animEngine := NewAnimationEngine()
	animState := animEngine.NewEntityAnimationState(x, y)

	// Calculate effective radius to match visual representation
	var effectiveSize float64
	switch size {
	case 1:
		effectiveSize = 0.8 // Tiny
	case 2:
		effectiveSize = 1.0 // Small
	case 3:
		effectiveSize = 1.3 // Medium
	case 4:
		effectiveSize = 1.6 // Large
	default:
		effectiveSize = float64(size) * 0.8
	}

	return &Sphere{
		BaseEntity: BaseEntity{
			ID:             generateID("sphere"),
			X:              x,
			Y:              y,
			VX:             0,
			VY:             0,
			Size:           size,
			Color:          color,
			Symbol:         "●",
			Type:           SphereType,
			Mass:           effectiveSize, // Mass proportional to effective size
			AnimationState: animState,
		},
		Radius: effectiveSize / 2.0,
	}
}

// Sphere-specific methods
func (s *Sphere) GetRadius() float64 {
	return s.Radius
}

func (s *Sphere) SetRadius(radius float64) {
	s.Radius = radius
	// Update size to match new radius (approximately)
	if radius <= 0.4 {
		s.Size = 1
	} else if radius <= 0.5 {
		s.Size = 2
	} else if radius <= 0.65 {
		s.Size = 3
	} else {
		s.Size = 4
	}
}

// Override GetBounds for circular collision using effective size
func (s *Sphere) GetBounds() (x, y, width, height float64) {
	return s.X - s.Radius, s.Y - s.Radius, s.Radius * 2, s.Radius * 2
}

// Sprite represents a custom character entity
type Sprite struct {
	BaseEntity
	CustomSymbol string
	Animation    []string
	CurrentFrame int
}

// NewSprite creates a new sprite entity
func NewSprite(x, y float64, size int, color lipgloss.Color, customSymbol string) *Sprite {
	symbol := customSymbol
	if symbol == "" {
		// Default sprite symbols
		symbols := []string{"◆", "◇", "★", "☆", "▲", "△", "♦", "♢"}
		symbol = symbols[rand.Intn(len(symbols))]
	}

	// Create animation engine for this entity
	animEngine := NewAnimationEngine()
	animState := animEngine.NewEntityAnimationState(x, y)

	// Calculate effective size to match visual representation
	var effectiveSize float64
	switch size {
	case 1:
		effectiveSize = 0.8 // Tiny
	case 2:
		effectiveSize = 1.0 // Small
	case 3:
		effectiveSize = 1.3 // Medium
	case 4:
		effectiveSize = 1.6 // Large
	default:
		effectiveSize = float64(size) * 0.8
	}

	return &Sprite{
		BaseEntity: BaseEntity{
			ID:             generateID("sprite"),
			X:              x,
			Y:              y,
			VX:             0,
			VY:             0,
			Size:           size,
			Color:          color,
			Symbol:         symbol,
			Type:           SpriteType,
			Mass:           effectiveSize * 0.8, // Sprites are slightly lighter than spheres
			AnimationState: animState,
		},
		CustomSymbol: symbol,
		Animation:    []string{symbol}, // Single frame by default
		CurrentFrame: 0,
	}
}

// Sprite-specific methods
func (s *Sprite) SetAnimation(frames []string) {
	if len(frames) > 0 {
		s.Animation = frames
		s.Symbol = frames[0]
	}
}

func (s *Sprite) NextFrame() {
	if len(s.Animation) > 1 {
		s.CurrentFrame = (s.CurrentFrame + 1) % len(s.Animation)
		s.Symbol = s.Animation[s.CurrentFrame]
	}
}

// Override Update to handle animation
func (s *Sprite) Update(deltaTime float64) {
	s.BaseEntity.Update(deltaTime)
	// Animate every few updates (simplified)
	if rand.Float64() < 0.1 { // 10% chance to animate per update
		s.NextFrame()
	}
}

// EntityManager manages a collection of entities with thread-safe operations
type EntityManager struct {
	mu       sync.RWMutex // Protects entities slice from concurrent access
	entities []Entity
	nextID   int
}

// NewEntityManager creates a new entity manager
func NewEntityManager() *EntityManager {
	return &EntityManager{
		entities: make([]Entity, 0),
		nextID:   1,
	}
}

// AddEntity adds an entity to the manager (thread-safe)
func (em *EntityManager) AddEntity(entity Entity) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.entities = append(em.entities, entity)
}

// RemoveEntity removes an entity by ID (thread-safe)
func (em *EntityManager) RemoveEntity(id string) bool {
	em.mu.Lock()
	defer em.mu.Unlock()
	for i, entity := range em.entities {
		if entity.GetID() == id {
			// Remove entity from slice
			em.entities = append(em.entities[:i], em.entities[i+1:]...)
			return true
		}
	}
	return false
}

// GetEntities returns a copy of all entities to prevent concurrent modification issues (thread-safe)
func (em *EntityManager) GetEntities() []Entity {
	em.mu.RLock()
	defer em.mu.RUnlock()
	// Return a copy to prevent concurrent modification
	entities := make([]Entity, len(em.entities))
	copy(entities, em.entities)
	return entities
}

// GetEntitiesByType returns entities of a specific type
func (em *EntityManager) GetEntitiesByType(entityType EntityType) []Entity {
	var result []Entity
	for _, entity := range em.entities {
		if entity.GetType() == entityType {
			result = append(result, entity)
		}
	}
	return result
}

// Clear removes all entities (thread-safe)
func (em *EntityManager) Clear() {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.entities = make([]Entity, 0)
}

// Count returns the number of entities (thread-safe)
func (em *EntityManager) Count() int {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return len(em.entities)
}

// CountByType returns the count of entities by type (thread-safe)
func (em *EntityManager) CountByType(entityType EntityType) int {
	em.mu.RLock()
	defer em.mu.RUnlock()
	count := 0
	for _, entity := range em.entities {
		if entity.GetType() == entityType {
			count++
		}
	}
	return count
}

// Update updates all entities (thread-safe)
func (em *EntityManager) Update(deltaTime float64) {
	em.mu.RLock()
	defer em.mu.RUnlock()
	for _, entity := range em.entities {
		entity.Update(deltaTime)
	}
}

// CheckCollisions checks for collisions between all entities (thread-safe)
func (em *EntityManager) CheckCollisions() []CollisionPair {
	em.mu.RLock()
	defer em.mu.RUnlock()
	var collisions []CollisionPair

	for i := 0; i < len(em.entities); i++ {
		for j := i + 1; j < len(em.entities); j++ {
			if em.entities[i].CheckCollision(em.entities[j]) {
				collisions = append(collisions, CollisionPair{
					Entity1: em.entities[i],
					Entity2: em.entities[j],
				})
			}
		}
	}

	return collisions
}

// CollisionPair represents two entities that are colliding
type CollisionPair struct {
	Entity1, Entity2 Entity
}

// Utility functions

// generateID generates a unique ID for entities
func generateID(prefix string) string {
	return fmt.Sprintf("%s_%d_%d", prefix, rand.Intn(10000), rand.Intn(10000))
}

// GetRandomColor returns a random color for entities using reliable hex colors
func GetRandomColor() lipgloss.Color {
	colors := []lipgloss.Color{
		lipgloss.Color("#00FF00"), // Green
		lipgloss.Color("#FFFF00"), // Yellow
		lipgloss.Color("#0000FF"), // Blue
		lipgloss.Color("#FF00FF"), // Magenta
		lipgloss.Color("#00FFFF"), // Cyan
		lipgloss.Color("#FF0000"), // Red
		lipgloss.Color("#FFFFFF"), // White
		lipgloss.Color("#FF6B6B"), // Light Red
		lipgloss.Color("#4ECDC4"), // Teal
		lipgloss.Color("#45B7D1"), // Sky Blue
		lipgloss.Color("#96CEB4"), // Light Green
		lipgloss.Color("#FECA57"), // Orange
		lipgloss.Color("#A29BFE"), // Purple
	}
	return colors[rand.Intn(len(colors))]
}
