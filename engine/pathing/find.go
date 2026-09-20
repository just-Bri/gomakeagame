package pathing

import "github.com/just-Bri/gomakeagame/engine/grid"

// FindPath returns a shortest 4-connected walkable path from (sx,sy) to (tx,ty).
// The path includes both the start and goal cells. Returns nil if unreachable
// or if either endpoint is not walkable.
//
// Useful for player unit move orders; creep traffic still uses the flow field.
func FindPath(g *grid.Grid, sx, sy, tx, ty int) []Cell {
	if g == nil || !g.Walkable(sx, sy) || !g.Walkable(tx, ty) {
		return nil
	}
	if sx == tx && sy == ty {
		return []Cell{{X: sx, Y: sy}}
	}

	w, h := g.Width, g.Height
	n := w * h
	visited := make([]bool, n)
	parent := make([]int, n)
	for i := range parent {
		parent[i] = -1
	}
	queue := make([]int, 0, n)
	start := sy*w + sx
	visited[start] = true
	queue = append(queue, start)

	found := false
	goal := ty*w + tx
	var neigh [][2]int

	for head := 0; head < len(queue); head++ {
		ci := queue[head]
		if ci == goal {
			found = true
			break
		}
		cx := ci % w
		cy := ci / w
		neigh = g.Neighbors4(cx, cy, neigh[:0])
		for _, nb := range neigh {
			nx, ny := nb[0], nb[1]
			if !g.Walkable(nx, ny) {
				continue
			}
			ni := ny*w + nx
			if visited[ni] {
				continue
			}
			visited[ni] = true
			parent[ni] = ci
			queue = append(queue, ni)
		}
	}
	if !found {
		return nil
	}

	var rev []Cell
	for ci := goal; ci >= 0; ci = parent[ci] {
		rev = append(rev, Cell{X: ci % w, Y: ci / w})
		if ci == start {
			break
		}
	}
	for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
		rev[i], rev[j] = rev[j], rev[i]
	}
	return rev
}
