# Ingest-Time Semantic Compilation: Accuracy, Provenance, and Token Savings in Dialogue RAG

**File:** `20260725_#TBD_ISCRAG_@Wild_EN_v4.md` (paper_template v0 structure, Markdown for Google Docs)
**Version:** v4 — 2026-07-25

---

## Part 1 — Initial Question Sheet

| # | Question | Entry |
|---|---|---|
| Q1 | Working title | Ingest-Time Semantic Compilation: Accuracy, Provenance, and Token Savings in Dialogue RAG |
| Q2 | Proposal type | **D — Empirical study** (controlled comparison of retrieval representations) |
| Q3 | Research phase | **P1 (position/pilot)** — development evidence only; upgrades to P3 after human calibration + document-disjoint replication |
| Q4 | Venue | ⚠️ **#TBD** — candidates: SIGDIAL (dialogue fit), EMNLP/ACL Findings (retrieval-granularity lineage: Dense X, PropRAG), ECIR/SIGIR short (IR framing). **Rationale needed** — see fulfillment question FQ-1 |
| Q5 | Skeleton | IMRAD (type D → always IMRAD) |
| Q6 | Language | EN |
| Q7 | Method name | **ISC-RAG** (Ingest-time Semantic Compilation for RAG); core distinction: **ISC vs. QSR** (query-time semantic reconstruction) |
| Q8 | Core claim (one sentence) | Compiling dialogue into provenance-validated, self-contained facts at ingest lets a RAG system answer more questions from substantially less read-time context than transcript chunking, at lower per-query token cost. |
| Q9 | Authors & affiliations | Kyle Wild (first author; ⚠️ affiliation/ORCID — FQ-2); Yusuke Takahashi — Faculty of Data Science, Asia AI Institute, Musashino University / AIx, Inc., Tokyo, Japan — https://orcid.org/0009-0006-8351-2280 |
| Q10 | Deadline & page limit | ⚠️ 〔 〕 — FQ-1 |
| Q11 | Attached material | `docs/preliminary-findings-and-abstract-20260725.md`; `docs/fact-development-answer-frontier-20260725.md`; `docs/fact-development-top-k-depth-full-20260725.md`; `docs/fact-development-frontier-xl-20260725.md`; `docs/literature-landscape.md` (33 verified refs); `docs/literature-deep-dive-isc-cost-20260725.md` (27 verified items); repo `README.md` status log |
| Q12 | Version & date | v4 / 2026-07-25 |

---

## Part 2 — Literature Review

Verified in `docs/literature-landscape.md` (33 entries, landing-page checked) and `docs/literature-deep-dive-isc-cost-20260725.md` (27 entries, arXiv-fetch verified). The subset below is carried into References; the full lists remain the working bibliography. Quality: >90% from 2020+, arXiv items marked preprint, zero unverified URLs.

| # | Work | Role in this paper | Verdict |
|---|---|---|---|
| 1 | Chen et al. 2024, Dense X Retrieval, EMNLP | Nearest predecessor: proposition units beat passages on Wikipedia | OK |
| 2 | Anthropic 2024, Contextual Retrieval (blog) | Ingest-time chunk contextualization — compiles the pointer, not the payload | OK |
| 3 | Hou et al. 2026, AtomicRAG (arXiv preprint) | Recent chunk-dilution argument; atoms as graph nodes, no provenance | OK |
| 4 | Wang 2025, PropRAG, EMNLP | Proposition paths for multi-hop; prompted extraction | OK |
| 5 | Zhang et al. 2024, SiReRAG, ICLR 2025 | Propositions as multihop index scaffolding | OK |
| 6 | Pollertlam & Kornsuwannawit 2026, Beyond the Context Window (arXiv preprint) | Only explicit ingest-amortization break-even analysis | OK |
| 7 | Chen et al. 2025, HaluMem (arXiv preprint) | Fact extraction hallucinates → motivates the validation gate | OK |
| 8 | Bao & Shi 2026, IndexRAG (arXiv preprint) | Nearest neighbor on index-time reasoning | OK |
| 9 | Taranukhin et al. 2026 〔verify authors〕, Token Tax (arXiv preprint) | Cost–quality frontier framing, chunk-level only | ⚠️ verify author list |
| 10 | Microsoft 2024, LazyGraphRAG (blog) | The counter-position: defer semantics to query time | OK |
| 11 | Lewis et al. 2020, RAG, NeurIPS | Canonical retrieve-then-generate (QSR) | OK |
| 12 | Karpukhin et al. 2020, DPR, EMNLP | Dense retrieval baseline lineage | OK |
| 13 | Choi et al. 2021, Decontextualization, TACL | Making sentences stand-alone: the rewriting contract | OK |
| 14 | Zhu et al. 2021, MediaSum, NAACL | Corpus | OK |
| 15 | Min et al. 2023, FActScore, EMNLP | Atomic-fact evaluation precedent | OK |
| 16 | Chhikara et al. 2025 〔verify authors〕, Mem0 (arXiv preprint) | Industry fact-extraction memory + 90% token-savings claim, no provenance | ⚠️ verify author list |
| 17 | Jiang et al. 2023 / PISCO 2025, prompt/context compression (arXiv preprints) | Query-time token reduction — lossy, provenance-destroying contrast | OK |
| 18 | Zahn & Chana 2026, Facts as First-Class Objects (arXiv preprint) | Symbolic knowledge-object memory; conceptual sibling | OK |

