package tilemap

import (
	"testing"

	"github.com/just-Bri/gomakeagame/engine/grid"
)

func TestFromWalkable(t *testing.T) {
	g, err := grid.New(3, 2)
	if err != nil {
		t.Fatal(err)
	}
	g.SetWalkable(1, 0, false)
	l := FromWalkable(g, RGB{0, 1, 0}, RGB{1, 0, 0})
	if l.At(0, 0) != 1 || l.At(1, 0) != 2 {
		t.Fatalf("tiles=%v", l.Tiles)
	}
}
