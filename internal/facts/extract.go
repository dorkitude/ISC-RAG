// Package facts extracts and validates atomic dialogue facts.
package facts

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/jackc/pgx/v5/pgxpool"
)

const SystemPrompt = `You extract explicit, atomic facts from human dialogue. Return ONLY JSON shaped as {"facts":[{"text":"standalone atomic claim","quote":"exact supporting quote"}]}. Extract only claims explicitly stated in the provided text. Do not speculate, summarize conversational dynamics, infer unstated context, or write meta-claims. Preserve attribution through the supporting quote. Every fact must be understandable on its own. Quotes must be verbatim, concise, and under 240 characters.`

type Config struct {
	Version, InputVersion, CorpusSplit, Model string
	Limit, MaxFacts                           int
}
type extracted struct {
	Facts []candidate `json:"facts"`
}
type candidate struct {
	Text  string `json:"text"`
	Quote string `json:"quote"`
}

type enrichmentItem struct {
	FactID        string `json:"fact_id"`
	Topic         string `json:"topic"`
	SubjectEntity string `json:"subject_entity"`
	RelationType  string `json:"relation_type"`
	TargetEntity  string `json:"target_entity"`
}
type enrichmentResponse struct {
	Items []enrichmentItem `json:"items"`
}

const enrichmentPrompt = `For each supplied immutable fact, return ONLY JSON shaped as {"items":[{"fact_id":"...","topic":"1-3 word title case category","subject_entity":"short entity","relation_type":"ALL_CAPS VERB PHRASE","target_entity":"short entity"}]}. Preserve every fact ID. Do not alter, repeat, or infer beyond the fact text. If an element is unclear, use an empty string.`

type inputUnit struct {
	ID, DocumentID, Text string
	CharStart            int
}

