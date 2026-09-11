package gfx_test

import (
	"math"
	"testing"

	"github.com/just-Bri/gomakeagame/engine/gfx"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func TestFitLetterboxPillarbox(t *testing.T) {
	// Virtual 16:10 into a wider window → bars on left/right.
	lb := gfx.Fit(1600, 800, 1280, 800)
	if !almostEqual(lb.Scale, 1) {
		t.Fatalf("scale=%v want 1", lb.Scale)
	}
	if !almostEqual(lb.OffsetX, 160) || !almostEqual(lb.OffsetY, 0) {
		t.Fatalf("offset=(%v,%v) want (160,0)", lb.OffsetX, lb.OffsetY)
	}
}

func TestFitLetterboxWindowbox(t *testing.T) {
	// Virtual into a taller window → bars on top/bottom.
	lb := gfx.Fit(1280, 1000, 1280, 800)
	if !almostEqual(lb.Scale, 1) {
		t.Fatalf("scale=%v want 1", lb.Scale)
	}
	if !almostEqual(lb.OffsetX, 0) || !almostEqual(lb.OffsetY, 100) {
		t.Fatalf("offset=(%v,%v) want (0,100)", lb.OffsetX, lb.OffsetY)
	}
}

func TestFitUniformDownscale(t *testing.T) {
	lb := gfx.Fit(640, 400, 1280, 800)
	if !almostEqual(lb.Scale, 0.5) {
		t.Fatalf("scale=%v want 0.5", lb.Scale)
	}
	if !almostEqual(lb.OffsetX, 0) || !almostEqual(lb.OffsetY, 0) {
		t.Fatalf("offset=(%v,%v) want (0,0)", lb.OffsetX, lb.OffsetY)
	}
}

func TestFitDegenerate(t *testing.T) {
	lb := gfx.Fit(0, 100, 1280, 800)
	if lb.Scale != 0 {
		t.Fatalf("scale=%v want 0", lb.Scale)
	}
}

func TestCoordRoundTrip(t *testing.T) {
	lb := gfx.Fit(1600, 800, 1280, 800)
	sx, sy := lb.ToScreen(100, 50)
	vx, vy := lb.ToVirtual(sx, sy)
	if !almostEqual(vx, 100) || !almostEqual(vy, 50) {
		t.Fatalf("round-trip got (%v,%v)", vx, vy)
	}
}
