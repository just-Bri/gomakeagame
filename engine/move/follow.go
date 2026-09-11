// Package move advances path-following agents across a shared flow field.
//
// Agents are intentionally generic: tower-defense enemies use them today;
// player-controlled units can reuse the same Transform + Agent stores later.
// Combat stats (health, reward, armor) stay in game-specific components.
package move

import (
	"math"

	"github.com/just-Bri/gomakeagame/engine/ecs"
	"github.com/just-Bri/gomakeagame/engine/grid"
	"github.com/just-Bri/gomakeagame/engine/pathing"
)

// Transform is world-space position. Shared by any moving unit.
type Transform struct {
	X, Y float64
}

// Agent is path-following state for enemies and later player units.
//
// CellX/CellY cache the tile the agent occupies; FollowFlow refreshes them
// when the world position crosses a tile boundary.
type Agent struct {
	Speed        float64
	CellX, CellY int
}

// FollowFlow moves every entity that has both Agent and Transform along field.
//
// Each tick:
//  1. Refresh the cell cache from world position.
//  2. Aim at the center of the next flow-field cell (or stop at the exit).
//  3. Integrate position by Speed * dt toward that target.
//
// Agents on unreachable tiles do not move. tileSize must match the grid used
// to build field. g may be nil if you only need field sampling (cell refresh
// still uses tileSize alone).
func FollowFlow(
	agents *ecs.Store[Agent],
	transforms *ecs.Store[Transform],
	field *pathing.Field,
	g *grid.Grid,
	tileSize float64,
	dt float64,
) {
	if agents == nil || transforms == nil || field == nil || tileSize <= 0 || dt == 0 {
		return
	}

	ecs.Join2(agents, transforms, func(_ ecs.Entity, agent *Agent, tr *Transform) {
		cx, cy := cellOf(tr.X, tr.Y, tileSize)
		agent.CellX, agent.CellY = cx, cy

		if g != nil && !g.InBounds(cx, cy) {
			return
		}

		nx, ny, ok := field.Next(cx, cy)
		if !ok {
			return
		}

		// At exit cell: ease toward center, then stop.
		tx, ty := cellCenter(nx, ny, tileSize)
		dx := tx - tr.X
		dy := ty - tr.Y
		dist := math.Hypot(dx, dy)
		if dist < 1e-6 {
			tr.X, tr.Y = tx, ty
			agent.CellX, agent.CellY = nx, ny
			return
		}

		step := agent.Speed * dt
		if step >= dist {
			tr.X, tr.Y = tx, ty
		} else {
			tr.X += dx / dist * step
			tr.Y += dy / dist * step
		}

		agent.CellX, agent.CellY = cellOf(tr.X, tr.Y, tileSize)
	})
}

func cellOf(wx, wy, tileSize float64) (cx, cy int) {
	return int(math.Floor(wx / tileSize)), int(math.Floor(wy / tileSize))
}

func cellCenter(cx, cy int, tileSize float64) (wx, wy float64) {
	return (float64(cx) + 0.5) * tileSize, (float64(cy) + 0.5) * tileSize
}
