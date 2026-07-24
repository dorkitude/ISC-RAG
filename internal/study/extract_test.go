package study

import "testing"

func TestValidateRequiresExactQuoteAndDeduplicatesText(t *testing.T) {
	transcript := "First source turn\nSecond source turn"
	turns := []turn{{ordinal: 0, start: 0, end: 17, speaker: "A"}, {ordinal: 1, start: 18, end: 36, speaker: "B"}}
	got := validate([]Candidate{
		{Text: "A made a first claim.", Quote: "First source turn"},
		{Text: "A made a first claim.", Quote: "First source turn"},
		{Text: "Unsupported claim.", Quote: "not in source"},
	}, transcript, turns)
	if len(got) != 1 {
		t.Fatalf("validated %d facts, want 1", len(got))
	}
	if got[0].turnStart != 0 || got[0].turnEnd != 0 || got[0].speaker != "A" {
		t.Fatalf("wrong evidence: %#v", got[0])
	}
}
