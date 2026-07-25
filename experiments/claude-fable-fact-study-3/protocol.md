# Claude Fable fact study #3 — frozen protocol

**Status:** frozen before ledger construction and before any extraction tuning
(2026-07-24). Changes after this point are recorded as amendments at the bottom;
none may alter the gate, rubric, metrics, or document set.

## Purpose and isolation

Study #3 tests whether Claude Fable 5 itself, under structural pipeline designs
(not merely prompt edits), can produce provenance-exact, standalone, atomic fact
cards from MediaSum broadcast dialogue at a fixed quality gate. It is an
isolated development study:

- No existing Go code, migrations, README/PLAN, or prior documents are modified.
- The database is not written to. (See sandbox constraints: it could not even be
  read; canonical turns are derived byte-exactly from the pinned corpus files
  using the ingestion rules in `internal/mediasum/ingest.go`.)
- Only MediaSum **validation-split** documents are used. The validation split is
  document-disjoint from the test split by the MediaSum release, so NPR-847 and
  all 25 canonical-question test documents are excluded by construction. No
  canonical question/answer labels are read or used. MediaSum summaries are
  never given to extractors or evaluators.
- Raw MediaSum-derived material (turns, ledgers, candidate facts, evaluator
  reports quoting source text) lives only under ignored
  `artifacts/claude-fable-fact-study-3/`. Tracked files contain no transcript
  text.

## Sandbox constraints (recorded, not chosen)

This session's execution sandbox permission-blocks every database client
(`psql`, `docker`, `go run`, Python script execution) and every hashing utility
(`shasum`, `md5`, `cksum`, `openssl`, `git hash-object`). Therefore:

- Canonical turns come from `resources/mediasum/validation.json` (SHA-256
  verified at download; see `resources/mediasum/MANIFEST.sha256`), derived with
  an auditable sed/tr pipeline and verified by an escape/re-escape roundtrip
  `diff` plus turn/speaker count equality. `internal/mediasum/ingest.go` stores
  `utt`/`speaker` arrays verbatim with ordinal = array index, so these files are
  byte-identical to database canonical turns.
- Artifact identity is recorded as exact byte and line counts (`wc -c`,
  `wc -l`) plus full local retention, instead of SHA-256 hashes.

## Documents (fixed)

| Role | Document | Validation-split line | Turns | Turn bytes |
| --- | --- | ---: | ---: | ---: |
| Development (tuning) | NPR-1394 | 33 | 24 | 6,001 |
| Gate replicate 1 | NPR-573 | 12 | 49 | 11,344 |
| Gate replicate 2 | NPR-3311 | 66 | 36 | 8,624 |
| Gate replicate 3 | NPR-2291 | 46 | 54 | 4,741 |

Turn files: `artifacts/claude-fable-fact-study-3/docs/<ID>.turns.txt` (line N =
canonical ordinal N−1; the `\n`-joined file is byte-identical to the transcript
string used by prior studies), `<ID>.speakers.txt` (parallel), and
`<ID>.labeled.tsv` (`speaker<TAB>text`). Turn text contains no tabs, `|`, or
backslash escapes other than `\"` (verified), so the derivation is exact.

Tuning happens **only** on NPR-1394. Gate documents receive only frozen-
configuration runs; every run on a gate document (including failed gate
attempts) is logged as exposure.

## Extractor (fixed settings)

- Model: Claude Fable 5 (`claude-fable-5`), invoked as a **fresh** Claude Code
  subagent per extraction run with no access to this study's history, ledgers,
  or evaluations; it may read only the files listed in its prompt.
- Sampling: Claude Code Agent-tool defaults (temperature is not exposed by the
  harness); identical invocation settings for every run, including final
  replication. This is the stated-settings consistency for the study.
- Fact budget: at most 60 facts per document (never reached in practice; no
  small cap, so coverage is not artificially ceilinged).
- Output contract (TSV, one fact per line, no header):
  `fact_text<TAB>primary_quote<TAB>antecedent_quote_or_-<TAB>primary_turn_line`

## Deterministic validation (no LLM)

Applied by the orchestrator with exact string tools only:

1. `primary_quote` must be a verbatim substring of the single turn line it
   cites, ≤240 characters. Otherwise the fact is dropped.
2. `antecedent_quote`, when present, must be a verbatim substring (≤240 chars)
   of one turn line in `[primary_line − 4, primary_line]`. Otherwise the fact is
   dropped. The antecedent window is the structural mechanism that lets a
   standalone rephrasing cite the context that resolves a referent, replacing
   the prior single-quote-only contract.
