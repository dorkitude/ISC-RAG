// Package study runs isolated, provenance-validated fact-rephrasing studies.
package study

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Namespace, Version, ExternalID, Model, Prompt, OutputDir string
	Approach                                                 string
	Iteration, MaxFacts                                      int
}

type Candidate struct {
	Text  string `json:"text"`
	Quote string `json:"quote"`
}
type response struct {
	Facts []Candidate `json:"facts"`
}
type result struct {
	RunID                     string
	RawCount, ValidatedCount  int
	PromptHash, RawOutputHash string
}
type turn struct {
	ordinal, start, end int
	speaker, text       string
}

// Run extracts one document only. It is deliberately separate from facts.Extract
// so a study cannot delete or supersede the primary fact condition.
func Run(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, cfg Config) (result, error) {
	if cfg.Namespace == "" || cfg.Version == "" || cfg.ExternalID == "" || cfg.Model == "" || cfg.Prompt == "" || cfg.Iteration < 1 {
		return result{}, fmt.Errorf("namespace, version, external-id, model, prompt, and positive iteration are required")
	}
	if cfg.MaxFacts <= 0 {
		cfg.MaxFacts = 60
	}
	var documentID string
	if err := pool.QueryRow(ctx, `SELECT id FROM documents WHERE corpus='mediasum' AND external_id=$1`, cfg.ExternalID).Scan(&documentID); err != nil {
		return result{}, fmt.Errorf("find MediaSum document %q: %w", cfg.ExternalID, err)
	}
	transcript, labeled, turns, err := source(ctx, pool, documentID)
	if err != nil {
		return result{}, err
	}
	user := fmt.Sprintf("<TRANSCRIPT>\n%s\n</TRANSCRIPT>\nReturn at most %d facts.", labeled, cfg.MaxFacts)
	var rawOutput string
	var completion fireworks.ChatResult
	if cfg.Approach == "quote-first" {
		selected, stageOne, selectionErr := selectQuotes(ctx, client, cfg.Model, labeled, transcript, turns, cfg.MaxFacts)
		if selectionErr != nil {
			return result{}, selectionErr
		}
		payload, _ := json.Marshal(map[string]any{"candidate_quotes": selected})
		completion, err = client.ChatJSONNoReasoning(ctx, cfg.Model, cfg.Prompt, "<CANDIDATE_QUOTES>\n"+string(payload)+"\n</CANDIDATE_QUOTES>\nReturn one card per defensible candidate, and no other claims.", 6000)
		rawOutput = fmt.Sprintf("{\"stage_1_quote_selection\":%s,\"stage_2_rephrasing\":%s}", stageOne, completion.Text)
	} else {
		completion, err = client.ChatJSONNoReasoning(ctx, cfg.Model, cfg.Prompt, user, 6000)
		rawOutput = completion.Text
	}
	if err != nil {
		return result{}, err
	}
	parsed, err := decode(completion.Text)
	if err != nil {
		return result{}, fmt.Errorf("decode model output: %w", err)
	}
	validated := validate(parsed.Facts, transcript, turns)
	promptMaterial := cfg.Prompt
	if cfg.Approach == "quote-first" {
		promptMaterial += "\n--- quote-selection stage ---\n" + quoteSelectionPrompt
	}
	effectivePrompt := promptMaterial
	promptHash, rawHash := hash(promptMaterial), hash(rawOutput)
	tx, err := pool.Begin(ctx)
	if err != nil {
		return result{}, err
	}
	defer tx.Rollback(ctx)
	var runID string
	err = tx.QueryRow(ctx, `INSERT INTO fact_rephrasing_study_runs(namespace,version,iteration,document_id,model,prompt_hash,config,raw_output,raw_output_hash,effective_prompt) VALUES($1,$2,$3,$4,$5,$6,jsonb_build_object('max_facts',$7::int,'context','full transcript with numbered speaker turns','approach',$8::text,'temperature',0,'reasoning_effort','none'),$9,$10,$11) ON CONFLICT(namespace,version,iteration,document_id) DO UPDATE SET model=EXCLUDED.model,prompt_hash=EXCLUDED.prompt_hash,config=EXCLUDED.config,raw_output=EXCLUDED.raw_output,raw_output_hash=EXCLUDED.raw_output_hash,effective_prompt=EXCLUDED.effective_prompt,created_at=now() RETURNING id`, cfg.Namespace, cfg.Version, cfg.Iteration, documentID, cfg.Model, promptHash, cfg.MaxFacts, cfg.Approach, rawOutput, rawHash, effectivePrompt).Scan(&runID)
	if err != nil {
		return result{}, err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM fact_rephrasing_study_facts WHERE run_id=$1`, runID); err != nil {
		return result{}, err
	}
	for i, f := range validated {
		if _, err = tx.Exec(ctx, `INSERT INTO fact_rephrasing_study_facts(run_id,ordinal,text,quote,quote_char_start,quote_char_end,turn_start,turn_end,speaker_label,content_hash) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, runID, i+1, f.Text, f.Quote, f.start, f.end, f.turnStart, f.turnEnd, f.speaker, hash(f.Text)); err != nil {
			return result{}, err
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return result{}, err
	}
	if cfg.OutputDir != "" {
		if err = writeArtifacts(cfg.OutputDir, cfg, effectivePrompt, rawOutput, validated, promptHash, rawHash); err != nil {
			return result{}, err
		}
	}
	return result{RunID: runID, RawCount: len(parsed.Facts), ValidatedCount: len(validated), PromptHash: promptHash, RawOutputHash: rawHash}, nil
}

