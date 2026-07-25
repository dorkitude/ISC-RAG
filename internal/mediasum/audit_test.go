package mediasum

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAuditFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sample.json")
	content := "{\"id\":\"NPR-1\",\"program\":\"Show\",\"date\":\"2020-01-01\",\"url\":\"https://example.test\",\"title\":\"Title\",\"summary\":\"Summary\",\"utt\":[\"Hello\",\"Hi\"],\"speaker\":[\"ALEX, host\",\"PAT\"]}\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	report, err := AuditFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !report.Valid() {
		t.Fatalf("expected valid audit, got %#v", report)
	}
	if report.Records != 1 || report.Turns != 2 || report.ExplicitSpeakerRoles != 1 {
		t.Fatalf("unexpected audit report: %#v", report)
	}
}
