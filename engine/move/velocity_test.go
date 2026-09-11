package move_test

import (
	"testing"

	"github.com/just-Bri/gomakeagame/engine/ecs"
	"github.com/just-Bri/gomakeagame/engine/move"
)

func TestIntegrate(t *testing.T) {
	w := ecs.NewWorld()
	tr := ecs.NewStore[move.Transform](w)
	vel := ecs.NewStore[move.Velocity](w)
	e := w.Spawn()
	tr.Set(e, move.Transform{X: 0, Y: 0})
	vel.Set(e, move.Velocity{VX: 10, VY: -5})
	move.Integrate(tr, vel, 0.5)
	got, _ := tr.Get(e)
	if got.X != 5 || got.Y != -2.5 {
		t.Fatalf("got %+v, want {5,-2.5}", got)
	}
}

func TestClampRect(t *testing.T) {
	tr := &move.Transform{X: -10, Y: 900}
	move.ClampRect(tr, 13, 17, move.Bounds{X: 0, Y: 0, W: 1280, H: 800})
	if tr.X != 0 {
		t.Fatalf("x = %v, want 0", tr.X)
	}
	if tr.Y != 800-17 {
		t.Fatalf("y = %v, want %v", tr.Y, 800-17)
	}
}
