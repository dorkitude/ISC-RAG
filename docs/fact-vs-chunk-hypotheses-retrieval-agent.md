# Fact vs. fixed-chunk diagnostic: retrieval, context, and agent behavior

**Status:** development diagnosis only. This is not a paper result and it does
not replace source-level answer adjudication.

## Scope and the important measurement distinction

Two evaluations answer different questions:

1. The initial agent probe contains 20 questions, three fresh sessions per
   model/condition. The Q&A model first rewrites/selects a search query and
   almost always asks for `top_k=1`. The answer model then sees that one result.
   Its runs are `4016c648-a2ad-44c7-9af9-4deae00cb902`,
   `993b131e-f7ce-4f0a-8a49-72af9cb65be2`,
   `ba81fbd3-8a3d-46ce-bdfa-b225737b739f`, and
   `cb54fc99-81bb-4dc7-8162-f4339fd6f8c5`.
2. The deterministic full-100 retrieval check embeds the original question and
   retrieves five units, without an answer agent. Its fact run is
   `79843416-81b0-4191-9fa3-5cb311a221a2`; its matched fixed-chunk run is
   `a9dec815-06d5-4523-b70f-8861fe538854`.

The CLI evaluator calls a hit when a retrieved unit overlaps a gold turn. The
agent harness currently calls a hit only when the unit *contains* the full gold
span. This is harmless for a one-turn gold label, but can differ for a
multi-turn label or an atomic fact attached to one of several gold turns. Do
not compare their percentages as though they were the same metric.

## Where facts are already better

On the full frozen-100 deterministic top-5 retrieval check, facts retrieve
gold-overlapping evidence for **95/100** questions and fixed chunks for
**87/100**. The fact advantage is concentrated, rather than universal:

| Question type | Facts | Fixed chunks | Fact-only wins | Chunk-only wins |
| --- | ---: | ---: | ---: | ---: |
| Pronoun/coreference | 20/20 | 13/20 | 7 | 0 |
| Direct control | 20/20 | 17/20 | 3 | 0 |
| Attribution | 19/20 | 19/20 | 1 | 1 |
| Decontextualization | 19/20 | 19/20 | 1 | 1 |
| Opinion/belief/judgment | 17/20 | 19/20 | 0 | 2 |
| **All** | **95/100** | **87/100** | **12** | **4** |

Representative fact-only retrieval wins:

- `fact-dev100-022` (pronoun): the top fact rewrites the reference as
  “Whitten expects **Takenaka** to solidify the mandates…”. None of the top
  five fixed chunks is from the source transcript.
- `fact-dev100-051` (direct): the top fact says “**Abedi killed 22 people** at
  the Ariana Grande concert…”. The fixed index instead retrieves chunks about
  the Las Vegas massacre and unrelated deaths.
- `fact-dev100-062` (pronoun): the top fact says “Senator **Maria Chappelle
  Nadal** ultimately removed the tweet.” The five fixed results are from other
  transcripts.
- `fact-dev100-096` (direct): the top fact states “Florida announced **nearly
  4,000 new cases** in the past 24 hours.” The chunk condition retrieves an
  unrelated later section of the correct document and other coronavirus
  broadcasts.

The smaller initial agent probe also contains stable fact wins. Both models
retrieve the fact (and neither retrieves the chunk) for the Madison Avenue /
marijuana attribution question (`fact-dev100-033`), the Maria Chappelle Nadal
pronoun question (`fact-dev100-062`), and the Florida case-count question
(`fact-dev100-096`). Flash additionally obtains all three fact hits for the
Joe Biden / *Will & Grace* opinion question while fixed chunks obtain none.

These are the intended strengths of self-contained facts: the query matches a
short, named assertion, while a long chunk's embedding is influenced by many
unrelated turns.

## Hypothesis 3 — query-to-fact embedding mismatch

**Hypothesis:** question wording or the answer agent's query rewrite is a
poor semantic match for the fact projection, causing fact retrieval to lose.

**Finding: partly plausible, but not the main explanation for the initial
agent deficit.** In the initial 20-question agent probe, Kimi issued the
*identical* selected query for fact and fixed conditions in 43/60 paired fresh
sessions; within those identical-query pairs fixed chunks won 15 evidence
comparisons and facts won 7. Flash used the identical query in 26/60 pairs;
fixed won 6 and facts won 3. Thus a different query rewrite cannot explain the
bulk of the fixed-chunk wins.

It still matters at the representation level. Facts are a much more precise
query target when the question names an entity and claim (as in the examples
above), while the same embedding sometimes fails to rank a useful opinion fact
highly enough. The full-100 result supports this mixed story: facts are ahead
on direct and pronoun questions but behind on opinion questions (17 vs. 19).

**Implication:** retain the fact-plus-quote projection, but compare an
original-question retrieval condition against an agent-query-rewrite condition
explicitly. Do not attribute the current agent result solely to query rewrite.

## Hypothesis 4 — top-k, ranking, or context-budget dilution

**Hypothesis:** the fact condition loses because too many small facts dilute
the ranking/context budget, whereas chunks retain enough surrounding context.

