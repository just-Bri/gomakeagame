package ecs

// Join2 calls fn for every entity that has both A and B.
// Iteration follows a's dense order; b is looked up by entity.
// Do not mutate either store's membership during iteration.
func Join2[A, B any](a *Store[A], b *Store[B], fn func(e Entity, a *A, b *B)) {
	for i := range a.dense {
		e := a.entities[i]
		bp, ok := b.GetMut(e)
		if !ok {
			continue
		}
		fn(e, &a.dense[i], bp)
	}
}

// Join3 calls fn for every entity that has A, B, and C.
func Join3[A, B, C any](a *Store[A], b *Store[B], c *Store[C], fn func(e Entity, a *A, b *B, c *C)) {
	for i := range a.dense {
		e := a.entities[i]
		bp, ok := b.GetMut(e)
		if !ok {
			continue
		}
		cp, ok := c.GetMut(e)
		if !ok {
			continue
		}
		fn(e, &a.dense[i], bp, cp)
	}
}
