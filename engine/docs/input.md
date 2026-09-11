# Input

Last updated: 2026-09-11

Package `engine/input` tracks keyboard state from gogpu `EventSource` callbacks.

Prefer this over polling `Native().Input()` — held-key polling is unreliable on
some Wayland compositors, while press/release events stay consistent.

```go
keys := input.NewKeyboard()
edges := input.NewEdges()
input.Wire(app.EventSource(), keys, edges)

// each sim tick:
if keys.Pressed(gpucontext.KeyW) { /* move */ }
if edges.ConsumeReleased(gpucontext.KeyEscape) { /* pause */ }
```

| Type | Role |
|------|------|
| `Keyboard` | Currently held keys |
| `Edges` | One-shot press/release latches until consumed |
| `Wire` | Registers both on an `EventSource` |
