package grid

import "fmt"

// Rect is an axis-aligned footprint in cell coordinates (top-left + size).
// Buildings larger than 1×1 use this to claim walkability.
type Rect struct {
	X, Y int
	W, H int
}

// NewRect returns a footprint. w and h must be positive.
func NewRect(x, y, w, h int) (Rect, error) {
	if w <= 0 || h <= 0 {
		return Rect{}, fmt.Errorf("grid: invalid rect size %dx%d", w, h)
	}
	return Rect{X: x, Y: y, W: w, H: h}, nil
}

// MaxX returns the exclusive right edge (X + W).
func (r Rect) MaxX() int { return r.X + r.W }

// MaxY returns the exclusive bottom edge (Y + H).
func (r Rect) MaxY() int { return r.Y + r.H }

// Contains reports whether cell (x, y) lies inside r.
func (r Rect) Contains(x, y int) bool {
	return x >= r.X && y >= r.Y && x < r.MaxX() && y < r.MaxY()
}

// RectInBounds reports whether every cell of r lies on the grid.
func (g *Grid) RectInBounds(r Rect) bool {
	if r.W <= 0 || r.H <= 0 {
		return false
	}
	return g.InBounds(r.X, r.Y) && g.InBounds(r.MaxX()-1, r.MaxY()-1)
}

// RectClear reports whether every cell of r is currently walkable.
// Out-of-bounds rects are not clear.
func (g *Grid) RectClear(r Rect) bool {
	if !g.RectInBounds(r) {
		return false
	}
	for y := r.Y; y < r.MaxY(); y++ {
		for x := r.X; x < r.MaxX(); x++ {
			if !g.Walkable(x, y) {
				return false
			}
		}
	}
	return true
}

// SetRectBlocked marks every cell in r blocked or clear.
// Out-of-bounds cells are skipped.
func (g *Grid) SetRectBlocked(r Rect, blocked bool) {
	for y := r.Y; y < r.MaxY(); y++ {
		for x := r.X; x < r.MaxX(); x++ {
			g.SetBlocked(x, y, blocked)
		}
	}
}

// EachCell calls fn for every cell in r that is in bounds.
func (g *Grid) EachCell(r Rect, fn func(x, y int)) {
	for y := r.Y; y < r.MaxY(); y++ {
		for x := r.X; x < r.MaxX(); x++ {
			if g.InBounds(x, y) {
				fn(x, y)
			}
		}
	}
}
