package evaluation

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GenerateDeterministicPilotQuestions builds a fast lexical retrieval stress
// set. Each probe contains a source phrase and has an exact gold turn, so it is
// useful for plumbing and storage comparisons but is intentionally not a
// natural-question or paper-gold benchmark.
func GenerateDeterministicPilotQuestions(ctx context.Context, pool *pgxpool.Pool, questionSplit, corpusSplit string, documents, perDocument int) (PilotQuestionReport, error) {
	if questionSplit == "" || corpusSplit == "" || documents <= 0 || perDocument <= 0 {
		return PilotQuestionReport{}, fmt.Errorf("question split, corpus split, documents, and per-document must be set")
	}
	rows, err := pool.Query(ctx, `SELECT id,external_id FROM documents WHERE corpus='mediasum' AND corpus_split=$1 ORDER BY external_id LIMIT $2`, corpusSplit, documents)
	if err != nil {
		return PilotQuestionReport{}, err
	}
	defer rows.Close()
	var docs []pilotDocument
	for rows.Next() {
		var doc pilotDocument
		if err := rows.Scan(&doc.id, &doc.externalID); err != nil {
			return PilotQuestionReport{}, err
		}
		docs = append(docs, doc)
	}
	if err := rows.Err(); err != nil {
		return PilotQuestionReport{}, err
	}
	if len(docs) != documents {
		return PilotQuestionReport{}, fmt.Errorf("requested %d documents but found %d", documents, len(docs))
	}
	report := PilotQuestionReport{}
	for _, doc := range docs {
		turns, err := pilotDocumentTurns(ctx, pool, doc.id)
		if err != nil {
			return report, err
		}
		items := deterministicItems(turns, perDocument)
		if len(items) != perDocument {
			return report, fmt.Errorf("%s supplied only %d/%d deterministic lexical probes", doc.externalID, len(items), perDocument)
		}
		if err = persistDeterministicPilotQuestions(ctx, pool, doc, items, questionSplit); err != nil {
			return report, err
		}
		report.Documents++
		report.Questions += len(items)
	}
	return report, nil
}

func deterministicItems(turns []pilotTurn, required int) []pilotItem {
	seen := map[string]bool{}
	var items []pilotItem
	for width := 3; width <= 8 && len(items) < required; width++ {
		for _, turn := range turns {
			words := strings.Fields(turn.Text)
			for start := 0; start+width <= len(words) && len(items) < required; start++ {
				phrase := strings.Join(words[start:start+width], " ")
				question := fmt.Sprintf("What did %s say about %q?", deterministicSpeaker(turn.Speaker), phrase)
				key := strings.ToLower(question)
				if seen[key] {
					continue
				}
				seen[key] = true
				items = append(items, pilotItem{Question: question, Answer: phrase, Turn: turn.Ordinal})
			}
		}
	}
	return items
}

func deterministicSpeaker(speaker string) string {
	speaker = strings.Join(strings.Fields(speaker), " ")
	if speaker == "" {
		return "the speaker"
	}
	return speaker
}

func persistDeterministicPilotQuestions(ctx context.Context, pool *pgxpool.Pool, document pilotDocument, items []pilotItem, questionSplit string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for index, item := range items {
		externalID := fmt.Sprintf("%s-%s-%03d", questionSplit, document.externalID, index+1)
		if _, err = tx.Exec(ctx, `INSERT INTO questions(external_id,corpus_split,question_type,text,reference_answer,metadata) VALUES ($1,$2,'deterministic_lexical_probe',$3,$4,jsonb_build_object('generator_model','deterministic','document_external_id',$5::text,'protocol','source-phrase lexical retrieval stress test')) ON CONFLICT(external_id) DO UPDATE SET corpus_split=EXCLUDED.corpus_split,question_type=EXCLUDED.question_type,text=EXCLUDED.text,reference_answer=EXCLUDED.reference_answer,metadata=EXCLUDED.metadata`, externalID, questionSplit, item.Question, item.Answer, document.externalID); err != nil {
			return err
		}
		var questionID string
		if err = tx.QueryRow(ctx, `SELECT id FROM questions WHERE external_id=$1`, externalID).Scan(&questionID); err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, `DELETE FROM gold_evidence WHERE question_id=$1`, questionID); err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO gold_evidence(question_id,document_id,turn_start,turn_end) VALUES($1,$2,$3,$3)`, questionID, document.id, item.Turn); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
