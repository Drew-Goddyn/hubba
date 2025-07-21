package main

import (
	"math"
	"math/rand"
)

// PhysicsEngine handles all physics calculations and simulations
type PhysicsEngine struct {
	// Physics constants
	Gravity        float64 // Gravity acceleration (pixels/second²)
	AirResistance  float64 // Air resistance coefficient (0-1)
	Restitution    float64 // Bounce factor for collisions (0-1)
	StaticFriction float64 // Static friction when entities are nearly at rest
	ContactDamping float64 // Damping when entities are in contact

	// Simulation bounds
	MinX, MinY float64
	MaxX, MaxY float64

	// Time tracking
	DeltaTime float64 // Time step for physics calculations

	// Performance settings
	MaxVelocity float64 // Maximum velocity cap
	MinVelocity float64 // Minimum velocity threshold (for stopping)

	// Collision precision
	ContactTolerance float64 // How close entities can get before being considered touching
}

// NewPhysicsEngine creates a new physics engine with default settings
func NewPhysicsEngine(boundsWidth, boundsHeight float64) *PhysicsEngine {
	// Validate and sanitize dimensions
	if boundsWidth <= 0 {
		boundsWidth = 10.0 // Default minimum width
	}
	if boundsHeight <= 0 {
		boundsHeight = 10.0 // Default minimum height
	}

	return &PhysicsEngine{
		Gravity:          25.0, // Reasonable gravity for terminal display
		AirResistance:    0.05, // Increased air resistance for better settling
		Restitution:      0.7,  // Bouncy but not perfectly elastic
		StaticFriction:   0.8,  // Strong static friction to prevent jittering
		ContactDamping:   0.9,  // Strong damping when entities touch
		MinX:             1.0,  // Keep entities away from borders
		MinY:             1.0,
		MaxX:             boundsWidth - 2.0,
		MaxY:             boundsHeight - 2.0,
		DeltaTime:        0.1,  // 100ms time steps
		MaxVelocity:      50.0, // Cap velocity for visual reasons
		MinVelocity:      0.05, // Lower threshold for stopping
		ContactTolerance: 0.1,  // Allow entities to touch more closely
	}
}

// UpdateBounds updates the simulation boundaries
func (pe *PhysicsEngine) UpdateBounds(width, height float64) {
	pe.MaxX = width - 2.0
	pe.MaxY = height - 2.0
}

// ApplyPhysics applies all physics calculations to entities
func (pe *PhysicsEngine) ApplyPhysics(entities []Entity) {
	for _, entity := range entities {
		pe.applyGravity(entity)
		pe.applyAirResistance(entity)
		pe.updatePosition(entity)
		pe.handleBoundaryCollisions(entity)
		pe.capVelocity(entity)
	}
}

// applyGravity applies downward gravitational force
func (pe *PhysicsEngine) applyGravity(entity Entity) {
	// Apply gravity force: F = mg (simplified to just g since mass is in the ApplyForce method)
	entity.ApplyForce(0, pe.Gravity*pe.DeltaTime)
}

// applyAirResistance applies air resistance to slow down entities
func (pe *PhysicsEngine) applyAirResistance(entity Entity) {
	vx, vy := entity.GetVelocity()

	// Air resistance opposes motion: F = -k * v
	resistanceX := -pe.AirResistance * vx
	resistanceY := -pe.AirResistance * vy

	entity.ApplyForce(resistanceX, resistanceY)
}

// updatePosition updates entity position based on velocity
func (pe *PhysicsEngine) updatePosition(entity Entity) {
	entity.Update(pe.DeltaTime)
}

