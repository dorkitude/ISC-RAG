// answer-judge adds an independent semantic correctness screen to stored RAG
// answers. It deliberately sees the question, reference answer, and generated
// answer only; it does not see which retrieval representation produced it.
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

const judgeSystem = `Judge semantic answer correctness. Return ONLY JSON: {"verdict":"correct"|"incorrect","reason":"brief"}. Mark correct when the generated answer directly answers the question consistently with the reference answer, even if it is shorter or differently worded. Mark incorrect if it changes the entity, number, attribution, polarity, or lacks an answer. Do not reward unsupported extra claims.`

type job struct{ answerID, question, reference, answer string }
type verdict struct {
	Verdict string `json:"verdict"`
	Reason  string `json:"reason"`
}

func main() {
	var runs, model, projectDir, key string
	var workers, limit int
	flag.StringVar(&runs, "runs", "", "Comma-separated experiment run UUIDs")
	flag.StringVar(&model, "model", "accounts/fireworks/models/kimi-k2p6", "Independent judge model")
	flag.StringVar(&projectDir, "project-dir", ".", "Project directory containing .env")
	flag.StringVar(&key, "key", "judge", "Usage key under which to store judgments")
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
	ids := strings.Split(runs, ",")
	rows, err := pool.Query(ctx, `SELECT a.id,q.text,COALESCE(q.reference_answer,''),a.text FROM answers a JOIN questions q ON q.id=a.question_id WHERE a.experiment_run_id=ANY($1) AND NOT (a.usage ? $2) ORDER BY a.id`, ids, key)
	if err != nil {
		panic(err)
	}
	var jobsList []job
	for rows.Next() {
		var j job
		if err = rows.Scan(&j.answerID, &j.question, &j.reference, &j.answer); err != nil {
			panic(err)
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
				payload, _ := json.Marshal(map[string]string{"question": j.question, "reference_answer": j.reference, "generated_answer": j.answer})
				response, callErr := client.ChatJSONNoReasoning(ctx, model, judgeSystem, string(payload), 180)
				var v verdict
				if callErr == nil {
					callErr = json.Unmarshal([]byte(response.Text), &v)
				}
				if callErr != nil || (v.Verdict != "correct" && v.Verdict != "incorrect") {
					mu.Lock()
					failed++
					if firstFailure == "" {
						firstFailure = fmt.Sprintf("call=%v response=%q verdict=%q", callErr, response.Text, v.Verdict)
					}
					mu.Unlock()
					continue
				}
				usage, _ := json.Marshal(map[string]any{"model": model, "verdict": v.Verdict, "reason": v.Reason, "prompt_tokens": response.PromptTokens, "completion_tokens": response.CompletionTokens, "total_tokens": response.TotalTokens})
				_, callErr = pool.Exec(ctx, `UPDATE answers SET usage=usage || jsonb_build_object($2::text,$3::jsonb) WHERE id=$1`, j.answerID, key, string(usage))
				mu.Lock()
				if callErr != nil {
					failed++
					if firstFailure == "" {
						firstFailure = callErr.Error()
					}
				} else {
					completed++
				}
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

var _ *pgxpool.Pool
