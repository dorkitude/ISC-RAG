// Command rejection-rate-sample measures the validation-gate rejection rate of
// the held-out fact extraction by re-running the frozen turnwise extraction on a
// deterministic sample of held-out documents and replaying the persistDocument
// validation logic in memory. It writes NOTHING to the database: no
// pipeline_runs row, no facts, no source_units. It exists because the
// production extractor drops gate failures without counting them.
//
// Sample: first N documents of the held-out 500 in the production processing
// order (external_id ascending). Model, prompts, windowing, max-facts cap, and
// every validation predicate replicate internal/facts/extract.go as of the
// heldout500-turnwise-opinion-v1 run.
//
// Usage:
//
//	set -a; source .env; set +a
//	go run ./experiments/rejection-rate-sample-20260803 --docs 20 --workers 10
package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/dorkitude/isc-rag/internal/facts"
	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	model        = "accounts/fireworks/models/kimi-k2p6"
	inputVersion = "heldout500-v1"
	maxFacts     = 30 // min(cfg.MaxFacts=50 default, 30) in extract.go
	maxTokens    = 2400
)

var genericSpeaker = regexp.MustCompile(`(?i)\bthe (speaker|host|guest)\b`)

type unit struct{ turnStart, turnEnd int }

type candidate struct {
	Text        string `json:"text"`
	Quote       string `json:"quote"`
	targetStart int
	targetEnd   int
}
type extracted struct {
	Facts []candidate `json:"facts"`
}

type tally struct {
	Docs             int `json:"docs"`
	Calls            int `json:"calls"`
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	Candidates       int `json:"candidates"`
	EmptyOrOverlong  int `json:"rejected_empty_or_overlong_quote"`
	BadRange         int `json:"rejected_bad_target_range"`
	QuoteNotFound    int `json:"rejected_quote_not_in_target"`
	GenericSpeaker   int `json:"rejected_generic_speaker"`
	Duplicate        int `json:"deduplicated"`
	Validated        int `json:"validated"`
}

func main() {
	docCount := flag.Int("docs", 20, "held-out documents to sample")
	workers := flag.Int("workers", 10, "concurrent documents")
	flag.Parse()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()
	client, err := fireworks.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := pool.Query(ctx, `SELECT su.document_id, d.external_id, su.turn_start, su.turn_end
		FROM source_units su JOIN documents d ON d.id=su.document_id
		WHERE su.representation='fact-input' AND su.version=$1 AND d.corpus='mediasum' AND d.corpus_split='test'
		ORDER BY d.external_id, su.char_start`, inputVersion)
	if err != nil {
		log.Fatalf("select units: %v", err)
	}
	units := map[string][]unit{}
	external := map[string]string{}
	var order []string
	for rows.Next() {
		var docID, ext string
		var u unit
		if err := rows.Scan(&docID, &ext, &u.turnStart, &u.turnEnd); err != nil {
			log.Fatal(err)
		}
		if _, ok := units[docID]; !ok {
			order = append(order, docID)
			external[docID] = ext
		}
		units[docID] = append(units[docID], u)
	}
	rows.Close()
	if len(order) > *docCount {
		order = order[:*docCount]
	}
	log.Printf("sampling %d of %d held-out documents (production order)", len(order), len(units))

	var mu sync.Mutex
	total := tally{}
	perDoc := map[string]*tally{}
	jobs := make(chan string)
	var wg sync.WaitGroup
	for range *workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for docID := range jobs {
				t, err := processDocument(ctx, pool, client, docID, units[docID])
				if err != nil {
					log.Printf("document %s FAILED: %v", external[docID], err)
					continue
				}
				mu.Lock()
				total.Docs++
				total.Calls += t.Calls
				total.PromptTokens += t.PromptTokens
				total.CompletionTokens += t.CompletionTokens
				total.Candidates += t.Candidates
				total.EmptyOrOverlong += t.EmptyOrOverlong
				total.BadRange += t.BadRange
				total.QuoteNotFound += t.QuoteNotFound
				total.GenericSpeaker += t.GenericSpeaker
				total.Duplicate += t.Duplicate
				total.Validated += t.Validated
				perDoc[external[docID]] = t
				log.Printf("%s: %d candidates -> %d validated (%d quote-miss, %d dupe)", external[docID], t.Candidates, t.Validated, t.QuoteNotFound, t.Duplicate)
				mu.Unlock()
			}
		}()
	}
	for _, docID := range order {
		jobs <- docID
	}
	close(jobs)
	wg.Wait()

	gateRejects := total.EmptyOrOverlong + total.BadRange + total.QuoteNotFound + total.GenericSpeaker
	out := map[string]any{
		"totals":  total,
		"per_doc": perDoc,
		"gate_rejection_rate_excl_dedup": rate(gateRejects, total.Candidates),
		"overall_drop_rate_incl_dedup":   rate(gateRejects+total.Duplicate, total.Candidates),
	}
	blob, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(blob))
}

func rate(n, d int) string {
	if d == 0 {
		return "n/a"
	}
	return fmt.Sprintf("%.2f%%", 100*float64(n)/float64(d))
}

type turn struct {
	ordinal, start, end int
	speaker, text       string
}

