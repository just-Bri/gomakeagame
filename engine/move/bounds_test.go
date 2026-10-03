package move

import "testing"

func TestClampPointXZ(t *testing.T) {
	b := AABBXZ{MinX: -5, MaxX: 5, MinZ: -5, MaxZ: 5}
	x, z := ClampPointXZ(9, -9, b)
	if x != 5 || z != -5 {
		t.Fatalf("got %v,%v", x, z)
	}
}
