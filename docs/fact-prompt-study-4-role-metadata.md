# Study #4 — deterministic speaker metadata prompt seeds

**Status:** completed exploratory development study; not a paper result, not a
canonical evaluation, and not a configuration to promote. This is one of three
parallel five-seed studies requested after deterministic `speaker_name` and
`speaker_role` annotation was added.

## Question and controls

Does explicitly supplying deterministic speaker name and role improve direct
fact-card extraction from a public-radio transcript? All seven calls used
Fireworks `accounts/fireworks/models/deepseek-v4-pro`, `temperature=0`,
`reasoning_effort=none`, JSON-object response mode, a 6,000-token completion
limit, and the same one-turn MediaSum development document, `NPR-1040`. The
document was selected by the coordinating experiment to remain disjoint from
the other parallel tracks and from canonical labels.

Its sole turn has deterministic `speaker_name=NOEL KING` and
`speaker_role=host`. This makes the study a narrow test of metadata propagation
and reporting-source language, **not** a general test of multi-speaker
interviews or of role diversity. The source-only six-claim ledger was fixed
before the first request in
[the tracked ledger](../experiments/fact-prompt-study-4-role-metadata/claim-ledger.md).

Prompts and the runner are tracked in
[`experiments/fact-prompt-study-4-role-metadata/`](../experiments/fact-prompt-study-4-role-metadata/).
The runner supplies exact source turns with the migrated metadata fields only;
it makes no dialogue-based role inference and writes no primary facts or study
tables. Source transcript copies, requests, provider completions, and response
metadata stay in ignored `artifacts/fact-prompt-study-4-role-metadata/` under
the MediaSum research-use restriction. The table records SHA-256 identities for
the exact local requests and provider responses.

## Fixed audit rubric

I manually audited every returned card against the locked ledger and source
turn. `Support` means the complete card text is entailed by its exact quote;
`attribution` means its name, role, and turn fields exactly match the quoted
turn; `qualifiers` includes reported-source markers such as “we're told”;
`standalone` means the card text can be understood without prior transcript
turns; and `atomic` means it makes only one independently checkable
proposition. Duplicate is an exact or near-duplicate card count. Ledger
coverage requires a source-supported card that preserves the ledger claim's
stance and qualifier. Empty outputs are shown as n/a for per-card percentages,
not as perfect scores.

## Five independent seeds

| Seed | Deliberate starting idea | Cards | Support | Attribution | Qualifiers | Standalone | Atomic | Dup. | Ledger coverage | Latency / tokens |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| A, role-aware direct | Explicit named-speaker attribution and role-aware direct facts | 0 | n/a | n/a | n/a | n/a | n/a | 0 | 0/6 | 1s / 380+7 |
| B, speech-act reporting | Conservative `says`/`defines` reporting verbs | 0 | n/a | n/a | n/a | n/a | n/a | 0 | 0/6 | 1s / 362+7 |
| C, metadata-minimal | Minimal schema and metadata copy instruction | 5 | 5/5 | 5/5 | 4/5 | 3/5 | 2/5 | 0 | 4/6 | 3s / 275+301 |
| D, question-aware | Reject question premises; retain speech acts and qualifiers | 8 | 8/8 | 8/8 | 8/8 | 2/8 | 6/8 | 0 | 6/6 | 5s / 360+474 |
| E, evidence-bundle-first | Quote selection before conservative reporting | 4 | 3/4 | 4/4 | 3/4 | 2/4 | 2/4 | 0 | 4/6 | 3s / 341+330 |

Token cells are provider `prompt_tokens + completion_tokens`; the provider
response did not contain a billable dollar amount, so cost is unavailable rather
than estimated. Cards and quotations are not reproduced here because they are
source-derived; their exact request/response hashes are below.

Seed D had the best ledger coverage and perfect support, metadata attribution,
and qualifier retention. Its weakness is central to the ISC-RAG question:
although it preserves the host metadata, six of its eight fact texts use
context-bound references such as an unnamed couple, “both,” or “they.” Metadata
identifies the reporter, but does not make other discourse entities standalone.

