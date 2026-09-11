# Go Make A Game

**Gomag** for short — a 2D+3D capable game engine written in Go.

**Website:** [gomakeagame.com](https://gomakeagame.com)

---

## What is Gomag?

Gomag is a game engine built from the ground up in Go. It aims to support both 2D and 3D game development with a long-term vision of a full IDE and editor — think Godot, Unity, or Unreal — but code-focused for now.

This project is wide open. The architecture, APIs, and tooling are all still taking shape.

## Why?

Popular game engines never quite clicked. Hand-writing everything — or leaning heavily on bindings like Raylib — gets cumbersome as projects grow. Gomag is an attempt to build something that feels right: a proper engine with an editor and GUI down the road, without giving up the joy of working in Go.

## Planned Stack

Nothing is set in stone yet. Libraries under consideration:

| Area | Candidates |
|------|------------|
| Graphics / GPU | [gogpu](https://github.com/gogpu/gogpu) + [gg](https://github.com/gogpu/gg) via `engine/gfx` (Pure Go, `CGO_ENABLED=0`). `g3d` later for 3D. |
| GUI / Editor (future) | [gogpu/ui](https://github.com/gogpu/ui), possibly Fyne |

More will be added as the engine takes shape.

## Using AI

Will AI help build this engine? Yes.

I really enjoy programming — actually writing code — and problem solving. But for boilerplate, docs, and other repetitive work (like this README), I'll lean on different AI models for convenience and speed. The interesting parts stay human-driven.

## Status

Very early. Engine has ECS, maze-TD arena/pathing, and a first `gfx` window on gogpu+gg. Expect breakage.

## Monorepo

| Path | Description |
|------|-------------|
| [`engine/`](engine/) | Game engine (pure Go) — ECS, arena TD, gfx (gogpu+gg) |
| [`gomakeagame.com/`](gomakeagame.com/) | Website (Go + templ + htmx + missing.css) |

## License

[MIT](LICENSE) — Copyright (c) 2026 just_Bri
