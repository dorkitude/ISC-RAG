# Fact prompt study 6: dual authoritative-evidence and retrieval-description fields

**Status:** exploratory development prompt study; not a paper condition, not a
configuration freeze, and not a primary fact representation.

## Objective and boundary

This study tests a deliberately split record:

- **Authoritative evidence:** `quote`, a short exact contiguous source span.
- **Retrieval description:** `text`, always prefixed `[RETRIEVAL-ONLY]` and
  explicitly non-authoritative. It can improve matching and display, but cannot
  establish, expand, or repair a claim.

The split is important: a readable resolved label is useful to retrieval, while
the exact source evidence remains the only field a downstream factual answer may
cite. All support and ledger-coverage scoring below therefore uses `quote`
alone. No generated description is counted as a fact.

## Fixed development material

- **Document:** `NPR-1964`, *Crews Continue Searching Through Remains In
  California's Camp Fire*.
- **Canonical corpus row:** MediaSum `validation` split (the supplied
  development document is recorded as validation in the local database).
- **Source-only ledger:** 20 atomic, source-attributed claims, frozen before
  model execution at
  [`claim-ledger-npr1964.md`](../experiments/fact-prompt-study-6-dual-field/claim-ledger-npr1964.md).
  No summary, QA label, or model output informed it.
- **Model/configuration:** Fireworks
  `accounts/fireworks/models/deepseek-v4-pro`, `temperature=0`,
  `reasoning_effort=none`, quote-first selection, canonical numbered speaker
  turns. The account-visible V4 Pro model was the project-configured option.
- **Isolation:** namespace `fact-prompt-study-6-dual-field`, version `v1`, in
  the existing isolated study tables. Primary `facts`, canonical labels,
  migrations, and shared implementation were untouched.

For every run, the runner first selected candidate exact spans, then generated
JSON records of `text` and `quote`. It discarded any record whose quote was not
an exact transcript substring or whose text duplicated another record. The
prompt further required the retrieval-only prefix. Raw completions, effective
prompts, and validated records are under ignored
`artifacts/fact-prompt-study-6-dual-field/iteration-XX/`.

## Locked scoring

The two fields were evaluated separately.

1. **Evidence validity:** the mechanical exact-contiguous validation rate.
2. **Evidence coverage:** number of frozen ledger claims explicitly supported
   by an authoritative quote. A generated description cannot add coverage.
3. **Retrieval-description utility:** a manual, non-factual screen for a clear
   retrieval hint with the visible marker, a topical/speaker cue, and no material
   premise absent from its stored quote. This is not an entailment or support
   score; a utility pass never promotes the text to evidence.

All evaluated records had the required prefix and no exact-quote duplicates.
The full worksheet is
[`evaluation.md`](../artifacts/fact-prompt-study-6-dual-field/evaluation.md).

## Five seeded runs

| Seed | Prompt idea | Raw / validated | Exact evidence | Ledger coverage | Retrieval utility | Prompt SHA-256 | Raw SHA-256 |
| ---: | --- | ---: | ---: | ---: | ---: | --- | --- |
| 01 | Literal labels using quote words | 14 / 14 | 14/14 | 10/20 | 12/14 | `91cd9ff29ab95ce1a79aba532268138015282d3a005334925689c00c4e62b450` | `d3871f458de494eaebf5c25d123c52c11294b46901535bb574201c9ee7c7b2be` |
| 02 | Supplied speaker plus local-context capsule | 11 / 11 | 11/11 | 9/20 | 8/11 | `b43582037d9052e9cc578374336c7466a6b86bbbe36fb9922e9e2c36af08af77` | `d038da776be81aed13acbfc25705f8cf10f5d11478bbbb7e35efdeb060133da7` |
| 03 | Explicitly resolved retrieval wording | 12 / 12 | 12/12 | 9/20 | 10/12 | `9f5de027c7de9dff1f3108e088d2dd2bb8f1a19ce2252714f873956abaf9481e` | `1b0ed4fc37f339c06abf8480ebabb76fadba7b5ab241ba68231825292ac29a0c` |
| 04 | Source-operator-first wording | 12 / 12 | 12/12 | 9/20 | 10/12 | `156e293ce77170744eebbd08a85c0f1c734c04d079f56245f3f6a1fac5c1f952` | `a64c3ef0e6c8a4477a9f3e7d0ea885c8bb95e7bb39f8c2ad8a9f510e714ede74` |
| 05 | Query-oriented speaker/topic description | 13 / 13 | 13/13 | 10/20 | 11/13 | `dd0e93f22ed22d99a1905ed090024f9022f3e9e2bc97c0ce86c7cd3e73aa3de5` | `8bd76d83449eca6d7c002cc0a7efcd82cec20f68ef92cf37fad57e20d3a5c1ea` |