---

## Part 3 — Paper Draft (v4)

### Abstract

*(7-sentence closed-world structure; sentences numbered for review only — remove numbering in submission version.)*

(1) When people answer a question about a conversation they took part in, they do not replay the recording; they consult remembered, already-interpreted claims — who said what, about whom, resolved into stand-alone meaning at the time of hearing. (2) Retrieval-augmented generation systems do not yet work this way: they store raw transcript chunks and force the answer model to re-derive speakers, antecedents, and claims from unprocessed dialogue at every query. (3) Making retrieval consult compiled, source-faithful claims rather than raw text — without sacrificing auditability — is therefore an important open problem for question answering over human dialogue. (4) We present **ISC-RAG**, an ingest-time semantic compilation pipeline that extracts factual and opinion claims from dialogue transcripts, validates every claim against an exact source quotation, attaches speaker and turn provenance, and rewrites claims so resolvable pronouns and local references name their antecedents. (5) The defining feature of ISC-RAG is that the compiled fact is the *retrieval payload itself* — the unit the answer model reads — rather than an index key pointing back into raw text, so the semantic work of interpretation is paid once at ingest instead of on every query. (6) This design makes it possible to answer dialogue questions from a fraction of the read-time context that transcript chunking requires — shrinking the input surface on which long-context degradation and distraction operate — while returning a verbatim supporting quote with every claim. (7) In controlled development experiments on 100 frozen MediaSum interview transcripts with 200 source-grounded questions, we compare facts against three chunking controls under identical embedding, retrieval, and answer models: facts reached the gold source turn at top-1 for 165/200 questions versus 134/200 for fixed chunks, and at a 512-token read budget produced 81/100 blind-judged correct answers versus 45–46/100 for all chunk controls — with 256-token facts exceeding every 512-token chunk condition while consuming 38–41% fewer answer-model tokens. In a confound-free extension where top-300 retrieval makes the budget the binding constraint at 2,048–16,384 tokens, facts at 2,048 tokens are 92/100 correct — above every chunk condition at every budget, whose ceiling is 75/100 at 7.3× the token cost — while fact accuracy *declines* at larger budgets, indicating that compiled context is subject to long-context degradation too and simply reaches its optimum at radically lower cost. These results support the feasibility of ingest-time compilation as a favorable accuracy-per-token frontier, pending human calibration and held-out replication.

**Keywords:** retrieval-augmented generation; retrieval granularity; dialogue; provenance; decontextualization; token efficiency

### 1. Introduction

Answering questions about conversations is a routine human competence. A person who attended a meeting does not re-read a transcript to recall who committed to a deadline; during the conversation itself they resolved pronouns, attributed claims to speakers, and stored compact, stand-alone facts. The interpretive work was paid once, at *ingest*.

Retrieval-augmented generation (RAG) systems for dialogue do not yet realize this pattern. The dominant design — which we call **query-time semantic reconstruction (QSR)** — stores raw transcript chunks, retrieves them by embedding similarity, and asks the answer model to reconstruct meaning from unprocessed dialogue on every query: to find the relevant claim inside a multi-turn window, distinguish the respondent from the questioner, and resolve "she" or "that decision" across turn boundaries. This re-paid interpretive work has three costs. First, *accuracy*: a specific claim competes with surrounding conversational material both at retrieval time (the chunk embedding dilutes it) and at reading time (the answer model must locate it). Second, *auditability*: nothing in the pipeline distinguishes what a speaker actually said from what the model inferred. Third, *economics*: every query pays for context tokens whose meaning was already available at ingest — a growing concern as agentic workloads inflate read-time context. The first and third costs compound, because context is not accuracy-neutral: model performance grows unreliable as input length increases even on simple tasks [20], relevant information is lost in the middle of long inputs [19], and irrelevant material measurably distracts reasoning [21]. Shrinking the read payload is therefore a reliability intervention, not only a cost one. Prior work has attacked pieces of this problem — finer index units [1], ingest-time chunk contextualization [2], query-time context compression [17] — but in each case the payload the answer model reads remains raw or unverified text. Realizing human-like *compiled* consultation of dialogue, with provenance strong enough to trust, remains an important open problem.

