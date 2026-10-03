// Package combat provides Classic-style melee helpers reusable across games.
package combat

// SwingTimer models a Classic WoW auto-attack swing clock.
//
// Speed is the full swing period in seconds (weapon speed).
// Remaining counts down to the next swing; when it hits 0 a swing is ready.
type SwingTimer struct {
	Speed     float64
	Remaining float64
	// OnNextSwingBonus is flat damage added to the next swing only (Heroic Strike).
	OnNextSwingBonus float64
	// Locked out while casting / on GCD extras can bump Remaining.
}

// NewSwing returns a timer that starts mid-swing (half a period remaining).
func NewSwing(speed float64) SwingTimer {
	if speed <= 0 {
		speed = 2.0
	}
	return SwingTimer{Speed: speed, Remaining: speed * 0.5}
}

// Tick advances the clock. ready is true when a swing should resolve this frame.
func (s *SwingTimer) Tick(dt float64) (ready bool) {
	if s == nil || s.Speed <= 0 {
		return false
	}
	s.Remaining -= dt
	if s.Remaining > 0 {
		return false
	}
	s.Remaining = 0
	return true
}

// ConsumeSwing resets the clock after a successful hit or miss resolution.
// Returns any OnNextSwingBonus and clears it.
func (s *SwingTimer) ConsumeSwing() (bonus float64) {
	if s == nil {
		return 0
	}
	bonus = s.OnNextSwingBonus
	s.OnNextSwingBonus = 0
	s.Remaining = s.Speed
	if s.Remaining < 0.1 {
		s.Remaining = 0.1
	}
	return bonus
}

// QueueNextSwingBonus stacks flat bonus damage onto the next auto-attack.
func (s *SwingTimer) QueueNextSwingBonus(bonus float64) {
	if s == nil {
		return
	}
	s.OnNextSwingBonus += bonus
}

// Progress returns 0..1 how close we are to the next swing (1 = ready).
func (s SwingTimer) Progress() float64 {
	if s.Speed <= 0 {
		return 1
	}
	if s.Remaining <= 0 {
		return 1
	}
	p := 1 - s.Remaining/s.Speed
	if p < 0 {
		return 0
	}
	if p > 1 {
		return 1
	}
	return p
}
