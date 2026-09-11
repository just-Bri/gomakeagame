package occupancy_test

import (
	"testing"

	"github.com/just-Bri/gomakeagame/engine/ecs"
	"github.com/just-Bri/gomakeagame/engine/grid"
	"github.com/just-Bri/gomakeagame/engine/occupancy"
)

func TestSetAtClear(t *testing.T) {
	m := occupancy.New(4, 3)
	e := ecs.Entity(42)
	m.Set(1, 1, e)
	got, ok := m.At(1, 1)
	if !ok || got != e {
		t.Fatalf("At = %v,%v want %v", got, ok, e)
	}
	m.Clear(1, 1)
	if _, ok := m.At(1, 1); ok {
		t.Fatal("expected empty after Clear")
	}
}

func TestSetRect(t *testing.T) {
	m := occupancy.New(5, 5)
	e := ecs.Entity(7)
	fp := grid.Rect{X: 1, Y: 2, W: 2, H: 1}
	m.SetRect(fp, e)
	for _, c := range [][2]int{{1, 2}, {2, 2}} {
		got, ok := m.At(c[0], c[1])
		if !ok || got != e {
			t.Fatalf("At(%d,%d)=%v,%v", c[0], c[1], got, ok)
		}
	}
	if _, ok := m.At(0, 2); ok {
		t.Fatal("outside rect should be empty")
	}
	m.ClearRect(fp, e)
	if _, ok := m.At(1, 2); ok {
		t.Fatal("expected cleared")
	}
}
