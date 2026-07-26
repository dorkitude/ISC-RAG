package cli

import (
	"fmt"

	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/spf13/cobra"
)

// newFactsCommand materializes immutable fact projections for retrieval
// ablations without changing the canonical fact or its evidence record.
func newFactsCommand() *cobra.Command {
	root := &cobra.Command{Use: "facts", Short: "Create versioned fact retrieval projections"}
	var projectDir, sourceVersion, targetVersion, split string
	quoteProjection := &cobra.Command{Use: "quote-projection", Short: "Copy facts into a version for fact-plus-quote embeddings", RunE: func(cmd *cobra.Command, _ []string) error {
		if sourceVersion == targetVersion {
			return fmt.Errorf("source and target versions must differ")
		}
		pool, err := openDatabase(projectDir)
		if err != nil {
			return err
		}
		defer pool.Close()
		if err = storage.Migrate(cmd.Context(), pool); err != nil {
			return err
		}
		result, err := pool.Exec(cmd.Context(), `INSERT INTO source_units(document_id,representation,version,turn_start,turn_end,char_start,char_end,text,content_hash,pipeline_run_id) SELECT su.document_id,'facts',$1,su.turn_start,su.turn_end,su.char_start,su.char_end,su.text,su.content_hash,su.pipeline_run_id FROM source_units su JOIN documents d ON d.id=su.document_id WHERE su.representation='facts' AND su.version=$2 AND d.corpus_split=$3 ON CONFLICT(document_id,representation,version,content_hash) DO NOTHING`, targetVersion, sourceVersion, split)
		if err != nil {
			return err
		}
		cmd.Printf("Materialized %d facts/%s source unit(s) from facts/%s.\n", result.RowsAffected(), targetVersion, sourceVersion)
		return nil
	}}
	quoteProjection.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	quoteProjection.Flags().StringVar(&sourceVersion, "source-version", "", "Existing canonical fact source-unit version")
	quoteProjection.Flags().StringVar(&targetVersion, "target-version", "", "New immutable fact source-unit version")
	quoteProjection.Flags().StringVar(&split, "split", "test", "Corpus split")
	quoteProjection.PreRunE = func(cmd *cobra.Command, _ []string) error {
		if sourceVersion == "" || targetVersion == "" {
			return fmt.Errorf("--source-version and --target-version are required")
		}
		return nil
	}
	root.AddCommand(quoteProjection)
	return root
}
