package spatial_test

import (
	"testing"

	"github.com/just-Bri/gomakeagame/engine/spatial"
)

func TestNearest(t *testing.T) {
	id, ok := spatial.Nearest(0, 0, 10, func(yield func(x, y float64, id uint64)) {
		yield(8, 0, 1)
		yield(3, 0, 2)
		yield(20, 0, 3) // out of range
	})
	if !ok || id != 2 {
		t.Fatalf("got id=%d ok=%v, want 2", id, ok)
	}
}

func TestNearestNone(t *testing.T) {
	_, ok := spatial.Nearest(0, 0, 1, func(yield func(x, y float64, id uint64)) {
		yield(5, 0, 1)
	})
	if ok {
		t.Fatal("expected no hit")
	}
}
