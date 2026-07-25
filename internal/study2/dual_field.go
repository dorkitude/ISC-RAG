package study2

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/jackc/pgx/v5/pgxpool"
)

const DualExpansionPrompt = `Return ONLY JSON: {"expansions":[{"fact_id":"uuid","text":"resolved retrieval description"}]}. Each input is an authoritative exact dialogue capsule. Produce a short resolved query-oriented retrieval description for every item. It is NON-AUTHORITATIVE: do not present it as a fact, evidence, quotation, or answer; preserve uncertainty and attribution. Do not add information. Return every supplied fact_id exactly once.`

type DualConfig struct {
	Version, Model string
	BatchSize      int
}
type dualItem struct {
	FactID, Text string
	DocumentID   string
	Start, End   int
}
type dualOut struct {
	Expansions []struct {
		FactID string `json:"fact_id"`
		Text   string `json:"text"`
	} `json:"expansions"`
}

// ExpandCapsules indexes only the generated retrieval text. Answer generation
// must dereference study_two_dual_field.capsule_fact_id for evidence.
func ExpandCapsules(ctx context.Context, p *pgxpool.Pool, c *fireworks.Client, cfg DualConfig) (int, error) {
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 20
	}
	rows, err := p.Query(ctx, `SELECT f.id,r.document_id,f.text,f.turn_start,f.turn_end FROM study_two_facts f JOIN study_two_runs r ON r.id=f.run_id WHERE r.candidate='quote-capsules-exact' AND r.study_version='v2-quote-capsules-exact' AND NOT EXISTS (SELECT 1 FROM study_two_dual_field d WHERE d.capsule_fact_id=f.id) ORDER BY f.id`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var items []dualItem
	for rows.Next() {
		var x dualItem
		if err = rows.Scan(&x.FactID, &x.DocumentID, &x.Text, &x.Start, &x.End); err != nil {
			return 0, err
		}
		items = append(items, x)
	}
	if err = rows.Err(); err != nil {
		return 0, err
	}
	n := 0
	for i := 0; i < len(items); i += cfg.BatchSize {
		end := i + cfg.BatchSize
		if end > len(items) {
			end = len(items)
		}
		b := items[i:end]
		payload, _ := json.Marshal(map[string]any{"capsules": b})
		r, e := studyChat(ctx, c, cfg.Model, DualExpansionPrompt, string(payload), 3000)
		if e != nil {
			return n, e
		}
		var out dualOut
		if e = json.Unmarshal([]byte(r.Text), &out); e != nil {
			return n, e
		}
		by := map[string]string{}
		for _, x := range out.Expansions {
			by[x.FactID] = x.Text
		}
		tx, e := p.Begin(ctx)
		if e != nil {
			return n, e
		}
		for _, x := range b {
			text := by[x.FactID]
			if text == "" {
				continue
			}
			var su string
			e = tx.QueryRow(ctx, `INSERT INTO source_units(document_id,representation,version,turn_start,turn_end,char_start,char_end,text,content_hash) VALUES($1,'study2-capsule-expansion',$2,$3,$4,0,0,$5,$6) RETURNING id`, x.DocumentID, cfg.Version, x.Start, x.End, text, hash(text)).Scan(&su)
			if e == nil {
				_, e = tx.Exec(ctx, `INSERT INTO study_two_dual_field(capsule_fact_id,retrieval_source_unit_id,retrieval_expansion,expansion_hash,expansion_model,expansion_prompt) VALUES($1,$2,$3,$4,$5,$6) ON CONFLICT(capsule_fact_id) DO NOTHING`, x.FactID, su, text, hash(text), cfg.Model, DualExpansionPrompt)
			}
			if e != nil {
				tx.Rollback(ctx)
				return n, e
			}
			n++
		}
		if e = tx.Commit(ctx); e != nil {
			return n, e
		}
	}
	return n, nil
}

var _ = fmt.Sprintf
