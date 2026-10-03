package input

import (
	"testing"

	"github.com/gogpu/gpucontext"
)

func TestPointerBasics(t *testing.T) {
	p := NewPointer()
	p.SetPos(12, 34)
	x, y := p.Pos()
	if x != 12 || y != 34 {
		t.Fatalf("pos=%v,%v", x, y)
	}
	p.SetButton(gpucontext.MouseButtonLeft, true)
	if !p.Pressed(gpucontext.MouseButtonLeft) {
		t.Fatal("expected left down")
	}
	edges := NewPointerEdges()
	edges.OnPress(gpucontext.MouseButtonRight)
	if !edges.ConsumePressed(gpucontext.MouseButtonRight) {
		t.Fatal("expected right press edge")
	}
	if edges.ConsumePressed(gpucontext.MouseButtonRight) {
		t.Fatal("edge should be consumed")
	}
}
