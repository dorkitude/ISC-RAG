# production-inspired fact study #2

**Status:** active, isolated document-disjoint development study. This report
records every attempted configuration through the current iteration; it does not
promote a primary ISC-RAG fact extractor unless the predeclared gate is met.
Raw MediaSum-derived provider completions, fact text, and evaluator worksheets
are local ignored artifacts/database rows, not tracked here.

## Question and predeclared gate

Can production-inspired transcript context construction and a dialogue-neutral
extraction contract produce provenance-backed, standalone dialogue fact cards
that meet all of the following on each of three fixed document-level
replicates, with fixed model, prompt, cap, chunk policy, and validation policy?

- support at least 98%, with zero unsupported cards in every replicate;
- speaker/attribution at least 98%; qualifier/negation/temporal preservation at
  least 95%; standalone clarity and atomicity at least 90%; and duplicate or
  near-duplicate rate at most 5%; and
- source-ledger coverage at least 75% on every document.

The test is deliberately conservative: exact substring location is not treated
as semantic entailment. A card passes support only when its *stored* evidence
bundle supports the complete rephrased text.

## Frozen development material and audit method

The three documents are MediaSum **validation** records `NPR-1781`, `NPR-2291`,
and `NPR-967`. They are document-disjoint from `NPR-847` (Study #1) and from
the MediaSum test/canonical-question set. Before prompt tuning, an independent
source-only claim ledger was written at
`experiments/fact-study-2/claim-ledger.md`: 20, 18, and 20
salient explicit atomic claims respectively (58 total). It cites only exact
source turns; MediaSum summaries are prohibited.

Every candidate receives a fresh blind evaluator. The evaluator sees the fixed
ledger, source turns, and only the candidate's current facts/evidence—not prior
scores. It evaluates every retained card for: (1) support/entailment, (2)
speaker/attribution, (3) qualifier/negation/temporal preservation, (4)
standalone clarity, (5) atomicity, (6) exact/near duplicates, and (7) claim-ledger
coverage. Local worksheets live below
`artifacts/fact-study-2/evaluations/`.

All runs use Fireworks
`accounts/fireworks/models/deepseek-v4-pro`, temperature zero, no reasoning,
and JSON-object output. This account did not expose a V4 Flash deployment; the
CLI retains `--model` for an explicit future Flash comparison. Complete effective
instructions, SHA-256 hashes, model/configuration, token counts, latency, raw
completion hash, and raw output are retained in isolated local PostgreSQL tables
`study_two_*`; no study uses the primary `facts` table.

## C4 structural cross-study replication (frozen before run)

`fable-c4-sbar` is a separate, cross-study structural replication of the
strongest natural-language antecedent-window condition in Fable Study #3. It
is not a ledger-targeted repair condition: the Study #2 ledger is withheld from
generation, mechanical validation, the span-blind auditor, and repair. It is
used only by the final blind evaluator.

The frozen pipeline is: (S1) one full-transcript generator emits a standalone,
atomic fact, one primary quote, an optional antecedent quote, and a cited primary
turn; (S2) exact mechanical validation requires the primary quote in that turn
and the antecedent in that turn or one of four previous turns; (S3) a fresh
span-blind auditor sees only fact text, the two quote spans, and their canonical
speaker labels and marks every card PASS or FAIL with exact ungrounded elements;
(S4) only failed cards go to a fresh repairer with transcript access; its sole
legal moves are `REANCHOR`, `RECUT`, quoted-`NAME_FORM`, or `DROP`; (S4b) every
repair is mechanically revalidated; (S5) a fresh span-blind auditor audits every
repaired card; and (S6) deterministic assembly keeps only cards passing S3 or
S5. The primary turn is frozen during repair. Reanchoring cannot change fact
text or primary quote, recutting cannot change fact text or the antecedent, and
name-form rewriting cannot alter either quote. This bounds repair blast radius
and prevents hidden predicate or qualifier rewrites.

All prompts, stage outputs, hashes, and token counts are persisted in the
isolated study tables. Raw dialogue-derived outputs remain ignored local
artifacts/database rows. A configuration cannot be called a paper condition
until all three fixed documents and subsequently two unchanged-config repeat
runs per document pass the predeclared gate under fresh final evaluation.

## Dialogue-act Q/A structural jump (frozen before run)

`dialogue-act-pairs` changes the source representation rather than adding a
stricter rephrasing rule. Before generation, a question turn is paired only
with its immediately following turn when that turn is declarative and spoken by
a different speaker. That exact Q/A pair becomes one source unit; every
declarative turn not consumed as an answer remains a standalone unit. The
pairing rule, source speaker, and unit ordering are deterministic and retained
with the run.

The generator emits atomic cards with a `unit_id`, an exact answer quote, and
an optional exact question quote. A paired answer must be phrased as an
attributed dialogue act when the question supplies essential referent or scope:
the answer speaker answered, explained, estimated, believed, or recalled the
claim. The question is context, not a newly asserted fact. Mechanical
validation rejects a quote outside its own answer/question turn, rejects a
question quote for a standalone unit, deterministically attaches the answer
speaker, and stores both exact Q and A spans. The claim ledger remains
evaluation-only. Every fixed document receives a fresh blind audit; no repeats
are allowed unless all first-pass document gates clear.

## Dialogue-act predicate-unit follow-up (frozen before run)

`dialogue-act-predicate-units` retains the identical deterministic Q/A pairing
and exact Q/A evidence contract, but changes the answer-side representation.
Each answer turn is deterministically split into lexical sentence, semicolon,
and `, and` predicate/coordination units before generation. The model sees one
such answer claim unit per source unit and may emit at most one card for it; it
cannot merge two answer units. Mechanical validation accepts an answer quote
only when it is an exact substring of that unit (not merely elsewhere in the
answer turn), while retaining the paired question quote where needed. This is a
new source-representation condition, not a post-hoc score-targeted rewrite.

## Quote-capsule structural condition (frozen before run)

`quote-capsules` removes free-form fact wording altogether. Given the same
deterministic Q/A and standalone source units, the model may output only a unit
ID and zero-based UTF-8 byte ranges for one atomic answer span and an optional
paired-question context span. The validator rejects out-of-range or non-rune
boundaries, question-like/acknowledgement spans, duplicate capsules, and a
question span outside its own paired unit. It then renders the stored fact
deterministically as an attributed capsule—for example, that the answer speaker
answered a quoted question with a quoted exact answer span. The question is
always context rather than an asserted premise. One answer span is one card;
the selector must split distinct propositions into distinct spans.

`quote-capsules-exact` is the boundary-selector fallback after the provider
returned placeholder numeric ranges. It retains the no-free-text-fact design:
the model returns only a source-unit ID, exact `answer_quote`, and optional
exact paired `question_quote`. Deterministic code locates each quote as a
contiguous substring inside its own unit, rejects any unlocatable string, and
renders the attributed capsule. The boundary condition remains recorded as a
non-result; this is a distinct exact-string selection condition.

## Dual-field RAG follow-up

The capsule result changes the RAG representation rather than relabeling a
capsule as a standalone fact. Each retrieval item has two strictly separated
fields: (1) an **authoritative capsule**, the exact Q/A quote spans and their
canonical speakers, which is the only field supplied to answer generation; and
(2) a **non-authoritative retrieval expansion**, a resolved query-oriented
description used solely for embedding/search. An expansion may help find a
capsule but is never stored as a fact, treated as provenance, or shown as answer
evidence. `010_study_two_dual_field.sql` stores expansions keyed to their
capsule fact IDs in an isolated table. Evaluation must report expansion-driven
capsule evidence recall separately from answer support.

## Concepts tried and results

| Candidate | Structural concept | Fixed first-run configuration | Audit result | Decision |
| --- | --- | --- | --- | --- |
| `direct` | Study-1-style one-call full-transcript extraction | 45 cap; concise contract | 49/68 support (72.1%); 46/58 coverage (79.3%); attribution 98.5%, qualifier 92.6%, standalone 94.1%, atomicity 91.2% | Reject: evidence does not establish resolved antecedents. |
| `speaker-packed` | Explicit dialogue-neutral the production system-style contract; 3k-character speaker-boundary packs, no carryover | 30 cap/pack | 100/105 support (95.2%); 53/58 coverage (91.4%); attribution 100%, qualifier 97.1%, standalone 92.4%, atomicity 94.3%, near duplicates 4.8% | Reject: five unsupported cards violates zero-error/98% gate. |
| `carryover` | Same speaker packing with two prior turns in prompt context | 3k packs; two-turn carryover | Runtime/structural diagnostic retained; not promoted without a blind aggregate audit | Separates prompt continuity from stored proof. |
| `bundle` | Closed-proof model-selected evidence bundles with explicit examples | 3k packs, two-turn prompt carryover, 20 cap | 45/77 support (58.4%); attribution 98.7%, qualifier 84.4%, standalone 98.7%, atomicity 94.8%, coverage 46/58 (79.3%) | Reject: model resolved context but stored inadequate short quotes. |
| `expanded-bundle` | Model-selected bundle plus deterministic preceding-turn proof expansion | 3k/two-turn/20 | 31 retained cards; blind audit retained as local diagnostic | Too low output volume to plausibly clear 75% coverage. |
| `expanded-packed` | Explicit the production system-style extraction plus deterministic proof expansion: retained evidence turn and two preceding canonical turns | 3k packs, two-turn carryover, 30 cap | 62/70 support (88.6%); attribution 52/70 (74.3%); qualifier 60/70 (85.7%); standalone 66/70 (94.3%); atomicity 68/70 (97.1%); no duplicates; coverage 45/58 (77.6%), but only 14/20 (70%) for `NPR-967` | Reject before replication: broadening evidence cannot repair unsupported wording or bare-world rephrasing. |
| `attributed-packed` | New peak-jump: represent each dialogue claim as a source-speaker-attributed report; expanded canonical proof bundle | 3k packs, two-turn carryover, 30 cap | 46/47 support (97.9%); attribution 100%; qualifier 91.5%; standalone 61.7%; atomicity 74.5%; no duplicates; coverage 47/58 (81.0%) | Reject: it nearly fixes support, but source-speaker wrappers do not by themselves resolve deixis or enforce one claim. |
| `atomic-attributed` | Peak-jump: hard self-check for a single post-attribution proposition and resolved entities before emission | 3k packs, two-turn carryover, 30 cap | Retained 22, 11, and 6 cards respectively before audit | Reject before audit: `NPR-967` cannot reach its 15/20 (75%) ledger gate with six retained cards. |
| `turnwise-staged` | Bounded-concurrent extraction from each declarative target turn with two prior turns context-only, followed by a separate document-level atomic rewrite-or-reject stage | 3 candidates/target turn; fixed DeepSeek model; exact expanded proof bundles | Support 67/68 (98.5%), attribution 67/68 (98.5%), but standalone 31/68 (45.6%), atomicity 59/68 (86.8%); `NPR-967` coverage 10/20 (50%) | Reject: the document-level rewrite retained dialogue deixis and undercovered multi-claim turns. |
| `turnwise-local-staged` | Peak-jump: same source-first target-turn candidate generation, but independently rewrites/rejects each target's candidates using only its local two-turn context; no global rewrite bottleneck | 3 candidates/target; bounded 8-way concurrency; deterministic expanded proof bundle | 85/89 support (95.5%), attribution 100%, qualifier 96.6%, standalone 66.3%, atomicity 91.0%, duplicates 1.1%; coverage 16/20, 9/18, and 11/20 | Reject: four unsupported cards, pervasive deixis, and coverage below 75% on two documents. |
| `clause-anchored` | Deterministic lexical source-clause enumeration, deterministic source-speaker assignment, and rewrite/reject only against an anchored clause | DeepSeek V4 Pro; 8-way bounded concurrency; 5 clauses/turn cap | 97/107 support (90.7%), attribution 97.2%, qualifier 94.4%, standalone 88.8%, atomicity 98.1%, duplicate 2.8%; coverage 48/58 (82.8%), but `NPR-2291` 12/18 | Reject: clause coverage improved but source support remains materially below gate. |
| `clause-anchored` / `v2-glm52` | Frozen model-family ablation of the same deterministic source-first policy | GLM 5.2; same corpus, prompt, segmentation, cap, and validation policy | 61/79 support (77.2%), attribution 100%, qualifier 91.1%, standalone 98.7%, atomicity 86.1%, duplicate 2.5%; supported coverage 33/58 (56.9%) | Reject: alternative model worsened bundle-grounded support and coverage. |
| `proof-obligation` | Requires a selected exact deterministic source clause for every resolved element; expands quoted substrings to their containing source clause and rejects deixis | DeepSeek V4 Pro; fixed lexical clauses and deterministic speaker assignment | 57/62 support (91.9%), qualifier 85.5%, standalone 77.4%, atomicity 98.4%, duplicate 3.2%; coverage 31/58 (53.4%) | Reject: self-censoring proof policy collapses coverage without resolving all references. |
| `highrecall-verified` | Speaker-packed high-recall candidates, then independent full-source bundle entailment/atomicity/decontextualization verifier with a per-turn coverage instruction | 3k speaker packs, no carryover, 10 candidate cap/pack; 2.5k verifier output cap | 57/62 support (91.9%), attribution 100%, qualifier 90.3%, standalone 98.4%, atomicity 91.9%, duplicate 1.6%; supported coverage 39/58 (67.2%) | Reject: distinct verifier still permits agent/scope/causal substitutions and misses two document coverage gates. |
| `highrecall-nli` | Same high-recall candidates followed by a distinct closed-world NLI role with explicit KEEP/DROP/minimal-REWRITE decisions and phrase-level evidence requirements | Same stage-one policy; NLI output cap 2.5k | 55/59 support (93.2%), attribution 98.3%, qualifier 94.9%, standalone 100%, atomicity 94.9%, no duplicates; coverage 33/58 (56.9%) | Reject: discrete NLI decisions reduce some surface failures but do not establish semantic support or coverage. |
| `prefilter-nli` | Current deterministic prefilter architecture: remove interrogative and broadcast-framing/promotion turns before high-recall extraction; source speaker assigned only from evidence turn; then closed-world NLI | 3k eligible-turn packs; no carryover; 10 cap/pack | 16, 15, and 17 retained cards; fresh blind audit in progress | Tests whether removing question/framing provenance before generation prevents the observed speaker-misattribution class. |
| `decomposed-nli` | Current scope/decomposition architecture: deterministic source proposition units plus mandatory modality/scope fields; KEEP/DROP-only NLI cannot merge, rewrite, or delete scope | Eligible declarative turns; 5 units/turn cap; exact source speaker | 53, 12, and 38 retained cards; fresh blind audit in progress | Tests whether separation of proposition coverage from NLI selection preserves the prefilter's literal-support advantage while restoring atomicity/scope. |
| `framed-render` | Eligible source units -> typed subject/relation/object/modifier/evidence frame -> deterministic renderer | Per-turn bounded frame calls; every nonempty frame field must occur verbatim in selected clause | 96/111 support (86.5%), qualifier 83.8%, standalone 12.6%, atomicity 86.5%, coverage 13/58 | Reject: open-text frame fields fragment names and grammar; lexical field validation does not yield usable claims. |
| `fable-c4-sbar` | Structural cross-study replication: frozen full-transcript antecedent-window generator -> mechanical validation -> span-blind support audit -> restricted repair -> revalidation -> fresh re-audit | 60-card cap; four-turn antecedent window; DeepSeek V4 Pro; temperature 0/no reasoning; all phase prompts persisted | `NPR-1781`: 12/27 support (44.4%), qualifier 21/27 (77.8%), strict ledger coverage 9/20. `NPR-2291`: 15/20 support (75.0%), qualifier 18/20, fully passing coverage 10/18. `NPR-967`: 20/28 support (71.4%), qualifier/standalone 25/28, atomicity 27/28, strict coverage 15/20. No run had a duplicate; attribution was 100% throughout. | Reject: independent final audits found unresolved referents/scope, unsupported specialization, dropped or strengthened hedges, a question rendered as a fact, and expansion from a superlative to a universal. C4's internal span-blind screen is not reliable enough for the zero-unsupported-card gate; do not replicate. |
| `dialogue-act-pairs` | Deterministic immediate interviewer-question + different-speaker declarative-answer units; answer cards cite both Q and A and use attributed dialogue-act form when Q supplies scope; leftover declaratives remain standalone | 60-card full-document source-unit pass; DeepSeek V4 Pro; temperature 0/no reasoning; exact unit-local Q/A validation | `NPR-1781`: 35/37 support, 37/37 attribution, 35/37 qualifiers, but 22/37 standalone and 30/37 atomic; coverage 18/20. `NPR-2291`: 12/12 support/attribution/unique, 10/12 qualifier/standalone/atomic, coverage 14/18. `NPR-967`: 35/35 support/attribution/qualifiers, 22/35 standalone, 31/35 atomic, coverage 20/20. | Reject: Q/A provenance substantially improves grounding (and eliminates question-premise-as-fact errors), but it does not force a standalone single proposition. `NPR-2291` also cannot meet its 15/18 coverage gate with twelve cards. Next: predicate-level answer units within the same Q/A evidence structure. |
| `dialogue-act-predicate-units` | Same Q/A source/evidence structure, but one deterministic lexical answer predicate/coordination unit per generation target and one card maximum per target | 80-card full-document predicate-unit pass; DeepSeek V4 Pro; temperature 0/no reasoning; answer quote restricted to the target unit | `NPR-1781`: 62/69 support, 67/69 attribution, 69/69 qualifier, 39/69 standalone, 61/69 atomic, 66/69 unique, coverage 17/20. `NPR-2291`: support/compatible coverage 40/60, attribution 60/60, qualifier 58/60, standalone 54/60, atomicity 52/60, unique 57/60. `NPR-967`: 45/46 support and 20/20 ledger coverage; one failure miscredited a question premise to a reply. | Reject: lexical splitting created more targets but reintroduced question-to-assertion conversion, outside-bundle antecedents, deictic/clipped evidence, bundled cards, duplicates, and cross-question premise miscrediting. |
| `quote-capsules` | Boundary-only selection of one atomic exact answer span plus optional exact paired-question context; deterministic attributed rendering | 80-card full-document Q/A-unit pass; DeepSeek V4 Pro; UTF-8 boundary validation and no model-generated text | `NPR-1781`: provider returned 42 syntactically valid selector objects, but every answer range was the placeholder `[0,0)`; exact validator retained 0 cards. JSON field-tag and optional-question-sentinel handling were independently checked, then a repeat returned the same all-placeholder output. | Non-result/reject: the model did not perform numeric boundary selection, so no source span can be fabricated and a blind fact audit would be meaningless. |
| `quote-capsules-exact` | Exact-string-only selection of one answer quote plus optional paired-question quote; deterministic attributed capsule renderer | 80-card full-document Q/A-unit pass; DeepSeek V4 Pro; unit-local contiguous-quote validation and no model-generated fact text | All runs had 100% support/attribution/qualifier: `NPR-1781` 5/39 standalone, 26/39 atomic, coverage 7/20; `NPR-2291` 9/18 standalone, 15/18 atomic, coverage 9/18; `NPR-967` 10/25 standalone, 17/25 atomic, coverage 9/20. Duplicates were 1, 1, and 0. Evaluator reuse fallback is disclosed in each worksheet because agent capacity was blocked. | Reject as standalone fact condition: capsules are provenance-exact evidence but preserve dialogue deixis and question dependence. Retain them only as a potential authoritative evidence field for a dual-field RAG condition. |

The first five rows are not a hill climb: they vary full-transcript versus
speaker-boundary construction, zero versus two-turn carryover, direct versus
closed-proof bundle generation, and deterministic proof-bundle expansion. The
last candidate changes the semantic representation itself to preserve the
speech-act attribution that broadcast dialogue requires.

## Implementation and reproducibility

`isc-rag study-two` runs one explicit candidate/document/replicate.
It requires a validation-split document, writes into an isolated database
namespace, and persists the complete effective instruction (not merely a hash).
For example:

```sh
go run ./cmd/isc-rag study-two \
  --candidate attributed-packed --version v2 --replicate 1 \
  --external-id NPR-1781 --chunk-chars 3000 --carryover-turns 2 \
  --max-facts 30
```

The C4 replication is run once per fixed document with no chunking parameter:

```sh
go run ./cmd/isc-rag study-two \
  --candidate fable-c4-sbar --version v2-c4-sbar --replicate 1 \
  --external-id NPR-1781 --max-facts 60
```

For `expanded-*` and `attributed-packed`, validation persists the selected source
turn plus its two preceding **canonical source turns** as separately stored
exact spans. This makes context that was supplied to the model inspectable at
audit time; it is not generated context. Facts with unlocatable model evidence
are discarded. The database migration is `008_study_two.sql`; focused
packing tests cover carryover behavior. `turnwise-staged` adds bounded
eight-way concurrency because sequential per-turn calls exceeded the local
command-run window; aggregation remains ordered and every provider response is
stored. Its first stage may use preceding turns only to resolve a target-turn
reference; the second stage sees the candidates and canonical source turns and
must rewrite or reject them under the atomicity/fidelity contract.

## Current stop decision

No configuration has met the gate, so none is designated as a paper condition.
In particular, it would be invalid to run three “replicates” of a candidate
that already has a nonzero unsupported-card audit. The next concept should keep
attributed proof bundles but use a separate clause-level rewrite/rejection stage:
one model card must either be split into independently supported claims with
resolved referents or be dropped. It is audited before any replication; only a
candidate that clears the full one-run gate on each document can advance to two
independently evaluated repeat runs under the unchanged frozen configuration.
Costs and latency are stored per run in the database; aggregate cost should be
reported from those rows after a successful configuration is frozen.

The next meaningful jump is not another rephrasing instruction variant. It must
make source coverage deterministic before generation: segment each source turn
into candidate declarative clauses, retain an auditable per-turn coverage ledger,
and ask the model only to rewrite or reject each independently anchored clause.
This separates “which claims were considered” from “how a claim is
decontextualized,” preventing a model-level coverage preference from silently
dropping substantive multi-claim turns.
