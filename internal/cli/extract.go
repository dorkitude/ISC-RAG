package cli

import (
	"fmt"

	"github.com/dorkitude/isc-rag/internal/facts"
	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/spf13/cobra"
)

func newExtractCommand() *cobra.Command {
	var projectDir, version, inputVersion, split, model string
	var limit, offset, maxFacts int
	command := &cobra.Command{Use: "extract facts", Short: "Extract verified atomic facts from fact-input units", RunE: func(cmd *cobra.Command, args []string) error {
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
		count, err := facts.Extract(cmd.Context(), pool, client, facts.Config{Version: version, InputVersion: inputVersion, CorpusSplit: split, Model: model, Limit: limit, Offset: offset, MaxFacts: maxFacts})
		if err != nil {
			return err
		}
		cmd.Printf("Extracted and provenance-validated facts for %d document(s).\n", count)
		return nil
	}}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	command.Flags().StringVar(&version, "version", "v1", "Fact representation version")
	command.Flags().StringVar(&inputVersion, "input-version", "v1", "fact-input representation version")
	command.Flags().StringVar(&split, "split", "test", "Corpus split to process")
	command.Flags().StringVar(&model, "model", "accounts/fireworks/models/gpt-oss-120b", "Fireworks extraction model")
	command.Flags().IntVar(&limit, "limit", 0, "Maximum documents (0 means all)")
	command.Flags().IntVar(&offset, "offset", 0, "Documents to skip after stable external-ID ordering")
	command.Flags().IntVar(&maxFacts, "max-facts", 50, "Maximum facts per extraction input unit")
	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		if model == "" {
			return fmt.Errorf("model is required")
		}
		return nil
	}
	return command
}

func newEnrichCommand() *cobra.Command {
	var projectDir, version, model string
	var limit int
	command := &cobra.Command{Use: "enrich", Short: "Add optional topic and triple metadata to extracted facts", RunE: func(cmd *cobra.Command, args []string) error {
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
		count, err := facts.Enrich(cmd.Context(), pool, client, version, model, limit)
		if err != nil {
			return err
		}
		cmd.Printf("Enriched %d fact(s).\n", count)
		return nil
	}}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	command.Flags().StringVar(&version, "version", "v1", "Fact representation version")
	command.Flags().StringVar(&model, "model", "accounts/fireworks/models/gpt-oss-120b", "Fireworks enrichment model")
	command.Flags().IntVar(&limit, "limit", 0, "Maximum facts (0 means all)")
	return command
}
