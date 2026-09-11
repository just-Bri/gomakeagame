package pathing

import (
	"fmt"
	"math"

	"github.com/just-Bri/gomakeagame/engine/grid"
)

// unreachable marks cells with no path to the exit.
const unreachable = -1

// Field is a reverse-BFS flow field toward one or more exit cells.
//
// For each reachable walkable tile, Next points one step closer to an exit.
// Exit cells point at themselves. Unreachable cells have no next step.
//
// When several downhill neighbors share the same Dist, Next prefers moving
// down (increasing Y) before sideways — so creeps drop straight toward a
// bottom exit strip when the maze opens up.
type Field struct {
	Width  int
	Height int
	Exits  []Cell

	nextX []int16
	nextY []int16
	dist  []int32 // steps to nearest exit; unreachable = -1
}

// Rebuild computes a flow field toward a single exit cell.
func Rebuild(g *grid.Grid, exitX, exitY int) (*Field, error) {
	return RebuildFromExits(g, []Cell{{X: exitX, Y: exitY}})
}

// RebuildFromExits computes a flow field toward any of the given exit cells
// (multi-source reverse BFS). Cost is O(width*height).
func RebuildFromExits(g *grid.Grid, exits []Cell) (*Field, error) {
	if g == nil {
		return nil, fmt.Errorf("pathing: nil grid")
	}
	if len(exits) == 0 {
		return nil, fmt.Errorf("pathing: no exit cells")
	}

	n := g.Width * g.Height
	f := &Field{
		Width:  g.Width,
		Height: g.Height,
		Exits:  append([]Cell(nil), exits...),
		nextX:  make([]int16, n),
		nextY:  make([]int16, n),
		dist:   make([]int32, n),
	}
	for i := range f.dist {
		f.dist[i] = unreachable
		f.nextX[i] = -1
		f.nextY[i] = -1
	}

	queue := make([]int, 0, n)
	for _, e := range exits {
		if !g.InBounds(e.X, e.Y) {
			return nil, fmt.Errorf("pathing: exit (%d,%d) out of bounds", e.X, e.Y)
		}
		if !g.Walkable(e.X, e.Y) {
			return nil, fmt.Errorf("pathing: exit (%d,%d) is not walkable", e.X, e.Y)
		}
		ei := e.Y*g.Width + e.X
		if f.dist[ei] == 0 {
			continue // duplicate exit
		}
		f.dist[ei] = 0
		f.nextX[ei] = int16(e.X)
		f.nextY[ei] = int16(e.Y)
		queue = append(queue, ei)
	}

	var neigh [][2]int
	for head := 0; head < len(queue); head++ {
		ci := queue[head]
		cx := ci % g.Width
		cy := ci / g.Width
		neigh = g.Neighbors4(cx, cy, neigh[:0])
		for _, nb := range neigh {
			nx, ny := nb[0], nb[1]
			if !g.Walkable(nx, ny) {
				continue
			}
			ni := ny*g.Width + nx
			if f.dist[ni] != unreachable {
				continue
			}
			f.dist[ni] = f.dist[ci] + 1
			queue = append(queue, ni)
		}
	}

	f.assignNextPreferDown(g)
	return f, nil
}

// assignNextPreferDown sets Next using Dist, preferring downhill (↑Y) steps.
func (f *Field) assignNextPreferDown(g *grid.Grid) {
	var neigh [][2]int
	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			i := f.index(x, y)
			d := f.dist[i]
			if d == unreachable {
				continue
			}
			if d == 0 {
				f.nextX[i] = int16(x)
				f.nextY[i] = int16(y)
				continue
			}

			bestScore := -1
			bx, by := x, y
			neigh = g.Neighbors4(x, y, neigh[:0])
			for _, nb := range neigh {
				nx, ny := nb[0], nb[1]
				if !g.Walkable(nx, ny) {
					continue
				}
				ni := f.index(nx, ny)
				if f.dist[ni] != d-1 {
					continue
				}
				// Prefer down, then sideways, then up.
				score := 1
				if ny > y {
					score = 3
				} else if ny == y {
					score = 2
				}
				if score > bestScore {
					bestScore = score
					bx, by = nx, ny
				}
			}
			if bestScore < 0 {
				continue
			}
			f.nextX[i] = int16(bx)
			f.nextY[i] = int16(by)
		}
	}
}

func (f *Field) index(x, y int) int {
	return y*f.Width + x
}

func (f *Field) inBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < f.Width && y < f.Height
}

// Reachable reports whether (x, y) has a path to an exit.
func (f *Field) Reachable(x, y int) bool {
	if !f.inBounds(x, y) {
		return false
	}
	return f.dist[f.index(x, y)] != unreachable
}

// Dist returns steps from (x, y) to the nearest exit, or -1 if unreachable.
func (f *Field) Dist(x, y int) int {
	if !f.inBounds(x, y) {
		return unreachable
	}
	return int(f.dist[f.index(x, y)])
}

// Next returns the neighboring cell one step closer to an exit.
// ok is false when (x, y) is out of bounds or unreachable.
// At an exit, Next returns the exit itself.
func (f *Field) Next(x, y int) (nx, ny int, ok bool) {
	if !f.Reachable(x, y) {
		return 0, 0, false
	}
	i := f.index(x, y)
	return int(f.nextX[i]), int(f.nextY[i]), true
}

// Direction returns a unit vector from cell (x, y) toward its next step.
// At an exit, or when unreachable, ok is false and dx, dy are 0.
func (f *Field) Direction(x, y int) (dx, dy float64, ok bool) {
	nx, ny, ok := f.Next(x, y)
	if !ok {
		return 0, 0, false
	}
	if nx == x && ny == y {
		return 0, 0, false // already at exit
	}
	dx = float64(nx - x)
	dy = float64(ny - y)
	len := math.Hypot(dx, dy)
	return dx / len, dy / len, true
}
