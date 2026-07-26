// Package embedding writes Fireworks embeddings into the canonical pgvector table.
package embedding

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Representation, Version, CorpusSplit, Model string
	Limit, BatchSize                            int
	FactAndQuote                                bool
}
type source struct{ ID, Text string }

func Embed(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, cfg Config) (int, error) {
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 32
	}
	if cfg.Representation == "" || cfg.Version == "" || cfg.Model == "" {
		return 0, fmt.Errorf("representation, version, and model are required")
	}
	if cfg.FactAndQuote && cfg.Representation != "facts" {
		return 0, fmt.Errorf("fact-and-quote embedding is only valid for facts")
	}
	statement := `SELECT su.id,su.text FROM source_units su JOIN documents d ON d.id=su.document_id LEFT JOIN embeddings e ON e.source_unit_id=su.id AND e.model=$1 WHERE su.representation=$2 AND su.version=$3 AND d.corpus_split=$4 AND e.id IS NULL ORDER BY su.id LIMIT $5`
	if cfg.FactAndQuote {
		statement = `SELECT su.id, 'Fact: ' || su.text || E'\n\nSupporting quote: ' || proof.quote FROM source_units su JOIN documents d ON d.id=su.document_id LEFT JOIN embeddings e ON e.source_unit_id=su.id AND e.model=$1 JOIN LATERAL (SELECT fe.quote FROM facts f JOIN fact_evidence fe ON fe.fact_id=f.id WHERE f.document_id=su.document_id AND f.content_hash=su.content_hash ORDER BY fe.id LIMIT 1) proof ON true WHERE su.representation=$2 AND su.version=$3 AND d.corpus_split=$4 AND e.id IS NULL ORDER BY su.id LIMIT $5`
	}
	rows, err := pool.Query(ctx, statement, cfg.Model, cfg.Representation, cfg.Version, cfg.CorpusSplit, limitSQL(cfg.Limit))
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var sources []source
	for rows.Next() {
		var s source
		if err := rows.Scan(&s.ID, &s.Text); err != nil {
			return 0, err
		}
		sources = append(sources, s)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	for start := 0; start < len(sources); start += cfg.BatchSize {
		end := start + cfg.BatchSize
		if end > len(sources) {
			end = len(sources)
		}
		batch := sources[start:end]
		input := make([]string, len(batch))
		for i, s := range batch {
			input[i] = s.Text
		}
		vectors, err := client.Embed(ctx, cfg.Model, input)
		if err != nil {
			return start, err
		}
		tx, err := pool.Begin(ctx)
		if err != nil {
			return start, err
		}
		for i, v := range vectors {
			if len(v) == 0 {
				_ = tx.Rollback(ctx)
				return start, fmt.Errorf("empty embedding for source unit %s", batch[i].ID)
			}
			if _, err = tx.Exec(ctx, `INSERT INTO embeddings(source_unit_id,model,dimensions,vector) VALUES ($1,$2,$3,$4::vector) ON CONFLICT(source_unit_id,model) DO UPDATE SET dimensions=EXCLUDED.dimensions,vector=EXCLUDED.vector,created_at=now()`, batch[i].ID, cfg.Model, len(v), VectorString(v)); err != nil {
				_ = tx.Rollback(ctx)
				return start, err
			}
		}
		if err = tx.Commit(ctx); err != nil {
			return start, err
		}
	}
	return len(sources), nil
}
func VectorString(values []float64) string {
	var b strings.Builder
	b.WriteByte('[')
	for i, v := range values {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.FormatFloat(v, 'g', -1, 64))
	}
	b.WriteByte(']')
	return b.String()
}
func limitSQL(n int) int {
	if n <= 0 {
		return 2147483647
	}
	return n
}
