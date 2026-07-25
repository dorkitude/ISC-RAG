package cli

import (
	"github.com/dorkitude/isc-rag/internal/mediasum"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/spf13/cobra"
)

func newAnnotateCommand() *cobra.Command {
	command := &cobra.Command{Use: "annotate", Short: "Add deterministic metadata to canonical source records"}
	command.AddCommand(newSpeakerMetadataCommand())
	return command
}

func newSpeakerMetadataCommand() *cobra.Command {
	var projectDir string
	command := &cobra.Command{
		Use:   "speaker-metadata",
		Short: "Parse corpus speaker labels into deterministic name and role fields",
		RunE: func(cmd *cobra.Command, args []string) error {
			pool, err := openDatabase(projectDir)
			if err != nil {
				return err
			}
			defer pool.Close()
			if err = storage.Migrate(cmd.Context(), pool); err != nil {
				return err
			}
			count, err := mediasum.BackfillSpeakerMetadata(cmd.Context(), pool)
			if err != nil {
				return err
			}
			cmd.Printf("Parsed deterministic speaker metadata for %d turn(s).\n", count)
			return nil
		},
	}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	return command
}
