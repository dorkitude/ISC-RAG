// Package retrieval performs representation-isolated pgvector similarity search.
package retrieval

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/dorkitude/isc-rag/internal/embedding"
	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Representation, Version, CorpusSplit, EmbeddingModel string
	TopK                                                 int
}
type Result struct {
	SourceUnitID, DocumentID, ExternalID, Text string
	Quote, SpeakerLabel                        string
	TurnStart, TurnEnd                         int
	Score                                      float64
}

func Search(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, query string, cfg Config) ([]Result, error) {
	if cfg.TopK <= 0 {
		cfg.TopK = 5
	}
	vectors, err := client.Embed(ctx, cfg.EmbeddingModel, []string{query})
	if err != nil {
		return nil, err
	}
	if len(vectors) != 1 {
		return nil, fmt.Errorf("expected one query embedding")
	}
	return SearchVector(ctx, pool, embedding.VectorString(vectors[0]), cfg)
}
func SearchVector(ctx context.Context, pool *pgxpool.Pool, vector string, cfg Config) ([]Result, error) {
	var dimension int
	if err := pool.QueryRow(ctx, `SELECT dimensions FROM embeddings WHERE model=$1 LIMIT 1`, cfg.EmbeddingModel).Scan(&dimension); err != nil {
		return nil, fmt.Errorf("find embedding dimensions for %q: %w", cfg.EmbeddingModel, err)
	}
	if dimension < 1 || dimension > 16000 {
		return nil, fmt.Errorf("invalid embedding dimension %d", dimension)
	}
	statement := fmt.Sprintf(`SELECT su.id,su.document_id,d.external_id,su.text,COALESCE(proof.quote,''),COALESCE(proof.speaker_label,''),su.turn_start,su.turn_end,1-(e.vector::vector(%d) <=> $1::vector(%d)) AS score FROM embeddings e JOIN source_units su ON su.id=e.source_unit_id JOIN documents d ON d.id=su.document_id LEFT JOIN LATERAL (SELECT fe.quote,fe.speaker_label FROM facts f JOIN fact_evidence fe ON fe.fact_id=f.id WHERE su.representation='facts' AND f.document_id=su.document_id AND f.content_hash=su.content_hash ORDER BY fe.id LIMIT 1) proof ON true WHERE e.model=$2 AND su.representation=$3 AND su.version=$4 AND d.corpus_split=$5 ORDER BY e.vector::vector(%d) <=> $1::vector(%d) LIMIT $6`, dimension, dimension, dimension, dimension)
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	// pgvector's HNSW scan returns at most hnsw.ef_search rows (default 40)
	// regardless of LIMIT; raise it to the requested depth so top-k is honored.
	if _, err = tx.Exec(ctx, fmt.Sprintf("SET LOCAL hnsw.ef_search = %d", clampEfSearch(cfg.TopK))); err != nil {
		return nil, err
	}
	rows, err := tx.Query(ctx, statement, vector, cfg.EmbeddingModel, cfg.Representation, cfg.Version, cfg.CorpusSplit, cfg.TopK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Result
	for rows.Next() {
		var r Result
		if err := rows.Scan(&r.SourceUnitID, &r.DocumentID, &r.ExternalID, &r.Text, &r.Quote, &r.SpeakerLabel, &r.TurnStart, &r.TurnEnd, &r.Score); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()
	return out, tx.Commit(ctx)
}

// clampEfSearch bounds hnsw.ef_search to pgvector's accepted range [1, 1000]
// while never lowering it beneath the default of 40.
func clampEfSearch(topK int) int {
	if topK < 40 {
		return 40
	}
	if topK > 1000 {
		return 1000
	}
	return topK
}

// EnsureHNSW creates the partial, model-specific index required for a dimensioned
// pgvector HNSW index. It refuses mixed dimensions rather than silently indexing
// an incomparable embedding space.
func EnsureHNSW(ctx context.Context, pool *pgxpool.Pool, model string) (string, error) {
	var minimum, maximum int
	if err := pool.QueryRow(ctx, `SELECT COALESCE(min(dimensions),0),COALESCE(max(dimensions),0) FROM embeddings WHERE model=$1`, model).Scan(&minimum, &maximum); err != nil {
		return "", err
	}
	if minimum == 0 {
		return "", fmt.Errorf("no embeddings exist for model %q", model)
	}
	if minimum != maximum {
		return "", fmt.Errorf("model %q has mixed dimensions (%d-%d)", model, minimum, maximum)
	}
	sum := sha256.Sum256([]byte(model))
	name := fmt.Sprintf("embeddings_hnsw_%x", sum[:6])
	// dimensions are queried from PostgreSQL; the model literal is safely quoted.
	statement := fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON embeddings USING hnsw ((vector::vector(%d)) vector_cosine_ops) WHERE model = '%s'", name, minimum, strings.ReplaceAll(model, "'", "''"))
	if _, err := pool.Exec(ctx, statement); err != nil {
		return "", err
	}
	return name, nil
}

// SearchVectorHybrid fuses pgvector similarity with PostgreSQL full-text
// ranking via reciprocal rank fusion (k=60), approximating the
// embedding+BM25 hybrid of contextual-retrieval systems. Both legs draw
// candidateK candidates; the fused top cfg.TopK are returned.
func SearchVectorHybrid(ctx context.Context, pool *pgxpool.Pool, vector, queryText string, cfg Config, candidateK int) ([]Result, error) {
	if candidateK < cfg.TopK {
		candidateK = cfg.TopK
	}
	var dimension int
	if err := pool.QueryRow(ctx, `SELECT dimensions FROM embeddings WHERE model=$1 LIMIT 1`, cfg.EmbeddingModel).Scan(&dimension); err != nil {
		return nil, fmt.Errorf("find embedding dimensions for %q: %w", cfg.EmbeddingModel, err)
	}
	if dimension < 1 || dimension > 16000 {
		return nil, fmt.Errorf("invalid embedding dimension %d", dimension)
	}
	statement := fmt.Sprintf(`WITH scoped AS (
  SELECT su.id, su.document_id, d.external_id, su.text, su.turn_start, su.turn_end, e.vector
  FROM embeddings e JOIN source_units su ON su.id=e.source_unit_id JOIN documents d ON d.id=su.document_id
  WHERE e.model=$2 AND su.representation=$3 AND su.version=$4 AND d.corpus_split=$5
), vec AS (
  SELECT id, row_number() OVER (ORDER BY vector::vector(%d) <=> $1::vector(%d)) AS r
  FROM scoped ORDER BY vector::vector(%d) <=> $1::vector(%d) LIMIT $7
), fts AS (
  SELECT id, row_number() OVER (ORDER BY ts_rank_cd(to_tsvector('english', text), websearch_to_tsquery('english', $6)) DESC) AS r
  FROM scoped WHERE to_tsvector('english', text) @@ websearch_to_tsquery('english', $6) LIMIT $7
), fused AS (
  SELECT COALESCE(v.id, f.id) AS id,
    COALESCE(1.0/(60+v.r), 0) + COALESCE(1.0/(60+f.r), 0) AS score
  FROM vec v FULL OUTER JOIN fts f USING (id)
)
SELECT s.id, s.document_id, s.external_id, s.text,
  COALESCE(proof.quote,''), COALESCE(proof.speaker_label,''), s.turn_start, s.turn_end, fused.score
FROM fused JOIN scoped s ON s.id=fused.id
LEFT JOIN LATERAL (SELECT fe.quote, fe.speaker_label FROM source_units su2 JOIN facts fct ON su2.representation='facts' AND fct.document_id=su2.document_id AND fct.content_hash=su2.content_hash JOIN fact_evidence fe ON fe.fact_id=fct.id WHERE su2.id=s.id ORDER BY fe.id LIMIT 1) proof ON true
ORDER BY fused.score DESC LIMIT $8`, dimension, dimension, dimension, dimension)
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, fmt.Sprintf("SET LOCAL hnsw.ef_search = %d", clampEfSearch(candidateK))); err != nil {
		return nil, err
	}
	rows, err := tx.Query(ctx, statement, vector, cfg.EmbeddingModel, cfg.Representation, cfg.Version, cfg.CorpusSplit, queryText, candidateK, cfg.TopK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Result
	for rows.Next() {
		var r Result
		if err := rows.Scan(&r.SourceUnitID, &r.DocumentID, &r.ExternalID, &r.Text, &r.Quote, &r.SpeakerLabel, &r.TurnStart, &r.TurnEnd, &r.Score); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()
	return out, tx.Commit(ctx)
}
