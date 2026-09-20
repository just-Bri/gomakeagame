package pathing

// Trace follows Field.Next from (sx, sy) until an exit (or a stuck cell).
// The returned slice includes the start cell and the exit. An unreachable
// start yields nil. Used by maze UIs to paint the forced corridor.
func Trace(f *Field, sx, sy int) []Cell {
	if f == nil || !f.Reachable(sx, sy) {
		return nil
	}

	maxSteps := f.Width*f.Height + 1
	out := make([]Cell, 0, f.Dist(sx, sy)+1)
	x, y := sx, sy
	for step := 0; step < maxSteps; step++ {
		out = append(out, Cell{X: x, Y: y})
		nx, ny, ok := f.Next(x, y)
		if !ok {
			break
		}
		if nx == x && ny == y {
			break // exit
		}
		x, y = nx, ny
	}
	return out
}

// CorridorFrom returns the union of Trace paths from every start cell,
// preserving first-seen order (starts in list order, then downhill).
// Unreachable starts are skipped. Nil field or empty starts yield nil.
func CorridorFrom(f *Field, starts []Cell) []Cell {
	if f == nil || len(starts) == 0 {
		return nil
	}
	seen := make([]bool, f.Width*f.Height)
	out := make([]Cell, 0, f.Width+f.Height)
	for _, s := range starts {
		for _, c := range Trace(f, s.X, s.Y) {
			i := c.Y*f.Width + c.X
			if i < 0 || i >= len(seen) || seen[i] {
				continue
			}
			seen[i] = true
			out = append(out, c)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// LongestDist returns the maximum Dist among the given cells that are
// reachable on f, or -1 if none are reachable. Useful for “path length”
// HUD readouts from a spawn strip.
func LongestDist(f *Field, cells []Cell) int {
	if f == nil || len(cells) == 0 {
		return unreachable
	}
	best := unreachable
	for _, c := range cells {
		d := f.Dist(c.X, c.Y)
		if d > best {
			best = d
		}
	}
	return best
}
