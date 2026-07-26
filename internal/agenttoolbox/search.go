// Package agenttoolbox provides deliberately capability-limited CLIs for RAG
// evaluation agents. A binary is compiled for exactly one representation.
package agenttoolbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/dorkitude/isc-rag/internal/retrieval"
	"github.com/dorkitude/isc-rag/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const (
	defaultVersion        = "v1"
	defaultCorpusSplit    = "test"
	defaultEmbeddingModel = "nomic-ai/nomic-embed-text-v1.5"
	maxTopK               = 5
)

// Capability is compiled into an individual binary, not selected by an agent
// flag. Keep versions frozen per capability so conditions cannot be mixed.
type Capability struct {
	Name           string
	Representation string
	Version        string
	CorpusSplit    string
	EmbeddingModel string
}

var (
	FixedChunkRAG = Capability{
		Name: "fixed-chunk-rag", Representation: "fixed", Version: defaultVersion,
		CorpusSplit: defaultCorpusSplit, EmbeddingModel: defaultEmbeddingModel,
	}
	// PilotFixedChunkRAG is frozen to the same 100-document CNN development
	// slice as NamedFactRAG. Keep it separate from the v1 binary: a model must
	// never be able to choose its representation version at runtime.
	PilotFixedChunkRAG = Capability{
		Name: "pilot-fixed-chunk-rag", Representation: "fixed", Version: "pilot-100x100-v1",
		CorpusSplit: defaultCorpusSplit, EmbeddingModel: defaultEmbeddingModel,
	}
	TurnChunkRAG = Capability{
		Name: "turn-chunk-rag", Representation: "turn", Version: defaultVersion,
		CorpusSplit: defaultCorpusSplit, EmbeddingModel: defaultEmbeddingModel,
	}
	SemanticChunkRAG = Capability{
		Name: "semantic-chunk-rag", Representation: "semantic", Version: defaultVersion,
		CorpusSplit: defaultCorpusSplit, EmbeddingModel: defaultEmbeddingModel,
	}
	// NamedFactRAG is a development-only capability for the locally matched
	// CNN pilot. It returns each rephrased fact together with its verified quote
	// and speaker label; agents cannot select a different fact version.
	NamedFactRAG = Capability{
		Name: "named-fact-rag", Representation: "facts", Version: "pilot-named-v3",
		CorpusSplit: defaultCorpusSplit, EmbeddingModel: defaultEmbeddingModel,
	}
	NamedFactQuoteRAG = Capability{
		Name: "named-fact-quote-rag", Representation: "facts", Version: "pilot-named-v3-quote-v1",
		CorpusSplit: defaultCorpusSplit, EmbeddingModel: defaultEmbeddingModel,
	}
	// TurnwiseOpinionFactRAG is the bounded 100-document development condition.
	// Its units are high-recall, two-pass claims (including explicit opinions),
	// retrieved with the fact plus its verified quote.  It remains a development
	// capability until its question set and extraction quality are frozen.
	TurnwiseOpinionFactRAG = Capability{
		Name: "turnwise-opinion-fact-rag", Representation: "facts", Version: "turnwise-opinion-v1",
		CorpusSplit: defaultCorpusSplit, EmbeddingModel: defaultEmbeddingModel,
	}
	// FactDevelopmentFixedRAG is the matched fixed-chunk control over the same
	// frozen 100 documents as TurnwiseOpinionFactRAG.
	FactDevelopmentFixedRAG = Capability{
		Name: "fact-development-fixed-rag", Representation: "fixed", Version: "fact-dev-sample-100-fixed-v1",
		CorpusSplit: defaultCorpusSplit, EmbeddingModel: defaultEmbeddingModel,
	}
)

type searchResponse struct {
	Capability     string         `json:"capability"`
	Representation string         `json:"representation"`
	Version        string         `json:"version"`
	CorpusSplit    string         `json:"corpus_split"`
	Query          string         `json:"query"`
	Results        []searchResult `json:"results"`
}

type searchResult struct {
	SourceUnitID string  `json:"source_unit_id"`
	DocumentID   string  `json:"document_id"`
	ExternalID   string  `json:"external_id"`
	Text         string  `json:"text"`
	Quote        string  `json:"quote"`
	SpeakerLabel string  `json:"speaker_label"`
	TurnStart    int     `json:"turn_start"`
	TurnEnd      int     `json:"turn_end"`
	Score        float64 `json:"score"`
}

// Execute starts a one-capability CLI. It intentionally exposes only search;
// there are no flags for representation, model, version, or corpus split.
func Execute(capability Capability) {
	root := NewCommand(capability)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func NewCommand(capability Capability) *cobra.Command {
	var projectDir, query string
	var topK int
	root := &cobra.Command{
		Use:           "agent-toolbox",
		Short:         "Capability-limited RAG search for evaluation agents",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.CompletionOptions.DisableDefaultCmd = true
	search := &cobra.Command{
		Use:   "search",
		Short: "Search the binary's fixed retrieval representation",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(query) == "" {
				return errors.New("--query is required")
			}
			if topK < 1 || topK > maxTopK {
				return fmt.Errorf("--top-k must be between 1 and %d", maxTopK)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			pool, err := openReadOnlyDatabase(cmd.Context(), projectDir)
			if err != nil {
				return err
			}
			defer pool.Close()
			client, err := fireworks.NewFromEnv()
			if err != nil {
				return err
			}
			results, err := retrieval.Search(cmd.Context(), pool, client, query, retrieval.Config{
				Representation: capability.Representation,
				Version:        capability.Version,
				CorpusSplit:    capability.CorpusSplit,
				EmbeddingModel: capability.EmbeddingModel,
				TopK:           topK,
			})
			if err != nil {
				return err
			}
			encoded, err := json.MarshalIndent(searchResponse{
				Capability: capability.Name, Representation: capability.Representation,
				Version: capability.Version, CorpusSplit: capability.CorpusSplit,
				Query: query, Results: normalizeResults(results),
			}, "", "  ")
			if err != nil {
				return err
			}
			cmd.Println(string(encoded))
			return nil
		},
	}
	search.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	search.Flags().StringVar(&query, "query", "", "Question to retrieve evidence for")
	search.Flags().IntVar(&topK, "top-k", maxTopK, "Results to return (1-5)")
	root.AddCommand(search)
	return root
}

func normalizeResults(results []retrieval.Result) []searchResult {
	out := make([]searchResult, len(results))
	for i, result := range results {
		out[i] = searchResult{
			SourceUnitID: result.SourceUnitID, DocumentID: result.DocumentID,
			ExternalID: result.ExternalID, Text: result.Text, Quote: result.Quote,
			SpeakerLabel: result.SpeakerLabel, TurnStart: result.TurnStart,
			TurnEnd: result.TurnEnd, Score: result.Score,
		}
	}
	return out
}

func openReadOnlyDatabase(ctx context.Context, projectDir string) (*pgxpool.Pool, error) {
	absoluteProjectDir, err := filepath.Abs(projectDir)
	if err != nil {
		return nil, fmt.Errorf("resolve project directory: %w", err)
	}
	envFile := filepath.Join(absoluteProjectDir, ".env")
	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err != nil {
			return nil, fmt.Errorf("load %s: %w", envFile, err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("inspect %s: %w", envFile, err)
	}
	url := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if url == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	return storage.OpenReadOnly(ctx, url)
}
