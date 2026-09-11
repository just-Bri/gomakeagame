package collide_test

import (
	"math"
	"testing"

	"github.com/just-Bri/gomakeagame/engine/collide"
)

func TestDistPointSegment(t *testing.T) {
	d := collide.DistPointSegment(0, 1, 0, 0, 10, 0)
	if math.Abs(d-1) > 1e-9 {
		t.Fatalf("dist = %v, want 1", d)
	}
	d = collide.DistPointSegment(5, 0, 0, 0, 10, 0)
	if math.Abs(d) > 1e-9 {
		t.Fatalf("on-segment dist = %v, want 0", d)
	}
}

func TestSegmentHitsCircle(t *testing.T) {
	if !collide.SegmentHitsCircle(0, 0, 10, 0, 5, 1, 1.5) {
		t.Fatal("expected hit")
	}
	if collide.SegmentHitsCircle(0, 0, 10, 0, 5, 3, 1) {
		t.Fatal("expected miss")
	}
}

func TestClosestHitAlongSegment(t *testing.T) {
	id, ok := collide.ClosestHitAlongSegment(0, 0, 10, 0, 2, func(yield func(x, y float64, id uint64)) {
		yield(5, 3, 1) // too far
		yield(5, 1, 2) // closer
		yield(5, 0.5, 3)
	})
	if !ok || id != 3 {
		t.Fatalf("got id=%d ok=%v, want 3", id, ok)
	}
}

func TestCircleHitsCircle(t *testing.T) {
	if !collide.CircleHitsCircle(0, 0, 1, 1.5, 0, 1) {
		t.Fatal("expected overlap")
	}
	if collide.CircleHitsCircle(0, 0, 1, 3, 0, 1) {
		t.Fatal("expected miss")
	}
}

func TestCircleHitsRect(t *testing.T) {
	if !collide.CircleHitsRect(5, 5, 1, 0, 0, 10, 10) {
		t.Fatal("center inside rect should hit")
	}
	if !collide.CircleHitsRect(11, 5, 1.5, 0, 0, 10, 10) {
		t.Fatal("near right edge should hit")
	}
	if collide.CircleHitsRect(20, 5, 1, 0, 0, 10, 10) {
		t.Fatal("far away should miss")
	}
}
