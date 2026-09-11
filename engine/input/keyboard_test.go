package input_test

import (
	"testing"

	"github.com/gogpu/gpucontext"
	"github.com/just-Bri/gomakeagame/engine/input"
)

func TestKeyboardHeld(t *testing.T) {
	k := input.NewKeyboard()
	if k.Pressed(gpucontext.KeyW) {
		t.Fatal("expected up")
	}
	k.Set(gpucontext.KeyW, true)
	if !k.Pressed(gpucontext.KeyW) {
		t.Fatal("expected down")
	}
	k.Set(gpucontext.KeyW, false)
	if k.Pressed(gpucontext.KeyW) {
		t.Fatal("expected up after release")
	}
}

func TestKeyboardNilSafe(t *testing.T) {
	var k *input.Keyboard
	k.Set(gpucontext.KeyA, true)
	if k.Pressed(gpucontext.KeyA) {
		t.Fatal("nil keyboard should report not pressed")
	}
}

func TestEdgesConsumeOnce(t *testing.T) {
	e := input.NewEdges()
	e.OnPress(gpucontext.KeyEscape)
	e.OnRelease(gpucontext.KeyEscape)

	if !e.ConsumePressed(gpucontext.KeyEscape) {
		t.Fatal("expected press edge")
	}
	if e.ConsumePressed(gpucontext.KeyEscape) {
		t.Fatal("press edge should be consumed")
	}
	if !e.ConsumeReleased(gpucontext.KeyEscape) {
		t.Fatal("expected release edge")
	}
	if e.ConsumeReleased(gpucontext.KeyEscape) {
		t.Fatal("release edge should be consumed")
	}
}

func TestEdgesNilSafe(t *testing.T) {
	var e *input.Edges
	e.OnPress(gpucontext.KeyH)
	e.OnRelease(gpucontext.KeyH)
	if e.ConsumePressed(gpucontext.KeyH) || e.ConsumeReleased(gpucontext.KeyH) {
		t.Fatal("nil edges should never consume")
	}
}