func Extract(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, cfg Config) (int, error) {
	if cfg.Version == "" || cfg.InputVersion == "" || cfg.Model == "" {
		return 0, fmt.Errorf("version, input version, and model are required")
	}
	if cfg.MaxFacts <= 0 {
		cfg.MaxFacts = 50
	}
	rows, err := pool.Query(ctx, `SELECT su.id,su.document_id,su.text,su.char_start FROM source_units su JOIN documents d ON d.id=su.document_id WHERE su.representation='fact-input' AND su.version=$1 AND su.document_id IN (SELECT id FROM documents WHERE corpus='mediasum' AND corpus_split=$2 ORDER BY external_id LIMIT $3) ORDER BY d.external_id,su.char_start`, cfg.InputVersion, cfg.CorpusSplit, limitSQL(cfg.Limit))
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	byDoc := map[string][]inputUnit{}
	var order []string
	for rows.Next() {
		var u inputUnit
		if err := rows.Scan(&u.ID, &u.DocumentID, &u.Text, &u.CharStart); err != nil {
			return 0, err
		}
		if _, ok := byDoc[u.DocumentID]; !ok {
			order = append(order, u.DocumentID)
		}
		byDoc[u.DocumentID] = append(byDoc[u.DocumentID], u)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	runID, err := startRun(ctx, pool, cfg)
	if err != nil {
		return 0, err
	}
	processed := 0
	for _, docID := range order {
		var all []candidate
		for _, u := range byDoc[docID] {
			maxFacts := min(cfg.MaxFacts, 50)
			response, err := client.ChatJSON(ctx, cfg.Model, SystemPrompt, fmt.Sprintf("<TEXT>\n%s\n</TEXT>\nReturn at most %d facts.", u.Text, maxFacts), 2000)
			if err != nil {
				failRun(ctx, pool, runID)
				return processed, fmt.Errorf("extract document %s: %w", docID, err)
			}
			facts, err := decode(response.Text)
			if err != nil {
				failRun(ctx, pool, runID)
				return processed, fmt.Errorf("parse document %s: %w", docID, err)
			}
			all = append(all, facts.Facts...)
		}
		if err := persistDocument(ctx, pool, docID, cfg, runID, all); err != nil {
			failRun(ctx, pool, runID)
			return processed, err
		}
		processed++
	}
	_, err = pool.Exec(ctx, `UPDATE pipeline_runs SET status='complete',completed_at=now() WHERE id=$1`, runID)
	return processed, err
}

func startRun(ctx context.Context, pool *pgxpool.Pool, cfg Config) (string, error) {
	var id string
	err := pool.QueryRow(ctx, `INSERT INTO pipeline_runs(kind,version,model,parameters) VALUES ('fact-extraction',$1,$2,jsonb_build_object('input_version',$3::text,'split',$4::text)) RETURNING id`, cfg.Version, cfg.Model, cfg.InputVersion, cfg.CorpusSplit).Scan(&id)
	return id, err
}
func failRun(ctx context.Context, pool *pgxpool.Pool, id string) {
	_, _ = pool.Exec(ctx, `UPDATE pipeline_runs SET status='failed',completed_at=now() WHERE id=$1`, id)
}

// Enrich adds optional search metadata without modifying fact text or provenance.
func Enrich(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, version, model string, limit int) (int, error) {
	rows, err := pool.Query(ctx, `SELECT f.id,f.text FROM facts f JOIN pipeline_runs p ON p.id=f.pipeline_run_id WHERE p.kind='fact-extraction' AND p.version=$1 ORDER BY f.id LIMIT $2`, version, limitSQL(limit))
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	type row struct{ id, text string }
	var facts []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.id, &r.text); err != nil {
			return 0, err
		}
		facts = append(facts, r)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	updated := 0
	for start := 0; start < len(facts); start += 25 {
		end := min(len(facts), start+25)
		batch := facts[start:end]
		items := make([]map[string]string, len(batch))
		for i, fact := range batch {
			items[i] = map[string]string{"fact_id": fact.id, "text": fact.text}
		}
		payload, _ := json.Marshal(map[string]any{"facts": items})
		response, err := client.ChatJSON(ctx, model, enrichmentPrompt, string(payload), 3000)
		if err != nil {
			return updated, err
		}
		parsed, err := decodeEnrichment(response.Text)
		if err != nil {
			return updated, err
		}
		permitted := map[string]bool{}
		for _, f := range batch {
			permitted[f.id] = true
		}
		tx, err := pool.Begin(ctx)
		if err != nil {
			return updated, err
		}
		for _, item := range parsed.Items {
			if !permitted[item.FactID] {
				continue
			}
			if _, err = tx.Exec(ctx, `UPDATE facts SET topic=$2,subject_entity=$3,relation_type=$4,target_entity=$5 WHERE id=$1`, item.FactID, trim(item.Topic), trim(item.SubjectEntity), trim(item.RelationType), trim(item.TargetEntity)); err != nil {
				_ = tx.Rollback(ctx)
				return updated, err
			}
			updated++
		}
		if err = tx.Commit(ctx); err != nil {
			return updated, err
		}
	}
	return updated, nil
}
func persistDocument(ctx context.Context, pool *pgxpool.Pool, documentID string, cfg Config, runID string, candidates []candidate) error {
	transcript, turns, err := documentText(ctx, pool, documentID)
	if err != nil {
		return err
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `DELETE FROM facts WHERE document_id=$1`, documentID); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM source_units WHERE document_id=$1 AND representation='facts' AND version=$2`, documentID, cfg.Version); err != nil {
		return err
	}
	seen := map[string]bool{}
	for _, c := range candidates {
		text := strings.TrimSpace(c.Text)
		quote := strings.TrimSpace(c.Quote)
		if text == "" || quote == "" || len(quote) > 240 {
			continue
		}
		key := hash(text)
		if seen[key] {
			continue
		}
		seen[key] = true
		at := strings.Index(transcript, quote)
		if at < 0 {
			continue
		}
		start, end := at, at+len(quote)
		first, last, speaker := evidenceTurns(turns, start, end)
		var factID string
		err = tx.QueryRow(ctx, `INSERT INTO facts(document_id,text,content_hash,pipeline_run_id) VALUES ($1,$2,$3,$4) RETURNING id`, documentID, text, key, runID).Scan(&factID)
		if err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO fact_evidence(fact_id,quote,char_start,char_end,turn_start,turn_end,validation_method,speaker_label) VALUES ($1,$2,$3,$4,$5,$6,'exact',$7)`, factID, quote, start, end, first, last, speaker); err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO source_units(document_id,representation,version,turn_start,turn_end,char_start,char_end,text,content_hash,pipeline_run_id) VALUES ($1,'facts',$2,$3,$4,$5,$6,$7,$8,$9)`, documentID, cfg.Version, first, last, start, end, text, key, runID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

type canonicalTurn struct {
	ordinal, start, end int
	speaker             string
}

func documentText(ctx context.Context, pool *pgxpool.Pool, id string) (string, []canonicalTurn, error) {
	rows, err := pool.Query(ctx, `SELECT ordinal,text,char_start,char_end,speaker_label FROM turns WHERE document_id=$1 ORDER BY ordinal`, id)
	if err != nil {
		return "", nil, err
	}
	defer rows.Close()
	var pieces []string
	var ts []canonicalTurn
	for rows.Next() {
		var t canonicalTurn
		var text string
		if err := rows.Scan(&t.ordinal, &text, &t.start, &t.end, &t.speaker); err != nil {
			return "", nil, err
		}
		pieces = append(pieces, text)
		ts = append(ts, t)
	}
	return strings.Join(pieces, "\n"), ts, rows.Err()
}
func evidenceTurns(ts []canonicalTurn, start, end int) (int, int, string) {
	if len(ts) == 0 {
		return 0, 0, ""
	}
	first, last := ts[0], ts[len(ts)-1]
	for _, t := range ts {
		if start >= t.start && start <= t.end {
			first = t
		}
		if end >= t.start && end <= t.end+1 {
			last = t
			break
		}
	}
	return first.ordinal, last.ordinal, first.speaker
}
func cleanJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)
	if start, end := strings.Index(raw, "{"), strings.LastIndex(raw, "}"); start >= 0 && end > start {
		return raw[start : end+1]
	}
	return raw
}
func decode(raw string) (extracted, error) {
	raw = cleanJSON(raw)
	var out extracted
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &out); err != nil {
		return out, err
	}
	return out, nil
}
func decodeEnrichment(raw string) (enrichmentResponse, error) {
	var out enrichmentResponse
	if err := json.Unmarshal([]byte(cleanJSON(raw)), &out); err != nil {
		return out, err
	}
	return out, nil
}
func trim(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return strings.TrimSpace(s)
}
func hash(s string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(s)))
	return hex.EncodeToString(sum[:])
}
func limitSQL(n int) int {
	if n <= 0 {
		return 2147483647
	}
	return n
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
