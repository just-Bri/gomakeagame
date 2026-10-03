package combat

import "math"

// Candidate is a tab-targetable unit in XZ space.
type Candidate struct {
	ID     uint64
	X, Z   float64
	Alive  bool
	Hostile bool
}

// TabNext cycles hostile living targets within maxRange of (ox,oz).
// Order is clockwise by angle from the current target (or from facingYaw if none).
// Returns 0 when no valid target exists.
func TabNext(ox, oz, facingYaw, maxRange float64, current uint64, list []Candidate) uint64 {
	type scored struct {
		id    uint64
		angle float64
	}
	var opts []scored
	maxRange2 := maxRange * maxRange
	for _, c := range list {
		if !c.Alive || !c.Hostile || c.ID == 0 {
			continue
		}
		dx := c.X - ox
		dz := c.Z - oz
		if dx*dx+dz*dz > maxRange2 {
			continue
		}
		ang := math.Atan2(dx, dz)
		opts = append(opts, scored{id: c.ID, angle: ang})
	}
	if len(opts) == 0 {
		return 0
	}

	// Sort by angle ascending (clockwise from -pi..pi via Atan2 order + wrap).
	for i := 0; i < len(opts); i++ {
		for j := i + 1; j < len(opts); j++ {
			if opts[j].angle < opts[i].angle {
				opts[i], opts[j] = opts[j], opts[i]
			}
		}
	}

	if current == 0 {
		// Prefer target closest to facing direction.
		best := 0
		bestDelta := math.MaxFloat64
		for i, o := range opts {
			d := angleDelta(facingYaw, o.angle)
			if d < bestDelta {
				bestDelta = d
				best = i
			}
		}
		return opts[best].id
	}

	idx := -1
	for i, o := range opts {
		if o.id == current {
			idx = i
			break
		}
	}
	if idx < 0 {
		return opts[0].id
	}
	return opts[(idx+1)%len(opts)].id
}

func angleDelta(a, b float64) float64 {
	d := b - a
	for d > math.Pi {
		d -= 2 * math.Pi
	}
	for d < -math.Pi {
		d += 2 * math.Pi
	}
	if d < 0 {
		d = -d
	}
	return d
}

// InMelee reports whether two XZ points are within reach (+ radii).
func InMelee(ax, az, ar, bx, bz, br, reach float64) bool {
	dx := ax - bx
	dz := az - bz
	r := ar + br + reach
	return dx*dx+dz*dz <= r*r
}
