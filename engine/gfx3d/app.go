// Package gfx3d boots a 3D window on the GoGPU + g3d stack.
//
// Layout:
//
//	gogpu — window, input, lifecycle
//	g3d  — scene graph, cameras, lights, forward renderer
//	wgpu — pulled in transitively (GPU API)
//
// Pin g3d to a release that agrees with Gomag's gogpu/gg wgpu line
// (see docs/graphics.md). Build with CGO disabled:
//
//	CGO_ENABLED=0 go run ./examples/hello3d
package gfx3d

import (
	"fmt"
	"log"
	"time"

	"github.com/gogpu/g3d"
	"github.com/gogpu/gogpu"
	"github.com/gogpu/gpucontext"
)

// Config configures the window. Zero values get sensible defaults.
type Config struct {
	Title  string
	Width  int
	Height int
	// Continuous keeps redrawing every frame (game loops).
	Continuous bool
}

// DefaultConfig returns a 1280×720 continuous game window titled "Gomag 3D".
func DefaultConfig() Config {
	return Config{
		Title:      "Gomag 3D",
		Width:      1280,
		Height:     720,
		Continuous: true,
	}
}

// UpdateFunc is called each frame before draw with delta time in seconds.
type UpdateFunc func(dt float64)

// DrawFunc renders the 3D scene. scene/camera may be mutated beforehand in Update.
type DrawFunc func(scene *g3d.Scene, camera g3d.Camera, w, h int)

// App is a thin gogpu + g3d host for 3D games and tools.
type App struct {
	inner    *gogpu.App
	cfg      Config
	update   UpdateFunc
	draw     DrawFunc
	scene    *g3d.Scene
	camera   g3d.Camera
	renderer *g3d.Renderer
	last     time.Time
}

// New creates an App. Call SetScene, OnUpdate/OnDraw, then Run.
func New(cfg Config) *App {
	if cfg.Title == "" {
		cfg.Title = "Gomag 3D"
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

// SetScene installs the scene graph and camera used by the renderer.
func (a *App) SetScene(scene *g3d.Scene, camera g3d.Camera) *App {
	a.scene = scene
	a.camera = camera
	return a
}

// OnUpdate sets the per-frame simulation callback (dt in seconds).
func (a *App) OnUpdate(fn UpdateFunc) *App {
	a.update = fn
	return a
}

// OnDraw sets an optional hook called each frame after update, before Render.
// Use this to sync meshes from game state. Pass nil to skip.
func (a *App) OnDraw(fn DrawFunc) *App {
	a.draw = fn
	return a
}

// Scene returns the active scene.
func (a *App) Scene() *g3d.Scene { return a.scene }

// Camera returns the active camera.
func (a *App) Camera() g3d.Camera { return a.camera }

// Native returns the underlying gogpu.App for advanced use.
func (a *App) Native() *gogpu.App { return a.inner }

// EventSource returns the gogpu event source (keys, pointer, …).
func (a *App) EventSource() gpucontext.EventSource {
	return a.inner.EventSource()
}

// RequestRedraw schedules another frame when Continuous is false.
func (a *App) RequestRedraw() { a.inner.RequestRedraw() }

// Run opens the window and blocks until it closes.
func (a *App) Run() error {
	if a.scene == nil || a.camera == nil {
		return fmt.Errorf("gfx3d: SetScene required before Run")
	}

	a.last = time.Now()
	a.inner.OnUpdate(func(dt float64) {
		// gogpu may pass 0 on the first tick; fall back to wall clock.
		if dt <= 0 {
			now := time.Now()
			dt = now.Sub(a.last).Seconds()
			a.last = now
		} else {
			a.last = time.Now()
		}
		if dt > 0.1 {
			dt = 0.1
		}
		if a.update != nil {
			a.update(dt)
		}
	})
	a.inner.OnResize(func(width, height int) {
		if height <= 0 {
			return
		}
		if cam, ok := a.camera.(*g3d.PerspectiveCamera); ok {
			cam.SetAspect(float32(width) / float32(height))
		}
	})
	a.inner.OnDraw(a.onDraw)
	a.inner.OnClose(func() {
		if a.renderer != nil {
			a.renderer.Release()
			a.renderer = nil
		}
	})
	if err := a.inner.Run(); err != nil {
		return fmt.Errorf("gfx3d: %w", err)
	}
	return nil
}

func (a *App) onDraw(dc *gogpu.Context) {
	w, h := dc.Width(), dc.Height()
	if w <= 0 || h <= 0 {
		return
	}

	if a.renderer == nil {
		provider := a.inner.GPUContextProvider()
		if provider == nil {
			return
		}
		r, err := g3d.NewRenderer(provider)
		if err != nil {
			log.Printf("gfx3d: create renderer: %v", err)
			return
		}
		a.renderer = r
	}

	fbW, fbH := dc.FramebufferSize()
	a.renderer.SetSize(uint32(fbW), uint32(fbH))

	if a.draw != nil {
		a.draw(a.scene, a.camera, w, h)
	}

	view := dc.SurfaceView()
	if view == nil {
		return
	}
	if err := a.renderer.Render(a.scene, a.camera, view); err != nil {
		log.Printf("gfx3d: render: %v", err)
	}
}
