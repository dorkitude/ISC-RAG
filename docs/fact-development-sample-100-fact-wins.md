# Where facts beat fixed chunks in the initial frozen-100 probe

**Status:** diagnostic analysis of the initial 20-question probe, not a paper
result.  This note identifies the places where the fact condition is genuinely
useful, as well as one false-positive "win" produced by the current
turn-overlap evidence metric.

## Scope and comparison

This analysis uses only the four run IDs recorded in
[`fact-development-sample-100-initial-probe.md`](fact-development-sample-100-initial-probe.md):

| Model | Fact run | Fixed-chunk run |
| --- | --- | --- |
| DeepSeek V4 Flash | `4016c648-a2ad-44c7-9af9-4deae00cb902` | `993b131e-f7ce-4f0a-8a49-72af9cb65be2` |
| Kimi K2.6 | `ba81fbd3-8a3d-46ce-bdfa-b225737b739f` | `cb54fc99-81bb-4dc7-8162-f4339fd6f8c5` |

Each comparison holds model, question, and fresh-session replicate constant.
Both conditions use `top_k=1`.  A “fact-only” metric means the fact attempt is
true and its matched fixed-chunk attempt is false.

Across the 60 matched pairs per model:

| Model | Fact-only evidence | Chunk-only evidence | Fact-only citation | Chunk-only citation | Fact-only answer match | Chunk-only answer match |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Flash | 15 | 15 | 2 | 14 | 3 | 0 |
| Kimi | 9 | 18 | 5 | 23 | 3 | 0 |

The aggregate does **not** say facts win.  It does show a small, repeatable
set of question shapes where the fact projection supplies a much cleaner
answering unit.

## Full-100 deterministic retrieval check

After the initial agent probe, the project ran a deterministic `top_k=5`
retrieval check over all 100 frozen questions, with no answer-generation
model.  The matched runs are facts
`79843416-81b0-4191-9fa3-5cb311a221a2` and fixed chunks
`a9dec815-06d5-4523-b70f-8861fe538854`.

Facts hit a gold turn for **95/100** questions and fixed chunks for **87/100**.
Twelve questions are fact-hit/fixed-miss (four fixed-hit/fact-miss).  The
top-ranked fact itself contains the answer or resolved entity in 11 of these
12 cases; `fact-dev100-037` is the exception, because its retrieved fact is
from the gold turn but does not explicitly resolve “he” to President Trump.

| ID | Type | Top retrieved fact in the fact-only hit |
| --- | --- | --- |
| 022 | pronoun | Whitten expects **Takenaka** to solidify the upper-house-election mandates. |
| 023 | attribution | **Scott Jennings** believes the Court decision will bait Democrats into saying dumb things. |
| 024 | decontextualization | Miller says Cornyn is about **ten points** behind Trump in that part of Texas. |
| 027 | pronoun | Jenkin does not think the President’s **payroll-tax executive order** will help. |
| 037 | pronoun | Schiff says **Giuliani** smeared Yovanovitch. *(Gold-turn overlap only; insufficient for the asked antecedent.)* |
| 047 | pronoun | A.J. Hammer believes **Sharon** left because it felt like work. |
| 051 | direct | Abedi killed **22** people at the Ariana Grande concert. |
| 061 | direct | The investigation resulted in federal indictment of **38** people. |
| 062 | pronoun | Senator **Maria Chappelle Nadal** removed the tweet. |
| 067 | pronoun | **Misty May Treanor and Kerri Walsh** might retire on a high note. |
| 082 | pronoun | Hill later asked **Sondland** who put him in charge. |
| 096 | direct | Florida announced **nearly 4,000** new cases in 24 hours. |

This is stronger retrieval evidence for the intended hypothesis than the
20-question agent probe, but it is still only evidence retrieval.  It must be
followed by answer-generation trials and a fact-level evidence metric before
being treated as a result.

## Genuine, repeatable fact advantages

### 1. Speaker attribution inside a long exchange

| Question | Matched results | Why the fact helps |
| --- | --- | --- |
| `fact-dev100-033` — *Which guest said Madison Avenue would push legal marijuana to children like alcohol and tobacco?* | Flash: fact-only evidence in r1–r3. Kimi: fact-only evidence in r1–r3 and fact-only correct citation in r2–r3. | The fact result is “Madison Avenue will push marijuana to kids …” with quote “They'll push it to kids …” and `speaker_label=EVANS`; both models answer `EVANS`. The fixed result is the opening 4,000-character debate chunk, which names both guests but does not reach the later claim, so both models abstain. |
| `fact-dev100-048` — *Who said Islam tells people how not to dress?* | Flash: answer match fact-only in r1–r3. Kimi: answer match fact-only in r1–r3; fact citation is correct in r3 (the fixed citation is correct in all three). | The fact returns the claim plus `speaker_label=JEAN ABINADER, MANAGING DIRECTOR, ARAB AMERICAN INSTITUTE`. The fixed chunk contains Liz Boyes’s question immediately before Abinader’s response. Flash answers `Liz Boyes`; Kimi either answers Liz Boyes or calls the respondent unnamed. The fact answers Jean Abinader every time. |

This is the clearest current use case: an isolated claim with an attached,
normalized speaker label prevents a model from mistaking a nearby questioner
for the claimant.  The 048 result is a substantive answering improvement even
where the fixed chunk retrieved the gold turn.

### 2. Pronoun resolution when semantic chunk retrieval lands elsewhere

| Question | Matched results | Why the fact helps |
| --- | --- | --- |
| `fact-dev100-062` — *When Berman said “she ultimately removed that tweet,” who is “she”?* | Flash: fact-only evidence in r1–r3 and fact-only correct citation in r1 and r3. Kimi: fact-only evidence in r1–r3. | The fact is “Senator Maria Chappelle Nadal ultimately removed the tweet,” derived from the quote “she ultimately removed that tweet,” with the antecedent resolved. Both models answer Maria Chappelle Nadal. Fixed retrieval instead returns the programme’s unrelated closing exchange, so it abstains. |

