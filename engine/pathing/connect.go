package pathing

import "github.com/just-Bri/gomakeagame/engine/grid"

// Cell is a grid coordinate used for spawn/exit lists.
type Cell struct {
	X, Y int
}

// CanReach reports whether a walkable path exists from (sx,sy) to (tx,ty)
// on the current grid (4-connected). Does not build a full flow field.
func CanReach(g *grid.Grid, sx, sy, tx, ty int) bool {
	if g == nil || !g.Walkable(sx, sy) || !g.Walkable(tx, ty) {
		return false
	}
	if sx == tx && sy == ty {
		return true
	}

	n := g.Width * g.Height
	visited := make([]bool, n)
	queue := make([]int, 0, n)
	start := sy*g.Width + sx
	visited[start] = true
	queue = append(queue, start)

	var neigh [][2]int
	for head := 0; head < len(queue); head++ {
		ci := queue[head]
		cx := ci % g.Width
		cy := ci / g.Width
		if cx == tx && cy == ty {
			return true
		}
		neigh = g.Neighbors4(cx, cy, neigh[:0])
		for _, nb := range neigh {
			nx, ny := nb[0], nb[1]
			if !g.Walkable(nx, ny) {
				continue
			}
			ni := ny*g.Width + nx
			if visited[ni] {
				continue
			}
			visited[ni] = true
			queue = append(queue, ni)
		}
	}
	return false
}

// SpawnsCanReachExit reports whether every spawn can reach at least one exit.
// Empty spawn or exit lists return false.
func SpawnsCanReachExit(g *grid.Grid, spawns []Cell, exits []Cell) bool {
	if len(spawns) == 0 || len(exits) == 0 {
		return false
	}
	for _, s := range spawns {
		ok := false
		for _, e := range exits {
			if CanReach(g, s.X, s.Y, e.X, e.Y) {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

// AllReachable reports whether every cell in cells is reachable on an
// already-built flow field.
func AllReachable(f *Field, cells []Cell) bool {
	if f == nil || len(cells) == 0 {
		return false
	}
	for _, c := range cells {
		if !f.Reachable(c.X, c.Y) {
			return false
		}
	}
	return true
}
