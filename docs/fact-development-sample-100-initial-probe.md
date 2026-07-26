# Frozen 100-document fact-development probe

**Status:** exploratory development evidence, not a paper result.

## Frozen inputs

- Documents: the 100 deterministic MediaSum IDs in
  `examples/fact-development-sample-100-20260725.txt`.
- Questions: 100 source-grounded questions in
  `experiments/fact-development-sample-100-20260725/questions.jsonl`. Claude
  Fable generated each one in an independent `claude -p --model fable` session.
  The raw outputs are ignored; all 100 promoted rows were mechanically checked
  against PostgreSQL for an existing document, cited turn range, and exact
  supporting quote. See the adjacent `authoring-notes.md`.
- Fact condition: `facts/turnwise-opinion-v1`, 13,064 units from a high-recall,
  per-turn harvest of explicit assertions and explicit opinions/beliefs/
  judgments, with exact target-turn quote validation. The retrieval projection
  embeds `Fact: <fact>\n\nSupporting quote: <quote>`.
- Matched control: `fixed/fact-dev-sample-100-fixed-v1`, 283 4,000-character
  chunks with 400-character overlap, built from the same documents.

## Initial fresh-session Q&A probe

This is a small integration and diagnostic run, deliberately limited to the
first 20 questions in the frozen JSONL (3 questions from direct control, 6
pronoun/coreference, 2 attribution, 4 decontextualization, and 5 opinion).
Each question-condition-model cell has three fresh sessions. The agent gets one
compiled, capability-limited search binary, submits one JSON search request,
then answers from the retrieved tool output. The full prompts, completions,
retrieval ranks, automatic screens, and latencies are in PostgreSQL.

| Model | Condition | Sessions | Valid tool calls | Gold-evidence hit | Correct citation | Literal answer match | Mean E2E |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: |
| DeepSeek V4 Flash | fact + quote | 60 | 60 | 33 (55.0%) | 3 (5.0%) | 6 (10.0%) | 2.68 s |
| DeepSeek V4 Flash | fixed chunk | 60 | 57 | 33 (55.0%) | 15 (25.0%) | 3 (5.0%) | 6.34 s |
| Kimi K2.6 | fact + quote | 60 | 60 | 30 (50.0%) | 20 (33.3%) | 6 (10.0%) | 1.79 s |
| Kimi K2.6 | fixed chunk | 60 | 60 | 39 (65.0%) | 38 (63.3%) | 3 (5.0%) | 2.25 s |

Run IDs: `4016c648-a2ad-44c7-9af9-4deae00cb902`,
`993b131e-f7ce-4f0a-8a49-72af9cb65be2`,
`ba81fbd3-8a3d-46ce-bdfa-b225737b739f`, and
`cb54fc99-81bb-4dc7-8162-f4339fd6f8c5`.

### Reading this responsibly

This probe is not evidence that the fact method wins. The high-recall
extractor now covers the entire frozen sample, and facts are faster to present
to the agent, but the matched fixed chunks lead Kimi on evidence and citation
coverage. Flash ties the two conditions on evidence hit, but cites fixed chunks
more reliably. The useful next diagnostic is not a larger paper run: inspect
the fact retrieval misses by question type and determine whether they are
claim-extraction loss, query-to-fact mismatch, overfragmentation, or answer
agent citation behavior. Only then should the 100-question run be expanded.

## Follow-up: direct retrieval over all 100 questions

The agent probe should not be confused with direct retrieval quality: the agent
first rewrites a question into one search query and receives only its selected
tool result. To separate that behavior from the representation, we also embedded
each original question text and retrieved top-5 units without answer generation.
The resulting immutable retrieval runs are `79843416-81b0-4191-9fa3-5cb311a221a2`
(facts) and `a9dec815-06d5-4523-b70f-8861fe538854` (fixed).

| Question type | Facts, gold hit@5 | Fixed chunks, gold hit@5 | Facts only | Fixed only |
| --- | ---: | ---: | ---: | ---: |
| Attribution | 19/20 | 19/20 | 1 | 1 |
| Decontextualization | 19/20 | 19/20 | 1 | 1 |
| Direct control | 20/20 | 17/20 | 3 | 0 |
| Opinion/belief/judgment | 17/20 | 19/20 | 0 | 2 |
| Pronoun/coreference | 20/20 | 13/20 | 7 | 0 |
| **Overall** | **95/100** | **87/100** | **12** | **4** |

Thus facts are currently strongest for coreference and compact direct lookups.
Concrete fact-only wins include the antecedent of Berman's “she ultimately
removed that tweet” (Maria Chappelle Nadal), Florida's “nearly 4,000” new cases,
and the `Will & Grace` claim that Joe Biden said it educated the American
public. The agent result therefore exposes a second-stage problem—query choice,
top-1 retrieval, and citation/answer behavior—not a general failure of facts to
retrieve the gold turn.
