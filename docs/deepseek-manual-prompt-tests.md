# Manual DeepSeek prompt diagnostic

**Status:** small exploratory diagnostic, 2026-07-24. This is independent of the ten-iteration study, not a paper result or frozen configuration. Raw MediaSum-derived completions and temporary prompts are ignored under `artifacts/deepseek-manual-prompt-tests/`.

## Question

Were the earlier failures plausibly caused by overly long instructions, rather than by document context or the provenance/standalone task itself?

## Existing prompt size

This table measures the current tracked instruction files exactly in characters. Token figures use `ceil(characters / 4)`: a deliberately rough comparison estimate, not a model-tokenizer count.

| Iteration | Instruction characters | Approx. tokens |
| ---: | ---: | ---: |
| 1 | 954 | 239 |
| 2 | 1,441 | 361 |
| 3 | 1,316 | 329 |
| 4 | 1,015 | 254 |
| 5 | 1,057 | 265 |
| 6 | 916 | 229 |
| 7 | 977 | 245 |
| 8 | 952 | 238 |
| 9 | 914 | 229 |
| 10 | 939 | 235 |

`NPR-847` has 114 turns. The labelled transcript sent to a direct extraction call was 33,231 characters (about 8,308 rough tokens); raw source text was 30,040 characters. Thus the largest system instruction was only about 4% of direct-call character input. Including transcript and its small user wrapper, direct calls were roughly 8.55k--8.68k rough input tokens.

For quote-first calls, stage 1 used the same full labelled transcript plus the 953-character (about 239-token) built-in quote-selection prompt. Stage 2 received selected quotes with local context, not the full transcript. The exact stage-2 payload was not persisted, so it cannot be reconstructed exactly from the run record.

### Reproducibility discrepancy

Recorded prompt hashes for iterations 1--4 match current tracked files. Recorded hashes for 5--10 do **not** match current tracked files (although those files were committed with the study). This does not affect the new runs below, but means current files alone cannot reproduce prior runs 5--10 exactly. The stored prompt/raw-output hashes are the authoritative identities for those runs; repair this before treating the prior study as fully reproducible.

## New isolated runs

All tests used account-visible `accounts/fireworks/models/deepseek-v4-pro`, `temperature=0`, `reasoning_effort=none`, `--max-facts 25`, and `NPR-847`. No account-visible DeepSeek V4 Flash model was available. Runs were isolated in namespace `deepseek-manual-prompt-diagnostic`, version `v1`, so they cannot alter primary facts or prior study runs.

Each followed `docker compose up -d postgres` and `.env` loading:

```sh
go run ./cmd/isc-rag deepseek-study \
  --namespace deepseek-manual-prompt-diagnostic --version v1 \
  --iteration N --approach APPROACH --max-facts 25 \
  --prompt-file PROMPT --output-dir artifacts/deepseek-manual-prompt-tests/runs
```

| Run | Approach / altered factor | Prompt chars (rough tokens) | Prompt hash used | Raw-output hash | Raw / validated facts |
| ---: | --- | ---: | --- | --- | ---: |
| 1 | Direct full transcript; compact schema/rules | 465 (117) | `5158654a2d38d4d69097339a19cba9d4633ad0710634044db043fee027c56b11` | `203e68f43934a90e50040d8c4f575a9c5ca06c8dd34b99e7a52adb912577be4a` | 25 / 25 |
| 2 | Direct full transcript; long iteration-2 control | 1,441 (361) | `d49a2d1b9c586ee8cb152201a51257c91ec26a65a33f6d546007c32106ffb70e` | `06ceefd914647b7ede50a61d8ef1f77cbcabe81d10a213b40d075f33af743e0f` | 25 / 24 |
| 3 | Quote-first; compact rephraser | 388 (97) | `5702cace844a5bfbec98e90e846349fdc382950eedb0b105eda42941ed224e85` | `0abc7ffcaeb6b98b82c1f1ee26e05788fe8c444f165a3c795ead594a9edb8d53` | 20 / 20 |
| 4 | Quote-first; concise strict schema | 516 (129) | `85ef956fde714bb71660a0c12858b0e55ec36e6a6242e183e8404e4995b3d65d` | `e012574ce4235e96ea6537535e9825b12b5c3ca1d53116aa6d6d5be44e6c6bb9` | 21 / 21 |

For runs 3--4, the stored effective prompt hash includes the built-in quote-selection prompt, so it intentionally differs from the rephraser-file hash used for its character count.

## Manual evaluation

One fresh manual evaluator (the author of this document) inspected a fixed spread of eight output ordinals per run: 1, 4, 7, 10, 13, 16, 19, and 20. The same binary five-part rubric was used: (1) explicit source support/entailment; (2) attribution or source-speaker preservation, including persisted speaker metadata; (3) preservation of qualifier, negation, modality, time, and reported source; (4) standalone clarity; and (5) atomicity (one predicate-level claim).

This is a limited, non-blind, non-independent diagnostic. Scores are not statistical evidence and must not be extrapolated to all facts or documents.

| Run | Support | Attribution | Qualifiers | Standalone | Atomic | Total |
| ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| 1: direct compact | 6/8 | 8/8 | 5/8 | 5/8 | 8/8 | 32/40 |
| 2: direct long | 5/8 | 8/8 | 7/8 | 6/8 | 8/8 | 34/40 |
| 3: quote-first compact | 8/8 | 8/8 | 8/8 | 2/8 | 5/8 | 31/40 |
| 4: quote-first schema | 8/8 | 8/8 | 8/8 | 4/8 | 7/8 | 35/40 |

Error patterns explain the scores. Compact direct dropped epistemic framing, repaired a pronoun using unquoted context ("it" to "the program"), and made an ungrammatical fifth-sanctions claim. The long control retained more stance wording, but added an unquoted comparison to U.N. sanctions and reported-source context absent from a selected quote. Quote-first raised literal support because output often copied selected text, but retained unresolved "this report," "it," "there was," and "the next question." Its selector also supplied multi-claim quotes, so concise schema instructions did not reliably enforce atomicity.

## Findings and recommendation

There is no diagnostic signal that earlier system prompts were too long. They are short relative to the 8.3k-token transcript context: shrinking a 361-token direct prompt to 117 tokens produced 25 rather than 24 validated facts, but did not improve the sample rubric (32/40 versus 34/40). Quote-first reduced second-stage context but did not solve standalone clarity (2/8 and 4/8).

Stronger hypotheses are structural:

- A 25-fact cap across 114 turns imposes a coverage ceiling independent of prompt length.
- A fact may need antecedent context to be standalone, while the current contract requires its rewritten text to be entailed by one short immutable quote. That is a genuine conflict, not merely a prompting issue.
- Transcript disfluencies and compound turns make LLM quote selection unreliable for one-predicate evidence spans.
- Speaker provenance should be attached deterministically from a source turn rather than delegated to rephrasing instructions.

For the next development experiment, keep a compact schema but change the pipeline rather than merely shorten prompts: deterministically segment candidate turns/clauses; attach speaker metadata programmatically; permit a clearly defined expanded evidence window when a rephrasing resolves an antecedent; and validate the card against that complete evidence window with an independent verifier. Use a fixed claim ledger and budget for coverage evaluation. Do not promote any of these configurations to the paper condition.
