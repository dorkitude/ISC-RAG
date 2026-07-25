# Study #3 run log

All extraction runs, in execution order. Metrics come from the fresh
`fact_evaluator` verdict for that run (see
`artifacts/claude-fable-fact-study-3/runs/<run id>/evaluation.txt`). Raw = fact
lines returned by the extractor; Valid = facts surviving deterministic
provenance validation. Rates are over validated facts; coverage is over the
frozen ledger for that document.

| Run | Config | Doc | Raw | Valid | Support | Attrib | Qualif | Standalone | Atomic | Dup | Coverage | Notes |
| --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| r01 | C1 direct baseline | NPR-1394 | 29 | 29 | 41.4% (17 unsup) | 96.6% | 100% | 100% | 93.1% | 0% | 47.5% (28/59) | FAIL. Single short quote lacks the named entity the standalone fact asserts → 17 support fails. Standalone bought at cost of support. Coverage low. |

| r02 | C2 contract+antecedent | NPR-1394 | 52 | 52 | 90.4% (5 unsup) | 100% | 100% | 98.1% | 94.2% | 0% | 93.2% (55/59) | Big jump. Passes all criteria EXCEPT support. 5 support fails = decontextualized name/descriptor not in the cited spans (over-reach). Antecedent window fixed most C1 support fails. |

| r03 | C3 resolve-then-report | NPR-1394 | 58 | 58 | 91.4% (5 unsup) | 96.6% | 94.8% | 69.0% | 100% | 0% | 72.9% (43/59) | FAIL. Regression vs r02 on every failing axis. 18 standalone fails: the Pass-A referent ledger was under-applied — extractor kept source-literal deictic phrases instead of substituting canonicals (the C3 prompt's "prefer dropping to guessing" rule pushed it to keep literals, tanking standalone). 2 attribution+qualifier fails from facts lifted out of reported speech and credited to the quoting speaker. Coverage fell below gate. Evaluator: fresh fact_evaluator (agentId aaf0672ee744fc970), blind, opaque id run-x7; verdict at `artifacts/claude-fable-fact-study-3/runs/r03-c3-npr1394/evaluation.txt` (verbatim). |

| r04 | C4 span-blind audit+repair (dev) | NPR-1394 | 50 | 50 | 100% (0 unsup) | 100% | 100% | 98.0% | 98.0% | 0% | 94.9% (56/59) | PASS on dev doc — first configuration to clear every criterion. Pipeline: C2 generator byte-identical → mechanical validation (50/50 kept) → span-blind audit (10 FAIL: 2 topic-word over-reach, 1 case-descriptor, 7 strict title-entailment) → targeted repair (3 REANCHOR, 1 antecedent re-cut, 6 REWRITE to quoted name form; 0 drops) → mechanical re-validation (10/10) → fresh blind re-audit (10/10 PASS) → mechanical assembly (50 rows, 0 re-check drops). Evaluator: fresh fact_evaluator (agentId a3a442e64ebe263c0), blind, opaque id run-k2; verdict at `runs/r04-c4-npr1394/evaluation.txt`. C4 frozen at this state for gate replicates. |

Gate: support ≥98% & 0 unsupported; attribution ≥98%; qualifiers ≥95%; standalone ≥90%; atomicity ≥90%; dup ≤5%; coverage ≥75%/doc.

Read of r01→r02: the structural conflict is real. C1 makes facts standalone but the lone quote lacks the named entity (support 41%). C2's antecedent window grounds most names (support 90%, coverage 93%) but ~10% of facts still substitute a referent not present in the cited spans. Remaining lever: force every substituted referent to be backed by a cited anchor span (C3), and/or add a support-repair pass that drops/reliterals any fact whose content isn't covered by its spans.

Read of r02→r03: C3 traded the wrong way. Making substitution conditional on a groundable anchor ("prefer dropping to guessing") caused the extractor to keep source-literal deictics, so standalone collapsed to 69% while support only reached 91% — worse than C2 on both axes it was meant to fix, plus new attribution leaks from reported speech. C2 remains the best configuration (fails only support, 5/52). The residual problem is narrow and asymmetric: ~10% of C2 facts assert one referent/descriptor not present in their cited spans, and C3 showed generation-time constraints degrade the generator. Next lever should therefore be a post-generation structural stage (verify/repair or filter on the support axis specifically) rather than another generation-time contract change.
