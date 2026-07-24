# DeepSeek fact extraction / rephrasing iteration study

**Status:** completed exploratory prompt study; not a paper result and not a
configuration freeze. This document records the ten requested, sequential
iterations so that the retained output and the decision not to promote it are
auditable.

## Scope and reproducibility

- **Corpus/document:** MediaSum test document `NPR-847`, *Talk of the Nation*;
  114 canonical turns and 29,927 source-text characters. It was selected before
  any prompt iteration because it has substantive multi-speaker discussion. The
  MediaSum summary was never supplied to extraction or evaluators as evidence.
- **Model request:** the requested “DeepSeek 4 Flash.” Fireworks account model
  discovery on 2026-07-24 exposed only
  `accounts/fireworks/models/deepseek-v4-pro`; no V4 Flash identifier was
  account-visible. Every run therefore used the accessible V4 Pro model, with
  `temperature=0` and `reasoning_effort=none`. `isc-rag deepseek-study --model`
  remains an explicit override for a future Flash deployment.
- **Isolation:** namespace `deepseek-fact-rephrasing`, version `v1`; raw output
  and validated facts live in `fact_rephrasing_study_runs` and
  `fact_rephrasing_study_facts`, not the primary `facts` table. No primary facts
  version or canonical labels was modified.
- **Artifacts:** tracked prompts are under
  `experiments/deepseek-fact-rephrasing-v1/prompts/`; local raw completions and
  validated quotations are under ignored `artifacts/deepseek-fact-rephrasing-v1/`
  because they derive from restricted MediaSum data. Each ignored iteration
  folder includes `raw-output.json`, `validated-facts.json`, and its short blind
  evaluation note. The database preserves the raw-output hash and exact quote
  offsets. **Post-study reproducibility correction:** the tracked prompt files
  reproduce execution hashes for iterations 1--4 only. Iterations 5--10 have
  correct execution prompt hashes below, but their current tracked prompt files
  differ and no exact executed-prompt copy was retained. Do not represent those
  six runs as exactly replayable from the repository; future runs must persist
  the complete effective instruction prompt alongside its hash before calling the
  provider. The follow-on `007_study_effective_prompt.sql` migration and CLI
  implementation now do so for future runs, both in the isolated database table
  and local ignored artifacts.

## Command

Run one controlled iteration after `docker compose up -d postgres`:

```sh
go run ./cmd/isc-rag deepseek-study \
  --iteration 10 --approach quote-first --max-facts 20 \
  --prompt-file experiments/deepseek-fact-rephrasing-v1/prompts/iteration-10.txt
```

The default document is `NPR-847`; `--external-id`, `--model`, `--namespace`,
`--version`, and `--output-dir` are explicit. `direct` sends the full numbered
transcript once. `quote-first` first selects exact contiguous candidate quotes
from the full transcript, then rephrases only those candidates; its raw output
records both stages. The command validates every persisted quote by exact lookup
in canonical source text and attaches source offsets, turn range, and speaker.

## Locked evaluation protocol

After each extraction, a **new child evaluator** was spawned. The evaluator was
asked to inspect only that iteration's validated facts and the cited canonical
turns; it received neither prior scores nor the MediaSum summary. It assessed
every fact on this fixed binary rubric:

1. explicit support / entailment;
2. attribution / speaker preservation;
3. qualifier, negation, and temporal preservation;
4. standalone clarity; and
5. atomicity.

It also counted unsupported facts, duplicate or near-duplicate facts, and
salient explicit claims missed. Missed-claim counts are conservative qualitative
audits, not a fixed gold-recall denominator. Only the immediately preceding
evaluation informed the next prompt/context/approach change. No canonical QA
labels or answer set were used for tuning.

## Results

| Iteration | Approach / material change | Validated | Support | Attribution | Qualifiers | Standalone | Atomic | Total | Unsupported | Near dup. | Salient missed* |
| ---: | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| 1 | Direct full transcript; baseline prompt | 48 | 47 | 48 | 33 | 47 | 39 | 214/240 (89.2%) | 1 | 1 | ≥9 |
| 2 | Coverage sweep; explicit hedge/scope instructions | 35 | 34 | 34 | 27 | 30 | 28 | 153/175 (87.4%) | 1 | 1 | ≥37 |
| 3 | Internal quote-first planning; deictic/compound checks | 37 | 36 | 37 | 28 | 36 | 26 | 163/185 (88.1%) | 1 | 0 | not enumerated |
| 4 | Actual two-call quote-first; 15-span cap | 15 | 14 | 15 | 9 | 13 | 8 | 59/75 (78.7%) | 1 | 0 | ≥45 |
| 5 | Coverage-first quote selection; strict card checklist | 24 | 24 | 24 | 23 | 21 | 19 | **111/120 (92.5%)** | 0 | 1 | 24/46 ledger |
| 6 | Higher 30-span target; clause/pronoun rules | 28 | 27 | 27 | 23 | 14 | 25 | 116/140 (82.9%) | 1 | 0 | ≥15 |
| 7 | Add source-turn + adjacent-turn context | 17 | 16 | 17 | 10 | 15 | 13 | 71/85 (83.5%) | 1 | 0 | ≥10 |
| 8 | Operator-complete quote selection | 18 | 17 | 18 | 17 | 8 | 15 | 75/90 (83.3%) | 1 | 0 | ≥18 |
| 9 | Antecedent-complete quote selection; reject unresolved cards | 18 | 17 | 18 | 17 | 7 | 12 | 71/90 (78.9%) | 1 | 0 | ≥16 |
| 10 | Quote-alone rephrase; named speaker; no context additions | 19 | 19 | 19 | 19 | 7 | 14 | 78/95 (82.1%) | 0 | 1 | ≥31 |

