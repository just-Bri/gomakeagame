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
| `Pointer` | Mouse position + held buttons + scroll |
| `PointerEdges` | One-shot mouse press/release latches |
| `WirePointer` | Registers pointer state on an `EventSource` |

```go
ptr := input.NewPointer()
pedges := input.NewPointerEdges()
input.WirePointer(app.EventSource(), ptr, pedges)

x, y := ptr.Pos()
if pedges.ConsumePressed(gpucontext.MouseButtonLeft) { /* click */ }
```
