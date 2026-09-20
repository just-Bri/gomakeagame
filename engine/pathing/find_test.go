package pathing_test

import (
	"testing"

	"github.com/just-Bri/gomakeagame/engine/grid"
	"github.com/just-Bri/gomakeagame/engine/pathing"
)

func TestFindPathOpen(t *testing.T) {
	g, err := grid.New(5, 5)
	if err != nil {
		t.Fatal(err)
	}
	p := pathing.FindPath(g, 0, 0, 4, 0)
	if len(p) != 5 {
		t.Fatalf("path len=%d want 5: %v", len(p), p)
	}
	if p[0] != (pathing.Cell{X: 0, Y: 0}) || p[len(p)-1] != (pathing.Cell{X: 4, Y: 0}) {
		t.Fatalf("endpoints: %v", p)
	}
}

func TestFindPathAroundWall(t *testing.T) {
	g, err := grid.New(5, 3)
	if err != nil {
		t.Fatal(err)
	}
	// Wall blocking the middle of the top row — must go around.
	g.SetBlocked(1, 0, true)
	g.SetBlocked(2, 0, true)
	g.SetBlocked(3, 0, true)
	p := pathing.FindPath(g, 0, 0, 4, 0)
	if p == nil {
		t.Fatal("expected path around wall")
	}
	for _, c := range p {
		if !g.Walkable(c.X, c.Y) {
			t.Fatalf("path stepped on blocked %v", c)
		}
	}
	if p[0].X != 0 || p[len(p)-1].X != 4 {
		t.Fatalf("bad path %v", p)
	}
}

func TestFindPathSealed(t *testing.T) {
	g, err := grid.New(3, 3)
	if err != nil {
		t.Fatal(err)
	}
	g.SetBlocked(1, 0, true)
	g.SetBlocked(1, 1, true)
	g.SetBlocked(1, 2, true)
	if p := pathing.FindPath(g, 0, 1, 2, 1); p != nil {
		t.Fatalf("expected nil across seal, got %v", p)
	}
}

func TestFindPathSameCell(t *testing.T) {
	g, err := grid.New(2, 2)
	if err != nil {
		t.Fatal(err)
	}
	p := pathing.FindPath(g, 1, 1, 1, 1)
	if len(p) != 1 || p[0] != (pathing.Cell{X: 1, Y: 1}) {
		t.Fatalf("got %v", p)
	}
}
