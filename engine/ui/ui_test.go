package ui

import "testing"

func TestHit(t *testing.T) {
	if !Hit(10, 10, 50, 20, 30, 15) {
		t.Fatal("expected hit")
	}
	if Hit(10, 10, 50, 20, 5, 15) {
		t.Fatal("expected miss")
	}
}
