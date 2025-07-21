package main

import (
	"math"
	"time"

	"github.com/charmbracelet/harmonica"
)

// AnimationEngine handles smooth animations using Harmonica springs
type AnimationEngine struct {
	// Animation settings
	SpringTension float64 // Spring stiffness (higher = faster convergence)
	SpringDamping float64 // Spring damping (higher = less oscillation)

	// Frame rate targeting
	TargetFPS     int
	LastFrameTime time.Time
	FrameDelta    time.Duration
}

// EntityAnimationState holds animation state for each entity
type EntityAnimationState struct {
	// Current visual position (what's displayed)
	DisplayX, DisplayY float64

	// Target position (from physics)
	TargetX, TargetY float64

	// Velocity for spring animation
	VelocityX, VelocityY float64

	// Spring animations for smooth movement
	SpringX harmonica.Spring
	SpringY harmonica.Spring

	// Animation tracking
	IsAnimating bool
	LastUpdate  time.Time
}

// NewAnimationEngine creates a new animation engine
func NewAnimationEngine() *AnimationEngine {
	return &AnimationEngine{
		SpringTension: 300.0, // Responsive but not too bouncy
		SpringDamping: 30.0,  // Well-damped
		TargetFPS:     60,    // 60 FPS for smooth animation
		LastFrameTime: time.Now(),
		FrameDelta:    time.Millisecond * 16, // ~60 FPS (16ms per frame)
	}
}

// NewEntityAnimationState creates animation state for an entity
func (ae *AnimationEngine) NewEntityAnimationState(x, y float64) *EntityAnimationState {
	return &EntityAnimationState{
		DisplayX:   x,
		DisplayY:   y,
		TargetX:    x,
		TargetY:    y,
		VelocityX:  0,
		VelocityY:  0,
		SpringX:    harmonica.NewSpring(harmonica.FPS(ae.TargetFPS), ae.SpringTension, ae.SpringDamping),
		SpringY:    harmonica.NewSpring(harmonica.FPS(ae.TargetFPS), ae.SpringTension, ae.SpringDamping),
		LastUpdate: time.Now(),
	}
}

// SetTarget updates the target position for smooth animation
func (eas *EntityAnimationState) SetTarget(x, y float64) {
	// Validate and sanitize inputs - reject NaN and infinite values
	if math.IsNaN(x) || math.IsInf(x, 0) {
		x = eas.DisplayX // Keep current position if invalid
	}
	if math.IsNaN(y) || math.IsInf(y, 0) {
		y = eas.DisplayY // Keep current position if invalid
	}

	eas.TargetX = x
	eas.TargetY = y
	eas.IsAnimating = true
}

// UpdateAnimation advances the spring animation
func (ae *AnimationEngine) UpdateAnimation(eas *EntityAnimationState) {
	now := time.Now()
	eas.LastUpdate = now

	// Update spring animations toward target positions
	// Harmonica Update(position, velocity, target) returns new position and velocity
	newX, newVX := eas.SpringX.Update(eas.DisplayX, eas.VelocityX, eas.TargetX)
	newY, newVY := eas.SpringY.Update(eas.DisplayY, eas.VelocityY, eas.TargetY)

	// Update display positions and velocities
	eas.DisplayX = newX
	eas.DisplayY = newY
	eas.VelocityX = newVX
	eas.VelocityY = newVY

	// Check if animation is essentially complete
	toleranceX := 0.01
	toleranceY := 0.01
	velocityThreshold := 0.01
	if abs(eas.DisplayX-eas.TargetX) < toleranceX && abs(eas.DisplayY-eas.TargetY) < toleranceY &&
		abs(eas.VelocityX) < velocityThreshold && abs(eas.VelocityY) < velocityThreshold {
		eas.IsAnimating = false
	}
}

// GetDisplayPosition returns the current animated position
func (eas *EntityAnimationState) GetDisplayPosition() (float64, float64) {
	return eas.DisplayX, eas.DisplayY
}

// GetTarget returns the target position
func (eas *EntityAnimationState) GetTarget() (float64, float64) {
	return eas.TargetX, eas.TargetY
}

// SetInitialPosition sets both display and target to the same position (no animation)
func (eas *EntityAnimationState) SetInitialPosition(x, y float64) {
	eas.DisplayX = x
	eas.DisplayY = y
	eas.TargetX = x
	eas.TargetY = y
	eas.VelocityX = 0
	eas.VelocityY = 0
	eas.IsAnimating = false
}

// IsStillAnimating returns whether the entity is still in motion
func (eas *EntityAnimationState) IsStillAnimating() bool {
	return eas.IsAnimating
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
