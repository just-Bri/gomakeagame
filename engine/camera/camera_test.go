package camera

import "testing"

func TestFollowXZMovesTowardTarget(t *testing.T) {
	c := FollowXZ{}
	c.Update(1, 10, 20)
	if c.X <= 0 || c.Z <= 0 {
		t.Fatalf("expected movement, got %+v", c)
	}
	c.Snap(5, 7)
	if c.X != 5 || c.Z != 7 {
		t.Fatalf("snap failed: %+v", c)
	}
}
