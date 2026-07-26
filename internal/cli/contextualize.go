package cli

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/storage"
)

const contextualizePrompt = `You situate document excerpts for search indexing. Given a full interview transcript and one chunk from it, write 50-100 tokens of succinct context that situates the chunk within the overall document: name the programme/topic, the speakers involved and their roles, and what is being discussed when the chunk occurs, so the chunk is understandable and searchable on its own. Name people that the chunk refers to only by pronouns. Return ONLY JSON: {"context":"..."}.`

// newContextualizeCommand builds a contextual-retrieval baseline: each chunk of
// a source representation gains an LLM-generated situating preamble and is
// stored as a new version, following the contextual-embeddings recipe.
func newContextualizeCommand() *cobra.Command {
	var projectDir, sourceRepresentation, sourceVersion, targetVersion, corpusSplit, model string
	var maxDocumentChars int
	command := &cobra.Command{Use: "contextualize", Short: "Create a contextual-retrieval chunk baseline by prepending LLM-generated situating context", RunE: func(cmd *cobra.Command, args []string) error {
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
		ctx := cmd.Context()
		rows, err := pool.Query(ctx, `SELECT su.id, su.document_id, su.turn_start, su.turn_end, su.char_start, su.char_end, su.text
			FROM source_units su JOIN documents d ON d.id=su.document_id
			WHERE su.representation=$1 AND su.version=$2 AND d.corpus_split=$3
			  AND NOT EXISTS (SELECT 1 FROM source_units t WHERE t.document_id=su.document_id AND t.representation=$1 AND t.version=$4 AND t.turn_start=su.turn_start AND t.turn_end=su.turn_end AND t.char_start=su.char_start)
			ORDER BY su.document_id, su.turn_start`, sourceRepresentation, sourceVersion, corpusSplit, targetVersion)
		if err != nil {
			return err
		}
		type unit struct {
			id, documentID, text                   string
			turnStart, turnEnd, charStart, charEnd int
		}
		var units []unit
		for rows.Next() {
			var u unit
			if err = rows.Scan(&u.id, &u.documentID, &u.turnStart, &u.turnEnd, &u.charStart, &u.charEnd, &u.text); err != nil {
				rows.Close()
				return err
			}
			units = append(units, u)
		}
		rows.Close()
		if err = rows.Err(); err != nil {
			return err
		}
		transcripts := map[string]string{}
		written := 0
		for _, u := range units {
			transcript, ok := transcripts[u.documentID]
			if !ok {
				turnRows, err := pool.Query(ctx, `SELECT COALESCE(speaker_label,''), text FROM turns WHERE document_id=$1 ORDER BY ordinal`, u.documentID)
				if err != nil {
					return err
				}
				var builder strings.Builder
				for turnRows.Next() {
					var speaker, text string
					if err = turnRows.Scan(&speaker, &text); err != nil {
						turnRows.Close()
						return err
					}
					fmt.Fprintf(&builder, "%s: %s\n", speaker, text)
				}
				turnRows.Close()
				if err = turnRows.Err(); err != nil {
					return err
				}
				transcript = builder.String()
				if len(transcript) > maxDocumentChars {
					transcript = transcript[:maxDocumentChars]
				}
				transcripts[u.documentID] = transcript
			}
			result, err := client.ChatJSONNoReasoning(ctx, model, contextualizePrompt, fmt.Sprintf("<DOCUMENT>\n%s\n</DOCUMENT>\n\n<CHUNK>\n%s\n</CHUNK>", transcript, u.text), 400)
			if err != nil {
				return fmt.Errorf("contextualize unit %s: %w", u.id, err)
			}
			var parsed struct {
				Context string `json:"context"`
			}
			if err = json.Unmarshal([]byte(result.Text), &parsed); err != nil || strings.TrimSpace(parsed.Context) == "" {
				return fmt.Errorf("unit %s: invalid context response %q", u.id, result.Text)
			}
			combined := strings.TrimSpace(parsed.Context) + "\n\n" + u.text
			hash := fmt.Sprintf("%x", sha256.Sum256([]byte(combined)))
			if _, err = pool.Exec(ctx, `INSERT INTO source_units(document_id,representation,version,turn_start,turn_end,char_start,char_end,text,content_hash) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT(document_id,representation,version,content_hash) DO NOTHING`, u.documentID, sourceRepresentation, targetVersion, u.turnStart, u.turnEnd, u.charStart, u.charEnd, combined, hash); err != nil {
				return err
			}
			written++
			if written%25 == 0 {
				cmd.Printf("contextualized %d/%d units\n", written, len(units))
			}
		}
		cmd.Printf("Contextualized %d unit(s) into %s/%s.\n", written, sourceRepresentation, targetVersion)
		return nil
	}}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	command.Flags().StringVar(&sourceRepresentation, "source-representation", "fixed", "Representation to contextualize")
	command.Flags().StringVar(&sourceVersion, "source-version", "", "Existing chunk version to read (required)")
	command.Flags().StringVar(&targetVersion, "version", "", "New version to write (required)")
	command.Flags().StringVar(&corpusSplit, "corpus-split", "test", "Source corpus split")
	command.Flags().StringVar(&model, "model", "accounts/fireworks/models/deepseek-v4-flash", "Fireworks model that writes the situating context")
	command.Flags().IntVar(&maxDocumentChars, "max-document-chars", 120000, "Transcript characters supplied to the context writer")
	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		if sourceVersion == "" || targetVersion == "" || sourceVersion == targetVersion {
			return fmt.Errorf("--source-version and --version are required and must differ")
		}
		return nil
	}
	return command
}
