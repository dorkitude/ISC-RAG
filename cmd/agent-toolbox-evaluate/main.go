// agent-toolbox-evaluate runs fresh, capability-limited Q&A sessions and
// records every attempt in PostgreSQL for later audit.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dorkitude/isc-rag/internal/evaluation"
	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const toolSelectionSystem = `You are a Q&A agent. You have exactly one executable tool:
%s search --query <text> --top-k <integer from 1 to 5>
It searches only this fixed RAG capability and returns JSON evidence. Before answering, choose one search request. Return ONLY this JSON object: {"query":"...","top_k":1}. Do not answer the question yet, and do not mention any other tool or command.`

const answerSystem = `You are a Q&A agent. Answer the user's question using only the supplied JSON tool output. Return ONLY JSON shaped as {"answer":"one concise answer or Insufficient retrieved evidence.","citation":"external_id turns start-end or empty"}. Cite the turn range shown in the retrieved evidence. Do not explain reasoning.`

type question struct {
	ID, Text, Reference string
	Gold                []evidence
}
type evidence struct {
	ExternalID string
	Start, End int
}
type toolRequest struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k"`
}
type answerResponse struct {
	Answer   string `json:"answer"`
	Citation string `json:"citation"`
}
type retrieved struct {
	SourceUnitID string  `json:"source_unit_id"`
	ExternalID   string  `json:"external_id"`
	TurnStart    int     `json:"turn_start"`
	TurnEnd      int     `json:"turn_end"`
	Score        float64 `json:"score"`
}
type toolOutput struct {
	Results []retrieved `json:"results"`
}
type attempt struct {
	Question                                                       question
	Replicate                                                      int
	ToolRequest                                                    toolRequest
	RawSelection                                                   string
	RawTool                                                        json.RawMessage
	RawAnswer                                                      string
	Answer, Citation, Error                                        string
	ToolRequestValid, EvidenceHit, CitationCorrect, ReferenceMatch bool
	SelectionMS, SearchMS, AnswerMS, EndMS                         int64
	Results                                                        []retrieved
}

var citationPattern = regexp.MustCompile(`(?i)([^\s]+)\s+turns\s+(\d+)-(\d+)`)

