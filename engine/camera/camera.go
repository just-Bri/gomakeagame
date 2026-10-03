// Package camera provides simple follow cameras for 2D and XZ ground games.
package camera

import (
	"math"

	"github.com/just-Bri/gomakeagame/engine/project"
)

// Follow2D smoothly tracks a point on the XY plane (top-down / side view).
type Follow2D struct {
	X, Y  float64
	Speed float64 // lerp rate; 0 defaults to 4
}

// Update eases the camera toward (tx, ty).
func (c *Follow2D) Update(dt, tx, ty float64) {
	if c == nil {
		return
	}
	sp := c.Speed
	if sp <= 0 {
		sp = 4
	}
	t := 1 - math.Exp(-sp*dt)
	c.X += (tx - c.X) * t
	c.Y += (ty - c.Y) * t
}

// FollowXZ smoothly tracks a ground-plane point for isometric / soft-3D views.
type FollowXZ struct {
	X, Z  float64
	Speed float64
}

// Update eases the camera toward (tx, tz).
func (c *FollowXZ) Update(dt, tx, tz float64) {
	if c == nil {
		return
	}
	sp := c.Speed
	if sp <= 0 {
		sp = 4
	}
	t := 1 - math.Exp(-sp*dt)
	c.X += (tx - c.X) * t
	c.Z += (tz - c.Z) * t
}

// ApplyIso writes CamX/CamZ onto an isometric projector.
func (c FollowXZ) ApplyIso(iso *project.Iso) {
	if iso == nil {
		return
	}
	iso.CamX, iso.CamZ = c.X, c.Z
}

// Snap jumps immediately to the target (call on load / teleport).
func (c *FollowXZ) Snap(tx, tz float64) {
	if c == nil {
		return
	}
	c.X, c.Z = tx, tz
}

// Snap2D jumps a Follow2D immediately to the target.
func (c *Follow2D) Snap(tx, ty float64) {
	if c == nil {
		return
	}
	c.X, c.Y = tx, ty
}
