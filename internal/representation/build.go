// Package representation creates the non-LLM retrieval baselines from canonical turns.
package representation

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	Fixed     = "fixed"
	Turn      = "turn"
	Semantic  = "semantic"
	FactInput = "fact-input"
)

type Config struct {
	Representation      string
	Version             string
	CorpusSplit         string
	Limit               int
	MaxChars            int
	OverlapChars        int
	SemanticMin         int
	FactInputMax        int
	DocumentExternalIDs []string
}

type turn struct {
	Ordinal int
	Text    string
	Start   int
	End     int
}

type unit struct {
	TurnStart int
	TurnEnd   int
	CharStart int
	CharEnd   int
	Text      string
}

// Build replaces one representation/version for selected documents. Canonical turns
// remain untouched, so reruns are safe and reproducible.
func Build(ctx context.Context, pool *pgxpool.Pool, cfg Config) (int, error) {
	if cfg.Version == "" || cfg.CorpusSplit == "" {
		return 0, fmt.Errorf("version and corpus split are required")
	}
	if cfg.MaxChars < 100 {
		return 0, fmt.Errorf("max chars must be at least 100")
	}
	if cfg.OverlapChars < 0 || cfg.OverlapChars >= cfg.MaxChars {
		return 0, fmt.Errorf("overlap chars must be non-negative and smaller than max chars")
	}
	representation := cfg.Representation
	if representation == "" {
		return 0, fmt.Errorf("representation is required")
	}
	statement, args := `SELECT id FROM documents WHERE corpus='mediasum' AND corpus_split=$1`, []any{cfg.CorpusSplit}
	if len(cfg.DocumentExternalIDs) > 0 {
		statement += ` AND external_id=ANY($2)`
		args = append(args, cfg.DocumentExternalIDs)
	}
	statement += ` ORDER BY external_id LIMIT $` + fmt.Sprint(len(args)+1)
	args = append(args, limitSQL(cfg.Limit))
	rows, err := pool.Query(ctx, statement, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	for _, id := range ids {
		if err := buildDocument(ctx, pool, id, cfg); err != nil {
			return 0, err
		}
	}
	return len(ids), nil
}

func limitSQL(limit int) int {
	if limit <= 0 {
		return 2147483647
	}
	return limit
}

func buildDocument(ctx context.Context, pool *pgxpool.Pool, documentID string, cfg Config) error {
	rows, err := pool.Query(ctx, `SELECT ordinal,text,char_start,char_end FROM turns WHERE document_id=$1 ORDER BY ordinal`, documentID)
	if err != nil {
		return err
	}
	defer rows.Close()
	var turns []turn
	for rows.Next() {
		var t turn
		if err := rows.Scan(&t.Ordinal, &t.Text, &t.Start, &t.End); err != nil {
			return err
		}
		turns = append(turns, t)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	var units []unit
	switch cfg.Representation {
	case Fixed:
		units = fixedUnits(turns, cfg.MaxChars, cfg.OverlapChars)
	case Turn:
		units = turnUnits(turns, cfg.MaxChars)
	case Semantic:
		units = semanticUnits(turns, cfg.MaxChars, cfg.SemanticMin)
	case FactInput:
		inputMax := cfg.FactInputMax
		if inputMax <= 0 {
			inputMax = 10000
		}
		units = turnUnits(turns, inputMax)
	default:
		return fmt.Errorf("unsupported representation %q", cfg.Representation)
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `DELETE FROM source_units WHERE document_id=$1 AND representation=$2 AND version=$3`, documentID, cfg.Representation, cfg.Version); err != nil {
		return err
	}
	for _, u := range units {
		if _, err = tx.Exec(ctx, `INSERT INTO source_units (document_id,representation,version,turn_start,turn_end,char_start,char_end,text,content_hash) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, documentID, cfg.Representation, cfg.Version, u.TurnStart, u.TurnEnd, u.CharStart, u.CharEnd, u.Text, hash(u.Text)); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func fixedUnits(turns []turn, max, overlap int) []unit {
	if len(turns) == 0 {
		return nil
	}
	full := joinTurns(turns)
	var out []unit
	for start := 0; start < len(full); {
		end := start + max
		if end > len(full) {
			end = len(full)
		}
		if end < len(full) {
			if cut := strings.LastIndex(full[start:end], "\n"); cut > max/2 {
				end = start + cut
			}
		}
		if end <= start {
			end = min(len(full), start+max)
		}
		out = append(out, makeUnit(turns, full, start, end))
		if end == len(full) {
			break
		}
		next := end - overlap
		if next <= start {
			next = end
		}
		start = next
	}
	return out
}

func turnUnits(turns []turn, max int) []unit {
	var out []unit
	for i := 0; i < len(turns); {
		start := i
		size := 0
		for i < len(turns) {
			addition := len(turns[i].Text)
			if i > start {
				addition++
			}
			if i == start && addition > max {
				break
			}
			if size > 0 && size+addition > max {
				break
			}
			size += addition
			i++
		}
		if i == start { // split an exceptional long turn rather than dropping it.
			t := turns[i]
			for p := 0; p < len(t.Text); p += max {
				end := min(len(t.Text), p+max)
				out = append(out, unit{t.Ordinal, t.Ordinal, t.Start + p, t.Start + end, t.Text[p:end]})
			}
			i++
			continue
		}
		out = append(out, turnsUnit(turns[start:i]))
	}
	return out
}

// semanticUnits is a deterministic lexical-cohesion baseline. It favours a low
// term-overlap boundary after the minimum size, while preserving whole turns.
func semanticUnits(turns []turn, max, minimum int) []unit {
	if minimum <= 0 {
		minimum = max / 2
	}
	var out []unit
	for start := 0; start < len(turns); {
		end, startSize := start, 0
		best := -1
		bestScore := 2.0
		for end < len(turns) {
			add := len(turns[end].Text)
			if end > start {
				add++
			}
			if startSize > 0 && startSize+add > max {
				break
			}
			startSize += add
			end++
			if startSize >= minimum && end < len(turns) {
				score := jaccard(words(joinTurns(turns[start:end])), words(turns[end].Text))
				if score < bestScore {
					bestScore, best = score, end
				}
			}
		}
		if end == start {
			end++
		}
		if best > start {
			end = best
		}
		out = append(out, turnsUnit(turns[start:end]))
		start = end
	}
	return out
}

func joinTurns(turns []turn) string {
	parts := make([]string, len(turns))
	for i, t := range turns {
		parts[i] = t.Text
	}
	return strings.Join(parts, "\n")
}
func turnsUnit(ts []turn) unit {
	return unit{ts[0].Ordinal, ts[len(ts)-1].Ordinal, ts[0].Start, ts[len(ts)-1].End, joinTurns(ts)}
}
func makeUnit(turns []turn, full string, start, end int) unit {
	first, last := turns[0], turns[len(turns)-1]
	for _, t := range turns {
		if start >= t.Start && start <= t.End {
			first = t
		}
		if end >= t.Start && end <= t.End+1 {
			last = t
			break
		}
	}
	return unit{first.Ordinal, last.Ordinal, start, end, strings.TrimSpace(full[start:end])}
}
func words(s string) map[string]struct{} {
	result := map[string]struct{}{}
	for _, w := range strings.FieldsFunc(strings.ToLower(s), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }) {
		if len(w) > 2 {
			result[w] = struct{}{}
		}
	}
	return result
}
func jaccard(a, b map[string]struct{}) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 1
	}
	n := 0
	for x := range a {
		if _, ok := b[x]; ok {
			n++
		}
	}
	return float64(n) / float64(len(a)+len(b)-n)
}
func hash(text string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(text)))
	return hex.EncodeToString(sum[:])
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