func main() {
	var projectDir, questionsPath, questionSplit, model, binary, capability, representation, version, runName string
	var repetitions, workers, timeoutSeconds, limit int
	root := &cobra.Command{Use: "agent-toolbox-evaluate", SilenceErrors: true, SilenceUsage: true}
	root.RunE = func(cmd *cobra.Command, _ []string) error {
		projectDir, _ = filepath.Abs(projectDir)
		if err := loadEnv(projectDir); err != nil {
			return err
		}
		if questionsPath != "" {
			pool, err := storage.Open(cmd.Context(), os.Getenv("DATABASE_URL"))
			if err != nil {
				return err
			}
			defer pool.Close()
			if err := storage.Migrate(cmd.Context(), pool); err != nil {
				return err
			}
			if _, err := evaluation.Import(cmd.Context(), pool, questionsPath); err != nil {
				return fmt.Errorf("import questions: %w", err)
			}
		}
		pool, err := storage.Open(cmd.Context(), os.Getenv("DATABASE_URL"))
		if err != nil {
			return err
		}
		defer pool.Close()
		if err := storage.Migrate(cmd.Context(), pool); err != nil {
			return err
		}
		questions, err := loadQuestions(cmd.Context(), pool, questionSplit)
		if err != nil {
			return err
		}
		if len(questions) == 0 {
			return fmt.Errorf("no questions with corpus_split %q", questionSplit)
		}
		if limit > 0 && limit < len(questions) {
			questions = questions[:limit]
		}
		absoluteBinary, err := filepath.Abs(binary)
		if err != nil {
			return err
		}
		if _, err = os.Stat(absoluteBinary); err != nil {
			return fmt.Errorf("inspect binary: %w", err)
		}
		client, err := fireworks.NewFromEnv()
		if err != nil {
			return err
		}
		selectionPrompt := fmt.Sprintf(toolSelectionSystem, filepath.Base(absoluteBinary))
		runID, err := createRun(cmd.Context(), pool, runName, questionSplit, model, capability, representation, version, repetitions, selectionPrompt)
		if err != nil {
			return err
		}
		fmt.Printf("Started agent evaluation run %s (%d questions × %d fresh sessions, model=%s, capability=%s)\n", runID, len(questions), repetitions, model, capability)
		jobs := make(chan attempt)
		var wg sync.WaitGroup
		var firstErr error
		var errMu sync.Mutex
		for range workers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for job := range jobs {
					ctx, cancel := context.WithTimeout(cmd.Context(), time.Duration(timeoutSeconds)*time.Second)
					row := runAttempt(ctx, client, model, absoluteBinary, projectDir, selectionPrompt, job)
					cancel()
					if err := persistAttempt(cmd.Context(), pool, runID, row); err != nil {
						errMu.Lock()
						if firstErr == nil {
							firstErr = err
						}
						errMu.Unlock()
					}
				}
			}()
		}
		for _, q := range questions {
			for replicate := 1; replicate <= repetitions; replicate++ {
				jobs <- attempt{Question: q, Replicate: replicate}
			}
		}
		close(jobs)
		wg.Wait()
		errMu.Lock()
		err = firstErr
		errMu.Unlock()
		status := "complete"
		if err != nil {
			status = "failed"
		}
		_, updateErr := pool.Exec(cmd.Context(), `UPDATE agent_evaluation_runs SET status=$2,completed_at=now() WHERE id=$1`, runID, status)
		if err != nil {
			return err
		}
		if updateErr != nil {
			return updateErr
		}
		fmt.Printf("Completed agent evaluation run %s\n", runID)
		return nil
	}
	root.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	root.Flags().StringVar(&questionsPath, "questions", "", "Optional JSONL questions to import before the run")
	root.Flags().StringVar(&questionSplit, "question-split", "agent-dev-20260725", "questions.corpus_split namespace")
	root.Flags().StringVar(&model, "model", "accounts/fireworks/models/deepseek-v4-flash", "Fireworks Q&A model")
	root.Flags().StringVar(&binary, "binary", "", "Capability-limited agent-toolbox binary")
	root.Flags().StringVar(&capability, "capability", "", "Immutable capability name recorded with the run")
	root.Flags().StringVar(&representation, "representation", "", "Representation recorded with the run")
	root.Flags().StringVar(&version, "version", "", "Representation version recorded with the run")
	root.Flags().StringVar(&runName, "name", "agent-toolbox-development-20260725", "Human-readable run name")
	root.Flags().IntVar(&repetitions, "repetitions", 10, "Fresh sessions per question")
	root.Flags().IntVar(&limit, "limit", 0, "Cap questions for a bounded integration check; 0 uses all")
	root.Flags().IntVar(&workers, "workers", 6, "Concurrent fresh sessions")
	root.Flags().IntVar(&timeoutSeconds, "timeout-seconds", 30, "Per-session timeout")
	root.PreRunE = func(cmd *cobra.Command, _ []string) error {
		if binary == "" || capability == "" || representation == "" || version == "" {
			return fmt.Errorf("--binary, --capability, --representation, and --version are required")
		}
		if repetitions < 1 || workers < 1 || timeoutSeconds < 1 {
			return fmt.Errorf("repetitions, workers, and timeout-seconds must be positive")
		}
		return nil
	}
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadEnv(projectDir string) error {
	path := filepath.Join(projectDir, ".env")
	if _, err := os.Stat(path); err == nil {
		if err := godotenv.Load(path); err != nil {
			return err
		}
	}
	if strings.TrimSpace(os.Getenv("DATABASE_URL")) == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	return nil
}

