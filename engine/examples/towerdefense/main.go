// Headless tower-defense smoke demo: arena build phase → wave → movement.
//
//	cd engine && mise exec -- go run ./examples/towerdefense
package main

import (
	"fmt"

	"github.com/just-Bri/gomakeagame/engine/arena"
	"github.com/just-Bri/gomakeagame/engine/ecs"
	"github.com/just-Bri/gomakeagame/engine/grid"
	"github.com/just-Bri/gomakeagame/engine/move"
)

// Enemy is game-specific combat data. Movement uses move.Agent separately.
type Enemy struct {
	Health float64
	Reward int
}

func main() {
	spawns, exit, err := arena.DefaultTopBottom(8, 5)
	if err != nil {
		panic(err)
	}
	a, err := arena.New(arena.Config{
		Width:       8,
		Height:      5,
		TileSize:    32,
		Spawns:      spawns,
		Exit:        exit,
		RequirePath: true,
		LivePathing: true,
	})
	if err != nil {
		panic(err)
	}

	// Build a simple maze with mixed footprint sizes.
	mustPlace := func(x, y, w, h int) {
		fp, err := grid.NewRect(x, y, w, h)
		if err != nil {
			panic(err)
		}
		if err := a.Place(fp); err != nil {
			panic(err)
		}
	}
	mustPlace(2, 1, 2, 1) // 2×1 bunker
	mustPlace(5, 2, 1, 2) // 1×2 wall
	mustPlace(1, 2, 1, 1)

	// Sell then buy elsewhere (game owns economy; arena only clears/blocks cells).
	sold, _ := grid.NewRect(1, 2, 1, 1)
	if err := a.Unplace(sold); err != nil {
		panic(err)
	}
	mustPlace(3, 3, 1, 1)

	if err := a.BeginWave(); err != nil {
		panic(err)
	}
	fmt.Printf("phase=%s path_ok=%v dirty=%v\n", a.Phase, a.SpawnsReachable(), a.Dirty())

	w := ecs.NewWorld()
	agents := ecs.NewStore[move.Agent](w)
	transforms := ecs.NewStore[move.Transform](w)
	enemies := ecs.NewStore[Enemy](w)

	spawnAt := func(cellX, cellY int, speed float64) {
		e := w.Spawn()
		wx, wy := a.Grid.CellCenter(cellX, cellY, a.TileSize())
		agents.Set(e, move.Agent{Speed: speed, CellX: cellX, CellY: cellY})
		transforms.Set(e, move.Transform{X: wx, Y: wy})
		enemies.Set(e, Enemy{Health: 100, Reward: 10})
	}
	spawnAt(0, 0, 64)
	spawnAt(4, 0, 48)

	fmt.Println("tick  cellX  cellY      x      y")
	for tick := 0; tick <= 180; tick++ {
		if tick%30 == 0 {
			ecs.Join2(agents, transforms, func(_ ecs.Entity, ag *move.Agent, tr *move.Transform) {
				fmt.Printf("%4d  %5d  %5d  %6.1f %6.1f\n", tick, ag.CellX, ag.CellY, tr.X, tr.Y)
			})
		}
		move.FollowFlow(agents, transforms, a.Field, a.Grid, a.TileSize(), 1.0/60.0)
	}
}
