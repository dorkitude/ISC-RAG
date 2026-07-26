package evaluation

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dorkitude/isc-rag/internal/fireworks"
)

// JudgeConfig controls one blind judging pass over stored generated answers.
type JudgeConfig struct {
	RunIDs         []string
	JudgeModel     string
	CorrectnessKey string
	SupportKey     string
	// SkipExisting leaves answers untouched when both keys already hold verdicts,
	// making the pass resumable after provider interruptions.
	SkipExisting bool
}

// JudgeMetrics summarizes one judging pass.
type JudgeMetrics struct {
	Answers   int
	Judged    int
	Skipped   int
	Correct   int
	Supported int
}

const correctnessSystemPrompt = `You judge whether a generated answer is semantically correct. You see only the question, the reference answer, and the generated answer. Judge whether the generated answer states the same essential fact(s) as the reference answer; wording may differ. If the generated answer is exactly "Insufficient source evidence." it is correct only when the reference answer also indicates the question cannot be answered. Return ONLY JSON with this exact schema: {"verdict":"correct"} or {"verdict":"incorrect"}.`

const supportSystemPrompt = `You judge whether the quoted source passages directly support a generated answer. You see only the generated answer and the passages it cited. "Supported" means the passages themselves establish the answer without outside knowledge or inference beyond simple reading. Return ONLY JSON with this exact schema: {"verdict":"supported"} or {"verdict":"unsupported"}.`

type judgeVerdict struct {
	Verdict string `json:"verdict"`
}

// Judge runs blind correctness and source-support screens for every stored
// answer in the supplied experiment runs, writing each verdict into the
// answers.usage JSONB under the configured keys.
func Judge(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, cfg JudgeConfig) (JudgeMetrics, error) {
	metrics := JudgeMetrics{}
	if len(cfg.RunIDs) == 0 {
		return metrics, fmt.Errorf("at least one run ID is required")
	}
	if strings.TrimSpace(cfg.JudgeModel) == "" {
		return metrics, fmt.Errorf("--judge-model is required")
	}
	if strings.TrimSpace(cfg.CorrectnessKey) == "" || strings.TrimSpace(cfg.SupportKey) == "" {
		return metrics, fmt.Errorf("--correctness-key and --support-key are required")
	}
	rows, err := pool.Query(ctx, `SELECT a.id, q.text, coalesce(q.reference_answer,''), a.text, a.citations,
			(a.usage ? $2) AND (a.usage ? $3) AS already
		FROM answers a JOIN questions q ON q.id = a.question_id
		WHERE a.experiment_run_id = ANY($1::uuid[]) ORDER BY a.created_at`,
		cfg.RunIDs, cfg.CorrectnessKey, cfg.SupportKey)
	if err != nil {
		return metrics, err
	}
	type pending struct {
		id        string
		question  string
		reference string
		answer    string
		citations []byte
	}
	var work []pending
	for rows.Next() {
		var item pending
		var already bool
		if err = rows.Scan(&item.id, &item.question, &item.reference, &item.answer, &item.citations, &already); err != nil {
			rows.Close()
			return metrics, err
		}
		metrics.Answers++
		if already && cfg.SkipExisting {
			metrics.Skipped++
			continue
		}
		work = append(work, item)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return metrics, err
	}
	for _, item := range work {
		correctness, err := judgeOnce(ctx, client, cfg.JudgeModel, correctnessSystemPrompt,
			fmt.Sprintf("Question: %s\n\nReference answer: %s\n\nGenerated answer: %s", item.question, item.reference, item.answer),
			"correct", "incorrect")
		if err != nil {
			return metrics, fmt.Errorf("correctness judgment for answer %s: %w", item.id, err)
		}
		supportUser, hasCitations, err := supportPrompt(item.answer, item.citations)
		if err != nil {
			return metrics, fmt.Errorf("citations for answer %s: %w", item.id, err)
		}
		support := "unsupported"
		if hasCitations {
			support, err = judgeOnce(ctx, client, cfg.JudgeModel, supportSystemPrompt, supportUser, "supported", "unsupported")
			if err != nil {
				return metrics, fmt.Errorf("support judgment for answer %s: %w", item.id, err)
			}
		}
		payload, _ := json.Marshal(map[string]any{
			cfg.CorrectnessKey: map[string]any{"verdict": correctness, "judge_model": cfg.JudgeModel, "judged_at": time.Now().UTC().Format(time.RFC3339)},
			cfg.SupportKey:     map[string]any{"verdict": support, "judge_model": cfg.JudgeModel, "judged_at": time.Now().UTC().Format(time.RFC3339), "cited_sources_seen": hasCitations},
		})
		if _, err = pool.Exec(ctx, `UPDATE answers SET usage = usage || $2::jsonb WHERE id = $1`, item.id, payload); err != nil {
			return metrics, err
		}
		metrics.Judged++
		if correctness == "correct" {
			metrics.Correct++
		}
		if support == "supported" {
			metrics.Supported++
		}
	}
	return metrics, nil
}

func judgeOnce(ctx context.Context, client *fireworks.Client, model, system, user, positive, negative string) (string, error) {
	result, err := client.ChatJSONNoReasoning(ctx, model, system, user, 200)
	if err != nil {
		return "", err
	}
	var parsed judgeVerdict
	if err = json.Unmarshal([]byte(result.Text), &parsed); err != nil {
		return "", fmt.Errorf("judge returned invalid structured verdict %q: %w", result.Text, err)
	}
	verdict := strings.ToLower(strings.TrimSpace(parsed.Verdict))
	if verdict != positive && verdict != negative {
		return "", fmt.Errorf("judge returned unknown verdict %q", parsed.Verdict)
	}
	return verdict, nil
}

func supportPrompt(answer string, rawCitations []byte) (string, bool, error) {
	var stored storedCitations
	if err := json.Unmarshal(rawCitations, &stored); err != nil {
		return "", false, err
	}
	if len(stored.CitedSourceIDs) == 0 {
		return "", false, nil
	}
	byID := make(map[int]storedSource, len(stored.Retrieved))
	for _, source := range stored.Retrieved {
		byID[source.ID] = source
	}
	var builder strings.Builder
	for _, id := range stored.CitedSourceIDs {
		source, ok := byID[id]
		if !ok {
			return "", false, fmt.Errorf("cited source S%d absent from stored retrieval", id)
		}
		fmt.Fprintf(&builder, "[S%d] %s\n", id, source.Text)
		if source.Quote != "" {
			fmt.Fprintf(&builder, "Quote: %s\n", source.Quote)
		}
	}
	return fmt.Sprintf("Generated answer: %s\n\nCited passages:\n%s", answer, builder.String()), true, nil
}
