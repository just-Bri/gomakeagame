package pathing_test

import (
	"testing"

	"github.com/just-Bri/gomakeagame/engine/pathing"
)

func TestTraceOpenCorridor(t *testing.T) {
	g := openGrid(t, 3, 1)
	f, err := pathing.Rebuild(g, 2, 0)
	if err != nil {
		t.Fatal(err)
	}
	got := pathing.Trace(f, 0, 0)
	want := []pathing.Cell{{0, 0}, {1, 0}, {2, 0}}
	if len(got) != len(want) {
		t.Fatalf("Trace len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Trace[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestTraceUnreachableNil(t *testing.T) {
	g := openGrid(t, 3, 1)
	g.SetBlocked(1, 0, true)
	f, err := pathing.Rebuild(g, 2, 0)
	if err != nil {
		t.Fatal(err)
	}
	if pathing.Trace(f, 0, 0) != nil {
		t.Fatal("unreachable start should Trace nil")
	}
	if pathing.Trace(nil, 0, 0) != nil {
		t.Fatal("nil field should Trace nil")
	}
}

func TestCorridorFromUnion(t *testing.T) {
	g := openGrid(t, 5, 3)
	// Wall with gaps at (0,1) left lane and (4,1) right lane — actually
	// simpler: open map with full bottom exits; corridor from two spawns.
	exits := []pathing.Cell{{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}, {X: 4, Y: 2}}
	f, err := pathing.RebuildFromExits(g, exits)
	if err != nil {
		t.Fatal(err)
	}
	starts := []pathing.Cell{{X: 0, Y: 0}, {X: 4, Y: 0}}
	cells := pathing.CorridorFrom(f, starts)
	if len(cells) < 4 {
		t.Fatalf("corridor too short: %v", cells)
	}
	seen := map[pathing.Cell]bool{}
	for _, c := range cells {
		if seen[c] {
			t.Fatalf("duplicate cell in corridor: %+v", c)
		}
		seen[c] = true
	}
	if !seen[pathing.Cell{X: 0, Y: 0}] || !seen[pathing.Cell{X: 4, Y: 0}] {
		t.Fatalf("starts missing from corridor: %v", cells)
	}
	if pathing.LongestDist(f, starts) != 2 {
		t.Fatalf("LongestDist = %d, want 2", pathing.LongestDist(f, starts))
	}
}

func TestCorridorFromChoke(t *testing.T) {
	g := openGrid(t, 5, 3)
	for y := 0; y < 3; y++ {
		g.SetBlocked(2, y, true)
	}
	g.SetBlocked(2, 1, false)
	f, err := pathing.Rebuild(g, 4, 1)
	if err != nil {
		t.Fatal(err)
	}
	cells := pathing.CorridorFrom(f, []pathing.Cell{{X: 0, Y: 1}})
	sawGap := false
	for _, c := range cells {
		if c.X == 2 && c.Y == 1 {
			sawGap = true
		}
	}
	if !sawGap {
		t.Fatalf("choke gap missing from corridor: %v", cells)
	}
}