func loadQuestions(ctx context.Context, pool *pgxpool.Pool, split string) ([]question, error) {
	rows, err := pool.Query(ctx, `SELECT q.id,q.text,COALESCE(q.reference_answer,''),COALESCE(jsonb_agg(jsonb_build_object('external_id',d.external_id,'turn_start',ge.turn_start,'turn_end',ge.turn_end)) FILTER (WHERE ge.id IS NOT NULL),'[]'::jsonb) FROM questions q LEFT JOIN gold_evidence ge ON ge.question_id=q.id LEFT JOIN documents d ON d.id=ge.document_id WHERE q.corpus_split=$1 GROUP BY q.id,q.text,q.reference_answer ORDER BY q.id`, split)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []question
	for rows.Next() {
		var q question
		var raw []byte
		if err := rows.Scan(&q.ID, &q.Text, &q.Reference, &raw); err != nil {
			return nil, err
		}
		var encoded []struct {
			ExternalID string `json:"external_id"`
			Start      int    `json:"turn_start"`
			End        int    `json:"turn_end"`
		}
		if err := json.Unmarshal(raw, &encoded); err != nil {
			return nil, err
		}
		for _, e := range encoded {
			q.Gold = append(q.Gold, evidence{ExternalID: e.ExternalID, Start: e.Start, End: e.End})
		}
		if len(q.Gold) == 0 {
			return nil, fmt.Errorf("question %s has no gold evidence", q.ID)
		}
		out = append(out, q)
	}
	return out, rows.Err()
}

func createRun(ctx context.Context, pool *pgxpool.Pool, name, split, model, capability, representation, version string, repetitions int, selectionPrompt string) (string, error) {
	parameters, _ := json.Marshal(map[string]any{"protocol": "fresh two-call JSON selection then answer; host executes capability binary", "tool_top_k_max": 5, "reasoning_effort": "none"})
	var id string
	err := pool.QueryRow(ctx, `INSERT INTO agent_evaluation_runs(name,question_split,model,capability,representation,representation_version,embedding_model,repetitions,tool_selection_system,answer_system,parameters) VALUES($1,$2,$3,$4,$5,$6,'nomic-ai/nomic-embed-text-v1.5',$7,$8,$9,$10) RETURNING id`, name, split, model, capability, representation, version, repetitions, selectionPrompt, answerSystem, parameters).Scan(&id)
	return id, err
}

func runAttempt(ctx context.Context, client *fireworks.Client, model, binary, projectDir, selectionPrompt string, row attempt) attempt {
	started := time.Now()
	t := time.Now()
	selection, err := client.ChatJSONNoReasoning(ctx, model, selectionPrompt, row.Question.Text, 120)
	row.SelectionMS = time.Since(t).Milliseconds()
	if err != nil {
		row.Error = "selection: " + err.Error()
		row.EndMS = time.Since(started).Milliseconds()
		return row
	}
	row.RawSelection = selection.Text
	if err := json.Unmarshal([]byte(selection.Text), &row.ToolRequest); err != nil || strings.TrimSpace(row.ToolRequest.Query) == "" || row.ToolRequest.TopK < 1 || row.ToolRequest.TopK > 5 {
		row.Error = "invalid tool request"
		row.EndMS = time.Since(started).Milliseconds()
		return row
	}
	row.ToolRequestValid = true
	t = time.Now()
	command := exec.CommandContext(ctx, binary, "search", "--project-dir", projectDir, "--query", row.ToolRequest.Query, "--top-k", strconv.Itoa(row.ToolRequest.TopK))
	command.Dir = projectDir
	raw, err := command.CombinedOutput()
	row.SearchMS = time.Since(t).Milliseconds()
	if err != nil {
		row.Error = "search: " + err.Error()
		row.EndMS = time.Since(started).Milliseconds()
		return row
	}
	row.RawTool = append(json.RawMessage(nil), raw...)
	var tool toolOutput
	if err := json.Unmarshal(raw, &tool); err != nil {
		row.Error = "decode search: " + err.Error()
		row.EndMS = time.Since(started).Milliseconds()
		return row
	}
	row.Results = tool.Results
	row.EvidenceHit = evidenceHit(row.Question.Gold, tool.Results)
	t = time.Now()
	completion, err := client.ChatJSONNoReasoning(ctx, model, answerSystem, fmt.Sprintf("Question: %s\n\nTool output:\n%s", row.Question.Text, raw), 120)
	row.AnswerMS = time.Since(t).Milliseconds()
	if err != nil {
		row.Error = "answer: " + err.Error()
		row.EndMS = time.Since(started).Milliseconds()
		return row
	}
	row.RawAnswer = completion.Text
	var answer answerResponse
	if err := json.Unmarshal([]byte(completion.Text), &answer); err != nil {
		row.Error = "decode answer"
		row.EndMS = time.Since(started).Milliseconds()
		return row
	}
	row.Answer, row.Citation = strings.TrimSpace(answer.Answer), strings.TrimSpace(answer.Citation)
	row.CitationCorrect = citationCorrect(row.Citation, row.Question.Gold)
	row.ReferenceMatch = referenceMatch(row.Answer, row.Question.Reference)
	row.EndMS = time.Since(started).Milliseconds()
	return row
}

