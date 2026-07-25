package cli

import (
	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/dorkitude/isc-rag/internal/study2"
	"github.com/spf13/cobra"
)

func newStudyTwoCommand() *cobra.Command {
	var dir, version, candidate, external, model string
	var replicate, max, chars, carry int
	c := &cobra.Command{Use: "study-two", Short: "Run one isolated production-inspired fact-study candidate", RunE: func(cmd *cobra.Command, args []string) error {
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
		r, e := study2.Run(cmd.Context(), p, cl, study2.Config{StudyVersion: version, Candidate: candidate, ExternalID: external, Model: model, Replicate: replicate, MaxFacts: max, ChunkChars: chars, Carryover: carry})
		if e != nil {
			return e
		}
		cmd.Printf("Study #2 %s replicate %d %s: %d raw, %d exact-bundle facts; run=%s prompt_sha256=%s tokens=%d/%d latency_ms=%d\n", candidate, replicate, external, r.Raw, r.Validated, r.RunID, r.PromptHash, r.PromptTokens, r.CompletionTokens, r.LatencyMS)
		return nil
	}}
	c.Flags().StringVar(&dir, "project-dir", ".", "ISC-RAG project directory")
	c.Flags().StringVar(&version, "version", "v2", "isolated study version")
	c.Flags().StringVar(&candidate, "candidate", "direct", "direct, speaker-packed, carryover, bundle, expanded-bundle, expanded-packed, attributed-packed, atomic-attributed, turnwise-staged, turnwise-local-staged, clause-anchored, proof-obligation, highrecall-verified, highrecall-nli, prefilter-nli, decomposed-nli, framed-render, quote-repair, predicate-repair, substantive-render, fable-c4-sbar, dialogue-act-pairs, dialogue-act-predicate-units, quote-capsules, quote-capsules-exact, or explicit")
	c.Flags().StringVar(&external, "external-id", "NPR-1781", "MediaSum validation document")
	c.Flags().StringVar(&model, "model", "accounts/fireworks/models/deepseek-v4-pro", "account-visible model; override explicitly")
	c.Flags().IntVar(&replicate, "replicate", 1, "positive replicate")
	c.Flags().IntVar(&max, "max-facts", 45, "per chunk cap")
	c.Flags().IntVar(&chars, "chunk-chars", 10000, "speaker-packed character budget")
	c.Flags().IntVar(&carry, "carryover-turns", 0, "preceding turns on each non-first chunk")
	return c
}
