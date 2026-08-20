# ISC-RAG I5c — Pre-registered acceptance criteria for the human calibration study (v3)

*Registered 20 August 2026, before rater recruitment and before any human label exists. Sole registration; supersedes uncommitted drafts of 4 and 18 August 2026.*

## 1. Registration record and disclosure of delay

| Field | Entry |
|---|---|
| Registered | 2026-08-20 (JST) |
| Registered by | Yusuke Takahashi — Associate Professor, Faculty of Data Science, Musashino University / Asia AI Institute / AIx, Inc. PI and corresponding author. |
| Drafting history | Drafts dated 4 August and 18 August 2026 were prepared but never committed. Neither is claimed as a registration. This document is the sole registration and carries its true commit date. |
| State at registration | No rater has been recruited or contracted. No human label exists: the committed response templates worksheet-ratings.csv (80 rows) and practice-ratings.csv (10 rows) contain zero filled judgments, which is machine-checkable. |
| Instrument provenance | The worksheet, practice set, verdict files, response templates, and scripts/kappa.py were committed on 2026-07-26 — before this registration and before any label existed. |
| Prior registration | Complements, and does not supersede, the I5a sampling protocol registered 2026-07-26 in experiments/i5a-calibration-20260726/PROTOCOL.md. |
| Target venue | ECIR 2027. Abstract 21 September 2026, full paper 5 October 2026, both 23:59 GMT. |
| Status | Frozen on commit. See section 11. |

This registration is sixteen days later than first drafted. The delay does not affect its standing, because the condition that gives a pre-registration its value — that it is fixed before any outcome data exists — holds and remains independently checkable. An earlier plan targeting a submission six days away had forced several methodological options to be closed for reasons of time; under the ECIR schedule those options are reachable and are reinstated here (sections 8 and 9).

## 2. What is being calibrated

This study calibrates two automated screens against human judgment: judge_kimi_frontier_xl_v1 for answer correctness (primary outcome) and source_support_kimi_frontier_xl_v1 for source support (secondary outcome).

The instrument is the 80-item blinded worksheet drawn from four cells. The composition below was verified against the raw database on 2026-08-20; it is measured, not asserted.

| Cell | Items | Screen correct / incorrect | Run ID (ledger) |
|---|---|---|---|
| facts @ 2,048 | 20 | 12 / 8 | `a00fb789…` |
| fixed @ 2,048 | 20 | 10 / 10 | `737b9df9…` |
| fixed @ 16,384 | 20 | 10 / 10 | `b13cc9e8…` |
| turn @ 2,048 | 20 | 10 / 10 | `f92ee367…` |

The 12/8 split in the facts cell is not a design choice but an inventory limit recorded in the 2026-07-26 protocol: that cell contains only eight screen-incorrect answers in total, and all eight are included. Pooled marginal balance is 42/38.

### Verification already completed

- Raw-data audit, 2026-08-20: all 56 rows of the committed McNemar tables reproduce exactly from the pg_dump snapshot (SHA-256 verified), and all 56 run IDs referenced by the ledgers exist in the snapshot.
- Calibration-key verification, 2026-08-20: all 90 recorded screen verdicts (80 worksheet, 10 practice) match answers.usage in the database, as do run IDs, question identifiers, representations, and budgets.
> Both checks were completed while zero human labels existed, so no recorded value can have been adjusted to suit an outcome.

## 3. Scope of transfer, and the check that will test it

The worksheet is drawn from frontier-XL cells over the frozen 100-document development sample, and calibrates the frontier-XL judge keys. The paper's headline results are on the held-out 500 sample under the separate keys judge_kimi_heldout500_v1 and source_support_kimi_heldout500_v1.

Calibration therefore transfers to the headline results only if the same judge model and prompt behave equivalently on the held-out corpus. Section 9 specifies the second worksheet that tests this assumption directly. Should that worksheet not be completed, the assumption must be stated as an untested limitation in Methods.

## 4. Raters, adjudicator, and blinding

- Two raters work independently. Eligibility: no affiliation with the ISC research programme, Endgame Labs, Inc., AIx, Inc., or Musashino University, and not a student supervised by any author. Compensated from the research budget.
- This eligibility requirement is not relaxable. It is the basis of the calibration's independence; a rater drawn from the project would make the exercise uninformative whatever kappa resulted.
- Adjudicator: Asako Uraki. The adjudicator sees only items on which the two raters disagree, and only after both have submitted. The adjudicator does not rate the full instrument and is not a third rater, which preserves the independence of the rater-rater statistic.
- Blinding: no participant sees the representation, the token budget, the run ID, or the screen's verdict at any stage. Item order is a fixed deterministic shuffle (seed i5a-20260726), identical for both raters.
- Judgment values are exactly Yes, No, and Unclear. Blanks are not permitted; a rater who cannot decide records Unclear.

## 5. Rater qualification via the practice set

Candidates complete the 10-item practice set before being contracted. Three properties of that set constrain what the screening can mean, and are recorded here.

