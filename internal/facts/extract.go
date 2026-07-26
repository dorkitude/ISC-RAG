// Package facts extracts and validates atomic dialogue facts.
package facts

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/jackc/pgx/v5/pgxpool"
)

const SystemPrompt = `You are a high-recall claim harvester for human dialogue. Return ONLY JSON shaped as {"facts":[{"text":"atomic attributed claim","quote":"exact supporting quote from TARGET TURN"}]}. Extract EVERY explicit atomic claim in the TARGET TURN: factual assertions AND explicitly expressed opinions, beliefs, judgments, assessments, recommendations, preferences, intentions, predictions, and hypotheses. Do not rank claims, summarize, merge claims, or omit a claim because it seems minor. CONTEXT TURNS are only for resolving names and references; never extract a claim or quote from context. Each quote must be a contiguous verbatim substring of the TARGET TURN, concise, and under 240 characters. Preserve modality, uncertainty, negation, and attribution. Use the supplied target speaker's name when attribution makes a claim clearer; never write "the speaker", "the host", or "the guest". Do not invent, strengthen, or infer a claim. If the target turn contains no extractable assertion or expressed opinion, return an empty facts array.`

var genericSpeaker = regexp.MustCompile(`(?i)\bthe (speaker|host|guest)\b`)

type Config struct {
	Version, InputVersion, CorpusSplit, Model string
	Limit, Offset, MaxFacts                   int
	DocumentExternalIDs                       []string
	Workers                                   int
}
type extracted struct {
	Facts []candidate `json:"facts"`
}
type candidate struct {
	Text        string `json:"text"`
	Quote       string `json:"quote"`
	TargetStart int    `json:"-"`
	TargetEnd   int    `json:"-"`
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
	CharStart, CharEnd   int
	TurnStart, TurnEnd   int
}

