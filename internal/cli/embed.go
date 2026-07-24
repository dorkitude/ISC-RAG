package cli

import (
	"fmt"

	"github.com/dorkitude/isc-rag/internal/embedding"
	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/spf13/cobra"
)

func newEmbedCommand() *cobra.Command {
	var projectDir, representationName, version, split, model string
	var limit, batchSize int
	command := &cobra.Command{Use: "embed", Short: "Embed a retrieval representation in pgvector", RunE: func(cmd *cobra.Command, args []string) error {
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
		count, err := embedding.Embed(cmd.Context(), pool, client, embedding.Config{Representation: representationName, Version: version, CorpusSplit: split, Model: model, Limit: limit, BatchSize: batchSize})
		if err != nil {
			return err
		}
		cmd.Printf("Embedded %d %s/%s source unit(s).\n", count, representationName, version)
		return nil
	}}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	command.Flags().StringVar(&representationName, "representation", "fixed", "Representation: fixed, turn, semantic, facts")
	command.Flags().StringVar(&version, "version", "v1", "Representation version")
	command.Flags().StringVar(&split, "split", "test", "Corpus split")
	command.Flags().StringVar(&model, "model", "nomic-ai/nomic-embed-text-v1.5", "Fireworks embedding model")
	command.Flags().IntVar(&limit, "limit", 0, "Maximum missing embeddings (0 means all)")
	command.Flags().IntVar(&batchSize, "batch-size", 32, "Inputs per embeddings request")
	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		if representationName == "fact-input" {
			return fmt.Errorf("fact-input is an internal extraction representation, not an experimental condition")
		}
		return nil
	}
	return command
}
