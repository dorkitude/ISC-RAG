package mediasum

import "testing"

func TestParseSpeakerLabel(t *testing.T) {
	tests := []struct {
		label, name, role string
	}{
		{"MADELELEINE BRAND, host", "MADELELEINE BRAND", "host"},
		{"  MICHELLE   SINGLETARY ", "MICHELLE SINGLETARY", ""},
		{"ALEX COHEN, Senior Correspondent", "ALEX COHEN", "senior correspondent"},
		{"", "", ""},
	}
	for _, test := range tests {
		got := ParseSpeakerLabel(test.label)
		if got.Name != test.name || got.Role != test.role {
			t.Errorf("ParseSpeakerLabel(%q) = %#v, want name=%q role=%q", test.label, got, test.name, test.role)
		}
	}
}
