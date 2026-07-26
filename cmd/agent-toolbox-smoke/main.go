// agent-toolbox-smoke runs an exploratory tool-use check. It is an evaluator,
// not an agent-facing binary.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const toolSelectionSystem = `You are a Q&A agent. You have exactly one executable tool:
agent-toolbox-fixed-rag search --query <text> --top-k <integer from 1 to 5>
It searches only the fixed-chunk RAG condition and returns JSON evidence. Before answering, choose one search request. Return ONLY this JSON object: {"query":"...","top_k":1}. Do not answer the question yet, and do not mention any other tool or command.`

const answerSystem = `You are a Q&A agent. Answer the user's question using only the supplied JSON tool output. Return ONLY JSON shaped as {"answer":"one concise answer or Insufficient retrieved evidence.","citation":"external_id turns start-end or empty"}. Do not explain your reasoning. Do not claim to have used any tool other than the supplied result.`

type question struct {
	ID, Text, ExpectedAnswer, ExpectedExternalID string
}

var smokeQuestions = []question{
	{ID: "toolbox-smoke-01", Text: "What kind of engineer did Ann Kiremidjian identify herself as?", ExpectedAnswer: "structural engineer", ExpectedExternalID: "NPR-454"},
	{ID: "toolbox-smoke-02", Text: "What did Ann Kiremidjian say non-ductile concrete frames have very little of?", ExpectedAnswer: "steel reinforcement", ExpectedExternalID: "NPR-454"},
}

type toolRequest struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k"`
}

type answerResponse struct {
	Answer   string `json:"answer"`
	Citation string `json:"citation"`
}

type result struct {
	QuestionID        string          `json:"question_id"`
	Question          string          `json:"question"`
	ToolSelectionRaw  string          `json:"tool_selection_raw"`
	ToolRequest       toolRequest     `json:"tool_request"`
	ToolOutput        json.RawMessage `json:"tool_output,omitempty"`
	Answer            string          `json:"answer,omitempty"`
	Error             string          `json:"error,omitempty"`
	ExpectedChunkHit  bool            `json:"expected_chunk_hit"`
	ExpectedAnswerHit bool            `json:"expected_answer_surface_hit"`
	ExpectedCitation  bool            `json:"expected_citation_surface_hit"`
}

