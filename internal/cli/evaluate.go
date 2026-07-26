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
	command.AddCommand(newQuestionImportCommand(), newQuestionExportCommand(), newPilotQuestionCommand(), newDeterministicPilotQuestionCommand(), newEvaluationRunCommand())
	return command
}
func newDeterministicPilotQuestionCommand() *cobra.Command {
	var projectDir, questionSplit, corpusSplit string
	var documents, perDocument int
	command := &cobra.Command{Use: "deterministic-pilot-questions", Short: "Create source-phrase retrieval stress-test probes", RunE: func(cmd *cobra.Command, args []string) error {
		pool, err := openDatabase(projectDir)
		if err != nil {
			return err
		}
		defer pool.Close()
		if err = storage.Migrate(cmd.Context(), pool); err != nil {
			return err
		}
		report, err := evaluation.GenerateDeterministicPilotQuestions(cmd.Context(), pool, questionSplit, corpusSplit, documents, perDocument)
		if err != nil {
			return err
		}
		cmd.Printf("Generated %d deterministic lexical retrieval probes across %d documents.\n", report.Questions, report.Documents)
		return nil
	}}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	command.Flags().StringVar(&questionSplit, "question-split", "pilot-100x100-v1", "Question-set namespace")
	command.Flags().StringVar(&corpusSplit, "corpus-split", "test", "Source corpus split")
	command.Flags().IntVar(&documents, "documents", 100, "Number of source documents")
	command.Flags().IntVar(&perDocument, "per-document", 100, "Probes per source document")
	return command
}
func newPilotQuestionCommand() *cobra.Command {
	var projectDir, questionSplit, corpusSplit, model string
	var documents, perDocument, batchSize, workers int
	command := &cobra.Command{Use: "pilot-questions", Short: "Generate exact-evidence retrieval probes for a bounded pilot", RunE: func(cmd *cobra.Command, args []string) error {
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
		report, err := evaluation.GeneratePilotQuestions(cmd.Context(), pool, client, evaluation.PilotQuestionConfig{QuestionSplit: questionSplit, CorpusSplit: corpusSplit, Model: model, Documents: documents, PerDocument: perDocument, BatchSize: batchSize, Workers: workers})
		if err != nil {
			return err
		}
		cmd.Printf("Generated %d exact-evidence diagnostic question(s) across %d document(s); rejected %d invalid candidate(s).\n", report.Questions, report.Documents, report.Rejected)
		return nil
	}}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	command.Flags().StringVar(&questionSplit, "question-split", "pilot-100x100-v1", "Question-set namespace")
	command.Flags().StringVar(&corpusSplit, "corpus-split", "test", "Source corpus split")
	command.Flags().StringVar(&model, "model", "accounts/fireworks/models/deepseek-v4-pro", "Fireworks question-generation model")
	command.Flags().IntVar(&documents, "documents", 100, "Number of source documents")
	command.Flags().IntVar(&perDocument, "per-document", 100, "Exact-evidence questions per source document")
	command.Flags().IntVar(&batchSize, "batch-size", 10, "Questions per structured provider response")
	command.Flags().IntVar(&workers, "workers", 1, "Concurrent document-generation workers")
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
func newQuestionExportCommand() *cobra.Command {
	var projectDir string
	command := &cobra.Command{Use: "export-questions SPLIT FILE.jsonl", Short: "Export a validated question split with canonical evidence", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
		pool, err := openDatabase(projectDir)
		if err != nil {
			return err
		}
		defer pool.Close()
		count, err := evaluation.Export(cmd.Context(), pool, args[0], args[1])
		if err != nil {
			return err
		}
		cmd.Printf("Exported %d question(s).\n", count)
		return nil
	}}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	return command
}
func newEvaluationRunCommand() *cobra.Command {
	var projectDir, representationName, version, questionSplit, corpusSplit, embedModel, generatorModel string
	var topK, embeddingBatchSize, contextChars, contextTokens int
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
		metrics, err := evaluation.Run(cmd.Context(), pool, client, evaluation.Config{Representation: representationName, Version: version, QuestionSplit: questionSplit, CorpusSplit: corpusSplit, EmbeddingModel: embedModel, GeneratorModel: generatorModel, TopK: topK, EmbeddingBatchSize: embeddingBatchSize, ContextChars: contextChars, ContextTokens: contextTokens, Generate: generate})
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
	command.Flags().IntVar(&embeddingBatchSize, "embedding-batch-size", 32, "Questions per Fireworks embeddings request")
	command.Flags().IntVar(&contextChars, "context-chars", 0, "Maximum retrieved-context characters passed to the answer model; 0 does not pack")
	command.Flags().IntVar(&contextTokens, "context-tokens", 0, "Maximum cl100k_base-estimated retrieved-context tokens passed to the answer model; 0 does not pack")
	command.Flags().BoolVar(&generate, "generate", false, "Generate source-cited answers and store them")
	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		if generate && generatorModel == "" {
			return fmt.Errorf("--generator-model is required with --generate")
		}
		return nil
	}
	return command
}
