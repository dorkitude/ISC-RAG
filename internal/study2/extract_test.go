package study2

import "testing"

func TestPackCarryover(t *testing.T) {
	ts := []turn{{1, 0, 3, "a", "aaa"}, {2, 4, 7, "b", "bbb"}, {3, 8, 11, "a", "ccc"}}
	got := pack(ts, 6, 2)
	if len(got) != 2 || got[1][0].ordinal != 1 || got[1][1].ordinal != 2 || got[1][2].ordinal != 3 {
		t.Fatalf("unexpected packs: %#v", got)
	}
}

func TestTurnWindowsExcludeQuestions(t *testing.T) {
	ts := []turn{{1, 0, 3, "a", "Who?"}, {2, 4, 7, "b", "A fact."}, {3, 8, 12, "a", "Thanks."}, {4, 13, 18, "b", "Another fact."}}
	got := turnWindows(ts)
	if len(got) != 2 || got[0][len(got[0])-1].ordinal != 2 || got[1][len(got[1])-1].ordinal != 4 {
		t.Fatalf("unexpected turnwise windows: %#v", got)
	}
}

func TestDialogueUnitsPairOnlyImmediateDifferentSpeakerAnswer(t *testing.T) {
	ts := []turn{
		{ordinal: 0, speaker: "host", text: "What happened?"},
		{ordinal: 1, speaker: "guest", text: "The bridge closed."},
		{ordinal: 2, speaker: "host", text: "Thanks."},
		{ordinal: 3, speaker: "guest", text: "Repairs start Monday."},
		{ordinal: 4, speaker: "host", text: "And why?"},
		{ordinal: 5, speaker: "host", text: "I will answer later."},
	}
	got := dialogueUnits(ts)
	if len(got) != 3 || !got[0].Paired || got[0].Question.ordinal != 0 || got[0].Answer.ordinal != 1 || got[1].Answer.ordinal != 3 || got[2].Answer.ordinal != 5 {
		t.Fatalf("unexpected dialogue units: %#v", got)
	}
}

func TestValidateDialogueStoresQuestionAndAnswerEvidence(t *testing.T) {
	q := turn{ordinal: 0, start: 0, end: 14, speaker: "host", text: "Which bridge?"}
	a := turn{ordinal: 1, start: 15, end: 33, speaker: "guest", text: "It closed Monday."}
	units := []dialogueUnit{{ID: 0, Question: &q, Answer: a, Paired: true}}
	got := validateDialogue([]dialogueCandidate{{UnitID: 0, Text: "The guest answered that the bridge closed Monday.", AnswerQuote: "It closed Monday.", QuestionQuote: "Which bridge?"}}, q.text+"\n"+a.text, units)
	if len(got) != 1 || got[0].speaker != "guest" || len(got[0].evidence) != 2 || got[0].evidence[0].speaker != "guest" || got[0].evidence[1].speaker != "host" {
		t.Fatalf("unexpected validated facts: %#v", got)
	}
}

func TestDialoguePredicateValidationRestrictsAnswerQuoteToItsUnit(t *testing.T) {
	q := turn{ordinal: 0, start: 0, end: 10, speaker: "host", text: "What changed?"}
	a := turn{ordinal: 1, start: 11, end: 50, speaker: "guest", text: "The bridge closed, and repairs start Monday."}
	d := dialogueUnit{ID: 0, Question: &q, Answer: a, Paired: true}
	units := dialoguePredicateUnits([]dialogueUnit{d})
	if len(units) != 2 {
		t.Fatalf("expected two predicate units, got %#v", units)
	}
	valid := validateDialoguePredicates([]dialogueCandidate{{UnitID: units[0].ID, Text: "The guest answered that the bridge closed.", AnswerQuote: "The bridge closed", QuestionQuote: "What changed?"}}, q.text+"\n"+a.text, units)
	if len(valid) != 1 || len(valid[0].evidence) != 2 {
		t.Fatalf("unexpected valid predicate facts: %#v", valid)
	}
	invalid := validateDialoguePredicates([]dialogueCandidate{{UnitID: units[0].ID, Text: "The guest answered that repairs start Monday.", AnswerQuote: "repairs start Monday", QuestionQuote: "What changed?"}}, q.text+"\n"+a.text, units)
	if len(invalid) != 0 {
		t.Fatalf("quote from a different predicate unit should be rejected: %#v", invalid)
	}
}
