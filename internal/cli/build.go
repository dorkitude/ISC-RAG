package cli

import (
	"fmt"
	"strings"

	"github.com/dorkitude/isc-rag/internal/representation"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/spf13/cobra"
)

func newBuildCommand() *cobra.Command {
	var projectDir, representationName, version, split string
	var limit, maxChars, overlap, semanticMin, factInputMax int
	command := &cobra.Command{
		Use: "build", Short: "Build reproducible retrieval representations from canonical turns",
		RunE: func(cmd *cobra.Command, args []string) error {
			pool, err := openDatabase(projectDir)
			if err != nil {
				return err
			}
			defer pool.Close()
			if err = storage.Migrate(cmd.Context(), pool); err != nil {
				return err
			}
			representations := []string{representationName}
			if representationName == "all" {
				representations = []string{representation.Fixed, representation.Turn, representation.Semantic, representation.FactInput}
			}
			for _, name := range representations {
				count, err := representation.Build(cmd.Context(), pool, representation.Config{Representation: name, Version: version, CorpusSplit: split, Limit: limit, MaxChars: maxChars, OverlapChars: overlap, SemanticMin: semanticMin, FactInputMax: factInputMax})
				if err != nil {
					return fmt.Errorf("build %s: %w", name, err)
				}
				cmd.Printf("Built %s/%s for %d %s document(s).\n", name, version, count, split)
			}
			return nil
		},
	}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	command.Flags().StringVar(&representationName, "representation", "all", "fixed, turn, semantic, fact-input, or all")
	command.Flags().StringVar(&version, "version", "v1", "Representation version")
	command.Flags().StringVar(&split, "split", "test", "Corpus split to process")
	command.Flags().IntVar(&limit, "limit", 0, "Maximum documents (0 means all)")
	command.Flags().IntVar(&maxChars, "max-chars", 4000, "Maximum source-unit size in UTF-8 bytes")
	command.Flags().IntVar(&overlap, "overlap-chars", 400, "Fixed-window overlap in UTF-8 bytes")
	command.Flags().IntVar(&semanticMin, "semantic-min-chars", 2000, "Minimum semantic chunk size in UTF-8 bytes")
	command.Flags().IntVar(&factInputMax, "fact-input-max-chars", 10000, "Maximum extraction-input size in UTF-8 bytes")
	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		valid := map[string]bool{representation.Fixed: true, representation.Turn: true, representation.Semantic: true, representation.FactInput: true, "all": true}
		if !valid[strings.ToLower(representationName)] {
			return fmt.Errorf("unknown representation %q", representationName)
		}
		representationName = strings.ToLower(representationName)
		return nil
	}
	return command
}