- practice-key.jsonl contains only the automated screen's verdicts; it is not a human gold standard. Screening candidates against it would select for agreement with the model under calibration and make the study circular. It must not be used for screening.
- A human gold key for the 10 items is therefore authored by Asako Uraki before any candidate sees the practice set, and is committed before screening. The adjudicator thus authors the standard the raters are trained against, so training and adjudication follow one interpretation. This is intentional and is disclosed in the paper.
- All 10 practice items come from a single run in the semantic @ 2,048 cell — a representation deliberately excluded from the main worksheet. The practice set is a training and attentiveness instrument, not a condition-balanced sample, and screening performance may not be reported as an accuracy estimate.
- Eight of the ten practice questions also appear in the main worksheet, with a different generated answer and the same reference answer: fact-dev100x-002, -018, -020, -031, -035, -047, -068, -099. Raters therefore meet those eight reference answers first. The instrument is frozen; this is disclosed as a limitation.

### Pre-specified screening criterion, fixed before any candidate is approached

- A candidate qualifies if, against the human practice key, agreement is at least 8 of 10 on correctness and at least 8 of 10 on source support, with at most 2 items marked Unclear, and a note supplied for every item.
- Candidates who fail are not coached and re-tested; a new candidate is approached. The number screened and the number failed is reported.
- Practice responses are never pooled with worksheet responses and enter no kappa.

## 6. Response instrument and the handling of flags

Responses are recorded on the committed templates worksheet-ratings.csv and practice-ratings.csv, whose columns are: item_id, rater_id, j1_correct, j1_note, j2_supported, j2_note, flag_reference_bad, flag_citation_broken, flag_evasive, flag_needs_discussion, minutes.

Only the first six columns feed scripts/kappa.py. The remaining five are observations the instrument's designer provided for, and their treatment is fixed here so that it cannot be chosen after the flags are seen.

- Flags are descriptive. They do not by themselves remove any item from kRR or from kCJ. The primary figures include every item.
- Sensitivity analysis, pre-specified: kappa is additionally recomputed with items excluded where BOTH raters set flag_reference_bad. Both the primary and the sensitivity figure are reported, whichever is more favourable. A divergence between them is itself a finding about the reference answers and is reported as such.
- flag_citation_broken, flag_evasive, and flag_needs_discussion are reported as counts by cell. They trigger no exclusion.
- minutes is reported as the median per rater, describing effort only. It never weights, filters, or excludes a judgment.
- Raters are not told that any flag has analytic consequences, so that flagging is not strategic.
- item_id must remain a three-digit zero-padded string (001–080); kappa.py matches on that form. Spreadsheet software that coerces it to an integer will silently break the join, and the returned files are checked for this before analysis.

## 7. Statistics

Both figures are produced by scripts/kappa.py, run unmodified.

- kRR — Cohen's kappa between rater A and rater B, before adjudication, over three categories, computed separately for j1_correct and j2_supported. This is the evidence that the human standard is itself reliable.
- kCJ — Cohen's kappa between the post-adjudication consensus and the automated screen, binary. Items whose consensus is Unclear are excluded, because the screen has no Unclear category. This is the calibration figure.
Raw agreement counts are reported alongside every kappa.

## 8. Acceptance criteria

Gate A is evaluated and reported before Gate B. Under the ECIR schedule a failed gate can be answered by repetition, and that is the specified response.

### Gate A — rater reliability (precondition), on j1_correct

| kRR | Outcome |
|---|---|
| at least 0.60 | Pass. Proceed to Gate B. |
| below 0.60 | Fail. kCJ is not interpretable as calibration and must not be reported as such. Required action: publish kRR, revise the rating rubric — never the sample — recruit two new raters, and repeat on the same 80 items. The failed round is disclosed in the paper. At most one repetition is performed; if the second round also fails, the paper reports both kRR values and treats the screen as uncalibrated. |

Rationale: if two trained humans do not agree with each other, there is no human standard against which an automated screen can be calibrated. The 0.60 boundary is the conventional threshold for substantial agreement (Landis and Koch, 1977).

### Gate B — judge calibration, on j1_correct (primary outcome)

| kCJ | Verdict | Consequence for the paper |
|---|---|---|
| at least 0.60 | Accept | The screen-judged accuracy figures stand as primary evidence. kappa is reported in Methods. |
| 0.40 to below 0.60 | Conditional | Screen-judged figures may stand, but must be reported together with the human-verdict accuracy on the 80-item sample, and the paper must state that the screen is only moderately calibrated. P3 is not claimed on the automated screen alone. |
| below 0.40 | Reject | The screen cannot serve as the primary outcome measure. Required action: expand human rating to a sample large enough to be the primary outcome, or withdraw answer accuracy as a primary claim and report retrieval and cost as the primary outcomes. Which of the two is taken is decided on feasibility at that date and recorded in the correction log. |

j2_supported is reported on the same scale but is a secondary outcome; on its own it does not trigger the Gate B consequences.

## 9. Confirmatory work reinstated under the ECIR schedule

The following two items were closed in an earlier draft solely because a six-day schedule could not reach them. They are reinstated as planned work, and their pre-specification is fixed here.

### 9.1 Second worksheet drawn from held-out cells

