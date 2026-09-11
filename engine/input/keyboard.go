// Package input tracks keyboard state from gogpu EventSource callbacks.
//
// Prefer this over polling Native().Input() — held-key polling is unreliable
// on some Wayland compositors, while press/release events stay consistent.
package input

import (
	"sync"

	"github.com/gogpu/gpucontext"
)

// Keyboard tracks which keys are currently held.
type Keyboard struct {
	mu   sync.Mutex
	down map[gpucontext.Key]bool
}

// NewKeyboard returns an empty held-key map.
func NewKeyboard() *Keyboard {
	return &Keyboard{down: make(map[gpucontext.Key]bool)}
}

// Set records whether key is currently down. Call from OnKeyPress/OnKeyRelease.
func (k *Keyboard) Set(key gpucontext.Key, pressed bool) {
	if k == nil {
		return
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	if pressed {
		k.down[key] = true
	} else {
		delete(k.down, key)
	}
}

// Pressed reports whether key is currently held.
func (k *Keyboard) Pressed(key gpucontext.Key) bool {
	if k == nil {
		return false
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.down[key]
}

// Edges latches press and release events until the sim consumes them.
// Useful for one-shot actions (pause, confirm) that must not repeat every frame.
type Edges struct {
	mu       sync.Mutex
	pressed  map[gpucontext.Key]bool
	released map[gpucontext.Key]bool
}

// NewEdges returns empty edge latches.
func NewEdges() *Edges {
	return &Edges{
		pressed:  make(map[gpucontext.Key]bool),
		released: make(map[gpucontext.Key]bool),
	}
}

// OnPress records a key-down edge. Call from EventSource OnKeyPress.
func (e *Edges) OnPress(key gpucontext.Key) {
	if e == nil {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.pressed[key] = true
}

// OnRelease records a key-up edge. Call from EventSource OnKeyRelease.
func (e *Edges) OnRelease(key gpucontext.Key) {
	if e == nil {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.released[key] = true
}

// ConsumePressed returns true once per press edge for key.
func (e *Edges) ConsumePressed(key gpucontext.Key) bool {
	if e == nil {
		return false
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.pressed[key] {
		return false
	}
	delete(e.pressed, key)
	return true
}

// ConsumeReleased returns true once per release edge for key.
func (e *Edges) ConsumeReleased(key gpucontext.Key) bool {
	if e == nil {
		return false
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.released[key] {
		return false
	}
	delete(e.released, key)
	return true
}

// Wire registers Keyboard and/or Edges on an EventSource.
// Pass nil for either side you do not need.
func Wire(es gpucontext.EventSource, keys *Keyboard, edges *Edges) {
	if es == nil {
		return
	}
	es.OnKeyPress(func(key gpucontext.Key, _ gpucontext.Modifiers) {
		if keys != nil {
			keys.Set(key, true)
		}
		if edges != nil {
			edges.OnPress(key)
		}
	})
	es.OnKeyRelease(func(key gpucontext.Key, _ gpucontext.Modifiers) {
		if keys != nil {
			keys.Set(key, false)
		}
		if edges != nil {
			edges.OnRelease(key)
		}
	})
}