In this paper we present **ISC-RAG** (Ingest-time Semantic Compilation for RAG). ISC-RAG extracts high-recall factual and opinion claims from dialogue transcripts, discards any claim whose exact supporting quotation cannot be found verbatim in the source, attaches speaker attribution and turn-span provenance to every surviving claim, and rewrites claims so that resolvable pronouns and local references name their antecedents. The compiled fact — claim plus verified quote — is the retrieval payload the answer model reads.

Our contributions are:

1. **A provenance-gated fact compilation pipeline for dialogue** in which every indexed unit carries an exact-match-validated source quotation, speaker attribution, and turn span (Section 3).
2. **A controlled representation comparison** on a frozen 100-document broadcast-interview corpus, holding embedding model, retriever, answer model, and question set constant across facts and three chunking controls (Section 4).
3. **An accuracy-per-token frontier measurement**: answer correctness at four matched read-time context budgets (256–2,048 tokens) with provider-reported token usage, showing development evidence that compiled facts dominate the frontier (Section 4.3).
4. **A distinction between compiled payloads and compiled pointers** that organizes prior work (propositions, contextual chunking, graph indexes) and identifies the unclaimed combination: validated provenance + fact-as-payload + budget-matched cost accounting (Sections 2, 5).

Section 2 positions ISC-RAG against related work. Section 3 describes the compilation pipeline and data model. Section 4 reports the controlled development experiments. Section 5 discusses limitations and the amortization argument. Section 6 concludes.

### 2. Related Work

We organize prior work along two axes: **(Axis 1) where the semantic work is paid** — at query time or at ingest — and **(Axis 2) what the answer model reads** — raw text or compiled units. Existing systems occupy three of the four cells; ISC-RAG targets the fourth (ingest-paid, compiled-read), with provenance validation as the property that makes the fourth cell trustworthy.

**Query-paid, raw-read (QSR).** Canonical RAG [11] with dense retrieval [12] retrieves passages and reconstructs meaning per query. Query-time context compression — LLMLingua-family token pruning, RECOMP summaries, and soft-compression successors such as PISCO [17] — reduces the token bill but re-pays compression on every query and destroys verbatim provenance. A separate measurement literature establishes that the raw-context tax is paid in accuracy as well as tokens: models use long contexts non-uniformly, with performance degrading as inputs grow ("context rot" [20]) and mid-context information under-used [19], and irrelevant context measurably harms reasoning [21]. This makes read-payload size a reliability parameter of RAG systems, one that unit choice directly controls.

**Ingest-paid, raw-read (compiled pointers).** Dense X Retrieval [1] showed that indexing Wikipedia by atomic propositions beats passage indexing, especially under fixed reading budgets; its successors extend propositions into multihop graph scaffolding (SiReRAG [5], PropRAG [4], AtomicRAG [3] — the last arguing, as we do, that chunks "rigidly group multiple atomic facts together"). Anthropic's Contextual Retrieval [2] prepends LLM-generated situating context to each chunk before embedding and BM25 indexing. IndexRAG [8] pre-generates cross-document bridging facts at index time. In all of these, the ingest-time product improves *where retrieval lands*; the reader still consumes raw or enlarged text, and no system validates its generated units against source spans. LazyGraphRAG [10] argues the counter-position — defer semantic work to query time because indexing spend is often wasted — making query volume per corpus the pivotal, unformalized variable (Section 5.3).

**Ingest-paid, compiled-read, unvalidated.** Agent-memory systems extract facts from dialogue at ingest and serve them as context: Mem0 [16] reports >90% token savings; a recent cost-performance analysis [6] gives the only explicit break-even result (fact memory undercuts long-context after ~10 turns). But HaluMem [7] demonstrates that memory-system fact extraction hallucinates and that these errors propagate into QA — and a 2026 survey of memory architectures finds no system that validates the provenance of extracted facts. The Token Tax analysis [9] measures the cost–quality frontier we target, but with chunk-level arms only.

**This work.** ISC-RAG is, to our knowledge, the first to combine (i) atomic claims as the *read payload*, (ii) exact-quote provenance validation answering the failure mode HaluMem measures, and (iii) budget-matched accuracy-per-token evaluation — on dialogue, where decontextualization pressure (speaker attribution, cross-turn coreference, opinions) is strongest [13, 15]. Relative to our own prior framing of ISC as substrate maintenance 〔cite ISC paper draft — FQ-3〕, this paper contributes the representation-level empirical comparison.

### 3. Method: the ISC-RAG compilation pipeline

*(Type D paper: methods here describe the studied artifact and data model; experimental design follows in Section 4.)*

**3.1 Data model.** The transcript is immutable; retrieval views are versioned projections over it. Canonical objects: `documents` (one transcript), `turns` (ordered utterances with corpus-provided speaker labels and character offsets), `source_units` (one retrievable unit per representation/version), `embeddings` (one vector per unit and model), and `facts` + `fact_evidence` (one standalone claim with exact quote, quote span, turn span, and speaker label). Every representation searches the same `source_units` path with the same query embedding and pgvector index, so conditions differ *only* in what a unit contains.

