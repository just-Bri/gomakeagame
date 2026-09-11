package move

import "github.com/just-Bri/gomakeagame/engine/ecs"

// Velocity is world-space linear velocity (units per second).
type Velocity struct {
	VX, VY float64
}

// Bounds is an axis-aligned rectangle (top-left + size).
type Bounds struct {
	X, Y, W, H float64
}

// Integrate advances every entity that has both Transform and Velocity by dt.
func Integrate(transforms *ecs.Store[Transform], velocities *ecs.Store[Velocity], dt float64) {
	if transforms == nil || velocities == nil || dt == 0 {
		return
	}
	ecs.Join2(transforms, velocities, func(_ ecs.Entity, tr *Transform, vel *Velocity) {
		tr.X += vel.VX * dt
		tr.Y += vel.VY * dt
	})
}

// ClampRect keeps the axis-aligned box (tr.X, tr.Y, width, height) inside bounds.
func ClampRect(tr *Transform, width, height float64, bounds Bounds) {
	if tr == nil {
		return
	}
	maxX := bounds.X + bounds.W - width
	maxY := bounds.Y + bounds.H - height
	if tr.X < bounds.X {
		tr.X = bounds.X
	} else if tr.X > maxX {
		tr.X = maxX
	}
	if tr.Y < bounds.Y {
		tr.Y = bounds.Y
	} else if tr.Y > maxY {
		tr.Y = maxY
	}
}