const quoteSelectionPrompt = `Return ONLY JSON shaped as {"quotes":[{"quote":"exact contiguous source text"}]}. Read the numbered dialogue transcript and select at most the requested number of short evidence spans. Coverage first: sweep every substantive declarative speaker turn in order, selecting at least one span from each new topic before selecting a second span from an earlier topic. Each selected quote must contain exactly one explicit, substantive speaker claim. Include the complete clause containing its reporting or epistemic operator (for example "I think", "it seems to me", "today it was reported", "according to") and all negation, conditions, and time limits; never cut the quote so these markers are lost. If the claim uses a pronoun, report, comparison, or ellipsis, extend the contiguous quote to include its explicit antecedent or comparator, up to 240 characters; otherwise do not select it. Do not select questions, duplicate claims, or background narration.`

type quoteSelection struct {
	Quotes []struct {
		Quote string `json:"quote"`
	} `json:"quotes"`
}

func selectQuotes(ctx context.Context, client *fireworks.Client, model, labeled, transcript string, turns []turn, max int) ([]map[string]any, string, error) {
	res, err := client.ChatJSONNoReasoning(ctx, model, quoteSelectionPrompt, fmt.Sprintf("<TRANSCRIPT>\n%s\n</TRANSCRIPT>\nSelect at most %d quotes.", labeled, max), 6000)
	if err != nil {
		return nil, "", err
	}
	var parsed quoteSelection
	raw := strings.TrimSpace(res.Text)
	if i, j := strings.Index(raw, "{"), strings.LastIndex(raw, "}"); i >= 0 && j > i {
		raw = raw[i : j+1]
	}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil, res.Text, fmt.Errorf("decode quote selection: %w", err)
	}
	seen := map[string]bool{}
	var out []map[string]any
	for _, q := range parsed.Quotes {
		quote := strings.TrimSpace(q.Quote)
		if quote == "" || len(quote) > 240 || seen[quote] {
			continue
		}
		at := strings.Index(transcript, quote)
		if at < 0 {
			continue
		}
		turnStart, _, speaker := evidence(turns, at, at+len(quote))
		contextText := localContext(turns, turnStart)
		seen[quote] = true
		out = append(out, map[string]any{"quote": quote, "speaker": speaker, "local_context": contextText})
		if len(out) >= max {
			break
		}
	}
	return out, res.Text, nil
}

