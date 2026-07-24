package cli

import (
	"fmt"

	"github.com/dorkitude/isc-rag/internal/evaluation"
	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/spf13/cobra"
)

func newEvaluateCommand() *cobra.Command {
	command := &cobra.Command{Use: "evaluate", Short: "Import evidence labels and run reproducible RAG evaluations"}
	command.AddCommand(newQuestionImportCommand(), newEvaluationRunCommand())
	return command
}
func newQuestionImportCommand() *cobra.Command {
	var projectDir, file string
	command := &cobra.Command{Use: "questions FILE.jsonl", Short: "Import independently authored questions and gold turn spans", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		pool, err := openDatabase(projectDir)
		if err != nil {
			return err
		}
		defer pool.Close()
		if err = storage.Migrate(cmd.Context(), pool); err != nil {
			return err
		}
		count, err := evaluation.Import(cmd.Context(), pool, args[0])
		if err != nil {
			return err
		}
		cmd.Printf("Imported %d question(s).\n", count)
		return nil
	}}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	_ = file
	return command
}
func newEvaluationRunCommand() *cobra.Command {
	var projectDir, representationName, version, questionSplit, corpusSplit, embedModel, generatorModel string
	var topK int
	var generate bool
	command := &cobra.Command{Use: "run", Short: "Retrieve every question and record evidence coverage (optionally generate answers)", RunE: func(cmd *cobra.Command, args []string) error {
		pool, err := openDatabase(projectDir)
		if err != nil {
			return err
		}
		defer pool.Close()
		if err = storage.Migrate(cmd.Context(), pool); err != nil {
			return err
		}
		client, err := fireworks.NewFromEnv()
		if err != nil {
			return err
		}
		metrics, err := evaluation.Run(cmd.Context(), pool, client, evaluation.Config{Representation: representationName, Version: version, QuestionSplit: questionSplit, CorpusSplit: corpusSplit, EmbeddingModel: embedModel, GeneratorModel: generatorModel, TopK: topK, Generate: generate})
		if err != nil {
			return err
		}
		recall := 0.0
		if metrics.EvidenceTotal > 0 {
			recall = float64(metrics.EvidenceHits) / float64(metrics.EvidenceTotal)
		}
		cmd.Printf("Run %s: %d questions, evidence recall@%d %.3f (%d/%d), questions with a hit %d/%d.\n", metrics.RunID, metrics.Questions, topK, recall, metrics.EvidenceHits, metrics.EvidenceTotal, metrics.QuestionsWithHit, metrics.Questions)
		return nil
	}}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	command.Flags().StringVar(&representationName, "representation", "fixed", "Representation")
	command.Flags().StringVar(&version, "version", "v1", "Representation version")
	command.Flags().StringVar(&questionSplit, "question-split", "dev", "Question split")
	command.Flags().StringVar(&corpusSplit, "corpus-split", "test", "Source corpus split")
	command.Flags().StringVar(&embedModel, "embedding-model", "nomic-ai/nomic-embed-text-v1.5", "Fireworks embedding model")
	command.Flags().StringVar(&generatorModel, "generator-model", "accounts/fireworks/models/gpt-oss-120b", "Fireworks answer model")
	command.Flags().IntVar(&topK, "top-k", 5, "Number of retrieval units")
	command.Flags().BoolVar(&generate, "generate", false, "Generate source-cited answers and store them")
	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		if generate && generatorModel == "" {
			return fmt.Errorf("--generator-model is required with --generate")
		}
		return nil
	}
	return command
}
