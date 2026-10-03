package input

import (
	"sync"

	"github.com/gogpu/gpucontext"
)

// Pointer tracks mouse / pointer position and button state.
type Pointer struct {
	mu     sync.Mutex
	x, y   float64
	down   map[gpucontext.MouseButton]bool
	scroll float64
}

// NewPointer returns empty pointer state.
func NewPointer() *Pointer {
	return &Pointer{down: make(map[gpucontext.MouseButton]bool)}
}

// SetPos records the current pointer position (logical DIP coords).
func (p *Pointer) SetPos(x, y float64) {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.x, p.y = x, y
}

// SetButton records whether button is currently down.
func (p *Pointer) SetButton(button gpucontext.MouseButton, pressed bool) {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if pressed {
		p.down[button] = true
	} else {
		delete(p.down, button)
	}
}

// AddScroll accumulates wheel delta (positive = away / up).
func (p *Pointer) AddScroll(delta float64) {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.scroll += delta
}

// Pos returns the last known pointer position.
func (p *Pointer) Pos() (x, y float64) {
	if p == nil {
		return 0, 0
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.x, p.y
}

// Pressed reports whether button is currently held.
func (p *Pointer) Pressed(button gpucontext.MouseButton) bool {
	if p == nil {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.down[button]
}

// ConsumeScroll returns and clears accumulated scroll delta.
func (p *Pointer) ConsumeScroll() float64 {
	if p == nil {
		return 0
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	d := p.scroll
	p.scroll = 0
	return d
}

// PointerEdges latches press/release until consumed.
type PointerEdges struct {
	mu       sync.Mutex
	pressed  map[gpucontext.MouseButton]bool
	released map[gpucontext.MouseButton]bool
}

// NewPointerEdges returns empty pointer edge latches.
func NewPointerEdges() *PointerEdges {
	return &PointerEdges{
		pressed:  make(map[gpucontext.MouseButton]bool),
		released: make(map[gpucontext.MouseButton]bool),
	}
}

// OnPress records a button-down edge.
func (e *PointerEdges) OnPress(button gpucontext.MouseButton) {
	if e == nil {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.pressed[button] = true
}

// OnRelease records a button-up edge.
func (e *PointerEdges) OnRelease(button gpucontext.MouseButton) {
	if e == nil {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.released[button] = true
}

// ConsumePressed returns true once per press edge for button.
func (e *PointerEdges) ConsumePressed(button gpucontext.MouseButton) bool {
	if e == nil {
		return false
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.pressed[button] {
		return false
	}
	delete(e.pressed, button)
	return true
}

// ConsumeReleased returns true once per release edge for button.
func (e *PointerEdges) ConsumeReleased(button gpucontext.MouseButton) bool {
	if e == nil {
		return false
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.released[button] {
		return false
	}
	delete(e.released, button)
	return true
}

// WirePointer registers Pointer and/or PointerEdges on an EventSource.
// Pass nil for either side you do not need.
func WirePointer(es gpucontext.EventSource, ptr *Pointer, edges *PointerEdges) {
	if es == nil {
		return
	}
	es.OnMouseMove(func(x, y float64) {
		if ptr != nil {
			ptr.SetPos(x, y)
		}
	})
	es.OnMousePress(func(button gpucontext.MouseButton, x, y float64) {
		if ptr != nil {
			ptr.SetPos(x, y)
			ptr.SetButton(button, true)
		}
		if edges != nil {
			edges.OnPress(button)
		}
	})
	es.OnMouseRelease(func(button gpucontext.MouseButton, x, y float64) {
		if ptr != nil {
			ptr.SetPos(x, y)
			ptr.SetButton(button, false)
		}
		if edges != nil {
			edges.OnRelease(button)
		}
	})
	es.OnScroll(func(_, dy float64) {
		if ptr != nil {
			// Invert so scroll-up (typically negative dy) yields a positive delta.
			ptr.AddScroll(-dy)
		}
	})
}
