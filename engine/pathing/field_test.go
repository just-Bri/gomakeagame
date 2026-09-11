package pathing_test

import (
	"testing"

	"github.com/just-Bri/gomakeagame/engine/grid"
	"github.com/just-Bri/gomakeagame/engine/pathing"
)

func openGrid(t *testing.T, w, h int) *grid.Grid {
	t.Helper()
	g, err := grid.New(w, h)
	if err != nil {
		t.Fatal(err)
	}
	return g
}

func TestOpenMapPathsToExit(t *testing.T) {
	g := openGrid(t, 5, 5)
	f, err := pathing.Rebuild(g, 4, 4)
	if err != nil {
		t.Fatal(err)
	}

	if !f.Reachable(0, 0) {
		t.Fatal("open map: (0,0) should reach exit")
	}
	if f.Dist(4, 4) != 0 {
		t.Fatalf("exit dist = %d, want 0", f.Dist(4, 4))
	}
	if f.Dist(0, 0) != 8 { // manhattan on open grid
		t.Fatalf("dist(0,0) = %d, want 8", f.Dist(0, 0))
	}

	// Each step should decrease distance.
	x, y := 0, 0
	for step := 0; step < 20; step++ {
		if x == 4 && y == 4 {
			return
		}
		nx, ny, ok := f.Next(x, y)
		if !ok {
			t.Fatalf("stuck at (%d,%d)", x, y)
		}
		if f.Dist(nx, ny) >= f.Dist(x, y) && !(nx == x && ny == y) {
			t.Fatalf("next (%d,%d) did not approach exit from (%d,%d)", nx, ny, x, y)
		}
		x, y = nx, ny
	}
	t.Fatal("did not reach exit in time")
}

func TestChokePoint(t *testing.T) {
	g := openGrid(t, 5, 3)
	// Wall with a single gap at (2,1).
	for y := 0; y < 3; y++ {
		g.SetBlocked(2, y, true)
	}
	g.SetBlocked(2, 1, false)

	f, err := pathing.Rebuild(g, 4, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !f.Reachable(0, 1) {
		t.Fatal("should path through choke")
	}

	// Path from left must pass through the gap.
	x, y := 0, 1
	sawGap := false
	for i := 0; i < 20; i++ {
		if x == 2 && y == 1 {
			sawGap = true
		}
		if x == 4 && y == 1 {
			break
		}
		nx, ny, ok := f.Next(x, y)
		if !ok {
			t.Fatalf("unreachable mid-path at (%d,%d)", x, y)
		}
		x, y = nx, ny
	}
	if !sawGap {
		t.Fatal("path did not go through choke gap (2,1)")
	}
}

func TestFullyBlockedUnreachable(t *testing.T) {
	g := openGrid(t, 4, 4)
	for y := 0; y < 4; y++ {
		g.SetBlocked(2, y, true) // solid wall splitting map
	}
	f, err := pathing.Rebuild(g, 3, 0)
	if err != nil {
		t.Fatal(err)
	}
	if f.Reachable(0, 0) {
		t.Fatal("left side should be unreachable")
	}
	if !f.Reachable(3, 3) {
		t.Fatal("right side should reach exit")
	}
	if _, _, ok := f.Direction(0, 1); ok {
		t.Fatal("unreachable cell must not yield a direction")
	}
}

func TestRebuildAfterPlacingBlocker(t *testing.T) {
	g := openGrid(t, 3, 1)
	f1, err := pathing.Rebuild(g, 2, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !f1.Reachable(0, 0) {
		t.Fatal("open corridor should be reachable")
	}

	g.SetBlocked(1, 0, true)
	f2, err := pathing.Rebuild(g, 2, 0)
	if err != nil {
		t.Fatal(err)
	}
	if f2.Reachable(0, 0) {
		t.Fatal("after blocking middle, start should be unreachable")
	}
	if !f2.Reachable(2, 0) {
		t.Fatal("exit should remain reachable")
	}
}

func TestRebuildRejectsBadExit(t *testing.T) {
	g := openGrid(t, 3, 3)
	g.SetBlocked(1, 1, true)
	if _, err := pathing.Rebuild(g, 1, 1); err == nil {
		t.Fatal("expected error for blocked exit")
	}
	if _, err := pathing.Rebuild(g, 9, 9); err == nil {
		t.Fatal("expected error for OOB exit")
	}
}

func TestDirectionUnitLength(t *testing.T) {
	g := openGrid(t, 3, 1)
	f, err := pathing.Rebuild(g, 2, 0)
	if err != nil {
		t.Fatal(err)
	}
	dx, dy, ok := f.Direction(0, 0)
	if !ok {
		t.Fatal("expected direction")
	}
	if dx != 1 || dy != 0 {
		t.Fatalf("Direction = (%v,%v), want (1,0)", dx, dy)
	}
	if _, _, ok := f.Direction(2, 0); ok {
		t.Fatal("exit should not report a move direction")
	}
}

func TestPreferDownOverSideways(t *testing.T) {
	// Open 5×5 with full bottom exit strip. From (0,2), Dist to exit is 2
	// via (0,3)→(0,4) or via sideways then down. Prefer straight down.
	g := openGrid(t, 5, 5)
	exits := make([]pathing.Cell, 5)
	for x := 0; x < 5; x++ {
		exits[x] = pathing.Cell{X: x, Y: 4}
	}
	f, err := pathing.RebuildFromExits(g, exits)
	if err != nil {
		t.Fatal(err)
	}
	nx, ny, ok := f.Next(0, 2)
	if !ok {
		t.Fatal("expected next")
	}
	if nx != 0 || ny != 3 {
		t.Fatalf("Next(0,2)=(%d,%d), want (0,3) straight down", nx, ny)
	}
	if f.Dist(0, 3) != 1 || f.Dist(2, 4) != 0 {
		t.Fatalf("unexpected dists: (0,3)=%d (2,4)=%d", f.Dist(0, 3), f.Dist(2, 4))
	}
}
