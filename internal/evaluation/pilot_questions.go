package evaluation

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PilotQuestionConfig defines a deliberately non-gold, model-generated
// retrieval probe set. Every accepted item still has an exact answer quote and
// a source-turn label verified mechanically before it is persisted.
type PilotQuestionConfig struct {
	QuestionSplit, CorpusSplit, Model          string
	Documents, PerDocument, BatchSize, Workers int
}

type PilotQuestionReport struct {
	Documents, Questions, Rejected int
}

type pilotTurn struct {
	Ordinal int
	Speaker string
	Text    string
}

type pilotDocument struct{ id, externalID string }

type pilotItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Turn     int    `json:"turn"`
}

type pilotResponse struct {
	Items []pilotItem `json:"items"`
}

const pilotQuestionSystem = `You create source-grounded retrieval probes from a human-dialogue transcript. Return ONLY JSON shaped as {"items":[{"question":"...","answer":"exact quote","turn":0}]}. Create factual, natural-language questions whose answer is directly established by the exact contiguous answer quote in the stated source turn. The answer must be 3-180 characters copied verbatim from that turn. Do not use titles, summaries, outside knowledge, inference, or a question premise as evidence. Questions must be semantically distinct; when a short transcript lacks 100 independent facts, create distinct question formulations about explicit claims, but never invent a claim. Do not mention turn numbers in a question.`

// GeneratePilotQuestions creates and validates a fixed-size diagnostic question
// set. It is explicitly not a substitute for independently authored paper gold
// questions; its purpose is to stress the retrieval/storage pipeline at scale.
func GeneratePilotQuestions(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, cfg PilotQuestionConfig) (PilotQuestionReport, error) {
	if cfg.QuestionSplit == "" || cfg.CorpusSplit == "" || cfg.Model == "" {
		return PilotQuestionReport{}, fmt.Errorf("question split, corpus split, and model are required")
	}
	if cfg.Documents <= 0 || cfg.PerDocument <= 0 {
		return PilotQuestionReport{}, fmt.Errorf("documents and per-document must be positive")
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 10
	}
	if cfg.Workers <= 0 {
		cfg.Workers = 1
	}
	rows, err := pool.Query(ctx, `SELECT id,external_id FROM documents WHERE corpus='mediasum' AND corpus_split=$1 ORDER BY external_id LIMIT $2`, cfg.CorpusSplit, cfg.Documents)
	if err != nil {
		return PilotQuestionReport{}, err
	}
	defer rows.Close()
	var documents []pilotDocument
	for rows.Next() {
		var doc pilotDocument
		if err := rows.Scan(&doc.id, &doc.externalID); err != nil {
			return PilotQuestionReport{}, err
		}
		documents = append(documents, doc)
	}
	if err := rows.Err(); err != nil {
		return PilotQuestionReport{}, err
	}
	if len(documents) != cfg.Documents {
		return PilotQuestionReport{}, fmt.Errorf("requested %d documents but found %d in %s", cfg.Documents, len(documents), cfg.CorpusSplit)
	}

	var report PilotQuestionReport
	var firstErr error
	var mutex sync.Mutex
	jobs := make(chan pilotDocument)
	var workers sync.WaitGroup
	for worker := 0; worker < minPilot(cfg.Workers, len(documents)); worker++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for doc := range jobs {
				questions, rejected, err := generatePilotDocument(ctx, pool, client, doc, cfg)
				mutex.Lock()
				if err != nil && firstErr == nil {
					firstErr = err
				}
				if err == nil {
					report.Documents++
					report.Questions += questions
					report.Rejected += rejected
				}
				mutex.Unlock()
			}
		}()
	}
	for _, doc := range documents {
		jobs <- doc
	}
	close(jobs)
	workers.Wait()
	if firstErr != nil {
		return report, firstErr
	}
	return report, nil
}

func generatePilotDocument(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, doc pilotDocument, cfg PilotQuestionConfig) (int, int, error) {
	turns, err := pilotDocumentTurns(ctx, pool, doc.id)
	if err != nil {
		return 0, 0, err
	}
	var accepted []pilotItem
	acceptedQuestions := map[string]bool{}
	reportRejected := 0
	attempts, maxAttempts := 0, (cfg.PerDocument/cfg.BatchSize+1)*4
	for len(accepted) < cfg.PerDocument {
		attempts++
		if attempts > maxAttempts {
			return 0, reportRejected, fmt.Errorf("%s produced only %d/%d valid distinct questions after %d bounded attempts", doc.externalID, len(accepted), cfg.PerDocument, maxAttempts)
		}
		needed := minPilot(cfg.BatchSize, cfg.PerDocument-len(accepted))
		response, err := client.ChatJSONNoReasoning(ctx, cfg.Model, pilotQuestionSystem, pilotQuestionPrompt(turns, needed, accepted), maxPilot(2400, needed*70))
		if err != nil {
			return 0, reportRejected, fmt.Errorf("generate pilot questions for %s: %w", doc.externalID, err)
		}
		var parsed pilotResponse
		if err = json.Unmarshal([]byte(response.Text), &parsed); err != nil {
			return 0, reportRejected, fmt.Errorf("parse pilot questions for %s: %w", doc.externalID, err)
		}
		batch, rejected := validatePilotItems(parsed.Items, turns, needed)
		reportRejected += rejected
		for _, item := range batch {
			key := strings.ToLower(item.Question)
			if acceptedQuestions[key] {
				reportRejected++
				continue
			}
			acceptedQuestions[key] = true
			accepted = append(accepted, item)
		}
	}
	if err = persistPilotQuestions(ctx, pool, doc, accepted, cfg); err != nil {
		return 0, reportRejected, err
	}
	return len(accepted), reportRejected, nil
}

