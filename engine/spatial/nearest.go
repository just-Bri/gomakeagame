// Package spatial provides simple range queries over 2D points.
package spatial

import "math"

// Nearest returns the id of the closest point to (ox,oy) within maxDist.
// each should call yield(x, y, id) for every candidate.
func Nearest(ox, oy, maxDist float64, each func(yield func(x, y float64, id uint64))) (id uint64, ok bool) {
	best := math.MaxFloat64
	each(func(x, y float64, cand uint64) {
		d := math.Hypot(x-ox, y-oy)
		if d <= maxDist && d < best {
			best = d
			id = cand
			ok = true
		}
	})
	return id, ok
}
