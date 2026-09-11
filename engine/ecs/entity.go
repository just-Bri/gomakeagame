package ecs

// Entity is a generational handle into a World.
//
// The low 32 bits are the slot index; the high 32 bits are the generation.
// Destroying an entity bumps the generation so old handles fail Alive checks.
type Entity uint64

const (
	entityIndexBits = 32
	entityIndexMask = 1<<entityIndexBits - 1
)

// Index returns the slot this entity occupies in the world.
func (e Entity) Index() uint32 {
	return uint32(e & entityIndexMask)
}

// Generation returns how many times this slot has been recycled.
func (e Entity) Generation() uint32 {
	return uint32(e >> entityIndexBits)
}

func makeEntity(index, generation uint32) Entity {
	return Entity(uint64(generation)<<entityIndexBits | uint64(index))
}

// Nil is the zero entity handle. It is never alive.
// Worlds start generations at 1 so a real entity is never equal to Nil.
const Nil Entity = 0