func processDocument(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, docID string, us []unit) (*tally, error) {
	t := &tally{}
	var all []candidate
	for _, u := range us {
		targets, err := targetedInputs(ctx, pool, docID, u.turnStart, u.turnEnd)
		if err != nil {
			return nil, err
		}
		for _, target := range targets {
			response, err := client.ChatJSONNoReasoning(ctx, model, facts.SystemPrompt, target.prompt+"\n\n"+fmt.Sprintf("Return at most %d claims; do not omit a claim for importance.", maxFacts), maxTokens)
			if err != nil {
				return nil, fmt.Errorf("turn %d: %w", target.ordinal, err)
			}
			t.Calls++
			t.PromptTokens += response.PromptTokens
			t.CompletionTokens += response.CompletionTokens
			parsed, err := decode(response.Text)
			if err != nil {
				return nil, fmt.Errorf("parse turn %d: %w", target.ordinal, err)
			}
			for i := range parsed.Facts {
				parsed.Facts[i].targetStart, parsed.Facts[i].targetEnd = target.charStart, target.charEnd
			}
			all = append(all, parsed.Facts...)
		}
	}

	transcript, turns, err := documentText(ctx, pool, docID)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for _, c := range all {
		t.Candidates++
		text := strings.TrimSpace(c.Text)
		quote := strings.TrimSpace(c.Quote)
		if text == "" || quote == "" || len(quote) > 240 {
			t.EmptyOrOverlong++
			continue
		}
		if c.targetStart < 0 || c.targetEnd > len(transcript) || c.targetStart >= c.targetEnd {
			t.BadRange++
			continue
		}
		local := strings.Index(transcript[c.targetStart:c.targetEnd], quote)
		if local < 0 {
			t.QuoteNotFound++
			continue
		}
		at := c.targetStart + local
		_, _, speaker := evidenceTurns(turns, at, at+len(quote))
		text = replaceGenericSpeaker(text, speaker)
		if genericSpeaker.MatchString(text) {
			t.GenericSpeaker++
			continue
		}
		key := hash(text)
		if seen[key] {
			t.Duplicate++
			continue
		}
		seen[key] = true
		t.Validated++
	}
	return t, nil
}

type targetInput struct {
	ordinal, charStart, charEnd int
	prompt                      string
}

func targetedInputs(ctx context.Context, pool *pgxpool.Pool, docID string, turnStart, turnEnd int) ([]targetInput, error) {
	rows, err := pool.Query(ctx, `SELECT ordinal,COALESCE(speaker_name,''),COALESCE(speaker_label,''),text,char_start,char_end FROM turns WHERE document_id=$1 AND ordinal BETWEEN GREATEST(0,$2-2) AND $3 ORDER BY ordinal`, docID, turnStart, turnEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var turns []turn
	for rows.Next() {
		var t turn
		var name, label string
		if err := rows.Scan(&t.ordinal, &name, &label, &t.text, &t.start, &t.end); err != nil {
			return nil, err
		}
		if name == "" {
			name = speakerName(label)
		}
		t.speaker = name
		turns = append(turns, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var out []targetInput
	for i, target := range turns {
		if target.ordinal < turnStart {
			continue
		}
		var b strings.Builder
		if i > 0 {
			b.WriteString("<CONTEXT_TURNS>\n")
			for _, context := range turns[:i] {
				fmt.Fprintf(&b, "[SPEAKER: %s]\n%s\n", context.speaker, context.text)
			}
			b.WriteString("</CONTEXT_TURNS>\n")
		}
		fmt.Fprintf(&b, "<TARGET_TURN>\n[SPEAKER: %s]\n%s\n</TARGET_TURN>", target.speaker, target.text)
		out = append(out, targetInput{ordinal: target.ordinal, charStart: target.start, charEnd: target.end, prompt: b.String()})
	}
	return out, nil
}

type canonicalTurn struct {
	ordinal, start, end int
	speaker             string
}

func documentText(ctx context.Context, pool *pgxpool.Pool, id string) (string, []canonicalTurn, error) {
	rows, err := pool.Query(ctx, `SELECT ordinal,text,char_start,char_end,speaker_label FROM turns WHERE document_id=$1 ORDER BY ordinal`, id)
	if err != nil {
		return "", nil, err
	}
	defer rows.Close()
	var pieces []string
	var ts []canonicalTurn
	for rows.Next() {
		var t canonicalTurn
		var text string
		if err := rows.Scan(&t.ordinal, &text, &t.start, &t.end, &t.speaker); err != nil {
			return "", nil, err
		}
		pieces = append(pieces, text)
		ts = append(ts, t)
	}
	return strings.Join(pieces, "\n"), ts, rows.Err()
}

func evidenceTurns(ts []canonicalTurn, start, end int) (int, int, string) {
	if len(ts) == 0 {
		return 0, 0, ""
	}
	first, last := ts[0], ts[len(ts)-1]
	for _, t := range ts {
		if start >= t.start && start <= t.end {
			first = t
		}
		if end >= t.start && end <= t.end+1 {
			last = t
			break
		}
	}
	return first.ordinal, last.ordinal, first.speaker
}

func speakerName(label string) string {
	name, _, _ := strings.Cut(strings.Join(strings.Fields(label), " "), ",")
	return strings.TrimSpace(name)
}

func replaceGenericSpeaker(text, label string) string {
	name := speakerName(label)
	if name == "" {
		return text
	}
	return genericSpeaker.ReplaceAllString(text, name)
}

func cleanJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)
	if start, end := strings.Index(raw, "{"), strings.LastIndex(raw, "}"); start >= 0 && end > start {
		return raw[start : end+1]
	}
	return raw
}

func decode(raw string) (extracted, error) {
	var out extracted
	if err := json.Unmarshal([]byte(cleanJSON(raw)), &out); err != nil {
		return out, err
	}
	return out, nil
}

func hash(s string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(s)))
	return hex.EncodeToString(sum[:])
}
