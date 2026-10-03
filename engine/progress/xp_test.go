package progress

import "testing"

func TestAddXPLevels(t *testing.T) {
	table := ClassicTable(8)
	lvl, xp, gained := AddXP(1, 0, table[1], table, 8)
	if lvl != 2 || xp != 0 || gained != 1 {
		t.Fatalf("got level=%d xp=%d gained=%d", lvl, xp, gained)
	}
	lvl, xp, gained = AddXP(7, table[7]-10, 100000, table, 8)
	if lvl != 8 || xp != 0 {
		t.Fatalf("cap failed: level=%d xp=%d gained=%d", lvl, xp, gained)
	}
}