Seed 05 was selected for hill climbing: it tied the best initial evidence
coverage, retained a large number of exact evidence bundles, and had the most
useful retrieval hints without treating them as evidence.

## Hill climb

| Run | Change from seed 05 | Raw / validated | Exact evidence | Ledger coverage | Retrieval utility | Prompt SHA-256 | Raw SHA-256 | Result |
| --- | --- | ---: | ---: | ---: | ---: | --- | --- | --- |
| R1 | Require retention of every selected declarative bundle; raise selection cap from 20 to 30 | 9 / 9 | 9/9 | 10/20 | 8/9 | `ed7bcb83dab4e2d9476cb6cdb176d505766a34f7572f517b2e01507c917aa246` | `143bcdacf730987fb5c87fd9dbe7d8ab71d1169d88a36a757abf9f23c00ab137` | Rejected: broader bundles became fewer and labels became generic; coverage did not improve. |
| R2 | Keep every selected bundle but prohibit context-only resolution in the retrieval text | 10 / 10 | 10/10 | 11/20 | 7/10 | `da35c7af3b4e8500d0d7c51106ac3a5ea57361444bc719eed1fc89b8f9ca9214` | `17ffd84cf9594f27568fdd307b278e26da51af372d21b80d9e11cffce10a1257` | Rejected: small coverage gain, but three retrieval descriptions still imported context outside their stored quote. |

The resulting contradiction is useful: exact evidence validation remains perfect
because it is mechanical, but the quote-first selector is low coverage (best
11/20 ledger claims), while richer descriptions tend to reach beyond an
individual evidence bundle. The authoritative field prevents that expansion
from becoming an asserted fact, but it does not solve retrieval-description
governance.

## Reproduction

Tracked prompts are in
`experiments/fact-prompt-study-6-dual-field/prompts/`; the executed effective
prompt is copied to each ignored iteration directory. For example:

```sh
go run ./cmd/isc-rag deepseek-study \
  --namespace fact-prompt-study-6-dual-field --version v1 \
  --external-id NPR-1964 \
  --model accounts/fireworks/models/deepseek-v4-pro \
  --iteration 5 --approach quote-first --max-facts 20 \
  --prompt-file experiments/fact-prompt-study-6-dual-field/prompts/seed-05-query-oriented.txt \
  --output-dir artifacts/fact-prompt-study-6-dual-field
```

The existing runner persists model, hashes, approach, cap, temperature, and
reasoning setting. It does **not** expose provider token usage, price, or request
latency in its result/schema, so cost and latency are unavailable rather than
reported as zero.

## Honest stop decision

Stop this prompt hill climb. The candidate satisfies an important engineering
invariant—every retained evidence field is exact and authoritative—but it does
not meet a useful coverage bar on this fixed 20-claim ledger. More importantly,
the richer retrieval text repeatedly needs context not present in the saved quote
and may echo question premises. Marking it retrieval-only prevents factual
credit, but does not make it a safe self-contained representation.

Do not call this a paper condition or promote any seed/refinement to the primary
fact condition. A future experiment would need a stored multi-span evidence
bundle with explicit role metadata (for example, question plus endorsed answer),
deterministic bundle boundaries, and a separate retrieval-text policy/audit. It
should then be frozen and evaluated on document-disjoint data before any paper
comparison.
