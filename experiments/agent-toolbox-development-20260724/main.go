// agent-toolbox-development is a one-off development harness. It keeps the
// model behind the fixed-only toolbox and records raw completions locally.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const selectSystem = `You are a Q&A agent. You have exactly one executable tool:
agent-toolbox-fixed-rag search --query <text> --top-k <integer from 1 to 5>
Before answering, choose one search request. Return ONLY this JSON object: {"query":"...","top_k":1}. Do not answer the question yet.`

const answerSystem = `Answer the user's question using only the supplied JSON tool output. Return ONLY JSON shaped as {"answer":"one concise answer or Insufficient retrieved evidence.","citation":"external_id turns start-end or empty"}. Cite the turn range shown in the retrieved evidence. Do not explain reasoning.`

type question struct {
	ID, Text, ReferenceAnswer string
	Aliases                   []string `json:"answer_aliases"`
	Evidence                  []evidence
}
type evidence struct {
	DocumentExternalID string `json:"document_external_id"`
	TurnStart          int    `json:"turn_start"`
	TurnEnd            int    `json:"turn_end"`
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
	ExternalID string `json:"external_id"`
	TurnStart  int    `json:"turn_start"`
	TurnEnd    int    `json:"turn_end"`
}
type toolResponse struct {
	Results []retrieved `json:"results"`
}
type result struct {
	QuestionID         string      `json:"question_id"`
	ToolRequestValid   bool        `json:"tool_request_valid"`
	EvidenceHit        bool        `json:"evidence_hit"`
	StrictAnswerHit    bool        `json:"strict_answer_hit"`
	CitationCorrect    bool        `json:"citation_correct"`
	SelectionLatencyMS int64       `json:"selection_latency_ms"`
	SearchLatencyMS    int64       `json:"search_latency_ms"`
	AnswerLatencyMS    int64       `json:"answer_latency_ms"`
	EndToEndLatencyMS  int64       `json:"end_to_end_latency_ms"`
	ToolRequest        toolRequest `json:"tool_request,omitempty"`
	Answer             string      `json:"answer,omitempty"`
	Citation           string      `json:"citation,omitempty"`
	RawToolOutput      string      `json:"raw_tool_output,omitempty"`
	Error              string      `json:"error,omitempty"`
}