func localContext(turns []turn, ordinal int) string {
	for i, t := range turns {
		if t.ordinal != ordinal {
			continue
		}
		var parts []string
		if i > 0 {
			parts = append(parts, "[PREVIOUS TURN] "+turns[i-1].text)
		}
		parts = append(parts, "[SOURCE TURN] "+t.text)
		if i+1 < len(turns) {
			parts = append(parts, "[NEXT TURN] "+turns[i+1].text)
		}
		return strings.Join(parts, "\n")
	}
	return ""
}

type fact struct {
	Candidate
	start, end, turnStart, turnEnd int
	speaker                        string
}

func validate(candidates []Candidate, transcript string, turns []turn) []fact {
	seen := map[string]bool{}
	out := make([]fact, 0, len(candidates))
	for _, c := range candidates {
		c.Text = strings.TrimSpace(c.Text)
		c.Quote = strings.TrimSpace(c.Quote)
		if c.Text == "" || c.Quote == "" || len(c.Quote) > 240 || seen[hash(c.Text)] {
			continue
		}
		at := strings.Index(transcript, c.Quote)
		if at < 0 {
			continue
		}
		first, last, speaker := evidence(turns, at, at+len(c.Quote))
		seen[hash(c.Text)] = true
		out = append(out, fact{Candidate: c, start: at, end: at + len(c.Quote), turnStart: first, turnEnd: last, speaker: speaker})
	}
	return out
}
func source(ctx context.Context, pool *pgxpool.Pool, id string) (string, string, []turn, error) {
	rows, err := pool.Query(ctx, `SELECT ordinal,speaker_label,text,char_start,char_end FROM turns WHERE document_id=$1 ORDER BY ordinal`, id)
	if err != nil {
		return "", "", nil, err
	}
	defer rows.Close()
	var raw, labeled []string
	var turns []turn
	for rows.Next() {
		var t turn
		if err := rows.Scan(&t.ordinal, &t.speaker, &t.text, &t.start, &t.end); err != nil {
			return "", "", nil, err
		}
		raw = append(raw, t.text)
		labeled = append(labeled, fmt.Sprintf("[TURN %d | %s]\n%s", t.ordinal, t.speaker, t.text))
		turns = append(turns, t)
	}
	return strings.Join(raw, "\n"), strings.Join(labeled, "\n\n"), turns, rows.Err()
}
func evidence(turns []turn, start, end int) (int, int, string) {
	if len(turns) == 0 {
		return 0, 0, ""
	}
	first, last := turns[0], turns[len(turns)-1]
	for _, t := range turns {
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
func decode(raw string) (response, error) {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	if i, j := strings.Index(raw, "{"), strings.LastIndex(raw, "}"); i >= 0 && j > i {
		raw = raw[i : j+1]
	}
	var r response
	err := json.Unmarshal([]byte(raw), &r)
	return r, err
}
func hash(s string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(s)))
	return hex.EncodeToString(sum[:])
}
func writeArtifacts(dir string, cfg Config, effectivePrompt, raw string, facts []fact, promptHash, rawHash string) error {
	iteration := fmt.Sprintf("iteration-%02d", cfg.Iteration)
	d := filepath.Join(dir, iteration)
	if err := os.MkdirAll(d, 0755); err != nil {
		return err
	}
	values := map[string]any{"namespace": cfg.Namespace, "version": cfg.Version, "iteration": cfg.Iteration, "external_id": cfg.ExternalID, "model": cfg.Model, "prompt_hash": promptHash, "raw_output_hash": rawHash, "validated_fact_count": len(facts), "validated_facts": facts}
	b, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(d, "validated-facts.json"), b, 0600); err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(d, "effective-prompt.txt"), []byte(effectivePrompt), 0600); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(d, "raw-output.json"), []byte(raw), 0600)
}
