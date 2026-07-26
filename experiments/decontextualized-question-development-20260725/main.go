// decontextualized-question-development is a retrieval-only diagnostic for
// source-grounded development questions. It is not a held-out paper evaluator.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/retrieval"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const embeddingModel = "nomic-ai/nomic-embed-text-v1.5"

type evidence struct {
	DocumentExternalID string `json:"document_external_id"`
	TurnStart          int    `json:"turn_start"`
	TurnEnd            int    `json:"turn_end"`
}

type question struct {
	ID   string     `json:"external_id"`
	Type string     `json:"question_type"`
	Text string     `json:"text"`
	Gold []evidence `json:"evidence"`
}

type condition struct {
	Name, Representation, Version string
}

type row struct {
	QuestionID   string `json:"question_id"`
	QuestionType string `json:"question_type"`
	EvidenceHit  bool   `json:"evidence_hit"`
	Error        string `json:"error,omitempty"`
}

type report struct {
	QuestionsPath string           `json:"questions_path"`
	TopK          int              `json:"top_k"`
	Conditions    []condition      `json:"conditions"`
	Results       map[string][]row `json:"results"`
}

var conditions = []condition{
	{Name: "fixed-pilot-100x100-v1", Representation: "fixed", Version: "pilot-100x100-v1"},
	{Name: "facts-pilot-named-v3", Representation: "facts", Version: "pilot-named-v3"},
}

func main() {
	var projectDir, questionsPath, output string
	var topK int
	root := &cobra.Command{Use: "decontextualized-question-development", SilenceErrors: true, SilenceUsage: true}
	root.RunE = func(cmd *cobra.Command, _ []string) error {
		projectDir, _ = filepath.Abs(projectDir)
		if err := loadEnv(projectDir); err != nil {
			return err
		}
		questions, err := loadQuestions(questionsPath)
		if err != nil {
			return err
		}
		pool, err := storage.OpenReadOnly(cmd.Context(), strings.TrimSpace(os.Getenv("DATABASE_URL")))
		if err != nil {
			return err
		}
		defer pool.Close()
		client, err := fireworks.NewFromEnv()
		if err != nil {
			return err
		}
		recorded := report{QuestionsPath: questionsPath, TopK: topK, Conditions: conditions, Results: map[string][]row{}}
		for _, c := range conditions {
			for _, q := range questions {
				found, err := retrieval.Search(cmd.Context(), pool, client, q.Text, retrieval.Config{Representation: c.Representation, Version: c.Version, CorpusSplit: "test", EmbeddingModel: embeddingModel, TopK: topK})
				r := row{QuestionID: q.ID, QuestionType: q.Type, EvidenceHit: containsEvidence(q.Gold, found)}
				if err != nil {
					r.Error = err.Error()
					r.EvidenceHit = false
				}
				recorded.Results[c.Name] = append(recorded.Results[c.Name], r)
			}
		}
		body, err := json.MarshalIndent(recorded, "", "  ")
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
			return err
		}
		return os.WriteFile(output, body, 0o600)
	}
	root.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	root.Flags().StringVar(&questionsPath, "questions", "experiments/decontextualized-question-development-20260725/questions.jsonl", "source-grounded development questions")
	root.Flags().StringVar(&output, "output", "artifacts/decontextualized-question-development-20260725/retrieval.json", "ignored output report")
	root.Flags().IntVar(&topK, "top-k", 5, "Number of retrieval results per question")
	returnOrExit(root.Execute())
}

func returnOrExit(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func loadEnv(projectDir string) error {
	path := filepath.Join(projectDir, ".env")
	if _, err := os.Stat(path); err == nil {
		return godotenv.Load(path)
	}
	if strings.TrimSpace(os.Getenv("DATABASE_URL")) == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	return nil
}

func loadQuestions(path string) ([]question, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var questions []question
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), 1<<20)
	for scanner.Scan() {
		var q question
		if err := json.Unmarshal(scanner.Bytes(), &q); err != nil {
			return nil, err
		}
		if q.ID == "" || q.Type == "" || q.Text == "" || len(q.Gold) == 0 {
			return nil, fmt.Errorf("invalid question record")
		}
		questions = append(questions, q)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return questions, nil
}

func containsEvidence(gold []evidence, results []retrieval.Result) bool {
	for _, want := range gold {
		for _, got := range results {
			if got.ExternalID == want.DocumentExternalID && got.TurnStart <= want.TurnStart && got.TurnEnd >= want.TurnEnd {
				return true
			}
		}
	}
	return false
}
