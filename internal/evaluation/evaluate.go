// Package evaluation imports human-authored evidence labels and runs RAG retrieval trials.
package evaluation

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dorkitude/isc-rag/internal/embedding"
	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/retrieval"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkoukk/tiktoken-go"
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

// Export writes a stable JSONL snapshot of a question split and its canonical
// evidence. It is deliberately the inverse of Import so generated development
// sets can be checked into the experiment record after database validation.
func Export(ctx context.Context, pool *pgxpool.Pool, split, path string) (int, error) {
	rows, err := pool.Query(ctx, `SELECT q.id,q.external_id,q.corpus_split,q.question_type,q.text,COALESCE(q.reference_answer,''),d.external_id,ge.turn_start,ge.turn_end FROM questions q LEFT JOIN gold_evidence ge ON ge.question_id=q.id LEFT JOIN documents d ON d.id=ge.document_id WHERE q.corpus_split=$1 ORDER BY q.external_id,ge.turn_start,ge.turn_end`, split)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	type collected struct {
		input QuestionInput
	}
	byID := map[string]*collected{}
	var order []string
	for rows.Next() {
		var id, documentID string
		var turnStart, turnEnd int
		var input QuestionInput
		if err := rows.Scan(&id, &input.ExternalID, &input.CorpusSplit, &input.QuestionType, &input.Text, &input.ReferenceAnswer, &documentID, &turnStart, &turnEnd); err != nil {
			return 0, err
		}
		item := byID[id]
		if item == nil {
			item = &collected{input: input}
			byID[id] = item
			order = append(order, id)
		}
		if documentID != "" {
			item.input.Evidence = append(item.input.Evidence, EvidenceInput{DocumentExternalID: documentID, TurnStart: turnStart, TurnEnd: turnEnd})
		}
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	file, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	for _, id := range order {
		line, err := json.Marshal(byID[id].input)
		if err != nil {
			return 0, err
		}
		if _, err = writer.Write(append(line, '\n')); err != nil {
			return 0, err
		}
	}
	return len(order), nil
}

type Config struct {
	Representation, Version, QuestionSplit, CorpusSplit, EmbeddingModel, GeneratorModel string
	TopK, EmbeddingBatchSize, ContextChars, ContextTokens                               int
	DocumentExternalIDs                                                                 []string
	ResumeRunID                                                                         string
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
	if cfg.EmbeddingBatchSize <= 0 {
		cfg.EmbeddingBatchSize = 32
	}
	var runID string
	if cfg.ContextChars > 0 && cfg.ContextTokens > 0 {
		return Metrics{}, fmt.Errorf("context character and token budgets are mutually exclusive")
	}
	parameters, _ := json.Marshal(map[string]any{"question_split": cfg.QuestionSplit, "corpus_split": cfg.CorpusSplit, "document_external_ids": cfg.DocumentExternalIDs, "top_k": cfg.TopK, "context_char_budget": cfg.ContextChars, "context_token_budget": cfg.ContextTokens, "context_tokenizer": tokenEncodingName(cfg.ContextTokens), "generate": cfg.Generate})
	if cfg.ResumeRunID == "" {
		err := pool.QueryRow(ctx, `INSERT INTO experiment_runs(kind,representation,representation_version,embedding_model,generator_model,parameters) VALUES('retrieval',$1,$2,$3,$4,$5) RETURNING id`, cfg.Representation, cfg.Version, cfg.EmbeddingModel, nullString(cfg.GeneratorModel), parameters).Scan(&runID)
		if err != nil {
			return Metrics{}, err
		}
	} else {
		runID = cfg.ResumeRunID
		var matches bool
		err := pool.QueryRow(ctx, `SELECT representation=$2 AND representation_version=$3 AND embedding_model=$4 AND COALESCE(generator_model,'')=COALESCE($5,'') AND parameters=$6::jsonb FROM experiment_runs WHERE id=$1`, runID, cfg.Representation, cfg.Version, cfg.EmbeddingModel, nullString(cfg.GeneratorModel), string(parameters)).Scan(&matches)
		if err != nil {
			return Metrics{}, fmt.Errorf("resume run %s: %w", runID, err)
		}
		if !matches {
			return Metrics{}, fmt.Errorf("resume run %s does not match this evaluation configuration", runID)
		}
		if _, err = pool.Exec(ctx, `UPDATE experiment_runs SET status='running',completed_at=NULL WHERE id=$1`, runID); err != nil {
			return Metrics{}, err
		}
	}
	metrics := Metrics{RunID: runID}
	questionSQL := `SELECT q.id,q.text FROM questions q WHERE q.corpus_split=$1`
	args := []any{cfg.QuestionSplit}
	if len(cfg.DocumentExternalIDs) > 0 {
		questionSQL += ` AND EXISTS (SELECT 1 FROM gold_evidence ge JOIN documents d ON d.id=ge.document_id WHERE ge.question_id=q.id AND d.external_id=ANY($2))`
		args = append(args, cfg.DocumentExternalIDs)
	}
	if cfg.ResumeRunID != "" && cfg.Generate {
		questionSQL += fmt.Sprintf(` AND NOT EXISTS (SELECT 1 FROM answers a WHERE a.experiment_run_id=$%d AND a.question_id=q.id)`, len(args)+1)
		args = append(args, runID)
	}
	questionSQL += ` ORDER BY q.external_id`
	rows, err := pool.Query(ctx, questionSQL, args...)
	if err != nil {
		return metrics, fail(ctx, pool, runID, err)
	}
	defer rows.Close()
	type questionRow struct{ id, text string }
	var questions []questionRow
	for rows.Next() {
		var qid, question string
		if err := rows.Scan(&qid, &question); err != nil {
			return metrics, fail(ctx, pool, runID, err)
		}
		questions = append(questions, questionRow{id: qid, text: question})
	}
	if err := rows.Err(); err != nil {
		return metrics, fail(ctx, pool, runID, err)
	}
	for start := 0; start < len(questions); start += cfg.EmbeddingBatchSize {
		end := start + cfg.EmbeddingBatchSize
		if end > len(questions) {
			end = len(questions)
		}
		batch := questions[start:end]
		texts := make([]string, len(batch))
		for i, question := range batch {
			texts[i] = question.text
		}
		vectors, err := client.Embed(ctx, cfg.EmbeddingModel, texts)
		if err != nil {
			return metrics, fail(ctx, pool, runID, err)
		}
		for i, question := range batch {
			metrics.Questions++
			results, err := retrieval.SearchVector(ctx, pool, embedding.VectorString(vectors[i]), retrieval.Config{Representation: cfg.Representation, Version: cfg.Version, CorpusSplit: cfg.CorpusSplit, EmbeddingModel: cfg.EmbeddingModel, TopK: cfg.TopK})
			if err != nil {
				return metrics, fail(ctx, pool, runID, err)
			}
			tx, err := pool.Begin(ctx)
			if err != nil {
				return metrics, fail(ctx, pool, runID, err)
			}
			for rank, r := range results {
				if _, err = tx.Exec(ctx, `INSERT INTO retrievals(experiment_run_id,question_id,source_unit_id,rank,score) VALUES($1,$2,$3,$4,$5)`, runID, question.id, r.SourceUnitID, rank+1, r.Score); err != nil {
					_ = tx.Rollback(ctx)
					return metrics, fail(ctx, pool, runID, err)
				}
			}
			var total, hits int
			err = tx.QueryRow(ctx, `SELECT count(*),count(*) FILTER (WHERE EXISTS (SELECT 1 FROM retrievals r JOIN source_units su ON su.id=r.source_unit_id WHERE r.experiment_run_id=$1 AND r.question_id=ge.question_id AND su.document_id=ge.document_id AND su.turn_start<=ge.turn_end AND su.turn_end>=ge.turn_start)) FROM gold_evidence ge WHERE ge.question_id=$2`, runID, question.id).Scan(&total, &hits)
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
				answer, citations, usage, err := generate(ctx, client, cfg.GeneratorModel, question.text, results, cfg.ContextChars, cfg.ContextTokens)
				if err != nil {
					_ = tx.Rollback(ctx)
					return metrics, fail(ctx, pool, runID, err)
				}
				_, err = tx.Exec(ctx, `INSERT INTO answers(experiment_run_id,question_id,text,citations,usage) VALUES($1,$2,$3,$4,$5)`, runID, question.id, answer, citations, usage)
				if err != nil {
					_ = tx.Rollback(ctx)
					return metrics, fail(ctx, pool, runID, err)
				}
			}
			if err = tx.Commit(ctx); err != nil {
				return metrics, fail(ctx, pool, runID, err)
			}
		}
	}
	_, err = pool.Exec(ctx, `UPDATE experiment_runs SET status='complete',completed_at=now() WHERE id=$1`, runID)
	return metrics, err
}

