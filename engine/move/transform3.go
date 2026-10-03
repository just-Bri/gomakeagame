package move

import "math"

// Transform3 is a world-space pose on the XZ ground plane (Y up).
// Yaw is radians: 0 faces +Z, increasing turns toward +X (clockwise from above).
type Transform3 struct {
	X, Y, Z float64
	Yaw     float64
}

// ForwardXZ returns the horizontal facing unit vector.
func (t Transform3) ForwardXZ() (fx, fz float64) {
	return math.Sin(t.Yaw), math.Cos(t.Yaw)
}

// MoveForward advances along facing by distance (can be negative).
func (t *Transform3) MoveForward(dist float64) {
	if t == nil || dist == 0 {
		return
	}
	fx, fz := t.ForwardXZ()
	t.X += fx * dist
	t.Z += fz * dist
}

// MoveStrafe advances along right vector by distance.
func (t *Transform3) MoveStrafe(dist float64) {
	if t == nil || dist == 0 {
		return
	}
	fx, fz := t.ForwardXZ()
	// Right = (fz, -fx) when forward is (fx, fz) on XZ.
	t.X += fz * dist
	t.Z += -fx * dist
}

// FaceToward sets yaw to look at (tx, tz) on the ground plane.
func (t *Transform3) FaceToward(tx, tz float64) {
	if t == nil {
		return
	}
	dx := tx - t.X
	dz := tz - t.Z
	if dx*dx+dz*dz < 1e-12 {
		return
	}
	t.Yaw = math.Atan2(dx, dz)
}

// DistanceXZ returns horizontal distance to another pose.
func (t Transform3) DistanceXZ(o Transform3) float64 {
	dx := t.X - o.X
	dz := t.Z - o.Z
	return math.Hypot(dx, dz)
}