func main() {
	var projectDir, model, binary, questionsPath string
	var limit int
	var seed int64
	root := &cobra.Command{Use: "agent-toolbox-smoke", Short: "Run an exploratory model tool-use smoke check", SilenceErrors: true, SilenceUsage: true}
	root.RunE = func(cmd *cobra.Command, args []string) error {
		absoluteProjectDir, err := filepath.Abs(projectDir)
		if err != nil {
			return fmt.Errorf("resolve --project-dir: %w", err)
		}
		absoluteBinary, err := filepath.Abs(binary)
		if err != nil {
			return fmt.Errorf("resolve --binary: %w", err)
		}
		if err := loadEnv(absoluteProjectDir); err != nil {
			return err
		}
		client, err := fireworks.NewFromEnv()
		if err != nil {
			return err
		}
		if _, err := os.Stat(absoluteBinary); err != nil {
			return fmt.Errorf("inspect --binary: %w", err)
		}
		questions, err := selectedQuestions(questionsPath, limit, seed)
		if err != nil {
			return err
		}
		var results []result
		for _, q := range questions {
			row := runQuestion(cmd.Context(), client, model, absoluteBinary, absoluteProjectDir, q)
			results = append(results, row)
		}
		output, err := json.MarshalIndent(struct {
			Model         string   `json:"model"`
			Binary        string   `json:"binary"`
			QuestionsPath string   `json:"questions_path,omitempty"`
			Seed          int64    `json:"seed"`
			Results       []result `json:"results"`
		}{Model: model, Binary: absoluteBinary, QuestionsPath: questionsPath, Seed: seed, Results: results}, "", "  ")
		if err != nil {
			return err
		}
		cmd.Println(string(output))
		return nil
	}
	root.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory containing .env")
	root.Flags().StringVar(&model, "model", "accounts/fireworks/models/deepseek-v4-flash", "Fireworks model acting as the Q&A agent")
	root.Flags().StringVar(&binary, "binary", "./bin/agent-toolbox-fixed-rag", "Capability-limited binary the agent may use")
	root.Flags().StringVar(&questionsPath, "questions", "", "Optional JSONL questions with reference answers and evidence")
	root.Flags().IntVar(&limit, "limit", 0, "Number of questions to sample; 0 uses the built-in two-question smoke set")
	root.Flags().Int64Var(&seed, "seed", 20260724, "Deterministic question-sampling seed")
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runQuestion(ctx context.Context, client *fireworks.Client, model, binary, projectDir string, q question) result {
	selection, err := chatJSON(ctx, client, model, toolSelectionSystem, q.Text, 120)
	if err != nil {
		return result{QuestionID: q.ID, Question: q.Text, Error: fmt.Sprintf("request tool selection: %v", err)}
	}
	row := result{QuestionID: q.ID, Question: q.Text, ToolSelectionRaw: selection.Text}
	var request toolRequest
	if err := json.Unmarshal([]byte(selection.Text), &request); err != nil {
		row.Error = fmt.Sprintf("decode tool request: %v", err)
		return row
	}
	request.Query = strings.TrimSpace(request.Query)
	if request.Query == "" || request.TopK < 1 || request.TopK > 5 {
		row.ToolRequest = request
		row.Error = fmt.Sprintf("invalid tool request %#v", request)
		return row
	}
	row.ToolRequest = request
	tool := exec.CommandContext(ctx, binary, "search", "--project-dir", projectDir, "--query", request.Query, "--top-k", strconv.Itoa(request.TopK))
	tool.Dir = projectDir
	toolOutput, err := tool.CombinedOutput()
	if err != nil {
		row.Error = fmt.Sprintf("execute capability-limited search: %v", err)
		return row
	}
	if !json.Valid(toolOutput) {
		row.Error = fmt.Sprintf("search did not return JSON: %q", string(toolOutput))
		return row
	}
	row.ToolOutput = json.RawMessage(toolOutput)
	answer, err := chatJSON(ctx, client, model, answerSystem, fmt.Sprintf("Question: %s\n\nTool output:\n%s", q.Text, toolOutput), 120)
	if err != nil {
		row.Error = fmt.Sprintf("request final answer: %v", err)
		return row
	}
	var parsed answerResponse
	if err := json.Unmarshal([]byte(answer.Text), &parsed); err != nil {
		row.Error = fmt.Sprintf("decode final answer: %v", err)
		return row
	}
	answerText := strings.TrimSpace(parsed.Answer)
	if citation := strings.TrimSpace(parsed.Citation); citation != "" {
		answerText += " [" + citation + "]"
	}
	toolText := string(toolOutput)
	row.Answer = answerText
	row.ExpectedChunkHit = strings.Contains(toolText, q.ExpectedExternalID)
	row.ExpectedAnswerHit = strings.Contains(strings.ToLower(answerText), strings.ToLower(q.ExpectedAnswer))
	row.ExpectedCitation = strings.Contains(answerText, q.ExpectedExternalID)
	return row
}

func selectedQuestions(path string, limit int, seed int64) ([]question, error) {
	if path == "" {
		if limit != 0 && limit != len(smokeQuestions) {
			return nil, fmt.Errorf("--limit requires --questions, or use %d for the built-in smoke set", len(smokeQuestions))
		}
		return smokeQuestions, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open --questions: %w", err)
	}
	defer file.Close()
	type evidence struct {
		DocumentExternalID string `json:"document_external_id"`
	}
	type record struct {
		ID              string     `json:"external_id"`
		Text            string     `json:"text"`
		ReferenceAnswer string     `json:"reference_answer"`
		Evidence        []evidence `json:"evidence"`
	}
	var questions []question
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), 1<<20)
	for scanner.Scan() {
		var record record
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, fmt.Errorf("decode --questions: %w", err)
		}
		if record.ID == "" || record.Text == "" || record.ReferenceAnswer == "" || len(record.Evidence) == 0 || record.Evidence[0].DocumentExternalID == "" {
			return nil, fmt.Errorf("invalid question record %q", record.ID)
		}
		questions = append(questions, question{ID: record.ID, Text: record.Text, ExpectedAnswer: record.ReferenceAnswer, ExpectedExternalID: record.Evidence[0].DocumentExternalID})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read --questions: %w", err)
	}
	if limit < 1 || limit > len(questions) {
		return nil, fmt.Errorf("--limit must be between 1 and %d for --questions", len(questions))
	}
	rng := rand.New(rand.NewSource(seed))
	rng.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })
	return questions[:limit], nil
}

func chatJSON(ctx context.Context, client *fireworks.Client, model, system, user string, maxTokens int) (fireworks.ChatResult, error) {
	if model == "accounts/fireworks/models/deepseek-v4-flash" || model == "accounts/fireworks/models/kimi-k2p6" {
		return client.ChatJSONNoReasoning(ctx, model, system, user, maxTokens)
	}
	return client.ChatJSON(ctx, model, system, user, maxTokens)
}

func loadEnv(projectDir string) error {
	abs, err := filepath.Abs(projectDir)
	if err != nil {
		return err
	}
	path := filepath.Join(abs, ".env")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	return godotenv.Load(path)
}
