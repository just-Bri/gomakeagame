package pathing_test

import (
	"testing"

	"github.com/just-Bri/gomakeagame/engine/grid"
	"github.com/just-Bri/gomakeagame/engine/pathing"
)

func TestCanReachOpenAndSealed(t *testing.T) {
	g, err := grid.New(5, 5)
	if err != nil {
		t.Fatal(err)
	}
	if !pathing.CanReach(g, 0, 0, 4, 4) {
		t.Fatal("open map should connect")
	}
	for x := 0; x < 5; x++ {
		g.SetBlocked(x, 2, true)
	}
	if pathing.CanReach(g, 0, 0, 4, 4) {
		t.Fatal("sealed map should not connect")
	}
}

func TestSpawnsCanReachExit(t *testing.T) {
	g, err := grid.New(4, 4)
	if err != nil {
		t.Fatal(err)
	}
	spawns := []pathing.Cell{{0, 0}, {3, 0}}
	exit := pathing.Cell{X: 1, Y: 3}
	if !pathing.SpawnsCanReachExit(g, spawns, exit) {
		t.Fatal("expected reachable")
	}
	for x := 0; x < 4; x++ {
		g.SetBlocked(x, 1, true)
	}
	if pathing.SpawnsCanReachExit(g, spawns, exit) {
		t.Fatal("expected sealed")
	}
}