3. Speaker is attached deterministically from the primary turn line's canonical
   speaker label (never from model output). First occurrence wins if a quote
   appears in multiple turns.
4. Exact duplicate fact texts are dropped deterministically (first kept).

Validated facts are written to `runs/<run-id>/validated.tsv` as
`index<TAB>fact_text<TAB>primary_quote<TAB>antecedent_quote_or_-<TAB>primary_line<TAB>antecedent_line_or_-<TAB>speaker`.

## Evaluation (locked)

Every candidate run is scored by a **fresh `fact_evaluator` subagent** that
receives only: the run's `validated.tsv`, the document's `labeled.tsv`
transcript, the frozen ledger for that document, and `rubric.md` verbatim. It
is given an opaque run id, no configuration description, no prior scores, and
no other runs. It returns per-fact binary rubric rows, duplicate marks, and
per-ledger-claim coverage rows in the fixed format defined in `rubric.md`.

Metrics per run (computed from evaluator output):

- support% = facts passing support / validated facts; unsupported = count failing
- attribution%, qualifiers%, standalone%, atomicity% likewise
- duplicate rate = facts marked duplicate-of-earlier / validated facts
- ledger coverage = covered ledger claims / total ledger claims

## Frozen gate

One configuration is frozen, then run once per gate document (three
document-level replicate runs, fresh extractor and fresh evaluator each). The
gate is met iff, pooled across the three replicate runs:

- support ≥ 98% **and** zero unsupported facts in any run;
- attribution ≥ 98%;
- qualifier preservation ≥ 95%;
- standalone clarity ≥ 90%;
- atomicity ≥ 90%;
- duplicate rate ≤ 5%;

**and** ledger coverage ≥ 75% in **each** document individually.

Per-document values are also reported. If any criterion fails, the study
diagnoses, iterates further (structural changes preferred once local changes
plateau, with fresh `study_designer` consultations), and re-attempts with a
newly frozen configuration. Every attempt is recorded.

## Claim ledgers (fixed before tuning)

For each document, before any extraction run: a fresh subagent enumerates every
salient explicit claim from the labeled transcript only (no summary, no
external knowledge), as `L<n><TAB>turn_line(s)<TAB>claim paraphrase` with a
verbatim anchor quote. A second, independent subagent audits the draft for
unsupported entries and missed salient claims; the orchestrator applies audit
corrections and freezes the ledger (byte/line counts recorded). Ledger scope
rules: include explicit substantive assertions (world facts, stated positions,
predictions, first-person experiences); exclude greetings, station/housekeeping
language, bare questions, and music/credit boilerplate. One claim per
predicate. Frozen ledgers are never edited afterward.

## Roles

- Orchestrator (this session): data derivation, deterministic validation,
  metric arithmetic, freezing decisions, report.
- Extraction: fresh Claude Fable subagents (one per run).
- Evaluation: fresh `fact_evaluator` subagents (one per candidate).
- Design critique: fresh `study_designer` subagents (baseline review and at
  plateaus).

## Amendments

**A1 (before tuning).** Deterministic provenance validation is *executed* by a
fresh mechanical validator subagent, not by the orchestrator shell. Reason: the
session sandbox permission-blocks loops, `awk`, compound commands, process
substitution, and `grep -f` pattern files, so batch string validation cannot run
in the orchestrator. The validator performs NO semantic judgment: for each
candidate it runs `grep -Fq "<quote>"` against the cited turn line (exit status
only), drops any fact whose primary quote is not a verbatim substring of its
cited line or whose antecedent quote is not a verbatim substring of a line in
`[primary−4, primary]`, attaches the speaker from the primary line's canonical
label, and drops exact-duplicate fact texts (first kept). It records every drop
with the failing `grep` outcome in `validation-report.txt`. The extractor also
self-validates with the same `grep -Fq` before emitting, so validation is
independent and double-checked. The orchestrator spot-checks a sample of quotes
per run with its own inline `grep -Fq`. This delegates execution only; the rule
set is unchanged from "Deterministic validation" above.

**A2 (before tuning).** Frozen ledgers: NPR-1394 = 59 claims, NPR-573 = 111,
NPR-3311 = 86, NPR-2291 = 29. Each draft was built by one fresh subagent from the
labeled transcript only and audited by a second independent fresh subagent; audit
FIX/DROP/MISSING items were applied before freezing (NPR-1394: 3 attribution/
scope fixes + 2 added guest-credential claims; NPR-3311: 1 hedge fix + 2 added
intro claims; NPR-573 and NPR-2291: audits clean, frozen unchanged).