**Finding: not supported as the cause in this probe.** The agent instruction
actually biases the model to `top_k=1`:

| Model | Condition | Mean requested k | Gold hits at rank 1 | Gold hits below rank 1 |
| --- | --- | ---: | ---: | ---: |
| Flash | facts | 1.00 | 33 | 0 |
| Flash | fixed | 0.95 (3 invalid) | 33 | 0 |
| Kimi | facts | 1.23 | 30 | 0 |
| Kimi | fixed | 1.33 | 36 | 3 |

The answer agent receives much *less* context for facts, not more: mean tool
JSON is 626 bytes (Flash) / 726 bytes (Kimi) for facts, versus 3,095 / 4,626
bytes for fixed chunks. Yet fixed chunks lead Kimi in the small agent probe.
There is therefore no evidence of fact-context dilution; if anything, the
top-1 policy gives facts their cleanest possible context. Allowing Kimi more
than one result helped chunks for three additional hits, not facts.

The remaining ranking explanation is finer-grained: one atom can lose its
semantic rank even when a neighboring chunk contains enough related language
to rank. This is a *claim-ranking/fact-selection* problem, not a top-k context
budget problem. The next controlled ablation should evaluate recall@1, 3, and
5 using the same original query and a fixed context-token budget (for example,
pack the highest-ranked facts until the chunk budget is reached).

## Hypothesis 5 — answer-agent and citation behavior create an apparent gap

**Hypothesis:** the facts are retrieved, but the answer agent handles their
small output or citation syntax differently, depressing apparent performance.

**Finding: definitely affects the citation screen; it cannot explain the
retrieval screen.** The evidence-hit flag is computed before the answer model,
so the Kimi fixed-versus-fact retrieval difference (39 vs. 30 in this small
probe) is real under the harness's containment criterion. But the reported
citation gap is substantially confounded by an overly strict parser.

The harness accepts only `external_id turns start-end`. Both models frequently
produce an otherwise correct citation such as `CNN-403178 3-3` or
`CNN-139486 turn 17-17`. Counting those as valid by a relaxed document/range
parser changes the probe as follows:

| Model | Facts strict → relaxed citations | Fixed strict → relaxed citations |
| --- | ---: | ---: |
| Flash | 3/60 → 22/60 | 15/60 → 30/60 |
| Kimi | 20/60 → 26/60 | 38/60 → 36/60 |

Flash's apparent fact citation score is therefore mostly formatting failure,
not missing evidence. Facts remain behind Kimi's fixed chunks on relaxed
citations, but the meaningful comparison should use a structured citation
object or a tolerant host parser—not prompt-dependent free text. The current
literal-reference-answer screen is also unsuitable as answer accuracy: it
matches only 3 evidence-hit answers per condition/model in this probe despite
many plainly useful answers.

**Implication:** make citations structured (document ID, start, end) and
validate them host-side. Human or rubric-based answer grading is needed before
claiming an answer-quality advantage or disadvantage.

## Reproduce the decisive full-100 slice

From the project directory, after loading `.env`:

```sh
set -a; source .env; set +a
psql "$DATABASE_URL" -P pager=off
```

Then run this query. It uses the evaluator's intended **overlap** condition:

```sql
WITH x AS (
  SELECT e.representation, q.external_id, q.question_type,
    EXISTS (
      SELECT 1
      FROM retrievals rr
      JOIN gold_evidence ge ON ge.question_id = rr.question_id
      JOIN source_units su ON su.id = rr.source_unit_id
      WHERE rr.experiment_run_id = e.id
        AND rr.question_id = q.id
        AND su.document_id = ge.document_id
        AND su.turn_start <= ge.turn_end
        AND su.turn_end >= ge.turn_start
    ) AS hit
  FROM experiment_runs e CROSS JOIN questions q
  WHERE e.id IN ('79843416-81b0-4191-9fa3-5cb311a221a2',
                 'a9dec815-06d5-4523-b70f-8861fe538854')
    AND q.corpus_split = 'fact-dev-sample-100-20260725'
), paired AS (
  SELECT external_id, max(question_type) AS question_type,
    max(hit::int) FILTER (WHERE representation = 'facts') AS fact_hit,
    max(hit::int) FILTER (WHERE representation = 'fixed') AS fixed_hit
  FROM x GROUP BY external_id
)
SELECT question_type, count(*) AS questions, sum(fact_hit) AS fact_hits,
       sum(fixed_hit) AS fixed_hits,
       sum((fact_hit = 1 AND fixed_hit = 0)::int) AS fact_only,
       sum((fact_hit = 0 AND fixed_hit = 1)::int) AS fixed_only
FROM paired
GROUP BY ROLLUP(question_type)
ORDER BY question_type NULLS LAST;
```

## Bottom line

“Fixed chunks outperform” is too broad. In the agent probe, fixed chunks are
better for Kimi's top-1, answer-agent-mediated evaluation, especially on the
five sampled opinion questions. In direct deterministic retrieval over all 100
questions, facts are already better overall and substantially better on
coreference. The next experiment should hold query, top-k/context budget, and
citation representation fixed, then isolate fact claim coverage and ranking on
the opinion slice.