func pilotDocumentTurns(ctx context.Context, pool *pgxpool.Pool, documentID string) ([]pilotTurn, error) {
	rows, err := pool.Query(ctx, `SELECT ordinal,COALESCE(speaker_label,''),text FROM turns WHERE document_id=$1 ORDER BY ordinal`, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var turns []pilotTurn
	for rows.Next() {
		var turn pilotTurn
		if err := rows.Scan(&turn.Ordinal, &turn.Speaker, &turn.Text); err != nil {
			return nil, err
		}
		turns = append(turns, turn)
	}
	return turns, rows.Err()
}

func pilotQuestionPrompt(turns []pilotTurn, perDocument int, existing []pilotItem) string {
	var source strings.Builder
	for _, turn := range turns {
		fmt.Fprintf(&source, "<TURN n=\"%d\" speaker=\"%s\">\n%s\n</TURN>\n", turn.Ordinal, turn.Speaker, turn.Text)
	}
	used := make([]string, len(existing))
	for i, item := range existing {
		used[i] = item.Question
	}
	return fmt.Sprintf("Create exactly %d items from these source turns. The `answer` field must be copied exactly from the selected turn. Do not repeat or trivially rephrase any existing question: %s\n<SOURCE_TURNS>\n%s</SOURCE_TURNS>", perDocument, strings.Join(used, " | "), source.String())
}

func minPilot(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxPilot(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func validatePilotItems(items []pilotItem, turns []pilotTurn, required int) ([]pilotItem, int) {
	byOrdinal := make(map[int]string, len(turns))
	for _, turn := range turns {
		byOrdinal[turn.Ordinal] = turn.Text
	}
	seen := map[string]bool{}
	var accepted []pilotItem
	rejected := 0
	for _, item := range items {
		item.Question = strings.Join(strings.Fields(item.Question), " ")
		item.Answer = strings.TrimSpace(item.Answer)
		if item.Question == "" || len(item.Answer) < 3 || len(item.Answer) > 180 || !strings.Contains(byOrdinal[item.Turn], item.Answer) {
			rejected++
			continue
		}
		key := strings.ToLower(item.Question)
		if seen[key] {
			rejected++
			continue
		}
		seen[key] = true
		accepted = append(accepted, item)
		if len(accepted) == required {
			break
		}
	}
	sort.Slice(accepted, func(i, j int) bool { return accepted[i].Question < accepted[j].Question })
	return accepted, rejected
}

func persistPilotQuestions(ctx context.Context, pool *pgxpool.Pool, document pilotDocument, items []pilotItem, cfg PilotQuestionConfig) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for index, item := range items {
		externalID := fmt.Sprintf("%s-%s-%03d", cfg.QuestionSplit, document.externalID, index+1)
		if _, err = tx.Exec(ctx, `INSERT INTO questions(external_id,corpus_split,question_type,text,reference_answer,metadata) VALUES ($1,$2,'model_generated_retrieval_probe',$3,$4,jsonb_build_object('generator_model',$5::text,'document_external_id',$6::text,'protocol','exact-answer-quote')) ON CONFLICT(external_id) DO UPDATE SET corpus_split=EXCLUDED.corpus_split,question_type=EXCLUDED.question_type,text=EXCLUDED.text,reference_answer=EXCLUDED.reference_answer,metadata=EXCLUDED.metadata`, externalID, cfg.QuestionSplit, item.Question, item.Answer, cfg.Model, document.externalID); err != nil {
			return err
		}
		var questionID string
		if err = tx.QueryRow(ctx, `SELECT id FROM questions WHERE external_id=$1`, externalID).Scan(&questionID); err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, `DELETE FROM gold_evidence WHERE question_id=$1`, questionID); err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO gold_evidence(question_id,document_id,turn_start,turn_end) VALUES($1,$2,$3,$3)`, questionID, document.id, item.Turn); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
