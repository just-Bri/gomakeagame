package ecs

// World owns entity lifetimes and notifies registered stores on destroy.
type World struct {
	generations []uint32
	alive       []bool
	free        []uint32
	stores      []storeRemover
}

type storeRemover interface {
	removeEntity(Entity)
}

// NewWorld creates an empty world.
func NewWorld() *World {
	return &World{}
}

// Spawn creates a new alive entity.
func (w *World) Spawn() Entity {
	var index uint32
	if n := len(w.free); n > 0 {
		index = w.free[n-1]
		w.free = w.free[:n-1]
		w.alive[index] = true
		return makeEntity(index, w.generations[index])
	}

	// Generations start at 1 so makeEntity(0, 0) stays reserved as Nil.
	index = uint32(len(w.generations))
	w.generations = append(w.generations, 1)
	w.alive = append(w.alive, true)
	return makeEntity(index, 1)
}

// Destroy marks an entity dead, recycles its slot, and strips components.
// Destroying a stale or already-dead handle is a no-op.
func (w *World) Destroy(e Entity) {
	if !w.Alive(e) {
		return
	}

	index := e.Index()
	for _, s := range w.stores {
		s.removeEntity(e)
	}

	w.alive[index] = false
	w.generations[index]++
	w.free = append(w.free, index)
}

// Alive reports whether e still refers to a living entity in this world.
func (w *World) Alive(e Entity) bool {
	index := e.Index()
	if int(index) >= len(w.alive) {
		return false
	}
	return w.alive[index] && w.generations[index] == e.Generation()
}

// Len returns how many entities are currently alive.
func (w *World) Len() int {
	return len(w.generations) - len(w.free)
}

// Cap returns how many entity slots have been allocated (alive + free).
func (w *World) Cap() int {
	return len(w.generations)
}

// Each calls fn for every alive entity. Order is by slot index, not spawn order.
func (w *World) Each(fn func(Entity)) {
	for i, alive := range w.alive {
		if !alive {
			continue
		}
		fn(makeEntity(uint32(i), w.generations[i]))
	}
}

func (w *World) register(s storeRemover) {
	w.stores = append(w.stores, s)
}