type generatedAnswer struct {
	Answer         string `json:"answer"`
	CitedSourceIDs []int  `json:"cited_source_ids"`
}

type storedSource struct {
	ID           int    `json:"id"`
	DocumentID   string `json:"document_id"`
	TurnStart    int    `json:"turn_start"`
	TurnEnd      int    `json:"turn_end"`
	Text         string `json:"text"`
	Quote        string `json:"quote,omitempty"`
	SpeakerLabel string `json:"speaker_label,omitempty"`
}

type storedCitations struct {
	CitedSourceIDs []int          `json:"cited_source_ids"`
	Retrieved      []storedSource `json:"retrieved"`
}

func generate(ctx context.Context, client *fireworks.Client, model, question string, results []retrieval.Result, contextChars, contextTokens int) (string, []byte, []byte, error) {
	if model == "" {
		return "", nil, nil, fmt.Errorf("--generator-model is required with --generate")
	}
	context, packed, estimatedTokens, err := packSources(results, contextChars, contextTokens)
	if err != nil {
		return "", nil, nil, err
	}
	system := `Answer only from the supplied transcript sources. Return ONLY JSON with this exact schema: {"answer":"...","cited_source_ids":[1,2]}. cited_source_ids must contain only source IDs that directly support the answer; use [] when the sources do not establish an answer. When the sources do not establish an answer, set answer exactly to "Insufficient source evidence." Do not mention source IDs in answer.`
	result, err := client.ChatJSONNoReasoning(ctx, model, system, fmt.Sprintf("Question: %s\n\nSources:\n%s", question, context), 600)
	if err != nil {
		return "", nil, nil, err
	}
	var parsed generatedAnswer
	if err = json.Unmarshal([]byte(result.Text), &parsed); err != nil {
		return "", nil, nil, fmt.Errorf("answer model returned invalid structured answer %q: %w", result.Text, err)
	}
	if strings.TrimSpace(parsed.Answer) == "" {
		parsed.Answer = "Insufficient source evidence."
	}
	valid := make(map[int]bool, len(packed))
	for _, source := range packed {
		valid[source.ID] = true
	}
	seen := map[int]bool{}
	for _, id := range parsed.CitedSourceIDs {
		if !valid[id] {
			return "", nil, nil, fmt.Errorf("answer model cited unavailable source S%d", id)
		}
		if seen[id] {
			return "", nil, nil, fmt.Errorf("answer model cited S%d more than once", id)
		}
		seen[id] = true
	}
	citations, _ := json.Marshal(storedCitations{CitedSourceIDs: parsed.CitedSourceIDs, Retrieved: packed})
	usage, _ := json.Marshal(map[string]any{"context_chars": len(context), "context_estimated_tokens": estimatedTokens, "context_token_budget": contextTokens, "context_tokenizer": tokenEncodingName(contextTokens), "prompt_tokens": result.PromptTokens, "completion_tokens": result.CompletionTokens, "total_tokens": result.TotalTokens})
	return parsed.Answer, citations, usage, nil
}

