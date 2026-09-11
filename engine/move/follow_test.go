package move_test

import (
	"math"
	"testing"

	"github.com/just-Bri/gomakeagame/engine/ecs"
	"github.com/just-Bri/gomakeagame/engine/grid"
	"github.com/just-Bri/gomakeagame/engine/move"
	"github.com/just-Bri/gomakeagame/engine/pathing"
)

const tile = 10.0

func setupCorridor(t *testing.T) (*grid.Grid, *pathing.Field) {
	t.Helper()
	g, err := grid.New(5, 1)
	if err != nil {
		t.Fatal(err)
	}
	f, err := pathing.Rebuild(g, 4, 0)
	if err != nil {
		t.Fatal(err)
	}
	return g, f
}

func spawnAgent(w *ecs.World, agents *ecs.Store[move.Agent], transforms *ecs.Store[move.Transform], x, y, speed float64) ecs.Entity {
	e := w.Spawn()
	agents.Set(e, move.Agent{Speed: speed})
	transforms.Set(e, move.Transform{X: x, Y: y})
	return e
}

func TestFollowFlowAdvancesTowardExit(t *testing.T) {
	g, f := setupCorridor(t)
	w := ecs.NewWorld()
	agents := ecs.NewStore[move.Agent](w)
	transforms := ecs.NewStore[move.Transform](w)

	// Start at center of cell (0,0).
	e := spawnAgent(w, agents, transforms, 5, 5, 20)
	_ = e

	for i := 0; i < 100; i++ {
		move.FollowFlow(agents, transforms, f, g, tile, 1.0/60.0)
	}

	tr, _ := transforms.Get(e)
	ag, _ := agents.Get(e)
	if tr.X <= 5 {
		t.Fatalf("expected movement toward exit, x=%v", tr.X)
	}
	if ag.CellX < 1 {
		t.Fatalf("expected cell cache to advance, cell=%d", ag.CellX)
	}
}

func TestFollowFlowReachesExitCenter(t *testing.T) {
	g, f := setupCorridor(t)
	w := ecs.NewWorld()
	agents := ecs.NewStore[move.Agent](w)
	transforms := ecs.NewStore[move.Transform](w)

	e := spawnAgent(w, agents, transforms, 5, 5, 200) // fast
	for i := 0; i < 500; i++ {
		move.FollowFlow(agents, transforms, f, g, tile, 1.0/60.0)
		tr, _ := transforms.Get(e)
		ex, ey := g.CellCenter(4, 0, tile)
		if math.Hypot(tr.X-ex, tr.Y-ey) < 1e-6 {
			ag, _ := agents.Get(e)
			if ag.CellX != 4 || ag.CellY != 0 {
				t.Fatalf("cell cache = (%d,%d), want (4,0)", ag.CellX, ag.CellY)
			}
			return
		}
	}
	tr, _ := transforms.Get(e)
	t.Fatalf("never reached exit; pos=(%v,%v)", tr.X, tr.Y)
}

func TestFollowFlowSkipsUnreachable(t *testing.T) {
	g, err := grid.New(4, 1)
	if err != nil {
		t.Fatal(err)
	}
	g.SetBlocked(1, 0, true)
	f, err := pathing.Rebuild(g, 3, 0)
	if err != nil {
		t.Fatal(err)
	}

	w := ecs.NewWorld()
	agents := ecs.NewStore[move.Agent](w)
	transforms := ecs.NewStore[move.Transform](w)
	e := spawnAgent(w, agents, transforms, 5, 5, 50) // cell (0,0), unreachable

	move.FollowFlow(agents, transforms, f, g, tile, 1.0)
	tr, _ := transforms.Get(e)
	if tr.X != 5 || tr.Y != 5 {
		t.Fatalf("unreachable agent moved to (%v,%v)", tr.X, tr.Y)
	}
}

func TestFollowFlowWorksWithoutEnemyComponent(t *testing.T) {
	// Proves the system is unit-agnostic: only Agent + Transform required.
	g, f := setupCorridor(t)
	w := ecs.NewWorld()
	agents := ecs.NewStore[move.Agent](w)
	transforms := ecs.NewStore[move.Transform](w)

	type PlayerTag struct{}
	players := ecs.NewStore[PlayerTag](w)

	e := spawnAgent(w, agents, transforms, 5, 5, 30)
	players.Set(e, PlayerTag{})

	move.FollowFlow(agents, transforms, f, g, tile, 0.5)
	tr, _ := transforms.Get(e)
	if tr.X <= 5 {
		t.Fatal("player-tagged agent should still follow the flow field")
	}
}
