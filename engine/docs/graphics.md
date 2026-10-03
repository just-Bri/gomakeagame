# Graphics (gogpu + gg)

Last updated: 2026-09-11

Gomag’s 2D windowing/rendering sits on the [GoGPU](https://github.com/gogpu) ecosystem.

| Package | Role in Gomag |
|---------|----------------|
| `gogpu` | Window, input, lifecycle |
| `gg` + `ggcanvas` | 2D draw → present to the window |
| `wgpu` | GPU API (transitive) |
| `g3d` | Wired via `engine/gfx3d` (pin `g3d@v0.1.7` to match wgpu 0.31.x) |

Engine wrappers: `engine/gfx` (2D gg) and `engine/gfx3d` (3D g3d).

## Version pins (important)

Latest `gogpu` (v0.54+) and latest tagged `gg` (v0.52.5) currently disagree on `wgpu` / `gputypes` APIs. Gomag pins a known-good set:

| Module | Version |
|--------|---------|
| `github.com/gogpu/gg` | `v0.52.5` |
| `github.com/gogpu/gogpu` | `v0.53.0` |
| `github.com/gogpu/wgpu` | `v0.31.6` (via gogpu/gg) |
| `github.com/gogpu/gputypes` | `v0.5.2` (via gogpu/gg) |

Do not `go get -u` these blindly. Bump only after checking both READMEs / `go.mod` files agree.

## Build rule

Pure Go backend requires **CGO off**:

```bash
cd engine
CGO_ENABLED=0 mise exec -- go run ./examples/hello2d
# or from repo root:
mise run hello2d
```

On some Linux/macOS setups CGO defaults on; always set `CGO_ENABLED=0` for GoGPU apps.

Optional overrides:

```bash
GOGPU_GRAPHICS_API=vulkan   # or software, gles, …
GOGPU_GRAPHICS_API=software # no GPU / CI-friendly
```

## Usage sketch

```go
app := gfx.New(gfx.DefaultConfig())
app.OnDraw(func(dc *gg.Context, w, h int) {
    dc.ClearWithColor(gg.RGB(0.1, 0.1, 0.12))
    dc.SetRGB(0.3, 0.8, 0.5)
    dc.DrawCircle(float64(w)/2, float64(h)/2, 40)
    dc.Fill()
})
app.Run()
```

For raw input / animation tokens / multi-window, use `app.Native()` (`*gogpu.App`).
For held keys and press/release edges, prefer `engine/input` (see `docs/input.md`).

### Letterbox (fixed virtual resolution)

```go
lb := gfx.Fit(float64(w), float64(h), 1280, 800)
lb.Apply(dc)
defer dc.Pop()
// draw in virtual 1280×800 coords…
```

`ToVirtual` / `ToScreen` convert pointer coordinates across the mapping.

### 3D window (`gfx3d`)

```go
scene := g3d.NewScene()
camera := g3d.NewPerspectiveCamera(60, 16.0/9.0, 0.1, 200)
app := gfx3d.New(gfx3d.DefaultConfig())
app.SetScene(scene, camera)
app.OnUpdate(func(dt float64) { /* sim */ })
app.Run()
```

`g3d` is pinned to **v0.1.7** so it shares the Gomag wgpu 0.31 / gputypes 0.5 line with `gg`. Do not bump to g3d ≥0.1.8 without also re-validating gg.

Smoke:

```bash
mise run hello3d
```

Soft-3D / isometric placeholders (no GPU mesh): `engine/project` projects AABB boxes for 2.5D games (Forever Bound).

## Next (later)

- Sprite/tile batcher owned by Gomag (still on wgpu; gg stays for UI/debug)
- Wire `let_none_through` to `gfx` once headless sim is solid
- Composite gg HUD over gfx3d surfaces
- Revisit version bumps when gg catches up to gogpu’s wgpu line
