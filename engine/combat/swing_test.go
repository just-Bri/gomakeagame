package combat

import "testing"

func TestSwingTimer(t *testing.T) {
	s := NewSwing(2.0)
	if s.Tick(0.9) {
		t.Fatal("should not be ready yet")
	}
	if !s.Tick(0.2) {
		t.Fatal("expected ready after remaining elapsed")
	}
	s.QueueNextSwingBonus(10)
	bonus := s.ConsumeSwing()
	if bonus != 10 {
		t.Fatalf("bonus=%v", bonus)
	}
	if s.Remaining != 2.0 {
		t.Fatalf("remaining=%v", s.Remaining)
	}
}

func TestTabNext(t *testing.T) {
	list := []Candidate{
		{ID: 1, X: 0, Z: 5, Alive: true, Hostile: true},
		{ID: 2, X: 5, Z: 0, Alive: true, Hostile: true},
		{ID: 3, X: -5, Z: 0, Alive: true, Hostile: true},
	}
	first := TabNext(0, 0, 0, 20, 0, list)
	if first == 0 {
		t.Fatal("expected a target")
	}
	second := TabNext(0, 0, 0, 20, first, list)
	if second == 0 || second == first {
		t.Fatalf("expected different next target, got %d after %d", second, first)
	}
}
