package mediasum

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Record struct {
	ID         string   `json:"id"`
	Program    string   `json:"program"`
	Date       string   `json:"date"`
	URL        string   `json:"url"`
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	Utterances []string `json:"utt"`
	Speakers   []string `json:"speaker"`
}

func Ingest(ctx context.Context, pool *pgxpool.Pool, path, split string, limit int) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	count := 0
	for scanner.Scan() {
		if limit > 0 && count >= limit {
			break
		}
		var record Record
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return count, fmt.Errorf("decode record %d: %w", count+1, err)
		}
		if err := persist(ctx, pool, record, split); err != nil {
			return count, fmt.Errorf("persist %s: %w", record.ID, err)
		}
		count++
	}
	return count, scanner.Err()
}

func persist(ctx context.Context, pool *pgxpool.Pool, record Record, split string) error {
	if record.ID == "" || len(record.Utterances) != len(record.Speakers) {
		return fmt.Errorf("invalid record")
	}
	var date *time.Time
	if parsed, err := time.Parse("2006-01-02", record.Date); err == nil {
		date = &parsed
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var documentID string
	err = tx.QueryRow(ctx, `INSERT INTO documents (external_id, corpus, corpus_split, programme, published_on, source_url, title, summary) VALUES ($1,'mediasum',$2,$3,$4,$5,$6,$7) ON CONFLICT (external_id) DO UPDATE SET corpus_split=EXCLUDED.corpus_split, programme=EXCLUDED.programme, published_on=EXCLUDED.published_on, source_url=EXCLUDED.source_url, title=EXCLUDED.title, summary=EXCLUDED.summary, updated_at=now() RETURNING id`, record.ID, split, record.Program, date, record.URL, record.Title, record.Summary).Scan(&documentID)
	if err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM turns WHERE document_id=$1`, documentID); err != nil {
		return err
	}
	offset := 0
	for ordinal, text := range record.Utterances {
		start := offset
		end := start + len(text)
		speaker := record.Speakers[ordinal]
		if _, err = tx.Exec(ctx, `INSERT INTO turns (document_id, ordinal, speaker_label, text, char_start, char_end) VALUES ($1,$2,$3,$4,$5,$6)`, documentID, ordinal, speaker, text, start, end); err != nil {
			return err
		}
		offset = end + 1
	}
	return tx.Commit(ctx)
}

func ContentHash(text string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(text)))
	return hex.EncodeToString(sum[:])
}
