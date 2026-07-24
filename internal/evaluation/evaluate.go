// Package evaluation imports human-authored evidence labels and runs RAG retrieval trials.
package evaluation

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/retrieval"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EvidenceInput struct {
	DocumentExternalID string `json:"document_external_id"`
	TurnStart          int    `json:"turn_start"`
	TurnEnd            int    `json:"turn_end"`
}
type QuestionInput struct {
	ExternalID      string          `json:"external_id"`
	CorpusSplit     string          `json:"corpus_split"`
	QuestionType    string          `json:"question_type"`
	Text            string          `json:"text"`
	ReferenceAnswer string          `json:"reference_answer"`
	Evidence        []EvidenceInput `json:"evidence"`
}

func Import(ctx context.Context, pool *pgxpool.Pool, path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	scan := bufio.NewScanner(file)
	scan.Buffer(make([]byte, 64*1024), 4<<20)
	n := 0
	for scan.Scan() {
		var input QuestionInput
		if err := json.Unmarshal(scan.Bytes(), &input); err != nil {
			return n, fmt.Errorf("decode question %d: %w", n+1, err)
		}
		if input.ExternalID == "" || input.Text == "" || input.QuestionType == "" {
			return n, fmt.Errorf("question %d needs external_id, text, and question_type", n+1)
		}
		if input.CorpusSplit == "" {
			input.CorpusSplit = "dev"
		}
		tx, err := pool.Begin(ctx)
		if err != nil {
			return n, err
		}
		var id string
		err = tx.QueryRow(ctx, `INSERT INTO questions(external_id,corpus_split,question_type,text,reference_answer) VALUES ($1,$2,$3,$4,$5) ON CONFLICT(external_id) DO UPDATE SET corpus_split=EXCLUDED.corpus_split,question_type=EXCLUDED.question_type,text=EXCLUDED.text,reference_answer=EXCLUDED.reference_answer RETURNING id`, input.ExternalID, input.CorpusSplit, input.QuestionType, input.Text, nullString(input.ReferenceAnswer)).Scan(&id)
		if err == nil {
			_, err = tx.Exec(ctx, `DELETE FROM gold_evidence WHERE question_id=$1`, id)
		}
		for _, e := range input.Evidence {
			if err != nil {
				break
			}
			var docID string
			err = tx.QueryRow(ctx, `SELECT id FROM documents WHERE external_id=$1`, e.DocumentExternalID).Scan(&docID)
			if err != nil {
				err = fmt.Errorf("question %s evidence document %s: %w", input.ExternalID, e.DocumentExternalID, err)
				break
			}
			_, err = tx.Exec(ctx, `INSERT INTO gold_evidence(question_id,document_id,turn_start,turn_end) VALUES($1,$2,$3,$4)`, id, docID, e.TurnStart, e.TurnEnd)
		}
		if err != nil {
			_ = tx.Rollback(ctx)
			return n, err
		}
		if err = tx.Commit(ctx); err != nil {
			return n, err
		}
		n++
	}
	return n, scan.Err()
}

type Config struct {
	Representation, Version, QuestionSplit, CorpusSplit, EmbeddingModel, GeneratorModel string
	TopK                                                                                int
	Generate                                                                            bool
}
type Metrics struct {
	RunID                                                    string
	Questions, EvidenceTotal, EvidenceHits, QuestionsWithHit int
}

