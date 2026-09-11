package aim_test

import (
	"math"
	"testing"

	"github.com/just-Bri/gomakeagame/engine/aim"
)

func TestLeadStationary(t *testing.T) {
	ax, ay := aim.LeadPoint(0, 0, 100, 0, 0, 0, 50)
	if math.Abs(ax-100) > 1e-6 || math.Abs(ay) > 1e-6 {
		t.Fatalf("stationary aim = (%v,%v), want (100,0)", ax, ay)
	}
}

func TestLeadHeadOn(t *testing.T) {
	ax, ay := aim.LeadPoint(0, 0, 100, 0, -10, 0, 50)
	if ay != 0 {
		t.Fatalf("ay = %v, want 0", ay)
	}
	if ax <= 0 || ax >= 100 {
		t.Fatalf("ax = %v, want in (0,100)", ax)
	}
	tb := ax / 50
	te := (100 - ax) / 10
	if math.Abs(tb-te) > 1e-4 {
		t.Fatalf("times mismatch bullet=%v enemy=%v", tb, te)
	}
}

func TestLeadFallbackWhenImpossible(t *testing.T) {
	ax, ay := aim.LeadPoint(0, 0, 10, 0, 100, 0, 20)
	if ax != 10 || ay != 0 {
		t.Fatalf("fallback = (%v,%v), want current (10,0)", ax, ay)
	}
}
