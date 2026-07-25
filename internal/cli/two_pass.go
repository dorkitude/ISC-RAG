package cli

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/spf13/cobra"
)

const twoPassSystem = `Rewrite one isolated dialogue claim as a self-contained fact. Return ONLY {"text":"..."}. Resolve a pronoun only when its antecedent is explicit in the supplied nearby labeled turns. Use supplied speaker names, never "the speaker", "the host", or "the guest". Do not add facts beyond the claim and context. If an antecedent is not knowable, retain the claim's original wording rather than guessing.`

const twoPassResolveReferencesSystem = `Rewrite one isolated dialogue claim as a self-contained fact. Return ONLY {"text":"..."}. Resolve a pronoun or role description only when its unique antecedent is explicit in the supplied nearby labeled turns, exact quote, or claim. For example, if the local evidence explicitly identifies “the president” as President Bush, use “President Bush”; otherwise retain the original wording. Use supplied speaker names, never "the speaker", "the host", or "the guest". Preserve attribution, modality, qualifiers, and time language. Do not add facts beyond the claim and context.`

func newTwoPassCommand() *cobra.Command {
	var projectDir, model, outputPath, promptVariant string
	var documents, sample int
	command := &cobra.Command{Use: "two-pass-test", Short: "Test contextual rephrasing over isolated fact claims", RunE: func(cmd *cobra.Command, args []string) error {
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
		var systemPrompt string
		switch promptVariant {
		case "conservative":
			systemPrompt = twoPassSystem
		case "resolve-references":
			systemPrompt = twoPassResolveReferencesSystem
		default:
			return fmt.Errorf("unknown --prompt-variant %q; use conservative or resolve-references", promptVariant)
		}
		rows, err := pool.Query(cmd.Context(), `SELECT su.id,su.document_id,su.text,fe.quote,fe.turn_start,fe.turn_end FROM source_units su JOIN documents d ON d.id=su.document_id JOIN facts f ON f.document_id=su.document_id AND f.content_hash=su.content_hash JOIN fact_evidence fe ON fe.fact_id=f.id WHERE su.representation='facts' AND su.version='pilot-named-v3' AND su.document_id IN (SELECT id FROM documents WHERE corpus='mediasum' AND corpus_split='test' ORDER BY external_id LIMIT $1) ORDER BY d.external_id,su.id`, documents)
		if err != nil {
			return err
		}
		defer rows.Close()
		type item struct {
			documentID, claim, quote string
			start, end               int
		}
		var items []item
		for rows.Next() {
			var item item
			var ignored string
			if err := rows.Scan(&ignored, &item.documentID, &item.claim, &item.quote, &item.start, &item.end); err != nil {
				return err
			}
			items = append(items, item)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		type result struct {
			DocumentID string `json:"document_id"`
			Claim      string `json:"isolated_claim"`
			Quote      string `json:"exact_quote"`
			TurnStart  int    `json:"turn_start"`
			TurnEnd    int    `json:"turn_end"`
			Text       string `json:"rephrased_fact"`
		}
		type failure struct {
			DocumentID string `json:"document_id"`
			Claim      string `json:"isolated_claim"`
			Error      string `json:"error"`
		}
		var output []result
		var failures []failure
		persist := func() error {
			if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
				return err
			}
			serialized, err := json.MarshalIndent(struct {
				Model         string    `json:"model"`
				PromptVariant string    `json:"prompt_variant"`
				Documents     int       `json:"documents"`
				Results       []result  `json:"results"`
				Failures      []failure `json:"failures"`
			}{Model: model, PromptVariant: promptVariant, Documents: documents, Results: output, Failures: failures}, "", "  ")
			if err != nil {
				return err
			}
			return os.WriteFile(outputPath, append(serialized, '\n'), 0o600)
		}
		for _, item := range items {
			contextRows, err := pool.Query(cmd.Context(), `SELECT COALESCE(speaker_name,speaker_label,''),text FROM turns WHERE document_id=$1 AND ordinal BETWEEN GREATEST(0,$2-2) AND $3+2 ORDER BY ordinal`, item.documentID, item.start, item.end)
			if err != nil {
				return err
			}
			var context strings.Builder
			for contextRows.Next() {
				var name, text string
				if err := contextRows.Scan(&name, &text); err != nil {
					contextRows.Close()
					return err
				}
				fmt.Fprintf(&context, "[SPEAKER: %s]\n%s\n", name, text)
			}
			contextRows.Close()
			user := fmt.Sprintf("<ISOLATED_CLAIM>\n%s\n</ISOLATED_CLAIM>\n<EXACT_QUOTE>\n%s\n</EXACT_QUOTE>\n<NEARBY_TURNS>\n%s\n</NEARBY_TURNS>", item.claim, item.quote, context.String())
			response, err := client.ChatJSONNoReasoning(cmd.Context(), model, systemPrompt, user, 300)
			if err != nil {
				failures = append(failures, failure{DocumentID: item.documentID, Claim: item.claim, Error: err.Error()})
				if err := persist(); err != nil {
					return err
				}
				continue
			}
			var parsed struct {
				Text string `json:"text"`
			}
			if err := json.Unmarshal([]byte(response.Text), &parsed); err != nil {
				failures = append(failures, failure{DocumentID: item.documentID, Claim: item.claim, Error: fmt.Sprintf("decode JSON response: %v", err)})
				if err := persist(); err != nil {
					return err
				}
				continue
			}
			text := strings.TrimSpace(parsed.Text)
			if text != "" {
				output = append(output, result{
					DocumentID: item.documentID,
					Claim:      item.claim,
					Quote:      item.quote,
					TurnStart:  item.start,
					TurnEnd:    item.end,
					Text:       text,
				})
			} else {
				failures = append(failures, failure{DocumentID: item.documentID, Claim: item.claim, Error: "empty text field in JSON response"})
			}
			if err := persist(); err != nil {
				return err
			}
		}
		if err := persist(); err != nil {
			return err
		}
		rand.New(rand.NewSource(time.Now().UnixNano())).Shuffle(len(output), func(i, j int) { output[i], output[j] = output[j], output[i] })
		if sample > len(output) {
			sample = len(output)
		}
		cmd.Printf("Two-pass test: %d isolated claims rephrased across %d documents (%d failures); random sample:\n", len(output), documents, len(failures))
		for _, text := range output[:sample] {
			cmd.Printf("- %s\n", text.Text)
		}
		cmd.Printf("Saved all outputs to %s\n", outputPath)
		return nil
	}}
	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	command.Flags().StringVar(&model, "model", "accounts/fireworks/models/deepseek-v4-flash", "Fireworks rephrasing model")
	command.Flags().StringVar(&promptVariant, "prompt-variant", "conservative", "Prompt variant: conservative or resolve-references")
	command.Flags().IntVar(&documents, "documents", 5, "Number of test documents")
	command.Flags().IntVar(&sample, "sample", 25, "Random facts to render")
	command.Flags().StringVar(&outputPath, "output", filepath.Join("artifacts", "two-pass-rephrasing", "flash-results.json"), "Ignored local path for complete rephrasing outputs")
	return command
}
