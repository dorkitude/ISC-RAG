package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/dorkitude/isc-rag/internal/study"
	"github.com/spf13/cobra"
)

// newDeepSeekStudyCommand is intentionally named after the requested provider
// family. The default is the account-visible model; --model remains explicit so
// a future DeepSeek V4 Flash deployment can be selected without code changes.
func newDeepSeekStudyCommand() *cobra.Command {
	var projectDir, namespace, version, externalID, model, promptFile, outputDir, approach string
	var iteration, maxFacts int
	cmd := &cobra.Command{Use: "deepseek-study", Short: "Run one isolated DeepSeek fact-rephrasing study iteration", RunE: func(cmd *cobra.Command, args []string) error {
		prompt, err := os.ReadFile(promptFile)
		if err != nil {
			return fmt.Errorf("read prompt file: %w", err)
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
		result, err := study.Run(cmd.Context(), pool, client, study.Config{Namespace: namespace, Version: version, ExternalID: externalID, Model: model, Prompt: string(prompt), OutputDir: outputDir, Approach: approach, Iteration: iteration, MaxFacts: maxFacts})
		if err != nil {
			return err
		}
		cmd.Printf("Study %s/%s iteration %d: %d raw candidates, %d exact-provenance facts; run=%s prompt_sha256=%s raw_sha256=%s\n", namespace, version, iteration, result.RawCount, result.ValidatedCount, result.RunID, result.PromptHash, result.RawOutputHash)
		return nil
	}}
	cmd.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	cmd.Flags().StringVar(&namespace, "namespace", "deepseek-fact-rephrasing", "Isolated study namespace")
	cmd.Flags().StringVar(&version, "version", "v1", "Study version")
	cmd.Flags().StringVar(&externalID, "external-id", "NPR-847", "One MediaSum document external ID")
	cmd.Flags().StringVar(&model, "model", "accounts/fireworks/models/deepseek-v4-flash", "Fireworks model for the isolated study")
	cmd.Flags().StringVar(&promptFile, "prompt-file", filepath.Join("experiments", "deepseek-fact-rephrasing-v1", "prompts", "iteration-01.txt"), "Tracked prompt file")
	cmd.Flags().StringVar(&outputDir, "output-dir", filepath.Join("artifacts", "deepseek-fact-rephrasing-v1"), "Ignored local raw-output directory")
	cmd.Flags().StringVar(&approach, "approach", "direct", "direct or quote-first")
	cmd.Flags().IntVar(&iteration, "iteration", 1, "Positive prompt iteration number")
	cmd.Flags().IntVar(&maxFacts, "max-facts", 60, "Maximum candidate facts")
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(model) == "" || iteration < 1 || maxFacts < 1 || (approach != "direct" && approach != "quote-first") {
			return fmt.Errorf("model and positive iteration/max-facts are required")
		}
		return nil
	}
	return cmd
}
