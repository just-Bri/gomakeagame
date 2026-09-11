// Package grid is a 2D tile map for walkability and world↔cell conversion.
//
// The grid is shared map state — not ECS entities. Pathing and movement
// sample it; towers mark cells blocked when placed or sold.
package grid

import (
	"fmt"
	"math"
)

// Grid is a rectangular tile map. Cells are indexed row-major: y*Width + x.
type Grid struct {
	Width  int
	Height int
	// walkable[i] is true when units may enter the cell.
	walkable []bool
}

// New creates a Width×Height grid. Every cell starts walkable.
func New(width, height int) (*Grid, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("grid: invalid size %dx%d", width, height)
	}
	n := width * height
	w := make([]bool, n)
	for i := range w {
		w[i] = true
	}
	return &Grid{Width: width, Height: height, walkable: w}, nil
}

// Index returns the dense index for (x, y). Panics if out of bounds.
func (g *Grid) Index(x, y int) int {
	return y*g.Width + x
}

// InBounds reports whether (x, y) lies inside the grid.
func (g *Grid) InBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < g.Width && y < g.Height
}

// Walkable reports whether units may enter (x, y).
// Out-of-bounds cells are never walkable.
func (g *Grid) Walkable(x, y int) bool {
	if !g.InBounds(x, y) {
		return false
	}
	return g.walkable[g.Index(x, y)]
}

// SetWalkable sets whether (x, y) may be entered.
// Out-of-bounds is a no-op.
func (g *Grid) SetWalkable(x, y int, walkable bool) {
	if !g.InBounds(x, y) {
		return
	}
	g.walkable[g.Index(x, y)] = walkable
}

// SetBlocked marks (x, y) non-walkable when blocked is true.
func (g *Grid) SetBlocked(x, y int, blocked bool) {
	g.SetWalkable(x, y, !blocked)
}

// FillWalkable sets every cell to the same walkability.
func (g *Grid) FillWalkable(walkable bool) {
	for i := range g.walkable {
		g.walkable[i] = walkable
	}
}

// WorldToCell converts a world position to a cell using floor division.
func (g *Grid) WorldToCell(wx, wy, tileSize float64) (cx, cy int) {
	return int(math.Floor(wx / tileSize)), int(math.Floor(wy / tileSize))
}

// CellCenter returns the world-space center of cell (cx, cy).
func (g *Grid) CellCenter(cx, cy int, tileSize float64) (wx, wy float64) {
	return (float64(cx) + 0.5) * tileSize, (float64(cy) + 0.5) * tileSize
}

// Neighbors4 appends the four cardinal neighbors of (x, y) that are in bounds.
// Order: N, E, S, W.
func (g *Grid) Neighbors4(x, y int, dst [][2]int) [][2]int {
	offs := [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	for _, o := range offs {
		nx, ny := x+o[0], y+o[1]
		if g.InBounds(nx, ny) {
			dst = append(dst, [2]int{nx, ny})
		}
	}
	return dst
}
