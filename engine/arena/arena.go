// Package arena is the shared playfield for maze-style tower defense.
//
// It owns the grid, flow field, spawn/exit cells, and build↔wave phase.
// Building *types*, combat, and economy stay in game code; this package
// only manages footprints on walkability and when pathing is rebuilt.
//
// Typical loop (Matrix Defense / Desktop TD style):
//
//	a.BeginBuild()
//	a.Place(footprint)   // or Unplace after a sell — marks dirty
//	a.BeginWave()        // rebuilds field if needed; requires open path
//	move.FollowFlow(...) // during wave
//	a.BeginBuild()       // after wave cleared
package arena

import (
	"errors"
	"fmt"

	"github.com/just-Bri/gomakeagame/engine/grid"
	"github.com/just-Bri/gomakeagame/engine/pathing"
)

// Phase is the high-level play cadence.
type Phase int

const (
	// PhaseBuild: place/pickup/move buildings; creeps are not marching.
	PhaseBuild Phase = iota
	// PhaseWave: flow field is live; footprint edits are locked.
	PhaseWave
)

func (p Phase) String() string {
	switch p {
	case PhaseBuild:
		return "build"
	case PhaseWave:
		return "wave"
	default:
		return fmt.Sprintf("Phase(%d)", int(p))
	}
}

var (
	// ErrNotBuildPhase is returned when Place/Unplace runs outside build.
	ErrNotBuildPhase = errors.New("arena: not in build phase")
	// ErrNotWavePhase is returned when wave-only ops run outside a wave.
	ErrNotWavePhase = errors.New("arena: not in wave phase")
	// ErrFootprintBlocked is returned when a footprint overlaps occupied cells.
	ErrFootprintBlocked = errors.New("arena: footprint not clear")
	// ErrFootprintOOB is returned when a footprint leaves the playable grid.
	ErrFootprintOOB = errors.New("arena: footprint out of bounds")
	// ErrPathBlocked is returned when a place/start would seal spawn→exit.
	ErrPathBlocked = errors.New("arena: placement blocks all paths to exit")
	// ErrNoSpawns is returned when the arena has no spawn cells.
	ErrNoSpawns = errors.New("arena: no spawn cells")
)

// Config seeds a new Arena.
type Config struct {
	Width, Height int
	TileSize      float64
	// Spawns are entry cells (typically a full top strip).
	Spawns []pathing.Cell
	// Exits are creep goals (typically the full bottom row).
	Exits []pathing.Cell
	// RequirePath rejects placements (and BeginWave) that seal every
	// spawn from every exit — Desktop TD / fair-maze style.
	// If false, players may wall off; creeps on sealed tiles simply stall.
	RequirePath bool
	// LivePathing rebuilds the flow field after every Place/Unplace in
	// build phase so UI can preview the maze path. If false, pathing
	// rebuilds on BeginWave (and SyncPathing) only.
	LivePathing bool
}

// Arena is the maze rectangle + phase machine.
type Arena struct {
	Grid   *grid.Grid
	Field  *pathing.Field
	Phase  Phase
	Config Config

	dirty bool
}

// New creates an arena in PhaseBuild with an initial flow field.
func New(cfg Config) (*Arena, error) {
	if cfg.TileSize <= 0 {
		return nil, fmt.Errorf("arena: TileSize must be positive")
	}
	if len(cfg.Spawns) == 0 {
		return nil, ErrNoSpawns
	}
	if len(cfg.Exits) == 0 {
		return nil, fmt.Errorf("arena: no exit cells")
	}
	g, err := grid.New(cfg.Width, cfg.Height)
	if err != nil {
		return nil, err
	}
	for _, s := range cfg.Spawns {
		if !g.InBounds(s.X, s.Y) {
			return nil, fmt.Errorf("arena: spawn (%d,%d) out of bounds", s.X, s.Y)
		}
	}
	for _, e := range cfg.Exits {
		if !g.InBounds(e.X, e.Y) {
			return nil, fmt.Errorf("arena: exit (%d,%d) out of bounds", e.X, e.Y)
		}
	}

	a := &Arena{
		Grid:   g,
		Phase:  PhaseBuild,
		Config: cfg,
		dirty:  true,
	}
	if err := a.SyncPathing(); err != nil {
		return nil, err
	}
	return a, nil
}

// DefaultTopBottom returns spawns along the full top row and exits along
// the full bottom row — the classic “enter top, escape bottom” layout.
func DefaultTopBottom(width, height int) (spawns, exits []pathing.Cell, err error) {
	return VerticalPlayfield(width, height, 1)
}

