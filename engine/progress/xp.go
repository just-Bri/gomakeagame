// Package progress provides Classic-flavored leveling helpers.
package progress

// ClassicTable returns XP required to go from level L to L+1 for L in [1, maxLevel).
// Tuned slow enough to feel like early Classic, not retail catch-up.
//
// maxLevel is the phase cap (inclusive). The returned slice is length maxLevel:
// index 0 unused, index L = XP to reach L+1 from L. Index maxLevel is 0 (capped).
func ClassicTable(maxLevel int) []int {
	if maxLevel < 1 {
		maxLevel = 1
	}
	t := make([]int, maxLevel+1)
	// Level 1→2 starts modest; curve steepens like early Azeroth.
	xp := 400
	for lvl := 1; lvl < maxLevel; lvl++ {
		t[lvl] = xp
		// ~1.35× per level keeps 1→8 feel like a real starter zone grind.
		xp = int(float64(xp)*1.35 + 50)
	}
	t[maxLevel] = 0
	return t
}

// AddXP applies gained XP with Classic leftover overflow into subsequent levels.
// Returns new level, new XP into current level, and how many levels gained.
func AddXP(level, xpIntoLevel, gained int, table []int, maxLevel int) (newLevel, newXP, levelsGained int) {
	if maxLevel < 1 {
		maxLevel = 1
	}
	if level < 1 {
		level = 1
	}
	if level > maxLevel {
		level = maxLevel
	}
	if gained < 0 {
		gained = 0
	}
	newLevel, newXP = level, xpIntoLevel+gained
	for newLevel < maxLevel {
		need := 0
		if newLevel >= 0 && newLevel < len(table) {
			need = table[newLevel]
		}
		if need <= 0 || newXP < need {
			break
		}
		newXP -= need
		newLevel++
		levelsGained++
	}
	if newLevel >= maxLevel {
		newLevel = maxLevel
		newXP = 0
	}
	return newLevel, newXP, levelsGained
}

// Fraction returns 0..1 progress toward the next level.
func Fraction(level, xpIntoLevel int, table []int, maxLevel int) float64 {
	if level >= maxLevel {
		return 1
	}
	if level < 0 || level >= len(table) || table[level] <= 0 {
		return 0
	}
	f := float64(xpIntoLevel) / float64(table[level])
	if f < 0 {
		return 0
	}
	if f > 1 {
		return 1
	}
	return f
}