func packSources(results []retrieval.Result, contextChars, contextTokens int) (string, []storedSource, int, error) {
	var context strings.Builder
	var packed []storedSource
	encoder, err := tokenEncoder(contextTokens)
	if err != nil {
		return "", nil, 0, err
	}
	usedTokens := 0
	for i, r := range results {
		source := storedSource{ID: i + 1, DocumentID: r.ExternalID, TurnStart: r.TurnStart, TurnEnd: r.TurnEnd, Text: r.Text, Quote: r.Quote, SpeakerLabel: r.SpeakerLabel}
		text := renderSource(source)
		if contextChars > 0 && len(context.String())+len(text) > contextChars {
			if context.Len() == 0 {
				context.WriteString(text[:min(len(text), contextChars)])
				packed = append(packed, source)
			}
			break
		}
		if contextTokens > 0 {
			tokens := encoder.Encode(text, nil, nil)
			if usedTokens+len(tokens) > contextTokens {
				if context.Len() == 0 {
					remaining := min(len(tokens), contextTokens)
					context.WriteString(encoder.Decode(tokens[:remaining]))
					usedTokens = remaining
					packed = append(packed, source)
				}
				break
			}
			usedTokens += len(tokens)
		}
		context.WriteString(text)
		packed = append(packed, source)
	}
	if contextTokens == 0 {
		usedTokens = len(encoder.Encode(context.String(), nil, nil))
	}
	return context.String(), packed, usedTokens, nil
}

func renderSource(source storedSource) string {
	var unit strings.Builder
	fmt.Fprintf(&unit, "[S%d document=%s turns=%d-%d]\nSOURCE TEXT: %s\n", source.ID, source.DocumentID, source.TurnStart, source.TurnEnd, source.Text)
	if source.Quote != "" {
		fmt.Fprintf(&unit, "VERIFIED QUOTE")
		if source.SpeakerLabel != "" {
			fmt.Fprintf(&unit, " (speaker: %s)", source.SpeakerLabel)
		}
		fmt.Fprintf(&unit, ": %s\n", source.Quote)
	}
	unit.WriteString("\n")
	return unit.String()
}

func tokenEncodingName(contextTokens int) string {
	if contextTokens == 0 {
		return ""
	}
	return "cl100k_base"
}

func tokenEncoder(contextTokens int) (*tiktoken.Tiktoken, error) {
	return tiktoken.GetEncoding("cl100k_base")
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
