# ECS + Grid Pathing — Agent Handoff

Last updated: 2026-09-11

This note captures where Gomag’s engine stands and the plan for maze-style
tower defense (Matrix Defense / Desktop TD inspired). Read this before extending `engine/`.

## Project context

- **Repo:** `github.com/just-Bri/gomakeagame` (monorepo)
- **Engine module:** `engine/` → `github.com/just-Bri/gomakeagame/engine`
- **First game:** `let_none_through` — top-down maze TD; creeps enter top → exit bottom; player builds multi-size buildings to shape the path; pick up / reposition buildings between waves.
- **Language:** Pure Go, **no CGO**.
- **Toolchain:** Go **1.27.1** via `mise.toml`.

## What exists today

| Package | Role |
|---------|------|
| `ecs/` | Generational entities, sparse-set `Store[T]`, `Join2`/`Join3` |
| `grid/` | Tile walkability, `Rect` footprints, world↔cell |
| `pathing/` | Flow field `Rebuild`, `Trace` / `CorridorFrom`, `CanReach` / `FindPath`, `SpawnsCanReachExit` |
| `move/` | Generic `Agent` + `Transform` + `FollowFlow` |
| `arena/` | Build↔wave phase, Place/Pickup/Move, dirty pathing sync |
| `examples/towerdefense/` | Headless maze + wave smoke demo |

**Run:**

```bash
cd engine && mise exec -- go test ./...
cd engine && mise exec -- go run ./examples/towerdefense
```

### Design choices (do not casually reverse)

1. **Fat game components** (`Enemy` health/reward) + shared `move.Agent` for pathing.
2. **Systems are plain functions.**
3. **Grid + flow field are shared resources**, not per-tile ECS entities.
4. **Arena owns phase + footprints on walkability.** Building defs, damage, economy stay in the game.
5. **`RequirePath`** (default for fair mazes) rejects placements that seal spawn→exit.

### Build / wave loop

```go
spawns, exit, _ := arena.DefaultTopBottom(24, 18)
a, _ := arena.New(arena.Config{
    Width: 24, Height: 18, TileSize: 32,
    Spawns: spawns, Exit: exit,
    RequirePath: true,
    LivePathing: true, // preview maze path while placing
})

a.BeginBuild()
a.Place(grid.Rect{X: 4, Y: 5, W: 2, H: 2}) // multi-tile bunker
a.Unplace(oldFootprint)                      // after game-side sell
a.BeginWave()                                // SyncPathing + lock edits

move.FollowFlow(agents, transforms, a.Field, a.Grid, a.TileSize(), dt)

a.BeginBuild() // after wave
```

## Engine vs game split

| Concern | Where |
|---------|--------|
| Footprint block/clear, path seal checks, build/wave lock | `engine/arena` + `grid` + `pathing` |
| Flow follow for any unit | `engine/move` |
| Building catalog (size, cost, damage, aura) | **game** |
| Sell / buy / “move” by refund+repurchase | **game** (calls `Unplace` then `Place`) |
| Waves, lives, cash, win/lose | **game** |
| Tower targeting / projectiles | **game** first; extract to engine when a second game needs it |
| Occupancy map (which entity owns a cell) | **game** (or thin engine helper later if needed) |

## Next steps

1. ~~grid / pathing / move / arena build phase~~ **done**
2. ~~`let_none_through` bootstrap~~ **done** (maze TD consuming arena)
3. ~~Path corridor helpers (`Trace` / `CorridorFrom`)~~ **done** — maze UI paints forced path
4. **Targeting** when shooting towers exist — range query over enemies (game or `engine/spatial`) — largely in game already
5. **Defer:** window/GPU polish, editor, fancy ECS

## Out of scope right now

- ECS library / CGO / Raylib as core
- Per-enemy A* as the default
- Full building/combat framework in the engine
