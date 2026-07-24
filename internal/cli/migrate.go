package cli

import (
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/spf13/cobra"
)

func newMigrateCommand() *cobra.Command {
	var projectDir string
	command := &cobra.Command{
		Use: "migrate", Short: "Apply PostgreSQL schema migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			pool, err := openDatabase(projectDir)
			if err != nil {
				return err
			}
			defer pool.Close()
			if err := storage.Migrate(cmd.Context(), pool); err != nil {
				return err
			}
			cmd.Println("PostgreSQL migrations are current.")
			return nil
		},
	}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	return command
}