## Hill climb from seed D

Two refinements were tested sequentially. Each retains the model, document,
schema, and all seed-D constraints; it changes one documented component.

| Run | Single change | Cards | Support | Attribution | Qualifiers | Standalone | Atomic | Dup. | Ledger coverage | Latency / tokens |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| R1, reporting source | Explicitly require retaining “we're told”/“according to” rather than treating a report as the host's firsthand claim | 7 | 7/7 | 7/7 | 7/7 | 2/7 | 6/7 | 0 | 5/6 | 5s / 412+421 |
| R2, clause boundaries | Add a clause-at-a-time, shortest-limiting-span instruction to R1 | 1 | 1/1 | 1/1 | 1/1 | 1/1 | 1/1 | 0 | 1/6 | 1s / 484+68 |

R1 correctly retained the reported-source marker and made cards more atomic,
but it dropped the ledger's “shutdown is personal” claim. R2 became severely
over-conservative, returning only the reported closure. Neither improves on
seed D's full coverage, and neither solves the standalone-reference problem, so
this hill climb stops rather than adding prompt complexity without a measured
gain.

## Reproducibility identities

| Run | Prompt file | Request SHA-256 | Provider response SHA-256 |
|---|---|---|---|
| A | `seed-a-role-aware-direct.md` | `0326efa351da33a91755cc3f00a885445f21022751ea73059b4a79a87f2c226f` | `0cc484798802ff5316cda9e192ededb6bd8e8850547c179597bfb2998862e22d` |
| B | `seed-b-speech-act-reporting.md` | `75092a6c40b2974770a32367d0e38ae181b20adc4ed8a869944b6630b5933653` | `8d9ea6454263796b1e7b19baed78d00af67352be323d3a5d7006e2fd34d3130e` |
| C | `seed-c-metadata-minimal.md` | `f3546a51c3575c2f5bd99f267b3a7448328bc6bed18515d964f5d029a9ae3b6f` | `64570694680548ee0e2c949bf3211c915c1551c488e4a75ec2eb7e3cfff71d3d` |
| D | `seed-d-question-aware.md` | `3c10b0adb79290cf8379db71d2ebad2cbe706efdacd2b0de19cb353e4da75684` | `c23334891956a3afa33913bdf2830c88d2bffcdb0aad855b9794c241b4f2aab1` |
| E | `seed-e-evidence-bundle-first.md` | `afbc9ac0d0487928031235c3b567e35e2be39a16127e9ecf1217846cb8f9a6c3` | `df15370a63b8615dbb26eebcfb4aae1aa1eb7dd113e9d63dc7e8d22c5b1bd27d` |
| R1 | `refinement-1-reporting-source.md` | `dbcb8c91608d57ad2e0c8644439ad85a99a7dad9e146b6702a7ccc93ee28c7b5` | `5cea2ea0ac9650e3e4162afdad26e6dedf741c4bb2fe1991a787d12cc4262eee` |
| R2 | `refinement-2-clause-boundaries.md` | `6736851fcae5f0e444d272d27494c8a7ab0cafd07e11d987c3c44a32e8d9c771` | `c938ca6d0e33a993f4e85685e0e40f9288ac0d6feabf580b755a076ff727decb` |

## Honest outcome

Deterministic role metadata worked mechanically: every nonempty card copied the
correct name, role, and source turn, and the reporting-source refinement kept
the one source-attribution qualifier. But this is insufficient evidence for a
production fact prompt. The development material has one reporter, one short
turn, no actual interviewer/answer exchange, and unnamed third parties. It
demonstrates that role metadata cannot by itself resolve third-party coreference
or make extracted propositions self-contained. The experiment therefore does
not promote a prompt, and it supports the existing evidence-capsule design over
free-form standalone fact rewriting.