- Purpose: to test the transfer assumption declared in section 3 rather than merely declaring it.
- Construction: the sampling rule registered on 2026-07-26 is applied unchanged to the held-out 500 cells facts @ 2,048 and ctx-full-stack @ 2,048, under the keys judge_kimi_heldout500_v1 and source_support_kimi_heldout500_v1. The deterministic shuffle uses seed `i5a-heldout-<commit date>`.
- Raters: the same two raters, after the main worksheet is submitted and adjudicated, so that the two instruments cannot contaminate each other.
- Reporting: kCJ is reported for both instruments. The two are compared descriptively; no threshold is attached to their difference, because this study is not powered to test one. If the second worksheet is not completed, section 3's limitation stands as written.

### 9.2 Re-judging the core held-out cells with a non-Kimi model

- Purpose: to address the extractor-judge overlap of section 10 by a direct method rather than by discounting.
- Design: facts @ 2,048 and ctx-full-stack @ 2,048 on the held-out 500 sample are re-judged with a judge model from a different family, chosen and recorded before the re-judging is run. Answer text, prompts, and retrieval payloads are unchanged; only the judge differs.
- Reporting: both judges' results are reported side by side. The paper's accuracy claim is stated under the non-Kimi judge where the two disagree in direction, and the disagreement is reported.

## 10. Pre-registered confound analysis: extractor-judge model overlap

The fact condition's retrieval payload was written by Kimi K2.6, and the automated screen is also Kimi K2.6. The chunk conditions' payloads are raw transcript text and receive no comparable benefit. The project's cross-judging rule, that no model judges its own answers, constrains the answer model only and does not address this overlap.

Directional hypothesis, registered before any label exists: if the overlap biases the screen, the screen will be more generous on the facts cell than on the chunk cells.

### Definitions, applied per item after adjudication

- Screen-generous: the screen says correct and the consensus says No.
- Screen-severe: the screen says incorrect and the consensus says Yes.

### Analysis

- Report the two-by-two table of screen against consensus separately for the facts cell (n = 20) and the pooled chunk cells (n = 60).
- Compare the screen-generous rate between facts and pooled chunks with a two-sided Fisher exact test.
> Power is declared low in advance. With 20 items against 60, only a large difference is detectable. This analysis is descriptive and directional, not confirmatory; section 9.2 carries the confirmatory burden.

### Pre-specified interpretation, fixed now so that it cannot be selected after the numbers are seen

- If the screen-generous rate on facts is not greater than on chunks: report as no evidence of overlap bias at this sample size, with the power limitation. A non-significant result may not be reported as evidence that the confound is absent.
- If the screen-generous rate on facts is greater than on chunks, at any p-value: report it as a possible confound. The primary remedy is the non-Kimi re-judging of section 9.2. If that has not been completed, discount the facts-versus-chunk accuracy gap by the observed generosity difference and say so.

## 11. Frozen elements

The following are fixed at registration and may not be changed once any human label has been seen: the thresholds and consequences in section 8; the screening criterion in section 5; the flag-handling rules in section 6; the number of raters; the adjudication rule; the exclusion of Unclear from kCJ; the definitions and interpretation rules in section 10; the pre-specifications in section 9; and the 80-item instrument together with its order.

Any deviation is entered in the protocol correction log with the date, the reason, and the person who decided, and is disclosed in the paper. This mirrors the correction-log rule already recorded for the held-out protocol.

## 12. Known implementation gaps

- scripts/kappa.py computes kRR and kCJ but does not produce the per-representation breakdown that section 10 requires. A supplementary script must be committed before any consensus CSV is produced, so that it cannot be tuned to the outcome. Until it exists, section 10 cannot be executed.
- No human gold key exists for the practice set. Section 5 cannot be executed until it is authored and committed.
- The second worksheet of section 9.1 does not yet exist and must be generated by the same committed sampling code, not by hand.

## 13. Schedule and owners

| Date | Item | Owner |
|---|---|---|
| 2026-08-20 | This document committed; audit record committed | Takahashi |
| 2026-08-26 | Human gold key for the 10 practice items | Uraki |
| 2026-08-28 | Ethics review filed with Musashino University | Takahashi |
| 2026-08-31 | Section 10 supplementary script committed | Takahashi |
| 2026-09-02 | Two external raters screened and contracted | Takahashi |
| 2026-09-08 – 09-12 | Main worksheet completed | Raters |
| 2026-09-14 | kRR computed; disagreement list generated | Takahashi |
| 2026-09-16 | Adjudication | Uraki |
| 2026-09-17 | kCJ; section 10 analysis; section 8 applied | Takahashi |
| 2026-09-18 | Second worksheet (section 9.1) | Takahashi, Raters |
| 2026-09-20 | Non-Kimi re-judging (section 9.2) | Wild |
| 2026-09-21 | ECIR abstract (23:59 GMT) | Takahashi |
| 2026-10-05 | ECIR full paper (23:59 GMT) | Takahashi |

> Both ECIR deadlines are stated in GMT, not anywhere-on-earth. 5 October 23:59 GMT is 6 October 08:59 JST.
