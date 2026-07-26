// human-calibration-export produces a blinded, stratified worksheet for a
// person to audit answer correctness and cited-source support.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/joho/godotenv"
)

type citationSource struct {
	ID           int    `json:"id"`
	DocumentID   string `json:"document_id"`
	TurnStart    int    `json:"turn_start"`
	TurnEnd      int    `json:"turn_end"`
	Text         string `json:"text"`
	Quote        string `json:"quote"`
	SpeakerLabel string `json:"speaker_label"`
}
type citationPayload struct {
	CitedSourceIDs []int            `json:"cited_source_ids"`
	Retrieved      []citationSource `json:"retrieved"`
}
type item struct {
	AnswerID, RunID, ExternalID, Representation, Budget, QuestionType string
	Question, Reference, Answer                                       string
	Correct, Support                                                  string
	Citations                                                         citationPayload
}

func main() {
	var projectDir, runs, out, correctnessKey, supportKey string
	flag.StringVar(&projectDir, "project-dir", ".", "Project directory containing .env")
	flag.StringVar(&runs, "runs", "", "Comma-separated experiment run UUIDs")
	flag.StringVar(&out, "out", "", "Output Markdown worksheet path")
	flag.StringVar(&correctnessKey, "correctness-key", "judge_kimi_v1", "Stored blind-correctness judgment key")
	flag.StringVar(&supportKey, "support-key", "source_support_kimi_v1", "Stored blind-source-support judgment key")
	flag.Parse()
	if runs == "" || out == "" {
		fmt.Fprintln(os.Stderr, "--runs and --out are required")
		os.Exit(2)
	}
	if err := godotenv.Load(filepath.Join(projectDir, ".env")); err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	ctx := context.Background()
	pool, err := storage.Open(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	if err = storage.Migrate(ctx, pool); err != nil {
		panic(err)
	}
	rows, err := pool.Query(ctx, `
WITH base AS (
  SELECT a.id AS answer_id,e.id AS run_id,q.external_id,q.question_type,q.text AS question,
         COALESCE(q.reference_answer,'') AS reference_answer,a.text AS answer,a.citations,e.representation,
         COALESCE(e.parameters->>'context_token_budget','') AS budget,
         COALESCE(a.usage->$2->>'verdict','') AS correct,
         COALESCE(a.usage->$3->>'verdict','') AS support
  FROM answers a
  JOIN experiment_runs e ON e.id=a.experiment_run_id
  JOIN questions q ON q.id=a.question_id
  WHERE e.id=ANY($1::uuid[])
), scored AS (
  SELECT b.*, EXISTS (
    SELECT 1 FROM base peer
    WHERE peer.external_id=b.external_id AND peer.budget=b.budget
      AND peer.correct <> b.correct
  ) AS contrast
  FROM base b
), ranked AS (
  SELECT *, row_number() OVER (
    PARTITION BY representation,budget,question_type
    ORDER BY contrast DESC, md5(answer_id::text)
  ) AS n
  FROM scored
)
SELECT answer_id,run_id,external_id,representation,budget,question_type,question,reference_answer,
       answer,citations,correct,support
FROM ranked
WHERE n=1
ORDER BY representation,budget::int,question_type`, strings.Split(runs, ","), correctnessKey, supportKey)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var items []item
	for rows.Next() {
		var x item
		var citations []byte
		if err = rows.Scan(&x.AnswerID, &x.RunID, &x.ExternalID, &x.Representation, &x.Budget, &x.QuestionType, &x.Question, &x.Reference, &x.Answer, &citations, &x.Correct, &x.Support); err != nil {
			panic(err)
		}
		if err = json.Unmarshal(citations, &x.Citations); err != nil {
			panic(fmt.Errorf("answer %s citations: %w", x.AnswerID, err))
		}
		items = append(items, x)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
	var worksheet strings.Builder
	worksheet.WriteString("# ISC-RAG blinded human-calibration worksheet\n\n")
	worksheet.WriteString("This sheet deliberately omits representation, budget, run ID, and model verdicts. For each item, independently label (a) whether the answer matches the reference answer and (b) whether the cited source directly supports the generated answer. Mark an answer incorrect if it changes an entity, number, attribution, polarity, scope, or time. Mark source support absent when a citation is missing, irrelevant, or only related.\n\n")
	worksheet.WriteString("There are 80 deterministic strata: 4 representations × 4 budgets × 5 question families. Within each stratum, the selection prefers a question where Kimi’s correctness verdict differs across representations at that budget; this makes disagreement cases available for calibration without revealing the condition to the reviewer.\n\n")
	for i, x := range items {
		fmt.Fprintf(&worksheet, "## Item %03d\n\n**Question:** %s\n\n**Reference answer:** %s\n\n**Generated answer:** %s\n\n", i+1, x.Question, x.Reference, x.Answer)
		cited := citedSources(x.Citations)
		if len(cited) == 0 {
			worksheet.WriteString("**Cited source(s):** None\n\n")
		} else {
			worksheet.WriteString("**Cited source(s):**\n\n")
			for _, source := range cited {
				fmt.Fprintf(&worksheet, "- Source %d (%s, turns %d–%d): %s\n", source.ID, source.DocumentID, source.TurnStart, source.TurnEnd, source.Text)
				if source.Quote != "" {
					fmt.Fprintf(&worksheet, "  - Verified quote%s: %s\n", speakerSuffix(source.SpeakerLabel), source.Quote)
				}
			}
			worksheet.WriteString("\n")
		}
		worksheet.WriteString("**Human correctness:** [ ] Correct  [ ] Incorrect  [ ] Unclear\n\n**Human source support:** [ ] Supported  [ ] Unsupported  [ ] Unclear\n\n**Notes:**\n\n---\n\n")
	}
	if len(items) != 80 {
		panic(fmt.Errorf("expected 80 calibration strata, got %d", len(items)))
	}
	if err = os.WriteFile(out, []byte(worksheet.String()), 0o644); err != nil {
		panic(err)
	}
	keyPath := strings.TrimSuffix(out, filepath.Ext(out)) + "-key.jsonl"
	key, _ := os.Create(keyPath)
	if key == nil {
		panic(fmt.Errorf("create %s", keyPath))
	}
	defer key.Close()
	for i, x := range items {
		record, _ := json.Marshal(map[string]any{"item": i + 1, "answer_id": x.AnswerID, "run_id": x.RunID, "question_external_id": x.ExternalID, "question_type": x.QuestionType, "representation": x.Representation, "context_token_budget": x.Budget, "kimi_correctness": x.Correct, "kimi_source_support": x.Support})
		if _, err = key.Write(append(record, '\n')); err != nil {
			panic(err)
		}
	}
	fmt.Printf("Wrote %d blinded calibration items to %s and hidden key to %s.\n", len(items), out, keyPath)
}

func citedSources(payload citationPayload) []citationSource {
	wanted := map[int]bool{}
	for _, id := range payload.CitedSourceIDs {
		wanted[id] = true
	}
	var result []citationSource
	for _, source := range payload.Retrieved {
		if wanted[source.ID] {
			result = append(result, source)
		}
	}
	return result
}

func speakerSuffix(label string) string {
	if label == "" {
		return ""
	}
	return " (speaker: " + label + ")"
}
