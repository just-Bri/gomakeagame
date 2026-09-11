package ecs

// Store is a sparse-set of fat component values of type T.
//
// Components live in a packed dense slice for iteration locality.
// Missing components are absent from the set — there is no zero-value sentinel.
//
// Prefer registering stores with NewStore so World.Destroy removes components.
type Store[T any] struct {
	world    *World
	sparse   []int // entity index -> dense index, -1 if absent
	dense    []T
	entities []Entity
}

// NewStore creates a component store bound to w.
// Destroying an entity in w automatically removes its T component.
func NewStore[T any](w *World) *Store[T] {
	s := &Store[T]{world: w}
	w.register(s)
	return s
}

// Has reports whether e currently has a T component.
func (s *Store[T]) Has(e Entity) bool {
	return s.denseIndex(e) >= 0
}

// Get returns a copy of the component and whether it exists.
func (s *Store[T]) Get(e Entity) (T, bool) {
	i := s.denseIndex(e)
	if i < 0 {
		var zero T
		return zero, false
	}
	return s.dense[i], true
}

// GetMut returns a pointer into dense storage for in-place mutation.
// The pointer is invalid after any Set/Remove/Clear on this store,
// or after World.Destroy of any entity that touches this store.
func (s *Store[T]) GetMut(e Entity) (*T, bool) {
	i := s.denseIndex(e)
	if i < 0 {
		return nil, false
	}
	return &s.dense[i], true
}

// Set inserts or overwrites the component for e.
// If e is not alive in the bound world, Set is a no-op and returns false.
func (s *Store[T]) Set(e Entity, c T) bool {
	if s.world != nil && !s.world.Alive(e) {
		return false
	}

	index := e.Index()
	s.ensureSparse(index)

	if di := s.sparse[index]; di >= 0 {
		s.dense[di] = c
		s.entities[di] = e
		return true
	}

	s.sparse[index] = len(s.dense)
	s.dense = append(s.dense, c)
	s.entities = append(s.entities, e)
	return true
}

// Remove deletes the component for e if present.
func (s *Store[T]) Remove(e Entity) {
	s.removeEntity(e)
}

// Len returns how many components are stored.
func (s *Store[T]) Len() int {
	return len(s.dense)
}

// Clear removes every component without affecting entity lifetimes.
func (s *Store[T]) Clear() {
	for i := range s.sparse {
		s.sparse[i] = -1
	}
	s.dense = s.dense[:0]
	s.entities = s.entities[:0]
}

// Each iterates dense storage. Mutating through c is allowed;
// do not Set/Remove/Clear this store during iteration.
func (s *Store[T]) Each(fn func(e Entity, c *T)) {
	for i := range s.dense {
		fn(s.entities[i], &s.dense[i])
	}
}

func (s *Store[T]) denseIndex(e Entity) int {
	index := e.Index()
	if int(index) >= len(s.sparse) {
		return -1
	}
	di := s.sparse[index]
	if di < 0 {
		return -1
	}
	// Stale handles: slot reused or generation mismatch.
	if s.entities[di] != e {
		return -1
	}
	return di
}

func (s *Store[T]) ensureSparse(index uint32) {
	for uint32(len(s.sparse)) <= index {
		s.sparse = append(s.sparse, -1)
	}
}

func (s *Store[T]) removeEntity(e Entity) {
	di := s.denseIndex(e)
	if di < 0 {
		return
	}

	last := len(s.dense) - 1
	lastEntity := s.entities[last]

	s.dense[di] = s.dense[last]
	s.entities[di] = lastEntity
	s.sparse[lastEntity.Index()] = di

	s.sparse[e.Index()] = -1
	s.dense = s.dense[:last]
	s.entities = s.entities[:last]
}
