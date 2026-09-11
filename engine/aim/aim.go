// Package aim provides ballistic helpers shared by games (lead prediction, etc.).
package aim

import "math"

// LeadPoint predicts where to shoot so a bullet of speed S from (ox,oy)
// meets a target at (px,py) moving with velocity (vx,vy), assuming constant V.
//
// If no positive intercept exists, it returns the target's current position.
func LeadPoint(ox, oy, px, py, vx, vy, bulletSpeed float64) (aimX, aimY float64) {
	if bulletSpeed <= 0 {
		return px, py
	}

	dx := px - ox
	dy := py - oy
	a := vx*vx + vy*vy - bulletSpeed*bulletSpeed
	b := 2 * (dx*vx + dy*vy)
	c := dx*dx + dy*dy

	const eps = 1e-6
	if math.Abs(a) < eps {
		if math.Abs(b) < eps {
			return px, py
		}
		t := -c / b
		if t > 0 {
			return px + vx*t, py + vy*t
		}
		return px, py
	}

	disc := b*b - 4*a*c
	if disc < 0 {
		return px, py
	}
	root := math.Sqrt(disc)
	t1 := (-b - root) / (2 * a)
	t2 := (-b + root) / (2 * a)

	t := math.Inf(1)
	if t1 > eps && t1 < t {
		t = t1
	}
	if t2 > eps && t2 < t {
		t = t2
	}
	if math.IsInf(t, 1) {
		return px, py
	}
	return px + vx*t, py + vy*t
}