func Extract(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, cfg Config) (int, error) {
	if cfg.Version == "" || cfg.InputVersion == "" || cfg.Model == "" {
		return 0, fmt.Errorf("version, input version, and model are required")
	}
	if cfg.MaxFacts <= 0 {
		cfg.MaxFacts = 50
	}
	if cfg.Offset < 0 {
		return 0, fmt.Errorf("offset must not be negative")
	}
	if cfg.Workers <= 0 {
		cfg.Workers = 1
	}
	statement := `SELECT su.id,su.document_id,su.text,su.char_start,su.char_end,su.turn_start,su.turn_end FROM source_units su JOIN documents d ON d.id=su.document_id WHERE su.representation='fact-input' AND su.version=$1 AND d.corpus='mediasum' AND d.corpus_split=$2`
	args := []any{cfg.InputVersion, cfg.CorpusSplit}
	if len(cfg.DocumentExternalIDs) > 0 {
		statement += ` AND d.external_id=ANY($3)`
		args = append(args, cfg.DocumentExternalIDs)
	}
	statement += ` ORDER BY d.external_id,su.char_start LIMIT $` + fmt.Sprint(len(args)+1) + ` OFFSET $` + fmt.Sprint(len(args)+2)
	args = append(args, limitSQL(cfg.Limit), cfg.Offset)
	rows, err := pool.Query(ctx, statement, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	byDoc := map[string][]inputUnit{}
	var order []string
	for rows.Next() {
		var u inputUnit
		if err := rows.Scan(&u.ID, &u.DocumentID, &u.Text, &u.CharStart, &u.CharEnd, &u.TurnStart, &u.TurnEnd); err != nil {
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
	workCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	jobs := make(chan string)
	var wg sync.WaitGroup
	var mu sync.Mutex
	processed := 0
	var firstErr error
	for range cfg.Workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for docID := range jobs {
				if err := extractDocument(workCtx, pool, client, cfg, runID, docID, byDoc[docID]); err != nil {
					mu.Lock()
					if firstErr == nil {
						firstErr = err
						cancel()
					}
					mu.Unlock()
					return
				}
				mu.Lock()
				processed++
				mu.Unlock()
			}
		}()
	}
	for _, docID := range order {
		select {
		case jobs <- docID:
		case <-workCtx.Done():
			break
		}
		if workCtx.Err() != nil {
			break
		}
	}
	close(jobs)
	wg.Wait()
	if firstErr != nil {
		failRun(ctx, pool, runID)
		return processed, firstErr
	}
	_, err = pool.Exec(ctx, `UPDATE pipeline_runs SET status='complete',completed_at=now() WHERE id=$1`, runID)
	return processed, err
}

func extractDocument(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, cfg Config, runID, docID string, units []inputUnit) error {
	var all []candidate
	for _, u := range units {
		targets, err := targetedInputs(ctx, pool, u)
		if err != nil {
			return err
		}
		for _, target := range targets {
			maxFacts := min(cfg.MaxFacts, 30)
			response, err := client.ChatJSONNoReasoning(ctx, cfg.Model, SystemPrompt, target.Prompt+"\n\n"+fmt.Sprintf("Return at most %d claims; do not omit a claim for importance.", maxFacts), 2400)
			if err != nil {
				return fmt.Errorf("extract document %s turn %d: %w", docID, target.Ordinal, err)
			}
			facts, err := decode(response.Text)
			if err != nil {
				return fmt.Errorf("parse document %s turn %d: %w", docID, target.Ordinal, err)
			}
			for i := range facts.Facts {
				facts.Facts[i].TargetStart, facts.Facts[i].TargetEnd = target.CharStart, target.CharEnd
			}
			all = append(all, facts.Facts...)
		}
	}
	if err := persistDocument(ctx, pool, docID, cfg, runID, all); err != nil {
		return err
	}
	return nil
}

type targetInput struct {
	Ordinal, CharStart, CharEnd int
	Prompt                      string
}

// targetedInputs makes every canonical turn in a fact-input unit an explicit
// extraction target. Two preceding turns are context-only, so the model has
// enough local dialogue to resolve references without silently harvesting a
// broad pack and omitting minor claims.
func targetedInputs(ctx context.Context, pool *pgxpool.Pool, input inputUnit) ([]targetInput, error) {
	rows, err := pool.Query(ctx, `SELECT ordinal,COALESCE(speaker_name,''),COALESCE(speaker_label,''),text,char_start,char_end FROM turns WHERE document_id=$1 AND ordinal BETWEEN GREATEST(0,$2-2) AND $3 ORDER BY ordinal`, input.DocumentID, input.TurnStart, input.TurnEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	type turn struct {
		ordinal, start, end int
		speaker, text       string
	}
	var turns []turn
	for rows.Next() {
		var t turn
		var name, label string
		if err := rows.Scan(&t.ordinal, &name, &label, &t.text, &t.start, &t.end); err != nil {
			return nil, err
		}
		if name == "" {
			name = speakerName(label)
		}
		t.speaker = name
		turns = append(turns, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var out []targetInput
	for i, target := range turns {
		if target.ordinal < input.TurnStart {
			continue
		}
		var b strings.Builder
		if i > 0 {
			b.WriteString("<CONTEXT_TURNS>\n")
			for _, context := range turns[:i] {
				fmt.Fprintf(&b, "[SPEAKER: %s]\n%s\n", context.speaker, context.text)
			}
			b.WriteString("</CONTEXT_TURNS>\n")
		}
		fmt.Fprintf(&b, "<TARGET_TURN>\n[SPEAKER: %s]\n%s\n</TARGET_TURN>", target.speaker, target.text)
		out = append(out, targetInput{Ordinal: target.ordinal, CharStart: target.start, CharEnd: target.end, Prompt: b.String()})
	}
	return out, nil
}

// labeledInput attaches only corpus-provided speaker metadata to the source
// turns overlapping an extraction unit. It does not infer a role or name from
// dialogue text, and preserves the source text used for quote validation.
func labeledInput(ctx context.Context, pool *pgxpool.Pool, input inputUnit) (string, error) {
	rows, err := pool.Query(ctx, `SELECT COALESCE(speaker_name,''),COALESCE(speaker_label,''),text FROM turns WHERE document_id=$1 AND char_end>$2 AND char_start<$3 ORDER BY ordinal`, input.DocumentID, input.CharStart, input.CharEnd)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var out strings.Builder
	for rows.Next() {
		var name, label, text string
		if err := rows.Scan(&name, &label, &text); err != nil {
			return "", err
		}
		if name == "" {
			name = speakerName(label)
		}
		fmt.Fprintf(&out, "[SPEAKER: %s]\n%s\n", name, text)
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	if out.Len() == 0 {
		return input.Text, nil
	}
	return strings.TrimSpace(out.String()), nil
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
		if c.TargetStart < 0 || c.TargetEnd > len(transcript) || c.TargetStart >= c.TargetEnd {
			continue
		}
		local := strings.Index(transcript[c.TargetStart:c.TargetEnd], quote)
		at := -1
		if local >= 0 {
			at = c.TargetStart + local
		}
		if at < 0 {
			continue
		}
		start, end := at, at+len(quote)
		first, last, speaker := evidenceTurns(turns, start, end)
		text = replaceGenericSpeaker(text, speaker)
		if genericSpeaker.MatchString(text) {
			continue
		}
		key := hash(text)
		if seen[key] {
			continue
		}
		seen[key] = true
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

func speakerName(label string) string {
	name, _, _ := strings.Cut(strings.Join(strings.Fields(label), " "), ",")
	return strings.TrimSpace(name)
}

func replaceGenericSpeaker(text, label string) string {
	name := speakerName(label)
	if name == "" {
		return text
	}
	return genericSpeaker.ReplaceAllString(text, name)
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
