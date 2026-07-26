package cli

import (
	"fmt"
	"strings"

	"github.com/dorkitude/isc-rag/internal/evaluation"
	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/spf13/cobra"
)

func newEvaluateCommand() *cobra.Command {
	command := &cobra.Command{Use: "evaluate", Short: "Import evidence labels and run reproducible RAG evaluations"}
	command.AddCommand(newQuestionImportCommand(), newQuestionExportCommand(), newPilotQuestionCommand(), newDeterministicPilotQuestionCommand(), newEvaluationRunCommand(), newDepthSummaryCommand(), newJudgeCommand())
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
	var projectDir, representationName, version, questionSplit, corpusSplit, embedModel, generatorModel, documentList, resumeRun string
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
		documents, err := readDocumentList(documentList)
		if err != nil {
			return err
		}
		metrics, err := evaluation.Run(cmd.Context(), pool, client, evaluation.Config{Representation: representationName, Version: version, QuestionSplit: questionSplit, CorpusSplit: corpusSplit, EmbeddingModel: embedModel, GeneratorModel: generatorModel, TopK: topK, EmbeddingBatchSize: embeddingBatchSize, ContextChars: contextChars, ContextTokens: contextTokens, DocumentExternalIDs: documents, ResumeRunID: resumeRun, Generate: generate})
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
	command.Flags().StringVar(&documentList, "document-list", "", "Optional newline-delimited source document IDs; evaluates only their questions")
	command.Flags().StringVar(&resumeRun, "resume-run", "", "Resume a failed matching run, preserving answers already stored")
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

func newJudgeCommand() *cobra.Command {
	var projectDir, runs, judgeModel, correctnessKey, supportKey string
	var rejudge bool
	command := &cobra.Command{Use: "judge", Short: "Run blind correctness and source-support screens over stored generated answers", RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(runs) == "" {
			return fmt.Errorf("--runs is required")
		}
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
		metrics, err := evaluation.Judge(cmd.Context(), pool, client, evaluation.JudgeConfig{RunIDs: strings.Split(runs, ","), JudgeModel: judgeModel, CorrectnessKey: correctnessKey, SupportKey: supportKey, SkipExisting: !rejudge})
		if err != nil {
			return err
		}
		cmd.Printf("Judged %d/%d answers (%d already judged, skipped): %d correct, %d supported.\n", metrics.Judged, metrics.Answers, metrics.Skipped, metrics.Correct, metrics.Supported)
		return nil
	}}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	command.Flags().StringVar(&runs, "runs", "", "Comma-separated experiment run UUIDs")
	command.Flags().StringVar(&judgeModel, "judge-model", "accounts/fireworks/models/kimi-k2p6", "Fireworks judge model; must differ from the answer model for an independent screen")
	command.Flags().StringVar(&correctnessKey, "correctness-key", "", "Usage key to store blind correctness verdicts (required)")
	command.Flags().StringVar(&supportKey, "support-key", "", "Usage key to store blind source-support verdicts (required)")
	command.Flags().BoolVar(&rejudge, "rejudge", false, "Re-judge answers that already hold verdicts under both keys")
	return command
}

func newDepthSummaryCommand() *cobra.Command {
	var projectDir, runs, correctnessKey, supportKey string
	command := &cobra.Command{Use: "depth-summary", Short: "Summarize answer success and minimum successful top-k across completed runs", RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(runs) == "" {
			return fmt.Errorf("--runs is required")
		}
		pool, err := openDatabase(projectDir)
		if err != nil {
			return err
		}
		defer pool.Close()
		if err = storage.Migrate(cmd.Context(), pool); err != nil {
			return err
		}
		ids := strings.Split(runs, ",")
		rows, err := pool.Query(cmd.Context(), `WITH runs AS (SELECT id,representation,(parameters->>'top_k')::int AS k FROM experiment_runs WHERE id=ANY($1::uuid[])), outcomes AS (SELECT r.representation,r.k,q.external_id,(a.usage->$2->>'verdict'='correct' AND a.usage->$3->>'verdict'='supported') AS success,(a.usage->>'total_tokens')::numeric AS total_tokens FROM answers a JOIN runs r ON r.id=a.experiment_run_id JOIN questions q ON q.id=a.question_id) SELECT representation,k,count(*) FILTER (WHERE success),count(*),round(avg(total_tokens),1) FROM outcomes GROUP BY 1,2 ORDER BY 1,2`, ids, correctnessKey, supportKey)
		if err != nil {
			return err
		}
		defer rows.Close()
		cmd.Println("Answerable-and-supported by top-k (mean provider total tokens):")
		for rows.Next() {
			var representation string
			var k, success, total int
			var tokens float64
			if err = rows.Scan(&representation, &k, &success, &total, &tokens); err != nil {
				return err
			}
			cmd.Printf("  %s k=%d: %d/%d (mean total tokens %.1f)\n", representation, k, success, total, tokens)
		}
		if err = rows.Err(); err != nil {
			return err
		}
		rows, err = pool.Query(cmd.Context(), `WITH runs AS (SELECT id,representation,(parameters->>'top_k')::int AS k FROM experiment_runs WHERE id=ANY($1::uuid[])), outcomes AS (SELECT r.representation,r.k,q.external_id,(a.usage->$2->>'verdict'='correct' AND a.usage->$3->>'verdict'='supported') AS success FROM answers a JOIN runs r ON r.id=a.experiment_run_id JOIN questions q ON q.id=a.question_id), first_success AS (SELECT representation,external_id,min(k) FILTER (WHERE success) AS minimum_k FROM outcomes GROUP BY 1,2) SELECT representation,external_id,minimum_k FROM first_success ORDER BY representation,external_id`, ids, correctnessKey, supportKey)
		if err != nil {
			return err
		}
		defer rows.Close()
		cmd.Println("Minimum successful top-k per question (never means no correct, supported answer in supplied runs):")
		for rows.Next() {
			var representation, externalID string
			var minimum *int
			if err = rows.Scan(&representation, &externalID, &minimum); err != nil {
				return err
			}
			value := "never"
			if minimum != nil {
				value = fmt.Sprintf("%d", *minimum)
			}
			cmd.Printf("  %s %s: %s\n", representation, externalID, value)
		}
		return rows.Err()
	}}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	command.Flags().StringVar(&runs, "runs", "", "Comma-separated experiment run UUIDs")
	command.Flags().StringVar(&correctnessKey, "correctness-key", "judge_kimi_depth_v1", "Usage key containing blind correctness verdicts")
	command.Flags().StringVar(&supportKey, "support-key", "source_support_kimi_depth_v1", "Usage key containing blind source-support verdicts")
	return command
}