**3.2 Compilation passes.** Pass 1 (*harvest*): for every turn, an extraction model proposes explicit factual assertions **and** opinions, beliefs, judgments, preferences, and predictions, each with a short verbatim supporting quotation. Pass 2 (*validation gate*): a candidate is discarded unless its quotation is found exactly in the canonical transcript; surviving claims store the quote span and turn span. Pass 3 (*contextual rewrite*): each retained claim is rewritten using local transcript context so that resolvable pronouns and references name their antecedents (following the decontextualization contract of [13]); rewriting is conservative — an antecedent is substituted only when knowable, and the verbatim quote is never altered. Speaker attribution uses deterministic parsing of corpus-supplied labels (e.g., `ALEX COHEN, host`), never inference from dialogue content.

**3.3 Retrieval projection.** The embedded representation of a fact is `Fact: <claim>\n\nSupporting quote: <quote>`; retrieval returns both fields. Chunk controls embed and return raw transcript windows: fixed windows (4,000 bytes, 400-byte overlap), turn-aware windows, and lexical-cohesion semantic windows.

**3.4 Research questions.**
- **RQ1 (retrieval):** Does a compiled fact index reach the gold source turn at top-k more often than matched chunk indexes, under identical query embeddings?
- **RQ2 (answering at matched budgets):** At equal read-time context budgets, does an answer model produce more correct, better-supported answers from facts than from chunks?
- **RQ3 (frontier):** Can a smaller fact budget match or exceed larger chunk budgets at lower provider-reported token cost?

### 4. Experiments

**4.1 Setup.** *Corpus:* a frozen, deterministically sampled 100-document subset of MediaSum [14] (NPR/CNN interview transcripts; full local corpus 463,596 documents, 13.9M turns, all stream-audited through the ingestion schema). *Fact index:* 13,064 validated fact-plus-quote units from the frozen sample (`facts/turnwise-opinion-v1`). *Controls:* fixed, turn-aware, and semantic chunk indexes built from exactly the same 100 documents. *Questions:* two independently authored 100-question sets (one per RQ round), each question source-grounded with a gold document, turn range, and verbatim supporting quote, all mechanically verified against the database before use. *Models:* one Fireworks-hosted embedding model (Nomic Embed Text v1.5) for all conditions; pgvector top-20 search; DeepSeek V4 Flash as the answer model with mandatory structured citations; two blind Kimi K2.6 sessions as development judges (semantic correctness; cited-source support). Judges never see the representation. 〔E2 settings table with exact prompts, temperatures, and packing algorithm — drafted from `fact-development-answer-frontier-20260725.md`; FQ-4〕

**4.2 RQ1: Retrieval.** Facts retrieved a gold source turn at top-1 for 165/200 pooled questions versus 134/200 for fixed chunks (82.5% vs 67.0%; per-set 83%/70% and 82%/64%; top-5 95%/87% and 95%/88%). The gold-turn-overlap metric is deliberately conservative in interpretation: a fact from the correct turn can be credited even when it states a different claim from that turn; fact-level evidence links are future work.

**4.3 RQ2/RQ3: Answer accuracy and token frontier.** All 16 conditions (4 representations × 4 budgets) use identical retrieval and answering machinery; retrieved sources are packed greedily to 256/512/1,024/2,048 estimated context tokens; Fireworks reports actual tokens.

| Representation | Budget | Correct | Cited support | Mean total tokens |
|---|---:|---:|---:|---:|
| Facts + quote | 256 | 74/100 | 77/100 | 389.8 |
| Fixed chunks | 256 | 23/100 | 24/100 | 419.3 |
| Turn-aware chunks | 256 | 21/100 | 23/100 | 420.3 |
| Semantic chunks | 256 | 28/100 | 29/100 | 404.9 |
| Facts + quote | 512 | 81/100 | 83/100 | 654.4 |
| Fixed chunks | 512 | 46/100 | 44/100 | 663.6 |
| Turn-aware chunks | 512 | 45/100 | 46/100 | 663.9 |
| Semantic chunks | 512 | 45/100 | 44/100 | 631.2 |
| Facts + quote | 1,024 | 87/100 | 86/100 | 1,158.8 |
| Fixed chunks | 1,024 | 60/100 | 59/100 | 906.1 |
| Turn-aware chunks | 1,024 | 56/100 | 57/100 | 921.1 |
| Semantic chunks | 1,024 | 61/100 | 61/100 | 870.1 |
| Facts + quote | 2,048 | 86/100 | 87/100 | 1,407.5 |
| Fixed chunks | 2,048 | 64/100 | 64/100 | 1,801.3 |
| Turn-aware chunks | 2,048 | 59/100 | 62/100 | 1,806.5 |
| Semantic chunks | 2,048 | 63/100 | 63/100 | 1,901.1 |

