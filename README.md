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
| Graphics / GPU | [gogpu](https://github.com/gogpu/gogpu) |
| GUI / Editor (future) | [Fyne](https://fyne.io/), possibly Raylib-based tooling |

More will be added as the engine takes shape.

## Using AI

Will AI help build this engine? Yes.

I really enjoy programming — actually writing code — and problem solving. But for boilerplate, docs, and other repetitive work (like this README), I'll lean on different AI models for convenience and speed. The interesting parts stay human-driven.

## Status

**Very early.** This repo is just getting started — no engine code yet, just the foundation. Expect things to move quickly (and break often) as the project grows.

## Monorepo

| Path | Description |
|------|-------------|
| [`gomakeagame.com/`](gomakeagame.com/) | Website (Go + templ + htmx + missing.css) |

## License

[MIT](LICENSE) — Copyright (c) 2026 just_Bri
