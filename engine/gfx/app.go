// Package gfx boots a 2D window on the GoGPU stack.
//
// Layout:
//
//	gogpu  — window, input, lifecycle
//	gg     — 2D drawing (via ggcanvas → surface)
//	wgpu   — pulled in transitively (GPU API)
//
// g3d is intentionally not wired here; add it when a 3D game needs it.
//
// Build with CGO disabled (Pure Go backend):
//
//	CGO_ENABLED=0 go run ./examples/hello2d
package gfx

import (
	"fmt"
	"log"

	"github.com/gogpu/gg"
	_ "github.com/gogpu/gg/gpu" // register GPU accelerator for gg
	"github.com/gogpu/gg/integration/ggcanvas"
	"github.com/gogpu/gogpu"
	"github.com/gogpu/gpucontext"
)

// Config configures the window. Zero values get sensible defaults.
type Config struct {
	Title  string
	Width  int
	Height int
	// Continuous keeps redrawing every frame (game loops).
	// When false, redraw only on RequestRedraw / input (power-friendly).
	Continuous bool
}

// DefaultConfig returns a 1280×720 continuous game window titled "Gomag".
func DefaultConfig() Config {
	return Config{
		Title:      "Gomag",
		Width:      1280,
		Height:     720,
		Continuous: true,
	}
}

// DrawFunc is called each frame with a gg drawing context and logical size.
type DrawFunc func(dc *gg.Context, width, height int)

// App is a thin gogpu + ggcanvas host for 2D games and tools.
type App struct {
	inner  *gogpu.App
	cfg    Config
	draw   DrawFunc
	canvas *ggcanvas.Canvas
}

// New creates an App. Call OnDraw, then Run.
func New(cfg Config) *App {
	if cfg.Title == "" {
		cfg.Title = "Gomag"
	}
	if cfg.Width <= 0 {
		cfg.Width = 1280
	}
	if cfg.Height <= 0 {
		cfg.Height = 720
	}

	gcfg := gogpu.DefaultConfig().
		WithTitle(cfg.Title).
		WithSize(cfg.Width, cfg.Height).
		WithContinuousRender(cfg.Continuous)

	return &App{
		inner: gogpu.NewApp(gcfg),
		cfg:   cfg,
	}
}

// OnDraw sets the per-frame 2D draw callback.
func (a *App) OnDraw(fn DrawFunc) *App {
	a.draw = fn
	return a
}

// Native returns the underlying gogpu.App for advanced use
// (input polling, animation tokens, multi-window, etc.).
func (a *App) Native() *gogpu.App {
	return a.inner
}

// EventSource returns the gogpu event source (keys, pointer, …).
func (a *App) EventSource() gpucontext.EventSource {
	return a.inner.EventSource()
}

// RequestRedraw schedules another frame when Continuous is false.
func (a *App) RequestRedraw() {
	a.inner.RequestRedraw()
}

// Run opens the window and blocks until it closes.
func (a *App) Run() error {
	a.inner.OnDraw(a.onDraw)
	a.inner.OnClose(func() {
		gg.CloseAccelerator()
	})
	if err := a.inner.Run(); err != nil {
		return fmt.Errorf("gfx: %w", err)
	}
	return nil
}

func (a *App) onDraw(dc *gogpu.Context) {
	w, h := dc.Width(), dc.Height()
	if w <= 0 || h <= 0 {
		return
	}

	if a.canvas == nil {
		provider := a.inner.GPUContextProvider()
		if provider == nil {
			return
		}
		canvas, err := ggcanvas.New(provider, w, h)
		if err != nil {
			log.Printf("gfx: create canvas: %v", err)
			return
		}
		a.canvas = canvas
	}

	if cw, ch := a.canvas.Size(); cw != w || ch != h {
		if err := a.canvas.Resize(w, h); err != nil {
			log.Printf("gfx: resize: %v", err)
		}
	}

	if a.draw != nil {
		if err := a.canvas.Draw(func(cc *gg.Context) {
			a.draw(cc, w, h)
		}); err != nil {
			log.Printf("gfx: draw: %v", err)
		}
	}

	if err := a.canvas.Render(dc.RenderTarget()); err != nil {
		log.Printf("gfx: present: %v", err)
	}
}