At the matched 512-token budget, facts produced 81/100 correct answers versus 45–46/100 for every chunk control at nearly identical actual cost. On the frontier reading, 256-token facts (74/100) exceeded every 512-token chunk condition while consuming 38–41% fewer provider-reported tokens, and 1,024-token facts (87/100) exceeded every 2,048-token chunk condition (59–64/100) at roughly two-thirds the cost.

**4.4 Extended frontier: large budgets with supply unconstrained.** The frontier above retrieved only the top 20 units before packing — enough for chunks, but 20 facts total only ~1,300 estimated tokens, so the fact condition's larger budgets were candidate-starved rather than budget-bound. A confound-free extension therefore retrieves the top 300 units for every condition — 300 facts ≈ 19k estimated tokens, 300 fixed chunks ≈ 210k, so greedy packing truncation binds at every cell — across budgets of 2,048/4,096/8,192/16,384 tokens, far below the answer model's 1M-token window. Judging uses a reproducible CLI command with frozen blind Kimi prompts.

| Representation | Budget | Correct | Corr+Supp | Mean total tokens |
|---|---:|---:|---:|---:|
| Facts + quote | 2,048 | 92/100 | 86/100 | 2,233.3 |
| Facts + quote | 4,096 | 87/100 | 83/100 | 4,332.6 |
| Facts + quote | 8,192 | 88/100 | 81/100 | 8,445.9 |
| Facts + quote | 16,384 | 86/100 | 81/100 | 13,094.1 |
| Fixed chunks | 2,048 | 62/100 | 59/100 | 1,801.2 |
| Fixed chunks | 4,096 | 70/100 | 67/100 | 3,946.8 |
| Fixed chunks | 8,192 | 68/100 | 66/100 | 8,002.2 |
| Fixed chunks | 16,384 | 75/100 | 73/100 | 16,243.4 |
| Turn-aware chunks | 2,048 | 59/100 | 55/100 | 1,804.9 |
| Turn-aware chunks | 4,096 | 59/100 | 55/100 | 3,872.7 |
| Turn-aware chunks | 8,192 | 64/100 | 61/100 | 8,003.9 |
| Turn-aware chunks | 16,384 | 75/100 | 71/100 | 16,271.2 |
| Semantic chunks | 2,048 | 63/100 | 63/100 | 1,901.2 |
| Semantic chunks | 4,096 | 61/100 | 61/100 | 3,964.0 |
| Semantic chunks | 8,192 | 65/100 | 65/100 | 8,049.9 |
| Semantic chunks | 16,384 | 70/100 | 70/100 | 16,361.4 |

Decomposing each cell into retrieval coverage (was any gold-overlapping source packed?) and reading accuracy (correctness among only those questions) isolates the representation effect from a needle-in-a-haystack account:

| Representation | Gold packed @2k/4k/8k/16k | Read-acc given gold @2k/4k/8k/16k |
|---|---|---|
| Facts + quote | 100 / 100 / 100 / 100 | 92.0 / 87.0 / 88.0 / 86.0 |
| Fixed chunks | 80 / 88 / 90 / 95 | 76.3 / 78.4 / 75.6 / 78.9 |
| Turn-aware chunks | 75 / 82 / 91 / 95 | 77.3 / 72.0 / 70.3 / 78.9 |
| Semantic chunks | 78 / 85 / 88 / 95 | 78.2 / 70.6 / 72.7 / 73.7 |

If large-budget chunk cells merely measured the answer model's haystack weakness, reading accuracy given packed gold would collapse with budget; instead it is roughly flat for fixed chunks (~76–79%) from 2,048 to 16,384 tokens, and chunk budget gains are mostly coverage gains (80%→95%). The dominant effect is a ~14–16-point reading-accuracy gap between compiled facts and raw transcript that is already present at 2,048 tokens, where haystack pressure is minimal; long-context degradation appears as a secondary effect on top (the fact conditions' own decline occurs at fixed 100% coverage and is therefore pure reading degradation with volume).

Three results stand out. First, the 2,048-token fact cell (92/100 correct at 2,233 actual tokens) exceeds every chunk cell at every budget; the best any chunk condition achieves is 75/100 at roughly 16,300 tokens — 7.3× the cost for 17 fewer correct answers — and facts' worst cell (86/100) still exceeds chunks' best. Second, with supply unconstrained, fact accuracy *declines* as budget grows (92→87→88→86; paired 2,048→16,384: 7 regressions vs 2 gains), so compiled context is itself subject to long-context degradation — its optimum simply sits near 2k tokens. Third, chunks convert budget into accuracy (turn-aware gains 20 paired questions from 2,048 to 16,384) but non-monotonically (fixed dips at 8,192), and a qualitative audit verified individual correct→insufficient flips with the gold evidence still packed. The prescription the table supports is not "facts tolerate large contexts better" but "facts make large contexts unnecessary."

