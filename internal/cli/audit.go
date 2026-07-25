package cli

import (
	"fmt"
	"path/filepath"

	"github.com/dorkitude/isc-rag/internal/mediasum"
	"github.com/spf13/cobra"
)

func newAuditCommand() *cobra.Command {
	command := &cobra.Command{Use: "audit", Short: "Run read-only corpus integrity checks"}
	command.AddCommand(newMediaSumAuditCommand())
	return command
}

func newMediaSumAuditCommand() *cobra.Command {
	var projectDir, split string
	command := &cobra.Command{
		Use:   "mediasum",
		Short: "Validate deterministic MediaSum source fields without writing records",
		RunE: func(cmd *cobra.Command, args []string) error {
			splits := []string{split}
			if split == "all" {
				splits = []string{"train", "validation", "test"}
			}
			for _, current := range splits {
				if current != "train" && current != "validation" && current != "test" {
					return fmt.Errorf("unknown MediaSum split %q (use train, validation, test, or all)", current)
				}
				report, err := mediasum.AuditFile(filepath.Join(projectDir, "resources", "mediasum", current+".json"))
				if err != nil {
					return err
				}
				report.Write(cmd.OutOrStdout())
				if !report.Valid() {
					return fmt.Errorf("MediaSum %s audit failed", current)
				}
			}
			return nil
		},
	}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	command.Flags().StringVar(&split, "split", "all", "MediaSum split to audit: train, validation, test, or all")
	return command
}
