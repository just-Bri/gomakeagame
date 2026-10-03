package project

import (
	"math"
	"testing"
)

func TestIsoYawRoundTrip(t *testing.T) {
	iso := DefaultIso()
	iso.CamX, iso.CamZ = 3, -2
	iso.Yaw = 1.1
	wx, wz := 8.0, 4.0
	sx, sy := iso.WorldToScreen(wx, 0, wz)
	gx, gz := iso.ScreenToGround(sx, sy)
	if math.Abs(gx-wx) > 1e-9 || math.Abs(gz-wz) > 1e-9 {
		t.Fatalf("round-trip got (%v,%v) want (%v,%v)", gx, gz, wx, wz)
	}
}

func TestIsoYawChangesScreen(t *testing.T) {
	iso := DefaultIso()
	sx0, sy0 := iso.WorldToScreen(5, 0, 0)
	iso.Yaw = math.Pi / 2
	sx1, sy1 := iso.WorldToScreen(5, 0, 0)
	if sx0 == sx1 && sy0 == sy1 {
		t.Fatal("yaw should move projected point")
	}
}