**4.5 Depth curve: coverage and cost at uncapped top-k.** To separate the budget manipulation from retrieval depth itself, a second controlled sweep retrieves exactly the top k ∈ {1, 3, 5, 10, 20, 30} units per question with *no packing cap*, passing everything retrieved to the answer model (24 conditions; 2,400 answers; 4,800 independent judgments; DeepSeek V4 Flash as both answerer and — a deliberate development shortcut flagged in §5.4 — automated judge). A question counts only when *both* screens pass: semantic correctness and direct source support.

| Representation | k=1 | k=3 | k=5 | k=10 | k=20 | k=30 |
|---|---:|---:|---:|---:|---:|---:|
| Facts + quote | 47/100 | 68/100 | 72/100 | 76/100 | 79/100 | 79/100 |
| Fixed chunks | 49/100 | 60/100 | 62/100 | 62/100 | 66/100 | 61/100 |
| Turn-aware chunks | 48/100 | 53/100 | 54/100 | 63/100 | 63/100 | 69/100 |
| Semantic chunks | 49/100 | 53/100 | 57/100 | 58/100 | 59/100 | 65/100 |

At top-1, chunks are marginally ahead (48–49 vs 47) — but not cost-neutrally, since one chunk is many times larger than one fact: facts reach 47/100 at a mean 233 total tokens versus 729–903 for chunks. From k=3 upward facts lead at every depth: 68/100 at 366 tokens versus 53–60/100 at 1,892–2,348 tokens (k=3), and 79/100 at 1,414 tokens versus 59–66/100 at 11,697–14,465 tokens (k=20) — roughly 8–10× cheaper while answering 13–20 more questions. Facts saturate at 79/100 by k=20; fixed chunks *fall* from 66/100 at k=20 to 61/100 at k=30, additional raw context distracting the answer model rather than helping. A per-question minimum-successful-depth analysis shows facts make 84/100 questions answerable-and-supported by k=30 versus 73–76/100 for the chunk controls (median first-success depth is 1 for all representations, reflecting simple direct lookups). Facts do not universally need fewer *units* — some direct questions are answered from one chunk; what they provide is a better joint coverage–depth–token curve.

**4.6 Where the advantage concentrates.** Qualitative audit of fact-win cases identifies four mechanisms: *attribution* (claims arrive with normalized speakers, so respondent vs. questioner needs no reconstruction), *coreference* ("she ultimately removed that tweet" is retrieved as "Senator Maria Chappelle-Nadal ultimately removed the tweet"), *claim salience* (a compact proposition matches its question without competing broadcast vocabulary), and *opinion coverage* (explicit judgments and predictions widen the answerable surface).

**4.7 Negative and cautionary results.** *(Results here are facts only; interpretation deferred to Section 5.)* (a) In an earlier 3,000-session agent-evaluation matrix where models queried each index through a constrained search tool, the fixed-chunk condition led the then-current fact index on evidence coverage — the fact advantage did not automatically survive the agent/tool layer. (b) A fact-plus-quote embedding ablation improved fact retrieval modestly but did not close that agent-layer gap by itself. (c) All results in 4.2–4.4 use same-document development data and LLM judges; an 80-item blinded human calibration worksheet exists but no human labels have been collected yet.

### 5. Discussion

**5.1 Claims and evidence.** Contribution 1 is supported by construction (13,064/13,064 indexed units carry validated quotes). Contributions 2–3 are supported by Sections 4.2–4.4 as *development evidence of feasibility and direction* — consistent with our P1 phase, we do not claim established effectiveness: the data are same-document, and correctness judgments await human calibration. Contribution 4 is an organizing claim supported by the literature analysis (Section 2).

**5.2 Compiled pointers vs. compiled payloads.** Dense X and Contextual Retrieval demonstrate that ingest-time work improves retrieval; our results suggest the larger prize is on the *reading* side. The 256-token fact condition outperforming 512-token chunks is precisely the outcome a payload-compilation account predicts and a pointer-compilation account does not: if propositions merely index better, accuracy should converge once chunks are found; instead the gap persists because the answer model no longer performs reconstruction. The depth sweep (4.5) adds a second signature: raw-chunk coverage *regresses* as depth grows (66→61/100 from k=20 to k=30) while compiled facts saturate flat — exactly what a reading-side account predicts, since each additional chunk adds distraction but each additional fact adds only a compact claim. This regression is a within-experiment instance of the context-rot and lost-in-the-middle phenomena [19, 20]: at k=30 the chunk conditions ask the answer model to find one claim inside ~21,000 tokens of transcript, precisely the regime where long-context reliability measurably decays. The extended frontier (4.4) completes the picture with a third signature: once candidate supply is unconstrained, even *compiled* context shows a mild budget penalty (92→86 as fact budgets grow eightfold, 7 paired regressions) — degradation is a property of context volume generally, and the advantage of compilation is that its accuracy optimum sits near 2k tokens where chunks have not yet reached 63/100. The implication runs beyond cost: because compiled facts hold the read payload small at any depth, they should also yield more *consistent* answers with fewer context-induced hallucinations — a prediction our fixed-budget and depth results support indirectly and that per-question answer-variance measurement should test directly. The agent-layer negative result (4.7a) sharpens rather than contradicts this: when a *tool-using model* mediates retrieval, single facts may under-supply context — motivating compiled-context remedies (sibling facts from adjacent turns) rather than falling back to raw windows, which would re-inflate the read budget the approach exists to shrink.