func evidenceHit(gold []evidence, got []retrieved) bool {
	for _, e := range gold {
		for _, r := range got {
			if r.ExternalID == e.ExternalID && r.TurnStart <= e.Start && r.TurnEnd >= e.End {
				return true
			}
		}
	}
	return false
}
func citationCorrect(citation string, gold []evidence) bool {
	for _, m := range citationPattern.FindAllStringSubmatch(citation, -1) {
		start, a := strconv.Atoi(m[2])
		end, b := strconv.Atoi(m[3])
		if a != nil || b != nil {
			continue
		}
		for _, e := range gold {
			if strings.EqualFold(m[1], e.ExternalID) && start <= e.Start && end >= e.End {
				return true
			}
		}
	}
	return false
}
func referenceMatch(answer, reference string) bool {
	return reference != "" && strings.Contains(strings.ToLower(answer), strings.ToLower(reference))
}

func persistAttempt(ctx context.Context, pool *pgxpool.Pool, runID string, row attempt) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var id string
	err = tx.QueryRow(ctx, `INSERT INTO agent_evaluation_attempts(run_id,question_id,replicate,tool_request,raw_tool_selection,raw_tool_output,raw_answer_completion,answer_text,citation,tool_request_valid,evidence_hit,citation_correct,reference_answer_match,selection_latency_ms,search_latency_ms,answer_latency_ms,end_to_end_latency_ms,error,completed_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,now()) RETURNING id`, runID, row.Question.ID, row.Replicate, jsonb(row.ToolRequest), nullIfEmpty(row.RawSelection), jsonb(row.RawTool), nullIfEmpty(row.RawAnswer), nullIfEmpty(row.Answer), nullIfEmpty(row.Citation), row.ToolRequestValid, row.EvidenceHit, row.CitationCorrect, row.ReferenceMatch, nullableInt(row.SelectionMS), nullableInt(row.SearchMS), nullableInt(row.AnswerMS), nullableInt(row.EndMS), nullIfEmpty(row.Error)).Scan(&id)
	if err != nil {
		return err
	}
	for rank, r := range row.Results {
		if _, err := tx.Exec(ctx, `INSERT INTO agent_evaluation_retrievals(attempt_id,source_unit_id,rank,score) VALUES($1,$2,$3,$4)`, id, r.SourceUnitID, rank+1, r.Score); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
func jsonb(value any) any {
	body, _ := json.Marshal(value)
	if string(body) == "null" {
		return nil
	}
	return body
}
func nullIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}
func nullableInt(value int64) any {
	if value == 0 {
		return nil
	}
	return value
}
