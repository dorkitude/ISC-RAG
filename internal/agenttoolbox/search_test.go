package agenttoolbox

import (
	"testing"

	"github.com/dorkitude/isc-rag/internal/retrieval"
)

func TestChunkCapabilityHasOnlySearchCommand(t *testing.T) {
	command := NewCommand(FixedChunkRAG)
	commands := command.Commands()
	if len(commands) != 1 || commands[0].Name() != "search" {
		t.Fatalf("expected exactly one search command, got %#v", commands)
	}
	if !command.CompletionOptions.DisableDefaultCmd {
		t.Fatal("capability binary must not expose the generated completion command")
	}
	for _, forbidden := range []string{"representation", "version", "split", "embedding-model"} {
		if commands[0].Flags().Lookup(forbidden) != nil {
			t.Fatalf("search must not expose --%s", forbidden)
		}
	}
}

func TestNormalizeResultsUsesStableJSONFields(t *testing.T) {
	result := normalizeResults([]retrieval.Result{{SourceUnitID: "unit", DocumentID: "doc", ExternalID: "NPR-1", Text: "text", TurnStart: 2, TurnEnd: 3, Score: 0.5}})
	if len(result) != 1 || result[0].SourceUnitID != "unit" || result[0].TurnEnd != 3 {
		t.Fatalf("unexpected normalized result: %#v", result)
	}
}

func TestChunkCapabilitiesAreDistinctAndFrozen(t *testing.T) {
	if FixedChunkRAG.Representation != "fixed" || TurnChunkRAG.Representation != "turn" || SemanticChunkRAG.Representation != "semantic" {
		t.Fatalf("unexpected representations: %#v %#v %#v", FixedChunkRAG, TurnChunkRAG, SemanticChunkRAG)
	}
	for _, capability := range []Capability{FixedChunkRAG, TurnChunkRAG, SemanticChunkRAG} {
		if capability.Version != defaultVersion || capability.CorpusSplit != defaultCorpusSplit || capability.EmbeddingModel != defaultEmbeddingModel {
			t.Fatalf("unfrozen capability: %#v", capability)
		}
	}
}