// handleBoundaryCollisions keeps entities within the simulation bounds
func (pe *PhysicsEngine) handleBoundaryCollisions(entity Entity) {
	x, y := entity.GetPosition()
	vx, vy := entity.GetVelocity()
	size := float64(entity.GetSize())

	// Calculate entity bounds
	entityMinX := x - size/2
	entityMaxX := x + size/2
	entityMinY := y - size/2
	entityMaxY := y + size/2

	// Horizontal boundary collisions
	if entityMinX <= pe.MinX {
		// Hit left wall
		newX := pe.MinX + size/2
		entity.SetImmediatePosition(newX, y) // Immediate position for crisp bounce
		entity.SetVelocity(-vx*pe.Restitution, vy)
		x = newX // Update position variable for subsequent collisions
	} else if entityMaxX >= pe.MaxX {
		// Hit right wall
		newX := pe.MaxX - size/2
		entity.SetImmediatePosition(newX, y) // Immediate position for crisp bounce
		entity.SetVelocity(-vx*pe.Restitution, vy)
		x = newX // Update position variable for subsequent collisions
	}

	// Vertical boundary collisions
	if entityMinY <= pe.MinY {
		// Hit top wall
		newY := pe.MinY + size/2
		entity.SetImmediatePosition(x, newY) // Use updated x position
		entity.SetVelocity(vx, -vy*pe.Restitution)
	} else if entityMaxY >= pe.MaxY {
		// Hit bottom wall
		newY := pe.MaxY - size/2
		entity.SetImmediatePosition(x, newY) // Immediate position for crisp bounce
		entity.SetVelocity(vx, -vy*pe.Restitution)
	}
}

// capVelocity ensures velocities don't become too extreme
func (pe *PhysicsEngine) capVelocity(entity Entity) {
	vx, vy := entity.GetVelocity()

	// Cap maximum velocity
	if math.Abs(vx) > pe.MaxVelocity {
		vx = math.Copysign(pe.MaxVelocity, vx)
	}
	if math.Abs(vy) > pe.MaxVelocity {
		vy = math.Copysign(pe.MaxVelocity, vy)
	}

	// Apply static friction when entities are nearly at rest
	totalSpeed := math.Sqrt(vx*vx + vy*vy)
	if totalSpeed < pe.MinVelocity*3 {
		// Apply static friction
		frictionFactor := pe.StaticFriction
		vx *= frictionFactor
		vy *= frictionFactor
	}

	// Stop very slow movement
	if math.Abs(vx) < pe.MinVelocity {
		vx = 0
	}
	if math.Abs(vy) < pe.MinVelocity {
		vy = 0
	}

	entity.SetVelocity(vx, vy)
}

// HandleEntityCollisions processes collisions between entities
func (pe *PhysicsEngine) HandleEntityCollisions(entities []Entity) {
	// Get all collisions
	collisions := pe.findCollisions(entities)

	// Resolve each collision
	for _, collision := range collisions {
		pe.resolveCollision(collision.Entity1, collision.Entity2)
	}
}

// findCollisions detects all entity-to-entity collisions
func (pe *PhysicsEngine) findCollisions(entities []Entity) []CollisionPair {
	var collisions []CollisionPair

	for i := 0; i < len(entities); i++ {
		for j := i + 1; j < len(entities); j++ {
			if pe.checkEntityCollision(entities[i], entities[j]) {
				collisions = append(collisions, CollisionPair{
					Entity1: entities[i],
					Entity2: entities[j],
				})
			}
		}
	}

	return collisions
}

// checkEntityCollision checks if two entities are colliding
func (pe *PhysicsEngine) checkEntityCollision(e1, e2 Entity) bool {
	x1, y1 := e1.GetPosition()
	x2, y2 := e2.GetPosition()

	// Use GetBounds to get the actual collision sizes
	_, _, w1, _ := e1.GetBounds()
	_, _, w2, _ := e2.GetBounds()

	// Calculate effective radii from bounds
	radius1 := w1 / 2
	radius2 := w2 / 2

	// Calculate distance between centers
	dx := x2 - x1
	dy := y2 - y1
	distance := math.Sqrt(dx*dx + dy*dy)

	// Check if distance is less than sum of radii plus contact tolerance
	// This allows entities to touch more closely
	minDistance := (radius1 + radius2) - pe.ContactTolerance
	return distance < minDistance
}

