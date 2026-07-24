package cli

import (
	"path/filepath"

	"github.com/dorkitude/isc-rag/internal/mediasum"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/spf13/cobra"
)

func newIngestCommand() *cobra.Command {
	command := &cobra.Command{Use: "ingest", Short: "Ingest canonical corpus records"}
	command.AddCommand(newMediaSumIngestCommand())
	return command
}

func newMediaSumIngestCommand() *cobra.Command {
	var projectDir, split, file string
	var limit int
	command := &cobra.Command{
		Use: "mediasum", Short: "Ingest a MediaSum JSONL split into PostgreSQL",
		RunE: func(cmd *cobra.Command, args []string) error {
			pool, err := openDatabase(projectDir)
			if err != nil {
				return err
			}
			defer pool.Close()
			if err := storage.Migrate(cmd.Context(), pool); err != nil {
				return err
			}
			if file == "" {
				file = filepath.Join(projectDir, "resources", "mediasum", split+".json")
			}
			count, err := mediasum.Ingest(cmd.Context(), pool, file, split, limit)
			if err != nil {
				return err
			}
			cmd.Printf("Ingested %d MediaSum %s record(s).\n", count, split)
			return nil
		},
	}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	command.Flags().StringVar(&split, "split", "test", "MediaSum split label")
	command.Flags().StringVar(&file, "file", "", "JSONL file path (defaults to resources/mediasum/<split>.json)")
	command.Flags().IntVar(&limit, "limit", 0, "Maximum records to ingest (0 means all)")
	return command
}