// VerticalPlayfield is the Matrix Defense–style map:
// a vertical rectangle, a horizontal spawn strip of spawnRows at the top,
// and a full-width exit strip on the bottom row.
//
// Players typically build horizontal walls left→right then right→left,
// leaving a one-cell gap on alternating ends so creeps snake downward.
func VerticalPlayfield(width, height, spawnRows int) (spawns, exits []pathing.Cell, err error) {
	if width <= 0 || height <= 0 {
		return nil, nil, fmt.Errorf("arena: invalid size %dx%d", width, height)
	}
	if spawnRows <= 0 {
		return nil, nil, fmt.Errorf("arena: spawnRows must be positive")
	}
	if spawnRows >= height {
		return nil, nil, fmt.Errorf("arena: spawnRows %d leaves no room for build/exit (height %d)", spawnRows, height)
	}
	spawns = make([]pathing.Cell, 0, width*spawnRows)
	for y := 0; y < spawnRows; y++ {
		for x := 0; x < width; x++ {
			spawns = append(spawns, pathing.Cell{X: x, Y: y})
		}
	}
	exits = make([]pathing.Cell, 0, width)
	for x := 0; x < width; x++ {
		exits = append(exits, pathing.Cell{X: x, Y: height - 1})
	}
	return spawns, exits, nil
}

// BeginBuild enters (or returns to) build phase. Safe to call at wave end.
func (a *Arena) BeginBuild() {
	a.Phase = PhaseBuild
}

// BeginWave locks footprint edits, rebuilds pathing if dirty, and starts
// the wave. Fails if RequirePath and any spawn cannot reach an exit.
func (a *Arena) BeginWave() error {
	if a.Phase == PhaseWave {
		return nil
	}
	if err := a.SyncPathing(); err != nil {
		return err
	}
	if a.Config.RequirePath && !a.SpawnsReachable() {
		return ErrPathBlocked
	}
	a.Phase = PhaseWave
	return nil
}

// MarkDirty notes that walkability changed outside Place/Pickup.
func (a *Arena) MarkDirty() {
	a.dirty = true
}

// Dirty reports whether the flow field may be stale.
func (a *Arena) Dirty() bool {
	return a.dirty
}

// SyncPathing rebuilds the flow field when dirty.
func (a *Arena) SyncPathing() error {
	if !a.dirty && a.Field != nil {
		return nil
	}
	f, err := pathing.RebuildFromExits(a.Grid, a.Config.Exits)
	if err != nil {
		return err
	}
	a.Field = f
	a.dirty = false
	return nil
}

// SpawnsReachable reports whether every configured spawn can reach an exit
// on the *current* grid (cheap BFS; does not require a fresh Field).
func (a *Arena) SpawnsReachable() bool {
	return pathing.SpawnsCanReachExit(a.Grid, a.Config.Spawns, a.Config.Exits)
}

// CanPlace reports whether footprint may be blocked in the current build.
func (a *Arena) CanPlace(fp grid.Rect) error {
	if a.Phase != PhaseBuild {
		return ErrNotBuildPhase
	}
	if !a.Grid.RectInBounds(fp) {
		return ErrFootprintOOB
	}
	if !a.Grid.RectClear(fp) {
		return ErrFootprintBlocked
	}
	if !a.Config.RequirePath {
		return nil
	}

	// Tentatively block, test connectivity, then restore.
	a.Grid.SetRectBlocked(fp, true)
	ok := a.SpawnsReachable()
	a.Grid.SetRectBlocked(fp, false)
	if !ok {
		return ErrPathBlocked
	}
	return nil
}

// Place blocks footprint cells. Build phase only.
func (a *Arena) Place(fp grid.Rect) error {
	if err := a.CanPlace(fp); err != nil {
		return err
	}
	a.Grid.SetRectBlocked(fp, true)
	a.dirty = true
	if a.Config.LivePathing {
		return a.SyncPathing()
	}
	return nil
}

// Unplace clears footprint cells (e.g. after the game sells a building).
// Build phase only. Game code owns economy and which entity occupied the cells.
func (a *Arena) Unplace(fp grid.Rect) error {
	if a.Phase != PhaseBuild {
		return ErrNotBuildPhase
	}
	if !a.Grid.RectInBounds(fp) {
		return ErrFootprintOOB
	}
	a.Grid.SetRectBlocked(fp, false)
	a.dirty = true
	if a.Config.LivePathing {
		return a.SyncPathing()
	}
	return nil
}

// TileSize returns the configured world units per cell.
func (a *Arena) TileSize() float64 {
	return a.Config.TileSize
}
