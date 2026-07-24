package representation

import "testing"

func TestBuildersPreserveTurnsAndBounds(t *testing.T) {
	turns := []turn{{0, "alpha speaks about rivers", 0, 24}, {1, "beta changes the subject", 25, 49}, {2, "alpha concludes clearly", 50, 72}}
	for name, units := range map[string][]unit{
		"fixed":    fixedUnits(turns, 35, 8),
		"turn":     turnUnits(turns, 35),
		"semantic": semanticUnits(turns, 35, 15),
	} {
		if len(units) == 0 {
			t.Fatalf("%s returned no units", name)
		}
		for _, u := range units {
			if u.Text == "" || u.TurnStart > u.TurnEnd || u.CharStart > u.CharEnd {
				t.Fatalf("%s produced invalid unit: %#v", name, u)
			}
		}
	}
}

func TestTurnBuilderSplitsOversizedTurn(t *testing.T) {
	units := turnUnits([]turn{{0, "abcdefghijklmnopqrstuvwxyz", 0, 26}}, 10)
	if len(units) != 3 || units[0].TurnStart != 0 || units[2].TurnEnd != 0 {
		t.Fatalf("unexpected oversized-turn split: %#v", units)
	}
}
