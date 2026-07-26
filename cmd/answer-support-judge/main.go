// answer-support-judge checks whether the sources selected by an answering
// model actually support its answer. It is deliberately distinct from answer
// correctness: an answer can match a reference while lacking source support.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const supportSystem = `Assess source support, not general-world truth. Return ONLY JSON: {"verdict":"supported"|"unsupported","reason":"brief"}. Mark supported only when the supplied cited transcript source(s), including any verified quote, directly establish the generated answer to the question. Mark unsupported for missing, irrelevant, contradictory, or merely related sources. Do not infer unstated identities, quantities, attribution, polarity, or time.`

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
type job struct {
	answerID, question, answer string
	citations                  citationPayload
}
type verdict struct {
	Verdict string `json:"verdict"`
	Reason  string `json:"reason"`
}

func main() {
	var runs, model, projectDir, key string
	var workers, limit int
	flag.StringVar(&runs, "runs", "", "Comma-separated experiment run UUIDs")
	flag.StringVar(&model, "model", "accounts/fireworks/models/deepseek-v4-flash", "Blind source-support judge model")
	flag.StringVar(&projectDir, "project-dir", ".", "Project directory containing .env")
	flag.StringVar(&key, "key", "source_support", "Usage key under which to store judgments")
	flag.IntVar(&workers, "workers", 8, "Concurrent judgments")
	flag.IntVar(&limit, "limit", 0, "Maximum answers to judge; 0 means all")
	flag.Parse()
	if runs == "" || workers < 1 {
		fmt.Fprintln(os.Stderr, "--runs and positive --workers are required")
		os.Exit(2)
	}
	if err := godotenv.Load(projectDir + "/.env"); err != nil && !os.IsNotExist(err) {
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
	client, err := fireworks.NewFromEnv()
	if err != nil {
		panic(err)
	}
	rows, err := pool.Query(ctx, `SELECT a.id,q.text,a.text,a.citations FROM answers a JOIN questions q ON q.id=a.question_id WHERE a.experiment_run_id=ANY($1) AND NOT (a.usage ? $2) ORDER BY a.id`, strings.Split(runs, ","), key)
	if err != nil {
		panic(err)
	}
	var jobsList []job
	for rows.Next() {
		var j job
		var raw []byte
		if err = rows.Scan(&j.answerID, &j.question, &j.answer, &raw); err != nil {
			panic(err)
		}
		if err = json.Unmarshal(raw, &j.citations); err != nil {
			panic(fmt.Errorf("answer %s has non-structured citations: %w", j.answerID, err))
		}
		jobsList = append(jobsList, j)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
	rows.Close()
	if limit > 0 && len(jobsList) > limit {
		jobsList = jobsList[:limit]
	}
	jobs := make(chan job)
	var wg sync.WaitGroup
	var mu sync.Mutex
	completed, failed := 0, 0
	var firstFailure string
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				if len(j.citations.CitedSourceIDs) == 0 {
					err := store(ctx, pool, j.answerID, key, map[string]any{"model": model, "verdict": "unsupported", "reason": "no cited source"})
					mu.Lock()
					if err != nil {
						failed++
					} else {
						completed++
					}
					mu.Unlock()
					continue
				}
				sources := citedSources(j.citations)
				payload, _ := json.Marshal(map[string]any{"question": j.question, "generated_answer": j.answer, "cited_sources": sources})
				response, callErr := client.ChatJSONNoReasoning(ctx, model, supportSystem, string(payload), 180)
				var v verdict
				if callErr == nil {
					callErr = json.Unmarshal([]byte(response.Text), &v)
				}
				if callErr != nil || (v.Verdict != "supported" && v.Verdict != "unsupported") {
					mu.Lock()
					failed++
					if firstFailure == "" {
						firstFailure = fmt.Sprintf("call=%v response=%q verdict=%q", callErr, response.Text, v.Verdict)
					}
					mu.Unlock()
					continue
				}
				if err := store(ctx, pool, j.answerID, key, map[string]any{"model": model, "verdict": v.Verdict, "reason": v.Reason, "prompt_tokens": response.PromptTokens, "completion_tokens": response.CompletionTokens, "total_tokens": response.TotalTokens}); err != nil {
					mu.Lock()
					failed++
					mu.Unlock()
					continue
				}
				mu.Lock()
				completed++
				mu.Unlock()
			}
		}()
	}
	for _, j := range jobsList {
		jobs <- j
	}
	close(jobs)
	wg.Wait()
	fmt.Printf("Judged %d answer(s); %d failed.\n", completed, failed)
	if firstFailure != "" {
		fmt.Printf("First failure: %s\n", firstFailure)
	}
}

func citedSources(payload citationPayload) []citationSource {
	allowed := map[int]bool{}
	for _, id := range payload.CitedSourceIDs {
		allowed[id] = true
	}
	var out []citationSource
	for _, source := range payload.Retrieved {
		if allowed[source.ID] {
			out = append(out, source)
		}
	}
	return out
}

func store(ctx context.Context, pool *pgxpool.Pool, answerID, key string, value map[string]any) error {
	raw, _ := json.Marshal(value)
	_, err := pool.Exec(ctx, `UPDATE answers SET usage=usage || jsonb_build_object($2::text,$3::jsonb) WHERE id=$1`, answerID, key, string(raw))
	return err
}
