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
