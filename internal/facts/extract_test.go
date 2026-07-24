package facts

import "testing"

func TestDecodeStrictJSONAndEvidenceTurns(t *testing.T) {
	decoded, err := decode("```json\n{\"facts\":[{\"text\":\"Alex lives in Seattle.\",\"quote\":\"I live in Seattle\"}]}\n```")
	if err != nil || len(decoded.Facts) != 1 {
		t.Fatalf("decode: %#v, %v", decoded, err)
	}
	start, end, speaker := evidenceTurns([]canonicalTurn{{0, 0, 10, "A"}, {1, 11, 22, "B"}}, 12, 18)
	if start != 1 || end != 1 || speaker != "B" {
		t.Fatalf("wrong evidence span: %d %d %q", start, end, speaker)
	}
}
