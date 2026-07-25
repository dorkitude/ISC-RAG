package mediasum

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SpeakerMetadata is parsed only from the corpus-provided speaker label. It
// deliberately does not infer roles from dialogue text or external sources.
type SpeakerMetadata struct {
	Name string
	Role string
}

// ParseSpeakerLabel separates MediaSum labels such as "ALEX COHEN, host" into
// a display name and the corpus-supplied role suffix. Labels without a comma
// retain their complete value as the name and have no role.
func ParseSpeakerLabel(label string) SpeakerMetadata {
	label = normalizeSpace(label)
	if label == "" {
		return SpeakerMetadata{}
	}
	name, role, found := strings.Cut(label, ",")
	if !found || normalizeSpace(name) == "" || normalizeSpace(role) == "" {
		return SpeakerMetadata{Name: label}
	}
	return SpeakerMetadata{Name: normalizeSpace(name), Role: strings.ToLower(normalizeSpace(role))}
}

func normalizeSpace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

// BackfillSpeakerMetadata populates deterministic fields for already-ingested
// corpus turns. Re-running it is safe and does not alter the original label.
func BackfillSpeakerMetadata(ctx context.Context, pool *pgxpool.Pool) (int, error) {
	rows, err := pool.Query(ctx, `SELECT id::text, COALESCE(speaker_label, '') FROM turns`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	type row struct{ id, label string }
	var pending []row
	for rows.Next() {
		var item row
		if err = rows.Scan(&item.id, &item.label); err != nil {
			return 0, err
		}
		pending = append(pending, item)
	}
	if err = rows.Err(); err != nil {
		return 0, err
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	for _, item := range pending {
		meta := ParseSpeakerLabel(item.label)
		if _, err = tx.Exec(ctx, `UPDATE turns SET speaker_name=$2, speaker_role=NULLIF($3, '') WHERE id=$1`, item.id, meta.Name, meta.Role); err != nil {
			return 0, fmt.Errorf("update speaker metadata: %w", err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return len(pending), nil
}