func Run(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, cfg Config) (Metrics, error) {
	if cfg.TopK <= 0 {
		cfg.TopK = 5
	}
	var runID string
	parameters, _ := json.Marshal(map[string]any{"question_split": cfg.QuestionSplit, "corpus_split": cfg.CorpusSplit, "top_k": cfg.TopK, "generate": cfg.Generate})
	err := pool.QueryRow(ctx, `INSERT INTO experiment_runs(kind,representation,representation_version,embedding_model,generator_model,parameters) VALUES('retrieval',$1,$2,$3,$4,$5) RETURNING id`, cfg.Representation, cfg.Version, cfg.EmbeddingModel, nullString(cfg.GeneratorModel), parameters).Scan(&runID)
	if err != nil {
		return Metrics{}, err
	}
	metrics := Metrics{RunID: runID}
	rows, err := pool.Query(ctx, `SELECT id,text FROM questions WHERE corpus_split=$1 ORDER BY external_id`, cfg.QuestionSplit)
	if err != nil {
		return metrics, fail(ctx, pool, runID, err)
	}
	defer rows.Close()
	for rows.Next() {
		var qid, question string
		if err := rows.Scan(&qid, &question); err != nil {
			return metrics, fail(ctx, pool, runID, err)
		}
		metrics.Questions++
		results, err := retrieval.Search(ctx, pool, client, question, retrieval.Config{Representation: cfg.Representation, Version: cfg.Version, CorpusSplit: cfg.CorpusSplit, EmbeddingModel: cfg.EmbeddingModel, TopK: cfg.TopK})
		if err != nil {
			return metrics, fail(ctx, pool, runID, err)
		}
		tx, err := pool.Begin(ctx)
		if err != nil {
			return metrics, fail(ctx, pool, runID, err)
		}
		for rank, r := range results {
			if _, err = tx.Exec(ctx, `INSERT INTO retrievals(experiment_run_id,question_id,source_unit_id,rank,score) VALUES($1,$2,$3,$4,$5)`, runID, qid, r.SourceUnitID, rank+1, r.Score); err != nil {
				_ = tx.Rollback(ctx)
				return metrics, fail(ctx, pool, runID, err)
			}
		}
		var total, hits int
		err = tx.QueryRow(ctx, `SELECT count(*),count(*) FILTER (WHERE EXISTS (SELECT 1 FROM retrievals r JOIN source_units su ON su.id=r.source_unit_id WHERE r.experiment_run_id=$1 AND r.question_id=ge.question_id AND su.document_id=ge.document_id AND su.turn_start<=ge.turn_end AND su.turn_end>=ge.turn_start)) FROM gold_evidence ge WHERE ge.question_id=$2`, runID, qid).Scan(&total, &hits)
		if err != nil {
			_ = tx.Rollback(ctx)
			return metrics, fail(ctx, pool, runID, err)
		}
		metrics.EvidenceTotal += total
		metrics.EvidenceHits += hits
		if hits > 0 {
			metrics.QuestionsWithHit++
		}
		if cfg.Generate {
			answer, err := generate(ctx, client, cfg.GeneratorModel, question, results)
			if err != nil {
				_ = tx.Rollback(ctx)
				return metrics, fail(ctx, pool, runID, err)
			}
			citations, _ := json.Marshal(results)
			_, err = tx.Exec(ctx, `INSERT INTO answers(experiment_run_id,question_id,text,citations) VALUES($1,$2,$3,$4)`, runID, qid, answer, citations)
			if err != nil {
				_ = tx.Rollback(ctx)
				return metrics, fail(ctx, pool, runID, err)
			}
		}
		if err = tx.Commit(ctx); err != nil {
			return metrics, fail(ctx, pool, runID, err)
		}
	}
	if err := rows.Err(); err != nil {
		return metrics, fail(ctx, pool, runID, err)
	}
	_, err = pool.Exec(ctx, `UPDATE experiment_runs SET status='complete',completed_at=now() WHERE id=$1`, runID)
	return metrics, err
}
func generate(ctx context.Context, client *fireworks.Client, model, question string, results []retrieval.Result) (string, error) {
	if model == "" {
		return "", fmt.Errorf("--generator-model is required with --generate")
	}
	var context strings.Builder
	for i, r := range results {
		fmt.Fprintf(&context, "[S%d document=%s turns=%d-%d]\n%s\n\n", i+1, r.ExternalID, r.TurnStart, r.TurnEnd, r.Text)
	}
	result, err := client.Chat(ctx, model, "Answer only from the supplied transcript sources. If the sources do not establish an answer, say so. Cite every substantive claim using the supplied [S#] labels.", fmt.Sprintf("Question: %s\n\nSources:\n%s", question, context.String()), 600)
	return result.Text, err
}
func fail(ctx context.Context, pool *pgxpool.Pool, id string, cause error) error {
	_, _ = pool.Exec(ctx, `UPDATE experiment_runs SET status='failed',completed_at=now() WHERE id=$1`, id)
	return cause
}
func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}
