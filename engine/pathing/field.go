// Package pathing builds shared flow fields over a grid.Grid.
//
// Rebuild once when walkability changes (build-phase end, tower place/sell).
// Units sample Direction or Next each tick instead of running A* per agent.
package pathing

import (
	"fmt"
	"math"

	"github.com/just-Bri/gomakeagame/engine/grid"
)

// unreachable marks cells with no path to the exit.
const unreachable = -1

// Field is a reverse-BFS flow field toward one exit cell.
//
// For each reachable walkable tile, Next points one step closer to the exit.
// The exit cell points at itself. Unreachable cells have no next step.
type Field struct {
	Width  int
	Height int
	ExitX  int
	ExitY  int

	nextX []int16
	nextY []int16
	dist  []int32 // steps to exit; unreachable = -1
}

// Rebuild computes a flow field from every walkable cell toward (exitX, exitY).
//
// Cost is O(width*height). Call when the map becomes dirty — not every frame.
func Rebuild(g *grid.Grid, exitX, exitY int) (*Field, error) {
	if g == nil {
		return nil, fmt.Errorf("pathing: nil grid")
	}
	if !g.InBounds(exitX, exitY) {
		return nil, fmt.Errorf("pathing: exit (%d,%d) out of bounds", exitX, exitY)
	}
	if !g.Walkable(exitX, exitY) {
		return nil, fmt.Errorf("pathing: exit (%d,%d) is not walkable", exitX, exitY)
	}

	n := g.Width * g.Height
	f := &Field{
		Width:  g.Width,
		Height: g.Height,
		ExitX:  exitX,
		ExitY:  exitY,
		nextX:  make([]int16, n),
		nextY:  make([]int16, n),
		dist:   make([]int32, n),
	}
	for i := range f.dist {
		f.dist[i] = unreachable
		f.nextX[i] = -1
		f.nextY[i] = -1
	}

	ei := exitY*g.Width + exitX
	f.dist[ei] = 0
	f.nextX[ei] = int16(exitX)
	f.nextY[ei] = int16(exitY)

	// BFS from exit outward; each newly reached cell steps toward its parent.
	queue := make([]int, 0, n)
	queue = append(queue, ei)

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
			f.nextX[ni] = int16(cx)
			f.nextY[ni] = int16(cy)
			queue = append(queue, ni)
		}
	}

	return f, nil
}

func (f *Field) index(x, y int) int {
	return y*f.Width + x
}

func (f *Field) inBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < f.Width && y < f.Height
}

// Reachable reports whether (x, y) has a path to the exit.
func (f *Field) Reachable(x, y int) bool {
	if !f.inBounds(x, y) {
		return false
	}
	return f.dist[f.index(x, y)] != unreachable
}

// Dist returns steps from (x, y) to the exit, or -1 if unreachable.
func (f *Field) Dist(x, y int) int {
	if !f.inBounds(x, y) {
		return unreachable
	}
	return int(f.dist[f.index(x, y)])
}

// Next returns the neighboring cell one step closer to the exit.
// ok is false when (x, y) is out of bounds or unreachable.
// At the exit, Next returns the exit itself.
func (f *Field) Next(x, y int) (nx, ny int, ok bool) {
	if !f.Reachable(x, y) {
		return 0, 0, false
	}
	i := f.index(x, y)
	return int(f.nextX[i]), int(f.nextY[i]), true
}

// Direction returns a unit vector from cell (x, y) toward its next step.
// At the exit, or when unreachable, ok is false and dx, dy are 0.
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