**5.3 The amortization argument.** Compilation is worthwhile when ingest cost is amortized across queries: with ingest cost c_c per document, query savings Δc_q per query, break-even occurs at R* = N·c_c / Δc_q reads 〔harmonize notation with the ISC substrate paper — FQ-3〕. LazyGraphRAG [10] correctly observes that low-query-volume corpora never reach R*; fact-memory break-even analysis [6] finds crossover within ~10 interactions in agent settings. Our measured 38–50% answer-model savings, plus distillation of the extractor as demonstrated by Dense X's Propositionizer [1], would drive R* down further. A dollar-denominated break-even estimate on this corpus is deliberately deferred until the accuracy results are human-calibrated (no cost claim should rest on uncalibrated correctness).

**5.4 Limitations.** (1) *Same-document development data*: questions were authored from the same 100 documents used to develop extraction; document-disjoint held-out replication is required before any effectiveness claim. (2) *LLM judges*: Kimi correctness/support screens are uncalibrated against humans (worksheet ready, labels pending), and the depth sweep (§4.5) uses DeepSeek Flash to judge DeepSeek Flash outputs — a self-model screen whose absolute values should not be compared across experiments; the confirmatory run will preserve the protocol with a different judge and human calibration. (3) *Single corpus and genre*: MediaSum is edited broadcast dialogue, not spontaneous workplace conversation. (4) *Gold-turn-overlap leniency* (4.2). (5) *Single answer/extraction model pairing*; a model-generality extension (Kimi K2.6 and GLM 5.2 as answer models at the 2,048 and 16,384 budgets, cross-judged so no model judges its own answers) is in progress. (6) *No statistical testing yet*: counts are reported without tests or intervals, consistent with the pilot phase; the P3 version requires tests, effect sizes, and an ablation over pipeline passes (validation gate on/off; rewrite on/off; quote-in-embedding on/off).

**5.5 Threats to validity.** *Internal:* question authorship by an LLM (albeit source-validated) may favor compact claims; the packing algorithm interacts with unit size. *External:* generalization beyond broadcast interviews and beyond the frozen sample is unestablished. *Construct:* "correct" is a blind semantic match to a reference answer, not human ground truth; "cited support" is likewise judged. Each threat has a planned mitigation in the P3 protocol (human labels; disjoint documents; second corpus; human-authored question subset).

### 6. Conclusion

We presented ISC-RAG, an ingest-time semantic compilation pipeline that turns dialogue transcripts into provenance-validated, self-contained facts serving as retrieval payloads, and reported controlled development evidence that this representation reaches gold evidence more often than transcript chunking and answers more questions correctly from a fraction of the read-time context, at lower per-query token cost. Future work, in order: collect human labels for the blinded 80-item calibration worksheet; freeze prompts, packing, and judgment protocol; replicate on document-disjoint MediaSum transcripts and a second dialogue corpus; add statistical testing and pipeline ablations; and estimate the extraction/indexing break-even against measured query-time savings.

### References

*(Working list; format to venue style on Q4 decision. arXiv items are preprints.)*