func main() {
	var projectDir, model, questionsPath, output string
	var limit int
	root := &cobra.Command{Use: "agent-toolbox-development", SilenceErrors: true, SilenceUsage: true}
	root.RunE = func(cmd *cobra.Command, _ []string) error {
		projectDir, _ = filepath.Abs(projectDir)
		if err := loadEnv(projectDir); err != nil {
			return err
		}
		client, err := fireworks.NewFromEnv()
		if err != nil {
			return err
		}
		questions, err := loadQuestions(questionsPath)
		if err != nil {
			return err
		}
		if limit > 0 && limit < len(questions) {
			questions = questions[:limit]
		}
		binary := filepath.Join(projectDir, "bin", "agent-toolbox-fixed-rag")
		if _, err := os.Stat(binary); err != nil {
			return fmt.Errorf("inspect fixed toolbox: %w", err)
		}
		rows := make([]result, 0, len(questions))
		for _, q := range questions {
			rows = append(rows, run(cmd.Context(), client, model, binary, projectDir, q))
		}
		payload, _ := json.MarshalIndent(map[string]any{"protocol": "one fixed-only search then answer; no-reasoning", "model": model, "results": rows}, "", "  ")
		if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
			return err
		}
		return os.WriteFile(output, payload, 0o600)
	}
	root.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	root.Flags().StringVar(&model, "model", "accounts/fireworks/models/deepseek-v4-flash", "Fireworks Q&A model")
	root.Flags().StringVar(&questionsPath, "questions", "artifacts/agent-toolbox-development-20260724/questions.jsonl", "local development question file")
	root.Flags().StringVar(&output, "output", "artifacts/agent-toolbox-development-20260724/results.json", "ignored raw result path")
	root.Flags().IntVar(&limit, "limit", 0, "Cap questions for a bounded diagnostic run; 0 uses all")
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, client *fireworks.Client, model, binary, projectDir string, q question) result {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	started := time.Now()
	row := result{QuestionID: q.ID}
	t := time.Now()
	selection, err := client.ChatJSONNoReasoning(ctx, model, selectSystem, q.Text, 120)
	row.SelectionLatencyMS = time.Since(t).Milliseconds()
	if err != nil {
		row.Error = "selection: " + err.Error()
		row.EndToEndLatencyMS = time.Since(started).Milliseconds()
		return row
	}
	var request toolRequest
	if err := json.Unmarshal([]byte(selection.Text), &request); err != nil || strings.TrimSpace(request.Query) == "" || request.TopK < 1 || request.TopK > 5 {
		row.Error = "invalid tool request"
		row.EndToEndLatencyMS = time.Since(started).Milliseconds()
		return row
	}
	row.ToolRequestValid = true
	row.ToolRequest = request
	t = time.Now()
	command := exec.CommandContext(ctx, binary, "search", "--project-dir", projectDir, "--query", request.Query, "--top-k", strconv.Itoa(request.TopK))
	command.Dir = projectDir
	raw, err := command.CombinedOutput()
	row.SearchLatencyMS = time.Since(t).Milliseconds()
	if err != nil {
		row.Error = "search: " + err.Error()
		row.EndToEndLatencyMS = time.Since(started).Milliseconds()
		return row
	}
	var retrievedResults toolResponse
	if err := json.Unmarshal(raw, &retrievedResults); err != nil {
		row.RawToolOutput = string(raw)
		row.Error = "decode search: " + err.Error()
		row.EndToEndLatencyMS = time.Since(started).Milliseconds()
		return row
	}
	row.EvidenceHit = evidenceHit(q.Evidence, retrievedResults.Results)
	t = time.Now()
	completion, err := client.ChatJSONNoReasoning(ctx, model, answerSystem, fmt.Sprintf("Question: %s\n\nTool output:\n%s", q.Text, raw), 120)
	row.AnswerLatencyMS = time.Since(t).Milliseconds()
	if err != nil {
		row.Error = "answer: " + err.Error()
		row.EndToEndLatencyMS = time.Since(started).Milliseconds()
		return row
	}
	var answer answerResponse
	if err := json.Unmarshal([]byte(completion.Text), &answer); err != nil {
		row.Error = "decode answer"
		row.EndToEndLatencyMS = time.Since(started).Milliseconds()
		return row
	}
	row.Answer = strings.TrimSpace(answer.Answer)
	row.Citation = strings.TrimSpace(answer.Citation)
	row.StrictAnswerHit = aliasHit(row.Answer, q.Aliases)
	row.CitationCorrect = citationCorrect(row.Citation, q.Evidence)
	row.EndToEndLatencyMS = time.Since(started).Milliseconds()
	return row
}

func evidenceHit(expected []evidence, got []retrieved) bool {
	for _, e := range expected {
		for _, r := range got {
			if r.ExternalID == e.DocumentExternalID && r.TurnStart <= e.TurnStart && r.TurnEnd >= e.TurnEnd {
				return true
			}
		}
	}
	return false
}
func citationCorrect(citation string, expected []evidence) bool {
	matches := citationPattern.FindAllStringSubmatch(citation, -1)
	for _, match := range matches {
		start, startErr := strconv.Atoi(match[2])
		end, endErr := strconv.Atoi(match[3])
		if startErr != nil || endErr != nil {
			continue
		}
		for _, e := range expected {
			if strings.EqualFold(match[1], e.DocumentExternalID) && start <= e.TurnStart && end >= e.TurnEnd {
				return true
			}
		}
	}
	return false
}

var citationPattern = regexp.MustCompile(`(?i)(NPR-\d+)\s+turns\s+(\d+)-(\d+)`)

func aliasHit(answer string, aliases []string) bool {
	normalized := strings.ToLower(answer)
	for _, a := range aliases {
		if strings.Contains(normalized, strings.ToLower(a)) {
			return true
		}
	}
	return false
}
func loadQuestions(path string) ([]question, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var out []question
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), 1<<20)
	for scanner.Scan() {
		var q question
		if err := json.Unmarshal(scanner.Bytes(), &q); err != nil {
			return nil, err
		}
		var raw struct {
			ID              string     `json:"external_id"`
			Text            string     `json:"text"`
			ReferenceAnswer string     `json:"reference_answer"`
			Aliases         []string   `json:"answer_aliases"`
			Evidence        []evidence `json:"evidence"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &raw); err != nil {
			return nil, err
		}
		q.ID, q.Text, q.ReferenceAnswer, q.Aliases, q.Evidence = raw.ID, raw.Text, raw.ReferenceAnswer, raw.Aliases, raw.Evidence
		out = append(out, q)
	}
	return out, scanner.Err()
}
func loadEnv(projectDir string) error {
	path := filepath.Join(projectDir, ".env")
	if _, err := os.Stat(path); err == nil {
		return godotenv.Load(path)
	}
	return nil
}
