# Why facts trail chunks: extraction and granularity audit

**Scope:** frozen `fact-dev-sample-100-20260725` development set only. This
is a diagnostic note, not a paper result. It tests two explanations for the
small 20-question / 240-session probe in
[`fact-development-sample-100-initial-probe.md`](fact-development-sample-100-initial-probe.md).

## Result in brief

The evidence does **not** support a simple claim that the high-recall extractor
failed to produce facts from the cited evidence. Every one of the 100 canonical
gold turn ranges has at least one `turnwise-opinion-v1` fact, and an exact
question-text retrieval run has higher turn-overlap recall for facts than for
matched fixed chunks at every tested depth:

| Units retrieved | facts | fixed chunks |
| --- | ---: | ---: |
| Recall@1 | 83/100 | 70/100 |
| Recall@3 | 92/100 | 82/100 |
| Recall@5 | 95/100 | 87/100 |

The agent probe's chunk advantage therefore cannot, by itself, be read as
evidence that facts are intrinsically worse at semantic retrieval. It is
consistent with an *agent/query-and-presentation* problem: agents normally ask
for one unit, and a one-fact result carries far less surrounding context than a
one-chunk result. A related measurement issue matters: the current automatic
"evidence hit" says only that the result overlaps a gold **turn**, not that the
returned fact itself expresses the answer. With many facts per turn, that can
call a wrong fact a hit.

## H1 — did claim extraction lose the gold evidence?

### Test A: existence in the canonical gold range

For every question, I joined its `gold_evidence` range to source units of each
condition. Facts cover all 100 questions, as do chunks. The extractor produces
an average of **13.7 facts** overlapping a gold range (minimum one) versus
**1.2 fixed chunks**. That rejects the narrow version of H1: no question lacks
an extracted fact from the source turn the question set cites.

| Question family | Questions with fact in gold turn | Mean facts in gold turn |
| --- | ---: | ---: |
| Attribution | 20/20 | 16.1 |
| Decontextualization | 20/20 | 15.3 |
| Direct control | 20/20 | 12.8 |
| Opinion / belief / judgment | 20/20 | 9.0 |
| Pronoun / coreference | 20/20 | 15.4 |
| **All** | **100/100** | **13.7** |

### Test B: direct semantic retrieval from the original question

I ran the same Fireworks embedding/search path on all 100 original question
texts, with no LLM query rewrite or answer model. Run IDs are:

| Condition | @1 | @3 | @5 |
| --- | --- | --- | --- |
| Facts (`turnwise-opinion-v1`) | `f32ae639-01e0-4ca4-9479-4900c92322ea` | `aa7f8765-4465-4a48-95e5-3887384f8e48` | `a7640b9c-3d52-49fc-b6db-240c0ddba9c6` |
| Fixed (`fact-dev-sample-100-fixed-v1`) | `ca2ba454-10da-48e2-92c5-2592ac8c2ed0` | `ad023981-2f7a-4426-bde6-ec5cc2863f59` | `6def7954-607e-4c4a-b16c-6e96f5193c75` |

Facts are especially better on coreference (18/20 at @1 versus chunks' 7/20;
20/20 at @5 versus 13/20). They are also better on direct controls (19/20
versus 15/20 at @1). At @1, facts-only successes outnumbered chunks-only
successes in three of five families:

| Question family | Facts only | Chunks only |
| --- | ---: | ---: |
| Attribution | 3 | 2 |
| Decontextualization | 3 | 4 |
| Direct control | 5 | 1 |
| Opinion / belief / judgment | 3 | 5 |
| Pronoun / coreference | 12 | 1 |

Representative fact-only **turn-overlap** @1 IDs: `fact-dev100-012`, `-017`,
`-022`, `-027`, `-047`, `-052`, `-057`, `-062`, `-067`, `-072`, `-077`, and
`-082` (coreference); `-016`, `-051`, `-061`, `-086`, and `-096` (direct
control). Several are clearly answer-bearing compact cards (for example,
`-016`, `-051`, `-061`, `-062`, and `-077`); do not assume that every listed
turn-overlap success is answer-bearing without fact-level adjudication.

### Important qualification: existence is not semantic sufficiency

The broad extractor can retain a related but under-specified claim. For
`fact-dev100-075`, the cited turn has the fact *“Harris states that it would
actually be an absolute bombshell if that were to be proven.”* It loses the
antecedent that the question names: compromise of Ptech software supplied to
government agencies. The original source turn also lacks that antecedent;
fixed chunks win because their wider window carries the prior setup. This is a
real form of context/claim loss, but it is not a failure to harvest the turn.

**H1 conclusion:** reject “the extractor broadly misses cited turns.” Retain a
narrower hypothesis for the next prompt study: self-contained rewriting must
carry the source antecedent or linked question/context when the target clause
contains `it`, `that`, a comparison, or an elliptical answer.

## H2 — are facts too fragmented or too context-poor for the agent?

### The unit-budget test

The two conditions have radically different information budgets:

