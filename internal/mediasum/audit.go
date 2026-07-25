package mediasum

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// AuditReport is a streaming, read-only integrity summary for one MediaSum
// split. Required means the fields ISC-RAG needs to construct canonical source
// turns: a unique document ID and equally sized utterance and speaker arrays.
// Optional document metadata is counted separately because MediaSum does not
// promise every descriptive field for every broadcast.
type AuditReport struct {
	Path                    string
	Records                 int
	Turns                   int
	RecordsMissingID        int
	DuplicateIDs            int
	MismatchedTurnArrays    int
	EmptyUtterances         int
	EmptySpeakerLabels      int
	InvalidUTF8Fields       int
	EmptyParsedSpeakerNames int
	ExplicitSpeakerRoles    int
	MissingProgram          int
	MissingDate             int
	MissingURL              int
	MissingTitle            int
	MissingSummary          int
}

// Valid reports whether every record has the required canonical-source fields.
func (r AuditReport) Valid() bool {
	return r.RecordsMissingID == 0 && r.DuplicateIDs == 0 && r.MismatchedTurnArrays == 0 && r.EmptyUtterances == 0 && r.EmptySpeakerLabels == 0 && r.InvalidUTF8Fields == 0 && r.EmptyParsedSpeakerNames == 0
}

// AuditFile validates a MediaSum JSONL split without loading it into memory or
// writing to PostgreSQL. It uses the same Record schema and ParseSpeakerLabel
// function as ingestion, so it checks the exact deterministic field path used
// by the experiment.
func AuditFile(path string) (AuditReport, error) {
	file, err := os.Open(path)
	if err != nil {
		return AuditReport{}, err
	}
	defer file.Close()

	report := AuditReport{Path: path}
	seenIDs := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		report.Records++
		var record Record
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return report, fmt.Errorf("decode record %d: %w", report.Records, err)
		}
		auditRecord(&report, seenIDs, record)
	}
	if err := scanner.Err(); err != nil {
		return report, err
	}
	return report, nil
}

func auditRecord(report *AuditReport, seenIDs map[string]struct{}, record Record) {
	if strings.TrimSpace(record.ID) == "" {
		report.RecordsMissingID++
	} else if _, exists := seenIDs[record.ID]; exists {
		report.DuplicateIDs++
	} else {
		seenIDs[record.ID] = struct{}{}
	}
	if len(record.Utterances) != len(record.Speakers) {
		report.MismatchedTurnArrays++
	}
	if strings.TrimSpace(record.Program) == "" {
		report.MissingProgram++
	}
	if strings.TrimSpace(record.Date) == "" {
		report.MissingDate++
	}
	if strings.TrimSpace(record.URL) == "" {
		report.MissingURL++
	}
	if strings.TrimSpace(record.Title) == "" {
		report.MissingTitle++
	}
	if strings.TrimSpace(record.Summary) == "" {
		report.MissingSummary++
	}

	for i, utterance := range record.Utterances {
		report.Turns++
		if !utf8.ValidString(utterance) {
			report.InvalidUTF8Fields++
		}
		if strings.TrimSpace(utterance) == "" {
			report.EmptyUtterances++
		}
		if i >= len(record.Speakers) {
			continue
		}
		label := record.Speakers[i]
		if !utf8.ValidString(label) {
			report.InvalidUTF8Fields++
		}
		if strings.TrimSpace(label) == "" {
			report.EmptySpeakerLabels++
		}
		meta := ParseSpeakerLabel(label)
		if meta.Name == "" {
			report.EmptyParsedSpeakerNames++
		}
		if meta.Role != "" {
			report.ExplicitSpeakerRoles++
		}
	}

	for _, value := range []string{record.ID, record.Program, record.Date, record.URL, record.Title, record.Summary} {
		if !utf8.ValidString(value) {
			report.InvalidUTF8Fields++
		}
	}
}

// Write writes a stable, human-readable report. Keeping presentation here lets
// CI or shell callers use AuditFile directly without parsing command output.
func (r AuditReport) Write(writer io.Writer) {
	fmt.Fprintf(writer, "MediaSum audit: %s\n", r.Path)
	fmt.Fprintf(writer, "  records: %d\n  turns: %d\n", r.Records, r.Turns)
	fmt.Fprintf(writer, "  required failures: missing_id=%d duplicate_id=%d mismatched_turn_arrays=%d empty_utterance=%d empty_speaker_label=%d invalid_utf8=%d empty_parsed_speaker_name=%d\n", r.RecordsMissingID, r.DuplicateIDs, r.MismatchedTurnArrays, r.EmptyUtterances, r.EmptySpeakerLabels, r.InvalidUTF8Fields, r.EmptyParsedSpeakerNames)
	fmt.Fprintf(writer, "  parsed explicit speaker roles: %d\n", r.ExplicitSpeakerRoles)
	fmt.Fprintf(writer, "  optional metadata missing: program=%d date=%d url=%d title=%d summary=%d\n", r.MissingProgram, r.MissingDate, r.MissingURL, r.MissingTitle, r.MissingSummary)
	if r.Valid() {
		fmt.Fprintln(writer, "  result: PASS")
	} else {
		fmt.Fprintln(writer, "  result: FAIL")
	}
}
