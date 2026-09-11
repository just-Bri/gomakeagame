// Package ecs is a hand-rolled Entity Component System for Gomag.
//
// Design goals:
//   - Pure Go — no CGO, no reflection in the hot path
//   - Fat value components stored densely (sparse-set)
//   - Generational entity IDs so stale handles are cheap to detect
//   - Systems are just functions; this package owns data, not game logic
//
// Components are plain structs. Prefer fewer, fatter types (e.g. one Enemy
// blob with health, path progress, and rewards) over micro-components.
// That fits tower-defense and most 2D games, and keeps Go code readable.
//
// Typical usage:
//
//	w := ecs.NewWorld()
//	enemies := ecs.NewStore[Enemy](w)
//	transforms := ecs.NewStore[Transform](w)
//
//	e := w.Spawn()
//	enemies.Set(e, Enemy{Health: 100, Speed: 40})
//	transforms.Set(e, Transform{X: 10, Y: 20})
//
//	enemies.Each(func(e ecs.Entity, enemy *Enemy) {
//	    // mutate in place
//	})
package ecs
