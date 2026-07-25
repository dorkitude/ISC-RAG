# Study #3 — Claude Fable fact extraction and rephrasing

## Status

**Incomplete; do not promote a paper condition.** Claude Fable ran four
independent configurations on the frozen development document `NPR-1394`,
using fresh `fact_evaluator` agents for the recorded audits. The fourth
configuration cleared the locked gate on that document, but the same frozen
pipeline was not scored on the other frozen validation documents before the
Fable worker exited. Empty run directories were created for those intended
replications (`NPR-573`, `NPR-3311`, and `NPR-2291`); they are not results.

The complete protocol, prompts, and run log are under
`experiments/claude-fable-fact-study-3/`. Source-derived artifacts are kept
only under the ignored `artifacts/claude-fable-fact-study-3/` directory.

## Locked gate

For every document: support at least 98% with zero unsupported cards;
attribution at least 98%; qualifier preservation at least 95%; standalone
clarity and atomicity at least 90%; duplicate rate at most 5%; and ledger
coverage at least 75%. A configuration must meet the gate on the development
document and on at least two document-disjoint frozen replications without
changing the pipeline.

## Audited development results

| Run | Configuration | Support | Attribution | Qualifier | Standalone | Atomic | Coverage | Decision |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| C1 | Direct baseline | 41.4% | 96.6% | 100% | 100% | 93.1% | 47.5% | Reject |
| C2 | Contract plus antecedent window | 90.4% | 100% | 100% | 98.1% | 94.2% | 93.2% | Reject: five unsupported cards |
| C3 | Resolve-then-report | 91.4% | 96.6% | 94.8% | 69.0% | 100% | 72.9% | Reject |
| C4 | Span-blind audit and targeted repair | 100% | 100% | 100% | 98.0% | 98.0% | 94.9% | Development pass only |

C4 kept the C2 generator fixed, mechanically validated its 50 evidence
bundles, performed a span-blind support audit, and repaired or re-anchored
only audited failures. A fresh evaluator then passed all 50 final cards. The
run log and the fresh evaluator worksheet contain the exact audit trail.

## Interpretation

The first three configurations reproduce the central trade-off observed in
Study #2: direct rephrasing is readable but overclaims beyond its evidence;
antecedent windows improve grounding but still permit a small number of
unsupported substitutions; conservative literal handling preserves support
but leaves dialogue deixis and harms standalone clarity. C4 is promising
because it isolates the problem to a post-generation, evidence-bound
audit-and-repair stage rather than asking generation-time prompts to solve all
constraints simultaneously.

The next required action is replication of the exact frozen C4 pipeline on
the three already-frozen document-disjoint validation documents, using a fresh
evaluator for each. Until those runs are complete, C4 is a development finding
and not an experimental condition for ISC-RAG.
