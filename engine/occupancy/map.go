// Package occupancy maps grid cells to owning entities (placed buildings, etc.).
package occupancy

import (
	"github.com/just-Bri/gomakeagame/engine/ecs"
	"github.com/just-Bri/gomakeagame/engine/grid"
)

// Map tracks which entity occupies each cell. Empty cells hold ecs.Nil.
type Map struct {
	Width  int
	Height int
	cells  []ecs.Entity
}

// New creates a Width×Height occupancy map (all empty).
func New(width, height int) *Map {
	if width <= 0 || height <= 0 {
		return &Map{}
	}
	return &Map{
		Width:  width,
		Height: height,
		cells:  make([]ecs.Entity, width*height),
	}
}

func (m *Map) index(x, y int) int {
	return y*m.Width + x
}

func (m *Map) inBounds(x, y int) bool {
	return m != nil && x >= 0 && y >= 0 && x < m.Width && y < m.Height && len(m.cells) == m.Width*m.Height
}

// At returns the entity at (x,y), if any.
func (m *Map) At(x, y int) (ecs.Entity, bool) {
	if !m.inBounds(x, y) {
		return ecs.Nil, false
	}
	e := m.cells[m.index(x, y)]
	if e == ecs.Nil {
		return ecs.Nil, false
	}
	return e, true
}

// Set assigns e to a single cell.
func (m *Map) Set(x, y int, e ecs.Entity) {
	if !m.inBounds(x, y) {
		return
	}
	m.cells[m.index(x, y)] = e
}

// Clear empties a single cell.
func (m *Map) Clear(x, y int) {
	m.Set(x, y, ecs.Nil)
}

// SetRect assigns e to every cell in fp.
func (m *Map) SetRect(fp grid.Rect, e ecs.Entity) {
	for y := fp.Y; y < fp.Y+fp.H; y++ {
		for x := fp.X; x < fp.X+fp.W; x++ {
			m.Set(x, y, e)
		}
	}
}

// ClearRect empties every cell in fp that currently holds e (or any if e is Nil).
func (m *Map) ClearRect(fp grid.Rect, e ecs.Entity) {
	for y := fp.Y; y < fp.Y+fp.H; y++ {
		for x := fp.X; x < fp.X+fp.W; x++ {
			if !m.inBounds(x, y) {
				continue
			}
			i := m.index(x, y)
			if e == ecs.Nil || m.cells[i] == e {
				m.cells[i] = ecs.Nil
			}
		}
	}
}
