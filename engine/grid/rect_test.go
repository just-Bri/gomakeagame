package grid_test

import (
	"testing"

	"github.com/just-Bri/gomakeagame/engine/grid"
)

func TestRectClearAndBlock(t *testing.T) {
	g, err := grid.New(6, 6)
	if err != nil {
		t.Fatal(err)
	}
	fp, err := grid.NewRect(2, 2, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if !g.RectClear(fp) {
		t.Fatal("fresh rect should be clear")
	}
	g.SetRectBlocked(fp, true)
	if g.RectClear(fp) {
		t.Fatal("blocked rect should not be clear")
	}
	if g.Walkable(2, 2) || g.Walkable(3, 3) {
		t.Fatal("footprint cells should be blocked")
	}
	if !g.Walkable(1, 2) {
		t.Fatal("outside footprint should stay walkable")
	}
	g.SetRectBlocked(fp, false)
	if !g.RectClear(fp) {
		t.Fatal("cleared rect should be walkable again")
	}
}

func TestRectOutOfBounds(t *testing.T) {
	g, err := grid.New(4, 4)
	if err != nil {
		t.Fatal(err)
	}
	fp, _ := grid.NewRect(3, 3, 2, 2)
	if g.RectInBounds(fp) {
		t.Fatal("overflow rect should be OOB")
	}
	if g.RectClear(fp) {
		t.Fatal("OOB rect must not be clear")
	}
}