// resolveCollision handles elastic collision between two entities
func (pe *PhysicsEngine) resolveCollision(e1, e2 Entity) {
	x1, y1 := e1.GetPosition()
	x2, y2 := e2.GetPosition()
	vx1, vy1 := e1.GetVelocity()
	vx2, vy2 := e2.GetVelocity()

	// Calculate collision normal
	dx := x2 - x1
	dy := y2 - y1
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance == 0 {
		// Entities are exactly on top of each other - separate them
		dx = 0.1 * (rand.Float64() - 0.5) // Small random separation
		dy = 0.1 * (rand.Float64() - 0.5)
		distance = math.Sqrt(dx*dx + dy*dy)
	}

	// Normalize collision vector
	nx := dx / distance
	ny := dy / distance

	// Separate entities if they're overlapping
	// Use GetBounds to get the actual collision sizes
	_, _, w1, _ := e1.GetBounds()
	_, _, w2, _ := e2.GetBounds()

	// Calculate effective radii from bounds
	radius1 := w1 / 2
	radius2 := w2 / 2

	minDistance := (radius1 + radius2) - pe.ContactTolerance
	overlap := minDistance - distance

	if overlap > 0 {
		// Move entities apart more gently to prevent energy injection
		separationFactor := 0.5 // Only separate by half the overlap to allow closer contact
		separationX := nx * overlap * separationFactor
		separationY := ny * overlap * separationFactor

		e1.SetImmediatePosition(x1-separationX, y1-separationY)
		e2.SetImmediatePosition(x2+separationX, y2+separationY)
	}

	// Calculate relative velocity in collision normal direction
	dvx := vx2 - vx1
	dvy := vy2 - vy1
	dvn := dvx*nx + dvy*ny

	// Do not resolve if velocities are separating
	if dvn > 0 {
		return
	}

	// Apply contact damping for entities that are barely moving
	relativeSpeed := math.Sqrt(dvx*dvx + dvy*dvy)
	if relativeSpeed < pe.MinVelocity*2 {
		// Apply strong damping when entities are moving slowly
		dampingFactor := pe.ContactDamping
		e1.SetVelocity(vx1*dampingFactor, vy1*dampingFactor)
		e2.SetVelocity(vx2*dampingFactor, vy2*dampingFactor)

		// If both entities are nearly at rest, stop them completely
		if relativeSpeed < pe.MinVelocity {
			e1.SetVelocity(0, 0)
			e2.SetVelocity(0, 0)
			return
		}
	}

	// Calculate collision impulse (simplified, assuming equal mass)
	impulse := 2 * dvn / 2 // Divided by 2 for equal mass distribution
	impulse *= pe.Restitution

	// Apply additional energy dissipation for more realistic settling
	energyLoss := 0.95 // Lose 5% energy on each collision
	impulse *= energyLoss

	// Apply impulse to velocities
	e1.SetVelocity(vx1+impulse*nx, vy1+impulse*ny)
	e2.SetVelocity(vx2-impulse*nx, vy2-impulse*ny)
}

// AddRandomVelocity adds some initial random velocity to an entity
func (pe *PhysicsEngine) AddRandomVelocity(entity Entity, maxVelocity float64) {
	// Add small random velocity for more interesting simulation
	vx := (rand.Float64() - 0.5) * maxVelocity
	vy := (rand.Float64() - 0.5) * maxVelocity

	currentVX, currentVY := entity.GetVelocity()
	entity.SetVelocity(currentVX+vx, currentVY+vy)
}

// SetGravity allows dynamic gravity adjustment
func (pe *PhysicsEngine) SetGravity(gravity float64) {
	// Validate input - reject infinite and NaN values
	if math.IsInf(gravity, 0) || math.IsNaN(gravity) {
		return // Reject invalid values
	}
	pe.Gravity = gravity
}

// GetGravity returns current gravity setting
func (pe *PhysicsEngine) GetGravity() float64 {
	return pe.Gravity
}

// SetRestitution allows dynamic bounce adjustment
func (pe *PhysicsEngine) SetRestitution(restitution float64) {
	if restitution >= 0 && restitution <= 1 {
		pe.Restitution = restitution
	}
}

// GetRestitution returns current bounce setting
func (pe *PhysicsEngine) GetRestitution() float64 {
	return pe.Restitution
}

// Pause stops physics calculations (sets deltaTime to 0)
func (pe *PhysicsEngine) Pause() {
	pe.DeltaTime = 0
}

// Resume restarts physics calculations
func (pe *PhysicsEngine) Resume() {
	pe.DeltaTime = 0.1
}

// IsRunning checks if physics is currently active
func (pe *PhysicsEngine) IsRunning() bool {
	return pe.DeltaTime > 0
}
