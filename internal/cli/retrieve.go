package cli

import (
	"encoding/json"
	"fmt"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/retrieval"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/spf13/cobra"
)

func newRetrieveCommand() *cobra.Command {
	var projectDir, representationName, version, split, model, query string
	var topK int
	command := &cobra.Command{Use: "retrieve", Short: "Run pgvector similarity search for one query", RunE: func(cmd *cobra.Command, args []string) error {
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
		results, err := retrieval.Search(cmd.Context(), pool, client, query, retrieval.Config{Representation: representationName, Version: version, CorpusSplit: split, EmbeddingModel: model, TopK: topK})
		if err != nil {
			return err
		}
		encoded, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return err
		}
		cmd.Println(string(encoded))
		return nil
	}}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	command.Flags().StringVar(&representationName, "representation", "fixed", "Representation")
	command.Flags().StringVar(&version, "version", "v1", "Representation version")
	command.Flags().StringVar(&split, "split", "test", "Corpus split")
	command.Flags().StringVar(&model, "embedding-model", "nomic-ai/nomic-embed-text-v1.5", "Fireworks embedding model")
	command.Flags().StringVar(&query, "query", "", "Question or search query")
	command.Flags().IntVar(&topK, "top-k", 5, "Number of source units")
	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		if query == "" {
			return fmt.Errorf("--query is required")
		}
		return nil
	}
	return command
}