| Condition | Units | Mean unit text | Median unit text | Mean turns/unit |
| --- | ---: | ---: | ---: | ---: |
| Facts | 13,064 | 73 chars | 68 chars | 1.0 |
| Fixed chunks | 283 | 3,212 chars | 3,608 chars | 11.5 |

In the actual fresh-agent runs, Flash requested exactly one fact on all 60
fact attempts, receiving 135 characters on average (fact plus quote), versus
2,620 characters from fixed chunks. Kimi requested 1.23 facts (162 characters)
versus 1.33 chunks (4,029 characters). Five facts would still be only roughly
365 fact-text characters before quotes. This is not a retrieval-depth matched
comparison: an agent selecting `top_k=1` gets one atomic proposition versus a
multi-turn local context window.

### The within-turn decoy test

Because there are 13.7 fact candidates in a typical gold range, turn-overlap
does not ensure that the answer-bearing candidate was returned. The clearest
case is `fact-dev100-034`. The correct fact exists in turn 3:

> Only one other time in history has there been a chance to have a majority leader from Kentucky.

All six fact-agent attempts instead retrieved the same *other* turn-3 fact:

> Mitch McConnell is playing the Senate majority leader card.

All six were automatically counted as gold-evidence hits, but both models
returned “Insufficient retrieved evidence.” The fixed chunk contains the full
turn and the number, so Kimi answered and cited it in all three repetitions.
This is retrieval competition among atomized facts, made invisible by a
turn-level overlap metric.

The same pattern appears in the agent probe summary:

| Model | Fact gold-turn hits | Hit but agent said insufficient | Chunk gold-turn hits | Hit but insufficient |
| --- | ---: | ---: | ---: | ---: |
| Flash | 33 | 8 | 33 | 0 |
| Kimi | 30 | 4 | 39 | 0 |

All six fact attempts for `fact-dev100-034` are in the first category. Other
examples are `fact-dev100-074` (three) and `fact-dev100-097` (three).

`fact-dev100-079` makes the second granularity failure concrete. The source
has the desired *“754 structures have been destroyed so far”* in turn 3, but
all six fact-agent calls retrieved the related turn-0 card *“The Thomas fire is
the state's 13th most destructive wildfire.”* Fixed chunks retrieve a window
containing both the topic and number. `fact-dev100-085` similarly has a useful
standalone opinion in turn 14, while the question's gold range is only turn 13;
the chunk spans both turns but a one-fact search selected the weaker turn-13
card. That last case also exposes a gold-range annotation defect and should be
re-adjudicated before using this set for scores.

**H2 conclusion:** strongly supported as an explanation for the *agent-probe*
gap, even though facts win the exact-query turn-overlap retrieval test. The
current fact condition needs either (a) a context-expanded fact result (adjacent
facts/turns), (b) substantially higher result count with a matched token
budget, or (c) a reranker that chooses an answer-bearing fact rather than a
topically related fact from the same source turn. Evaluation should add
fact-level answer-support judgment; raw turn overlap is an upper-bound
retrieval screen, not answer evidence.

## Reproduce

Load project credentials only locally, then run the six direct-query retrieval
conditions:

```sh
set -a; source .env; set +a
for top_k in 1 3 5; do
  go run ./cmd/isc-rag evaluate run \
    --representation facts --version turnwise-opinion-v1 \
    --question-split fact-dev-sample-100-20260725 --corpus-split test \
    --top-k "$top_k" --embedding-batch-size 32
  go run ./cmd/isc-rag evaluate run \
    --representation fixed --version fact-dev-sample-100-fixed-v1 \
    --question-split fact-dev-sample-100-20260725 --corpus-split test \
    --top-k "$top_k" --embedding-batch-size 32
done
```

Gold-turn unit coverage:

```sql
WITH q AS (
  SELECT q.id, q.question_type
  FROM questions q
  WHERE q.corpus_split = 'fact-dev-sample-100-20260725'
)
SELECT q.question_type,
       count(*) FILTER (WHERE EXISTS (
         SELECT 1 FROM gold_evidence ge JOIN source_units su
           ON su.document_id = ge.document_id
          AND su.representation = 'facts' AND su.version = 'turnwise-opinion-v1'
          AND su.turn_start <= ge.turn_end AND su.turn_end >= ge.turn_start
         WHERE ge.question_id = q.id
       )) AS questions_with_fact_overlap
FROM q GROUP BY q.question_type;
```

The first probe's gold-turn-hit-but-insufficient diagnostic:

```sql
SELECT r.model, r.capability,
       count(*) FILTER (WHERE a.evidence_hit) AS turn_hits,
       count(*) FILTER (WHERE a.evidence_hit
                          AND a.answer_text = 'Insufficient retrieved evidence.')
         AS hit_but_insufficient
FROM agent_evaluation_attempts a
JOIN agent_evaluation_runs r ON r.id = a.run_id
WHERE r.question_split = 'fact-dev-sample-100-20260725'
GROUP BY r.model, r.capability;
```