\*Not comparable as a formal recall metric; different evaluators produced
conservative, uncertainty-qualified audits. Percentages must not be interpreted
as cross-iteration statistical evidence because output count and span cap changed.

The iteration-5 percentage is the best observed rubric score, but it represents
only 24 facts and 22 of a reviewer’s 46-claim ledger. It is a useful diagnostic,
not evidence that prompt 5 is globally superior. Iteration 10 provides the
strongest exact-provenance/fidelity cells but still fails standalone clarity for
12 of 19 facts and misses much of the transcript.

## Run identities

All runs used `accounts/fireworks/models/deepseek-v4-pro`. Prompt/raw-output
SHA-256 hashes and fact counts are below; the prompt hash includes the quote
selection stage when present. The hashes identify what was executed. As noted
above, only the iteration 1--4 hashes currently match the corresponding tracked
prompt files; iteration 5--10 prompt files are useful development records but
not exact replay artifacts.

| Iteration | Prompt SHA-256 | Raw output SHA-256 | Facts |
| ---: | --- | --- | ---: |
| 1 | `3b1eaa7ebc965b3983d9920e441a9d8839780cfa3ba349572f09db58cb0bd457` | `f5072d2884c09b4bf422dfadb21a0d5d1c2877e54c3e6e6497527e583a6cf03a` | 48 |
| 2 | `d49a2d1b9c586ee8cb152201a51257c91ec26a65a33f6d546007c32106ffb70e` | `3d012678d9380bc92c79914eec523dd8c40fafb5cc0ac9987564eec82b7c22fc` | 35 |
| 3 | `41dbb2f8b83b07d319a3aace00d434a1948407a60f52eee77561746861e11ee5` | `a68ac0734e0ca0fb0e6311badd9a43df85a1d62e83957d1e5549604c7c39e049` | 37 |
| 4 | `642439961abd38b491e011f4bf5f85aac2a78772dcda548b423d2d070c4ab433` | `f094d484689fc61af50b837c1655c57aac288266d80a85ffebb4dddbdf4f01b8` | 15 |
| 5 | `f8dcf36b842b667580db25c377b426f72c1f364b90e86741708fd75db403d026` | `1a892372d951d0aec5dcc75d76dc4aba8ded6cf513d3e0c725b4d73b97f4a81d` | 24 |
| 6 | `a77fa6d1aded41efc72164137a81ccf97ee9d60af023eb525dc5d5782b5545c9` | `70ae7fa355d24f4e08b669e7dbff2b26381c9e388b496062c976041b4d5d8324` | 28 |
| 7 | `97860ae2b626ef89486f12ec43b82a3966b208f043eb17336e81daaa55ba709e` | `02f87db4a17ca3276c096b78ef642327fc5519a34fad676f73332f9947db1d67` | 17 |
| 8 | `be504218035fb7bc2a34b4b18061f1de5a476da9a0c066b895b6bc365570d966` | `f35101dec261a541389813fec05dd771c1c99fba38af3c753b4ca04d85c071d2` | 18 |
| 9 | `a2de916ebf67774ba1b655f44614ad3d085e32eec11c612315b35baddced21af` | `9763e5f66187736222c864f80a1b137325a1a9c8b8f6e2db261e5307d6e307e4` | 18 |
| 10 | `8a084a7da533ef18680ad39c2e5f16c87036b408fa50384de05abcb4465cbf08` | `71e79ab88936301151490ed95b48b2a9abdf1928d3f87a7e73db098e8376cf17` | 19 |

## Final recommendation

Do **not** promote any of these prompts as the paper's primary fact condition.
Exact quote validation plus DeepSeek V4 Pro can yield high local support and
speaker attribution, especially in iteration 5 and iteration 10. However, the
same study repeatedly found three unresolved problems: (1) output facts often
remain context-dependent despite exact provenance; (2) enforcing quality through
small quote caps sharply reduces document coverage; and (3) per-iteration
evaluators and changing output counts cannot establish a robust winner.

For a next, preregistered development experiment, separate **quote selection**,
**decontextualizing rephrasing**, and **deterministic/independent validation**;
set a fixed extraction budget and a fixed human-audited claim ledger; then compare
against the direct baseline on document-disjoint development documents. Require
an explicit-source/standalone validator to reject unresolved referents,
comparatives, incomplete predicates, and multi-claim cards. Only after that
configuration is frozen should it enter the four-condition ISC-RAG evaluation.
