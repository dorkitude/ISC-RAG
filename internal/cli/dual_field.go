package cli

import (
	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/dorkitude/isc-rag/internal/study2"
	"github.com/spf13/cobra"
)

func newDualFieldCommand() *cobra.Command {
	var dir, version, model string
	var batch int
	c := &cobra.Command{Use: "study2-dual-expand", Short: "Generate non-authoritative retrieval expansions for exact capsules", RunE: func(cmd *cobra.Command, _ []string) error {
		p, e := openDatabase(dir)
		if e != nil {
			return e
		}
		defer p.Close()
		if e = storage.Migrate(cmd.Context(), p); e != nil {
			return e
		}
		cl, e := fireworks.NewFromEnv()
		if e != nil {
			return e
		}
		n, e := study2.ExpandCapsules(cmd.Context(), p, cl, study2.DualConfig{Version: version, Model: model, BatchSize: batch})
		if e == nil {
			cmd.Printf("Indexed %d non-authoritative capsule expansions; answer generation must use linked exact capsules.\n", n)
		}
		return e
	}}
	c.Flags().StringVar(&dir, "project-dir", ".", "project directory")
	c.Flags().StringVar(&version, "version", "v2-dual-field", "source unit version")
	c.Flags().StringVar(&model, "model", "accounts/fireworks/models/deepseek-v4-pro", "model")
	c.Flags().IntVar(&batch, "batch-size", 20, "capsules per expansion call")
	return c
}
