package main

import (
	"testing"
	"time"
)

func TestNewAnimationEngine(t *testing.T) {
	ae := NewAnimationEngine()

	if ae == nil {
		t.Fatal("NewAnimationEngine returned nil")
	}

	if ae.TargetFPS != 60 {
		t.Errorf("Expected TargetFPS 60, got %d", ae.TargetFPS)
	}

	if ae.SpringTension != 300.0 {
		t.Errorf("Expected SpringTension 300.0, got %f", ae.SpringTension)
	}

	if ae.SpringDamping != 30.0 {
		t.Errorf("Expected SpringDamping 30.0, got %f", ae.SpringDamping)
	}
}

func TestNewEntityAnimationState(t *testing.T) {
	ae := NewAnimationEngine()
	eas := ae.NewEntityAnimationState(10.0, 20.0)

	if eas == nil {
		t.Fatal("NewEntityAnimationState returned nil")
	}

	if eas.DisplayX != 10.0 || eas.DisplayY != 20.0 {
		t.Errorf("Expected DisplayX=10.0, DisplayY=20.0, got %f, %f", eas.DisplayX, eas.DisplayY)
	}

	if eas.TargetX != 10.0 || eas.TargetY != 20.0 {
		t.Errorf("Expected TargetX=10.0, TargetY=20.0, got %f, %f", eas.TargetX, eas.TargetY)
	}

	if eas.VelocityX != 0.0 || eas.VelocityY != 0.0 {
		t.Errorf("Expected VelocityX=0.0, VelocityY=0.0, got %f, %f", eas.VelocityX, eas.VelocityY)
	}
}

func TestSetTarget(t *testing.T) {
	ae := NewAnimationEngine()
	eas := ae.NewEntityAnimationState(0.0, 0.0)

	eas.SetTarget(50.0, 75.0)

	if eas.TargetX != 50.0 || eas.TargetY != 75.0 {
		t.Errorf("Expected TargetX=50.0, TargetY=75.0, got %f, %f", eas.TargetX, eas.TargetY)
	}

	if !eas.IsAnimating {
		t.Error("Expected IsAnimating to be true after SetTarget")
	}
}

func TestUpdateAnimation(t *testing.T) {
	ae := NewAnimationEngine()
	eas := ae.NewEntityAnimationState(0.0, 0.0)
	eas.SetTarget(10.0, 10.0)

	// Store initial position
	initialX := eas.DisplayX
	initialY := eas.DisplayY

	// Update animation - should move toward target
	ae.UpdateAnimation(eas)

	// Position should have changed (moved toward target)
	if eas.DisplayX == initialX && eas.DisplayY == initialY {
		t.Error("Animation should have moved position toward target")
	}

	// Should be animating toward the target
	if !eas.IsAnimating {
		t.Error("Should still be animating toward target")
	}
}

func TestGetDisplayPosition(t *testing.T) {
	ae := NewAnimationEngine()
	eas := ae.NewEntityAnimationState(15.5, 25.5)

	x, y := eas.GetDisplayPosition()

	if x != 15.5 || y != 25.5 {
		t.Errorf("Expected GetDisplayPosition to return 15.5, 25.5, got %f, %f", x, y)
	}
}

func TestGetTarget(t *testing.T) {
	ae := NewAnimationEngine()
	eas := ae.NewEntityAnimationState(0.0, 0.0)
	eas.SetTarget(100.0, 200.0)

	x, y := eas.GetTarget()

	if x != 100.0 || y != 200.0 {
		t.Errorf("Expected GetTarget to return 100.0, 200.0, got %f, %f", x, y)
	}
}

func TestSetInitialPosition(t *testing.T) {
	ae := NewAnimationEngine()
	eas := ae.NewEntityAnimationState(0.0, 0.0)
	eas.SetTarget(50.0, 50.0) // Start animation

	// Reset to new initial position
	eas.SetInitialPosition(30.0, 40.0)

	if eas.DisplayX != 30.0 || eas.DisplayY != 40.0 {
		t.Errorf("Expected DisplayX=30.0, DisplayY=40.0, got %f, %f", eas.DisplayX, eas.DisplayY)
	}

	if eas.TargetX != 30.0 || eas.TargetY != 40.0 {
		t.Errorf("Expected TargetX=30.0, TargetY=40.0, got %f, %f", eas.TargetX, eas.TargetY)
	}

	if eas.VelocityX != 0.0 || eas.VelocityY != 0.0 {
		t.Errorf("Expected VelocityX=0.0, VelocityY=0.0, got %f, %f", eas.VelocityX, eas.VelocityY)
	}

	if eas.IsAnimating {
		t.Error("Expected IsAnimating to be false after SetInitialPosition")
	}
}

func TestAnimationConvergence(t *testing.T) {
	ae := NewAnimationEngine()
	eas := ae.NewEntityAnimationState(0.0, 0.0)
	eas.SetTarget(1.0, 1.0) // Small target for faster convergence

	// Simulate multiple animation frames
	maxFrames := 1000
	for i := 0; i < maxFrames; i++ {
		ae.UpdateAnimation(eas)

		// Add small delay to simulate real timing
		time.Sleep(time.Microsecond)

		// Check if animation has converged
		if !eas.IsAnimating {
			t.Logf("Animation converged after %d frames", i+1)

			// Verify final position is close to target
			if abs(eas.DisplayX-eas.TargetX) > 0.02 || abs(eas.DisplayY-eas.TargetY) > 0.02 {
				t.Errorf("Animation converged but position too far from target. Display: (%f, %f), Target: (%f, %f)",
					eas.DisplayX, eas.DisplayY, eas.TargetX, eas.TargetY)
			}
			return
		}
	}

	t.Error("Animation did not converge within reasonable time")
}
