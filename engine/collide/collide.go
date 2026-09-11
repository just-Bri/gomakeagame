// Package collide provides 2D swept / proximity hit tests.
package collide

import "math"

// DistPointSegment returns distance from P to the closest point on segment AB.
func DistPointSegment(px, py, ax, ay, bx, by float64) float64 {
	abx := bx - ax
	aby := by - ay
	len2 := abx*abx + aby*aby
	if len2 < 1e-12 {
		return math.Hypot(px-ax, py-ay)
	}
	t := ((px-ax)*abx + (py-ay)*aby) / len2
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	qx := ax + abx*t
	qy := ay + aby*t
	return math.Hypot(px-qx, py-qy)
}

// SegmentHitsCircle reports whether segment AB comes within radius of (cx,cy).
func SegmentHitsCircle(ax, ay, bx, by, cx, cy, radius float64) bool {
	if radius < 0 {
		return false
	}
	return DistPointSegment(cx, cy, ax, ay, bx, by) <= radius
}

// ClosestHitAlongSegment finds the candidate whose center is closest to
// segment AB among those within hitRadius. each should call yield(x, y, id)
// for every candidate. Returns the winning id.
func ClosestHitAlongSegment(ax, ay, bx, by, hitRadius float64, each func(yield func(x, y float64, id uint64))) (id uint64, ok bool) {
	best := math.MaxFloat64
	each(func(x, y float64, cand uint64) {
		d := DistPointSegment(x, y, ax, ay, bx, by)
		if d <= hitRadius && d < best {
			best = d
			id = cand
			ok = true
		}
	})
	return id, ok
}

// CircleHitsCircle reports whether two disks overlap (edges touching counts).
func CircleHitsCircle(ax, ay, ar, bx, by, br float64) bool {
	if ar < 0 || br < 0 {
		return false
	}
	dx := ax - bx
	dy := ay - by
	r := ar + br
	return dx*dx+dy*dy <= r*r
}

// CircleHitsRect reports whether a disk overlaps an axis-aligned rectangle
// defined by top-left (rx,ry) and size (rw,rh).
func CircleHitsRect(cx, cy, radius, rx, ry, rw, rh float64) bool {
	if radius < 0 || rw < 0 || rh < 0 {
		return false
	}
	// Closest point on the rect to the circle center.
	closestX := cx
	if closestX < rx {
		closestX = rx
	} else if closestX > rx+rw {
		closestX = rx + rw
	}
	closestY := cy
	if closestY < ry {
		closestY = ry
	} else if closestY > ry+rh {
		closestY = ry + rh
	}
	dx := cx - closestX
	dy := cy - closestY
	return dx*dx+dy*dy <= radius*radius
}