The reference answer includes extra descriptive material (“Missouri Democratic
state senator … who had written …”), so the current literal answer-match
screen marks the correct short answer false.  This is an evaluator sensitivity,
not a fact failure.

### 3. Small, specific numeric claims buried in a broadcast

| Question | Matched results | Why the fact helps |
| --- | --- | --- |
| `fact-dev100-096` — *How many new coronavirus cases did Florida announce in the past 24 hours?* | Flash: fact-only evidence in r1–r3. Kimi: fact-only evidence and correct citation in r1–r3. | The fact is exactly “Florida announced nearly 4,000 new cases in the past 24 hours,” attributed to Keilar. Both models answer nearly 4,000. Fixed retrieval returns the beginning of the broadcast—Keilar’s introduction and unrelated COVID news—and both abstain. |

Again, the answer-match screen is overly literal here: “nearly 4,000 new
cases” and “nearly 4,000” are semantically identical to the reference but did
not receive a match.

### 4. Explicit opinion / judgment wording

| Question | Matched results | Why the fact helps |
| --- | --- | --- |
| `fact-dev100-005` — *What did Joe Biden think “Will & Grace” had done?* | Flash: fact-only evidence in r1–r3. | The fact preserves both the judgment and the content: “Joe Biden believes that ‘Will & Grace’ probably did more to educate the American public than almost anything …,” with Biden’s exact quote. Flash answers “Educate the American public.” Fixed retrieval selects an earlier chunk of the broadcast rather than Biden’s later remarks and abstains. |

This is evidence that the opinion-inclusive first pass can retain a subjective
claim in a self-contained form.  The literal answer screen nevertheless marks
the short, correct answer false because the reference is a full sentence.

## An apparent win that should *not* be counted as a fact advantage

`fact-dev100-034` asks how many other times Kentucky had a chance to have a
Senate majority leader.  Flash has fact-only `evidence_hit` in r1–r3, but the
retrieved fact merely says “Mitch McConnell is playing the Senate majority
leader card.”  Flash correctly responds “Insufficient retrieved evidence.”

The gold answer—“only one other time in history”—is in the same source turn,
but this fact does not contain it.  The current evidence screen regards any
retrieved source unit that overlaps the gold turn as a hit.  That is too coarse
for fact units: a turn may yield many facts.  Treat these three cells as
**turn-overlap false positives**, not wins.

## What the cases say about the mechanism

The real wins are not simply “facts are shorter.”  In these examples, facts
make one of three pieces of latent transcript structure explicit:

1. **Claimant identity** — speaker metadata is adjacent to the extracted claim
   (033, 048).
2. **Resolved reference** — an otherwise ambiguous pronoun is rewritten to its
   antecedent (062).
3. **Claim salience** — a compact proposition matches a specific question,
   whereas one 4,000-character chunk gets ranked for broad programme/topic
   vocabulary (005, 033, 096).

The facts also expose two evaluation defects before any larger experiment:

- evidence for facts should require a validated quote or fact-level gold link,
  not only turn overlap; and
- answer matching should allow a reference-answer substring / normalized
  semantic equivalence, since correct concise answers are currently discarded.

## Reproduce from PostgreSQL

Run from the project root after loading `.env`:

```sh
set -a; source .env; set +a
psql "$DATABASE_URL" -P pager=off -c '
WITH paired AS (
  SELECT fr.model, q.external_id, q.question_type, f.replicate,
         f.evidence_hit AS fact_evidence, x.evidence_hit AS fixed_evidence,
         f.citation_correct AS fact_citation, x.citation_correct AS fixed_citation,
         f.reference_answer_match AS fact_answer, x.reference_answer_match AS fixed_answer,
         f.answer_text AS fact_answer_text, x.answer_text AS fixed_answer_text
  FROM agent_evaluation_attempts f
  JOIN agent_evaluation_runs fr ON fr.id = f.run_id
  JOIN agent_evaluation_attempts x
    ON x.question_id = f.question_id AND x.replicate = f.replicate
  JOIN questions q ON q.id = f.question_id
  WHERE (f.run_id = $$4016c648-a2ad-44c7-9af9-4deae00cb902$$
         AND x.run_id = $$993b131e-f7ce-4f0a-8a49-72af9cb65be2$$)
     OR (f.run_id = $$ba81fbd3-8a3d-46ce-bdfa-b225737b739f$$
         AND x.run_id = $$cb54fc99-81bb-4dc7-8162-f4339fd6f8c5$$)
)
SELECT * FROM paired
WHERE (fact_evidence AND NOT fixed_evidence)
   OR (fact_citation AND NOT fixed_citation)
   OR (fact_answer AND NOT fixed_answer)
ORDER BY model, external_id, replicate;'
```

For the audited retrieval payload for any fact win (including fact text, quote,
speaker label, tool query, answer, and citation):

```sql
SELECT q.external_id, a.replicate, a.tool_request,
       a.raw_tool_output, a.answer_text, a.citation,
       a.evidence_hit, a.citation_correct, a.reference_answer_match
FROM agent_evaluation_attempts a
JOIN questions q ON q.id = a.question_id
WHERE a.run_id = '4016c648-a2ad-44c7-9af9-4deae00cb902'
  AND q.external_id IN
    ('fact-dev100-005', 'fact-dev100-033', 'fact-dev100-034',
     'fact-dev100-048', 'fact-dev100-062', 'fact-dev100-096')
ORDER BY q.external_id, a.replicate;
```
