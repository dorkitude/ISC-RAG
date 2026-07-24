package cli

import (
	"github.com/dorkitude/isc-rag/internal/retrieval"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/spf13/cobra"
)

func newIndexCommand() *cobra.Command {
	var projectDir, model string
	command := &cobra.Command{Use: "index", Short: "Create a model-specific pgvector HNSW index", RunE: func(cmd *cobra.Command, args []string) error {
		pool, err := openDatabase(projectDir)
		if err != nil {
			return err
		}
		defer pool.Close()
		if err = storage.Migrate(cmd.Context(), pool); err != nil {
			return err
		}
		name, err := retrieval.EnsureHNSW(cmd.Context(), pool, model)
		if err != nil {
			return err
		}
		cmd.Printf("HNSW index ready: %s\n", name)
		return nil
	}}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	command.Flags().StringVar(&model, "embedding-model", "nomic-ai/nomic-embed-text-v1.5", "Embedding model to index")
	return command
}