[1] T. Chen et al., "Dense X Retrieval: What Retrieval Granularity Should We Use?" EMNLP, 2024.
[2] Anthropic, "Introducing Contextual Retrieval," engineering blog, 2024. https://www.anthropic.com/engineering/contextual-retrieval
[3] Y. Hou et al., "AtomicRAG: Atom–Entity Graphs for Retrieval-Augmented Generation," arXiv:2604.20844, 2026 (preprint).
[4] J. Wang, "PropRAG: Guiding Retrieval with Beam Search over Proposition Paths," EMNLP, 2025.
[5] N. Zhang et al., "SiReRAG: Indexing Similar and Related Information for Multihop Reasoning," ICLR, 2025.
[6] N. Pollertlam and W. Kornsuwannawit, "Beyond the Context Window: A Cost-Performance Analysis of Fact-Based Memory vs. Long-Context LLMs for Persistent Agents," arXiv:2603.04814, 2026 (preprint).
[7] D. Chen et al., "HaluMem: Evaluating Hallucinations in Memory Systems of Agents," arXiv:2511.03506, 2025 (preprint).
[8] Z. Bao and Y. Shi, "IndexRAG: Bridging Facts for Cross-Document Reasoning at Index Time," arXiv:2603.16415, 2026 (preprint).
[9] 〔verify authors〕, "The Token Tax of Epistemic Accuracy: RAG vs Long-Context," arXiv:2606.20898, 2026 (preprint).
[10] Microsoft Research, "LazyGraphRAG: Setting a New Standard for Quality and Cost," blog, 2024.
[11] P. Lewis et al., "Retrieval-Augmented Generation for Knowledge-Intensive NLP Tasks," NeurIPS, 2020.
[12] V. Karpukhin et al., "Dense Passage Retrieval for Open-Domain Question Answering," EMNLP, 2020.
[13] E. Choi et al., "Decontextualization: Making Sentences Stand-Alone," TACL, 2021.
[14] C. Zhu et al., "MediaSum: A Large-scale Media Interview Dataset for Dialogue Summarization," NAACL, 2021.
[15] S. Min et al., "FActScore: Fine-grained Atomic Evaluation of Factual Precision in Long Form Text Generation," EMNLP, 2023.
[16] 〔verify authors〕, "Mem0: Building Production-Ready AI Agents with Scalable Long-Term Memory," arXiv:2504.19413, 2025 (preprint).
[17] M. Louis et al. 〔verify〕, "PISCO: Pretty Simple Compression for Retrieval-Augmented Generation," arXiv:2501.16075, 2025 (preprint); H. Jiang et al., "LLMLingua," EMNLP, 2023.
[18] O. Zahn and S. Chana, "Facts as First Class Objects: Knowledge Objects for Persistent LLM Memory," arXiv:2603.17781, 2026 (preprint).
[19] N. F. Liu et al., "Lost in the Middle: How Language Models Use Long Contexts," TACL, 2024.
[20] K. Hong et al., "Context Rot: How Increasing Input Tokens Impacts LLM Performance," Chroma technical report, 2025. https://www.trychroma.com/research/context-rot
[21] F. Shi et al., "Large Language Models Can Be Easily Distracted by Irrelevant Context," ICML, 2023.

---

## Part 4 — Review History

| vX (date) | Review summary |
|---|---|
| v4 (2026-07-25) | • ✅ Praised: new §4.4 extended frontier removes the top-20 supply confound and lands the paper's strongest table — facts@2,048 (92/100) above every chunk cell at every budget through 16,384, chunks' ceiling 75/100 at 7.3× cost; facts' own budget decline (92→86) honestly reported and integrated into §5.2 as a third context-rot signature; judging now via a repo-frozen reproducible CLI command. (Retrieval-depth defect and fix documented in the repo experiment record, not the paper.) • ❌ 2 / ⚠️ 6 (unchanged). • HIGH: FQ-1 venue, FQ-2 author block still open. • Changes: abstract sentence 7; new §4.4 (old 4.4→4.5, 4.5→4.6, 4.6→4.7 with cross-refs); §5.2 third signature; §5.4 limitations (2) renumbered ref, (5) model-generality in progress; Q11/Q12. |
| v3 (2026-07-25) | • ✅ Praised: the token-savings result is now double-framed as cost *and* reliability — context rot [20], lost-in-the-middle [19], and distraction [21] wired into Introduction, Related Work (QSR paragraph), and Discussion; the depth-sweep chunk regression is identified as a within-experiment instance of context rot; consistency/hallucination prediction stated as a testable follow-up (per-question answer variance). • ❌ 2 / ⚠️ 6 (unchanged). • HIGH: FQ-1 venue, FQ-2 author block still open. • Changes: abstract sentence 6; Introduction cost paragraph; Related Work QSR paragraph; §5.2; References [19]–[21]; Q12. |
| v2 (2026-07-25) | • ✅ Praised: new §4.4 depth curve strengthens the reading-side thesis with an independent manipulation (uncapped top-k) and adds the chunk-distraction regression (66→61/100) as direct payload-account evidence; abstract updated with the depth headline; self-judging caveat honestly added to Limitations. • ❌ 2 / ⚠️ 6 (unchanged — no fulfillment answers received this round). • HIGH: FQ-1 venue, FQ-2 author block still open. • Changes: abstract sentence 7; new §4.4 (old 4.4→4.5, 4.5→4.6); §5.1 evidence pointer; §5.2 distraction paragraph; §5.4 limitation (2); Q11/Q12. |
| v1 (2026-07-25) | • ✅ Praised: closed-world abstract compiles cleanly from the ISC-vs-QSR framing; E-block numbers fully sourced from frozen-100 experiments; Related Work 2-axis positioning (where work is paid × what is read) is the paper's organizing asset; negative results included honestly (4.5). • ❌ 2 / ⚠️ 6 (first compile). • HIGH: venue undecided (FQ-1); first-author affiliation/ORCID missing (FQ-2). • Key fulfillment questions: FQ-1 venue+deadline; FQ-2 author block; FQ-3 relationship to the ISC substrate paper (cite as prior work or merge?); FQ-4 settings table (prompts/temps/packing) transfer; FQ-5 human-calibration timeline. |
