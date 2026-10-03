package move

// AABB is an inclusive-ish axis-aligned bounds on the XY plane (min/max).
type AABB struct {
	MinX, MaxX float64
	MinY, MaxY float64
}

// AABBXZ is ground-plane bounds (Y up worlds).
type AABBXZ struct {
	MinX, MaxX float64
	MinZ, MaxZ float64
}

// ClampPoint keeps (x,y) inside b.
func ClampPoint(x, y float64, b AABB) (float64, float64) {
	if x < b.MinX {
		x = b.MinX
	}
	if x > b.MaxX {
		x = b.MaxX
	}
	if y < b.MinY {
		y = b.MinY
	}
	if y > b.MaxY {
		y = b.MaxY
	}
	return x, y
}

// ClampPointXZ keeps (x,z) inside b.
func ClampPointXZ(x, z float64, b AABBXZ) (float64, float64) {
	if x < b.MinX {
		x = b.MinX
	}
	if x > b.MaxX {
		x = b.MaxX
	}
	if z < b.MinZ {
		z = b.MinZ
	}
	if z > b.MaxZ {
		z = b.MaxZ
	}
	return x, z
}

// ClampTransform3 clamps a ground pose into b.
func ClampTransform3(t *Transform3, b AABBXZ) {
	if t == nil {
		return
	}
	t.X, t.Z = ClampPointXZ(t.X, t.Z, b)
}
