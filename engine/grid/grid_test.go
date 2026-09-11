package grid_test

import (
	"testing"

	"github.com/just-Bri/gomakeagame/engine/grid"
)

func TestNewRejectsInvalidSize(t *testing.T) {
	if _, err := grid.New(0, 5); err == nil {
		t.Fatal("expected error for zero width")
	}
	if _, err := grid.New(5, -1); err == nil {
		t.Fatal("expected error for negative height")
	}
}

func TestDefaultsWalkable(t *testing.T) {
	g, err := grid.New(3, 2)
	if err != nil {
		t.Fatal(err)
	}
	for y := 0; y < 2; y++ {
		for x := 0; x < 3; x++ {
			if !g.Walkable(x, y) {
				t.Fatalf("(%d,%d) should start walkable", x, y)
			}
		}
	}
	if g.Walkable(-1, 0) || g.Walkable(3, 0) || g.Walkable(0, 2) {
		t.Fatal("out of bounds must not be walkable")
	}
}

func TestSetBlocked(t *testing.T) {
	g, err := grid.New(4, 4)
	if err != nil {
		t.Fatal(err)
	}
	g.SetBlocked(1, 1, true)
	if g.Walkable(1, 1) {
		t.Fatal("blocked cell should not be walkable")
	}
	g.SetBlocked(1, 1, false)
	if !g.Walkable(1, 1) {
		t.Fatal("unblocked cell should be walkable")
	}
}

func TestWorldCellRoundTrip(t *testing.T) {
	g, err := grid.New(10, 10)
	if err != nil {
		t.Fatal(err)
	}
	const tile = 32.0
	cx, cy := g.WorldToCell(32.5, 64.1, tile)
	if cx != 1 || cy != 2 {
		t.Fatalf("WorldToCell = (%d,%d), want (1,2)", cx, cy)
	}
	wx, wy := g.CellCenter(1, 2, tile)
	if wx != 48 || wy != 80 {
		t.Fatalf("CellCenter = (%v,%v), want (48,80)", wx, wy)
	}
}

func TestNeighbors4(t *testing.T) {
	g, err := grid.New(3, 3)
	if err != nil {
		t.Fatal(err)
	}
	n := g.Neighbors4(1, 1, nil)
	if len(n) != 4 {
		t.Fatalf("center neighbors = %d, want 4", len(n))
	}
	corner := g.Neighbors4(0, 0, nil)
	if len(corner) != 2 {
		t.Fatalf("corner neighbors = %d, want 2", len(corner))
	}
}
