package ecs_test

import (
	"testing"

	"github.com/just-Bri/gomakeagame/engine/ecs"
)

type Transform struct {
	X, Y float64
}

// Enemy is intentionally fat — health, movement, and rewards in one blob.
type Enemy struct {
	Health    float64
	MaxHealth float64
	Speed     float64
	PathT     float64
	Reward    int
	Armor     float64
}

type Projectile struct {
	Damage   float64
	Speed    float64
	Target   ecs.Entity
	Lifetime float64
}

func TestSpawnDestroyAlive(t *testing.T) {
	w := ecs.NewWorld()
	e := w.Spawn()
	if e == ecs.Nil {
		t.Fatal("spawned entity must not equal Nil")
	}
	if !w.Alive(e) {
		t.Fatal("fresh entity should be alive")
	}
	if w.Alive(ecs.Nil) {
		t.Fatal("Nil must never be alive")
	}
	if w.Len() != 1 {
		t.Fatalf("Len = %d, want 1", w.Len())
	}

	w.Destroy(e)
	if w.Alive(e) {
		t.Fatal("destroyed entity should be dead")
	}
	if w.Len() != 0 {
		t.Fatalf("Len = %d, want 0", w.Len())
	}
}

func TestGenerationalReuse(t *testing.T) {
	w := ecs.NewWorld()
	old := w.Spawn()
	w.Destroy(old)

	fresh := w.Spawn()
	if old.Index() != fresh.Index() {
		t.Fatalf("expected slot reuse, old=%d fresh=%d", old.Index(), fresh.Index())
	}
	if old.Generation() == fresh.Generation() {
		t.Fatal("generation should bump on reuse")
	}
	if w.Alive(old) {
		t.Fatal("stale handle must not be alive")
	}
	if !w.Alive(fresh) {
		t.Fatal("fresh handle should be alive")
	}
}

func TestStoreSetGetRemove(t *testing.T) {
	w := ecs.NewWorld()
	enemies := ecs.NewStore[Enemy](w)

	e := w.Spawn()
	if !enemies.Set(e, Enemy{Health: 100, MaxHealth: 100, Speed: 40, Reward: 15}) {
		t.Fatal("Set should succeed for alive entity")
	}

	got, ok := enemies.Get(e)
	if !ok || got.Health != 100 || got.Reward != 15 {
		t.Fatalf("Get = %+v, ok=%v", got, ok)
	}

	mut, ok := enemies.GetMut(e)
	if !ok {
		t.Fatal("GetMut missing")
	}
	mut.Health = 50

	got, _ = enemies.Get(e)
	if got.Health != 50 {
		t.Fatalf("Health = %v, want 50", got.Health)
	}

	enemies.Remove(e)
	if enemies.Has(e) {
		t.Fatal("component should be gone after Remove")
	}
}

func TestDestroyStripsComponents(t *testing.T) {
	w := ecs.NewWorld()
	transforms := ecs.NewStore[Transform](w)
	enemies := ecs.NewStore[Enemy](w)

	e := w.Spawn()
	transforms.Set(e, Transform{X: 1, Y: 2})
	enemies.Set(e, Enemy{Health: 10})

	w.Destroy(e)
	if transforms.Has(e) || enemies.Has(e) {
		t.Fatal("destroy should strip all registered stores")
	}
	if transforms.Len() != 0 || enemies.Len() != 0 {
		t.Fatal("stores should be empty after destroy")
	}
}

func TestStaleHandleIgnored(t *testing.T) {
	w := ecs.NewWorld()
	enemies := ecs.NewStore[Enemy](w)

	old := w.Spawn()
	enemies.Set(old, Enemy{Health: 1})
	w.Destroy(old)

	fresh := w.Spawn()
	enemies.Set(fresh, Enemy{Health: 99})

	if enemies.Has(old) {
		t.Fatal("stale handle must not see recycled slot component")
	}
	got, ok := enemies.Get(fresh)
	if !ok || got.Health != 99 {
		t.Fatalf("fresh Get = %+v ok=%v", got, ok)
	}
}

func TestSetOnDeadIsNoop(t *testing.T) {
	w := ecs.NewWorld()
	enemies := ecs.NewStore[Enemy](w)
	e := w.Spawn()
	w.Destroy(e)

	if enemies.Set(e, Enemy{Health: 1}) {
		t.Fatal("Set on dead entity should fail")
	}
	if enemies.Len() != 0 {
		t.Fatal("store should stay empty")
	}
}

func TestEachAndJoin(t *testing.T) {
	w := ecs.NewWorld()
	transforms := ecs.NewStore[Transform](w)
	enemies := ecs.NewStore[Enemy](w)
	projectiles := ecs.NewStore[Projectile](w)

	a := w.Spawn()
	b := w.Spawn()
	c := w.Spawn()

	transforms.Set(a, Transform{X: 1})
	enemies.Set(a, Enemy{Health: 10, Speed: 5})

	transforms.Set(b, Transform{X: 2})
	enemies.Set(b, Enemy{Health: 20, Speed: 5})

	transforms.Set(c, Transform{X: 3})
	projectiles.Set(c, Projectile{Damage: 5}) // no enemy

	joined := 0
	ecs.Join2(enemies, transforms, func(e ecs.Entity, enemy *Enemy, tr *Transform) {
		joined++
		tr.X += enemy.Speed
	})
	if joined != 2 {
		t.Fatalf("Join2 count = %d, want 2", joined)
	}

	ta, _ := transforms.Get(a)
	tb, _ := transforms.Get(b)
	tc, _ := transforms.Get(c)
	if ta.X != 6 || tb.X != 7 || tc.X != 3 {
		t.Fatalf("transforms after join: a=%v b=%v c=%v", ta.X, tb.X, tc.X)
	}
}

func TestSwapRemoveCompacts(t *testing.T) {
	w := ecs.NewWorld()
	enemies := ecs.NewStore[Enemy](w)

	e0 := w.Spawn()
	e1 := w.Spawn()
	e2 := w.Spawn()
	enemies.Set(e0, Enemy{Health: 1})
	enemies.Set(e1, Enemy{Health: 2})
	enemies.Set(e2, Enemy{Health: 3})

	enemies.Remove(e0)
	if enemies.Len() != 2 {
		t.Fatalf("Len = %d, want 2", enemies.Len())
	}
	if !enemies.Has(e1) || !enemies.Has(e2) {
		t.Fatal("remaining entities should still have components")
	}

	seen := map[float64]bool{}
	enemies.Each(func(_ ecs.Entity, enemy *Enemy) {
		seen[enemy.Health] = true
	})
	if !seen[2] || !seen[3] || seen[1] {
		t.Fatalf("unexpected dense contents: %v", seen)
	}
}

func TestWorldEach(t *testing.T) {
	w := ecs.NewWorld()
	a := w.Spawn()
	b := w.Spawn()
	w.Destroy(a)
	c := w.Spawn()

	var got []ecs.Entity
	w.Each(func(e ecs.Entity) { got = append(got, e) })
	if len(got) != 2 {
		t.Fatalf("Each len = %d, want 2", len(got))
	}
	alive := map[ecs.Entity]bool{b: true, c: true}
	for _, e := range got {
		if !alive[e] {
			t.Fatalf("unexpected entity %v", e)
		}
	}
}
