package arena_test

import (
	"errors"
	"testing"

	"github.com/just-Bri/gomakeagame/engine/arena"
	"github.com/just-Bri/gomakeagame/engine/grid"
	"github.com/just-Bri/gomakeagame/engine/pathing"
)

func newTestArena(t *testing.T, requirePath, live bool) *arena.Arena {
	t.Helper()
	spawns, exits, err := arena.DefaultTopBottom(8, 6)
	if err != nil {
		t.Fatal(err)
	}
	a, err := arena.New(arena.Config{
		Width:       8,
		Height:      6,
		TileSize:    32,
		Spawns:      spawns,
		Exits:       exits,
		RequirePath: requirePath,
		LivePathing: live,
	})
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func TestPlaceUnplaceBuildPhase(t *testing.T) {
	a := newTestArena(t, true, false)
	fp, _ := grid.NewRect(2, 2, 2, 2)
	if err := a.Place(fp); err != nil {
		t.Fatal(err)
	}
	if a.Grid.Walkable(2, 2) {
		t.Fatal("placed footprint should block")
	}
	if !a.Dirty() {
		t.Fatal("place should mark dirty when LivePathing is off")
	}
	if err := a.Unplace(fp); err != nil {
		t.Fatal(err)
	}
	if !a.Grid.Walkable(2, 2) {
		t.Fatal("unplace should clear footprint")
	}
}

func TestRequirePathRejectsSeal(t *testing.T) {
	a := newTestArena(t, true, false)
	// Wall off the entire width on a middle row with 1×1 blocks.
	for x := 0; x < 7; x++ {
		fp, _ := grid.NewRect(x, 3, 1, 1)
		if err := a.Place(fp); err != nil {
			t.Fatalf("place %d: %v", x, err)
		}
	}
	last, _ := grid.NewRect(7, 3, 1, 1)
	if err := a.Place(last); !errors.Is(err, arena.ErrPathBlocked) {
		t.Fatalf("sealing place = %v, want ErrPathBlocked", err)
	}
	if !a.Grid.Walkable(7, 3) {
		t.Fatal("rejected place must not leave the cell blocked")
	}
}

func TestBeginWaveRebuildsAndLocks(t *testing.T) {
	a := newTestArena(t, true, false)
	fp, _ := grid.NewRect(3, 2, 2, 1)
	if err := a.Place(fp); err != nil {
		t.Fatal(err)
	}
	if err := a.BeginWave(); err != nil {
		t.Fatal(err)
	}
	if a.Phase != arena.PhaseWave {
		t.Fatalf("phase = %s, want wave", a.Phase)
	}
	if a.Dirty() {
		t.Fatal("BeginWave should sync pathing")
	}
	if a.Field == nil || !pathing.AllReachable(a.Field, a.Config.Spawns) {
		t.Fatal("spawns should be reachable on field after BeginWave")
	}
	if err := a.Place(fp); !errors.Is(err, arena.ErrNotBuildPhase) {
		t.Fatalf("place during wave = %v, want ErrNotBuildPhase", err)
	}

	a.BeginBuild()
	if a.Phase != arena.PhaseBuild {
		t.Fatalf("phase = %s, want build", a.Phase)
	}
}

func TestBeginWaveFailsWhenSealedWithoutRequireOnPlace(t *testing.T) {
	// RequirePath false allows sealing during place; BeginWave with
	// RequirePath true would need a different config — here we seal with
	// RequirePath false, then manually check SpawnsReachable.
	a := newTestArena(t, false, false)
	for x := 0; x < 8; x++ {
		fp, _ := grid.NewRect(x, 3, 1, 1)
		if err := a.Place(fp); err != nil {
			t.Fatal(err)
		}
	}
	if a.SpawnsReachable() {
		t.Fatal("full wall should seal spawns")
	}
	// Flip policy for start check by using a RequirePath arena after seal:
	a.Config.RequirePath = true
	if err := a.BeginWave(); !errors.Is(err, arena.ErrPathBlocked) {
		t.Fatalf("BeginWave = %v, want ErrPathBlocked", err)
	}
	if a.Phase != arena.PhaseBuild {
		t.Fatal("failed BeginWave must stay in build")
	}
}

func TestUnplaceDuringWaveRejected(t *testing.T) {
	a := newTestArena(t, true, false)
	fp, _ := grid.NewRect(3, 2, 2, 1)
	if err := a.Place(fp); err != nil {
		t.Fatal(err)
	}
	if err := a.BeginWave(); err != nil {
		t.Fatal(err)
	}
	if err := a.Unplace(fp); !errors.Is(err, arena.ErrNotBuildPhase) {
		t.Fatalf("unplace during wave = %v, want ErrNotBuildPhase", err)
	}
}

func TestLivePathingClearsDirty(t *testing.T) {
	a := newTestArena(t, true, true)
	fp, _ := grid.NewRect(2, 2, 1, 1)
	if err := a.Place(fp); err != nil {
		t.Fatal(err)
	}
	if a.Dirty() {
		t.Fatal("LivePathing place should sync immediately")
	}
	if a.Field == nil {
		t.Fatal("expected field after live place")
	}
}

func TestVerticalPlayfieldSpawnStrip(t *testing.T) {
	spawns, exits, err := arena.VerticalPlayfield(11, 24, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(spawns) != 22 {
		t.Fatalf("spawns = %d, want 22", len(spawns))
	}
	if len(exits) != 11 {
		t.Fatalf("exits = %d, want 11", len(exits))
	}
	for _, e := range exits {
		if e.Y != 23 {
			t.Fatalf("exit outside bottom row: %+v", e)
		}
	}
	for _, s := range spawns {
		if s.Y >= 2 {
			t.Fatalf("spawn outside strip: %+v", s)
		}
	}
}

func TestOverlappingPlaceRejected(t *testing.T) {
	a := newTestArena(t, true, false)
	fp, _ := grid.NewRect(2, 2, 2, 2)
	if err := a.Place(fp); err != nil {
		t.Fatal(err)
	}
	overlap, _ := grid.NewRect(3, 3, 2, 2)
	if err := a.Place(overlap); !errors.Is(err, arena.ErrFootprintBlocked) {
		t.Fatalf("overlap = %v, want ErrFootprintBlocked", err)
	}
}
