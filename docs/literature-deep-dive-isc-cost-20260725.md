# Literature deep-dive: ISC-like compilation, fact units, and token economics (mid-2024 – July 2026)

**Timestamp:** 2026-07-25 (America/Los_Angeles)

**Status:** research-agent survey; every arXiv ID below was verified by fetching
its landing page at survey time. Industry items are blogs and are cited by URL.
This document records candidate related work; nothing here is promoted into
[the literature landscape](literature-landscape.md) until independently
re-verified and positioned.

## Section 1 — ISC-like concepts: semantic work at ingest/index time

1. **IndexRAG: Bridging Facts for Cross-Document Reasoning at Index Time** — arXiv, Mar 17 2026 — https://arxiv.org/abs/2603.16415
   Moves multi-hop reasoning from online inference to offline indexing: identifies shared bridge entities across documents and pre-generates "bridging facts" as retrievable units, so query time needs only single-pass retrieval and one LLM call. The closest published articulation of "do the reasoning once at index time" for QA.
   *Relation to ISC-RAG: compiles facts at ingest, but no provenance validation and no dialogue domain; cost framing is latency, not token economics.*

2. **Knowledge Capsules: Structured Nonparametric Memory Units for LLMs** — arXiv, Apr 22 2026 — https://arxiv.org/abs/2604.20487
   Distills corpora into structured "capsules" of normalized relational knowledge using a frozen LLM, then compiles them into attention-compatible key-value representations injected directly into attention (KVI), bypassing token competition.
   *Relation to ISC-RAG: pay-once compilation, but into latent KV — payload is not human-readable text and carries no source quotes.*

3. **Don't Retrieve, Navigate (Corpus2Skill)** — arXiv, Apr 16 2026 — https://arxiv.org/abs/2604.14572
   One-time offline "compile phase" distills an enterprise corpus into a hierarchical skill directory; at serve time an agent navigates progressively finer summaries down to documents. Explicit compile/serve vocabulary; notes a "moderate cost tradeoff."
   *Relation to ISC-RAG: compiles pointers (navigation structure), not payload — the read is still raw documents.*

4. **Mem0: Building Production-Ready AI Agents with Scalable Long-Term Memory** — arXiv, Apr 28 2025 — https://arxiv.org/abs/2504.19413
   Extracts salient facts from ongoing conversations at ingest with LLM-driven ADD/UPDATE/DELETE/NOOP consolidation; graph variant adds entity-relation structure. Claims >90% token cost savings and 91% lower p95 latency vs full-context on LOCOMO.
   *Relation to ISC-RAG: fact extraction from dialogue with a token-cost claim, but no provenance validation (no quote/span check) and no budget-matched accuracy comparison.*

5. **Zep: A Temporal Knowledge Graph Architecture for Agent Memory (Graphiti)** — arXiv, Jan 20 2025 — https://arxiv.org/abs/2501.13956
   Temporally-aware knowledge graph memory: fact edges with validity intervals, extracted from dialogue and business data. Graphiti's data model keeps EpisodicNodes storing the original input each fact was extracted from — a provenance *pointer*, not a validated quote.
   *Relation to ISC-RAG: facts carry a source-episode link but provenance is untested/unvalidated; facts are edges in an index, retrieval often returns them plus context.*

6. **A-MEM: Agentic Memory for LLM Agents** — arXiv, Feb 2025 — https://arxiv.org/abs/2502.12110
   Zettelkasten-style atomic notes with LLM-driven construction, linking, and evolution into a semantic graph; reports 2–6× multi-hop efficiency gains with token savings.
   *Relation to ISC-RAG: atomic self-contained notes as ingest-time units, but notes are free-form and carry no verified source spans.*

7. **AtomMem: Learnable Dynamic Agentic Memory with Atomic Memory Operation** — arXiv, Jan 13 2026 — https://arxiv.org/abs/2601.08323
   Learns memory management (SFT+RL) as a decision problem over atomic CRUD operations, beating static memory workflows on long-context benchmarks.
   *Relation to ISC-RAG: "atomic" refers to operations, not fact units; no provenance.*

8. **SGMem: Sentence Graph Memory for Long-Term Conversational Agents** — arXiv, Sep 25 2025 — https://arxiv.org/abs/2509.21212
   Represents dialogue as sentence-level graphs across turn/round/session granularity; answer generation combines retrieved raw dialogue with generated summaries/facts/insights.
   *Relation to ISC-RAG: hybrid — compiled facts supplement rather than replace raw-dialogue payload; provenance implicit, unvalidated.*

9. **xMemory: Beyond RAG for Agent Memory — Retrieval by Decoupling and Aggregation** — arXiv, Feb 2 2026 — https://arxiv.org/abs/2602.02007
   Ingest-time hierarchy (messages → segments → decoupled fact/update/detail components → groups); query time reads a compact backbone of components and expands to raw messages only when reader uncertainty warrants. Explicitly targets context quality vs inference efficiency.
   *Relation to ISC-RAG: structurally the closest agent-memory cousin (compiled facts as the default read payload, raw text as fallback), but no quote-level provenance validation.*

10. **TurboRAG: Precomputed KV Caches for Chunked Text** — arXiv, Oct 2024 — https://arxiv.org/abs/2410.07590
    Pre-computes document KV caches offline and reuses them at query time; 8.6–9.4× TTFT reduction at comparable accuracy. Related systems: CacheBlend, CacheClip (https://arxiv.org/pdf/2510.10129).
    *Relation to ISC-RAG: the "pay once at ingest" idea applied at the compute level, not the semantic level — amortizes prefill FLOPs, not meaning.*

11. **LazyGraphRAG** — Microsoft Research blog, Nov 2024 (updated 2025) — https://www.microsoft.com/en-us/research/blog/lazygraphrag-setting-a-new-standard-for-quality-and-cost/
    Deliberately *defers* all LLM summarization to query time, keeping indexing near-free; global queries ~700× cheaper than full GraphRAG at comparable quality.
    *Relation to ISC-RAG: the explicit counter-position — argues ingest-time LLM spend is often wasted; ISC-RAG must answer why fact compilation earns back its ingest cost (answer: high query volume per corpus + provenance value).*

12. **From RAG to Memory: Non-Parametric Continual Learning (HippoRAG 2)** — arXiv, Feb 2025 — https://arxiv.org/abs/2502.14802
    *Relation to ISC-RAG: ingest-time KG construction, but retrieval payload remains passages; graph is index-only.*

## Section 2 — Fact-like retrieval units: provenance and payload status

13. **Clue-RAG** — arXiv, Jul 11 2025 — https://arxiv.org/abs/2507.08445
    Multi-partite graph over chunks, disambiguated "knowledge units," and entities, with query-driven iterative retrieval; cuts indexing cost 72.58% vs graph-RAG baselines while improving accuracy.
    *Relation to ISC-RAG: knowledge units are index keys routing back to chunks — pointer, not payload; no provenance validation.*

14. **E²RAG: Entity-Event Knowledge Graphs for Temporal-Causal RAG** — arXiv, Jun 6 2025 — https://arxiv.org/abs/2506.05939
    Dual entity/event subgraphs built at ingest preserve chronology and causality for narrative QA; introduces ChronoQA benchmark.
    *Relation to ISC-RAG: event units are ISC-adjacent for narrative/dialogue-like text, but no quote provenance and units serve retrieval, not reading.*

15. **GenProve: Learning to Generate Text with Fine-Grained Provenance** — arXiv, Jan 8 2026 — https://arxiv.org/abs/2601.04932
    Trains models (SFT+RL) to emit sentence-level provenance triples distinguishing quotation vs compression vs inference, with the expert-verified ReFInE dataset. Finding: models handle surface quotation well, inference-provenance poorly.
    *Relation to ISC-RAG: validates exactly the quote-anchored provenance ISC facts carry — but at generation time, not as a compiled corpus artifact.*

16. **PaperTrail: A Claim-Evidence Interface for LLM Scholarly Q&A** — arXiv, Feb 24 2026 — https://arxiv.org/abs/2602.21045
    Decomposes answers and sources into claims/evidence and maps them, surfacing supported, unsupported, and omitted assertions; user study (n=26) shows calibrated distrust.
    *Relation to ISC-RAG: claim–evidence mapping as UI at query time; ISC does it once at ingest.*

17. **From Fluent to Verifiable: Claim-Level Auditability for Deep Research Agents** — arXiv, Feb 14 2026 — https://arxiv.org/abs/2602.13855
    Proposes claim-level auditability with "provenance coverage" and "provenance soundness" metrics and semantic provenance graphs built during synthesis.
    *Relation to ISC-RAG: supplies ready-made vocabulary/metrics for evaluating ISC's provenance validation step.*

18. **HaluMem: Evaluating Hallucinations in Memory Systems of Agents** — arXiv, Nov 5 2025 — https://arxiv.org/abs/2511.03506
    Operation-level benchmark (~15k memory points, >1M-token conversations) showing memory systems hallucinate during *extraction and updating*, and these errors propagate into QA.
    *Relation to ISC-RAG: the strongest empirical motivation for ISC's provenance-validation step — extraction without validation is a measured failure mode.*

19. **GETALP@AutoMin 2025: RAG over Meeting Transcripts** — arXiv, Aug 1 2025 — https://arxiv.org/abs/2508.00476
    SIGDial AutoMin Task B system: dense retrieval over meeting transcripts plus AMR-derived graphs verbalized into descriptions; AMR notably helps "who said/did" questions.
    *Relation to ISC-RAG: rare meeting-transcript RAG datapoint; structured semantics help speaker attribution, but no atomic facts and no provenance.*

20. **Memory in the LLM Era: Modular Architectures in a Unified Framework** — arXiv, Apr 2 2026 — https://arxiv.org/abs/2604.01707
    Survey + benchmark unifying agent-memory methods (Mem0, A-MEM, Zep lineage) with comparative evaluation.
    *Relation to ISC-RAG: useful related-work map; confirms no surveyed system validates provenance of extracted facts.*

Note: the survey found no 2025–26 work adding *validated* provenance to
proposition units; Dense X's extensions concentrate in PropRAG/AtomicRAG and
the memory papers above.

## Section 3 — Token cost / economics of ingest vs query

21. **Beyond the Context Window: Cost-Performance Analysis of Fact-Based Memory vs. Long-Context LLMs** — arXiv, Mar 5 2026 — https://arxiv.org/abs/2603.04814
    Directly compares fact-based memory (Mem0-style) vs long-context on accuracy *and cumulative API cost*, with a cost model including prompt caching. Explicit break-even/amortization result: memory becomes cheaper after ~10 interaction turns at 100k context, improving with context length.
    *Relation to ISC-RAG: the closest anyone has come to the ISC-vs-QSR amortization argument — but memory-vs-long-context, not compiled-facts-vs-retrieved-chunks at matched read budgets.*

22. **The Token Tax of Epistemic Accuracy: RAG vs Long-Context** — arXiv, Jun 18 2026 — https://arxiv.org/abs/2606.20898
    Manufacturing-safety benchmark: long-context beats semantic RAG on accuracy (73.1% vs 65.4%) but at 26× per-query token cost — frames the "token tax" of broad evidentiary access.
    *Relation to ISC-RAG: establishes the cost-quality frontier framing, but chunk-level only; compiled facts don't appear as an arm.*

23. **Cost-Aware Query Routing in RAG (CA-RAG)** — arXiv, Mar 26 2026 — https://arxiv.org/abs/2606.02581
    Per-query utility-based routing across retrieval depths; 26% fewer billed tokens than always-heavy retrieval at equal quality.
    *Relation to ISC-RAG: measures cost but chunk-level only; optimizes query-time depth rather than ingest-time representation.*

24. **RAGRouter-Bench** — arXiv, Jan 30 2026 — https://arxiv.org/abs/2602.00296
    Benchmark for routing across RAG paradigms measuring quality and resource consumption; search-surfaced material describes token accounting split into one-time construction cost vs per-query inference cost, though the abstract itself doesn't foreground the decomposition.
    *Relation to ISC-RAG: rare instance of construction-vs-inference token bookkeeping across paradigms — useful evaluation template.*

25. **PISCO: Pretty Simple Compression for RAG** — arXiv, Jan 27 2025 — https://arxiv.org/abs/2501.16075
    xRAG-successor soft compression: 16× compression with 0–3% accuracy loss via sequence-level distillation, no annotations needed. (Family: LLMLingua/LLMLingua-2, RECOMP arXiv:2310.04408, xRAG arXiv:2405.13792.)
    *Relation to ISC-RAG: reduces read tokens but per-query, lossy, latent, and provenance-destroying — QSR with a smaller bill, not ISC.*

26. **Chunking Methods on RAG: Effectiveness vs Computational Cost** — arXiv, May 30 2026 — https://arxiv.org/abs/2606.00881
    First systematic evaluation of chunking strategies against compute cost; finds chunking choice materially impacts performance and many proposed strategies lack cross-scenario evidence.
    *Relation to ISC-RAG: documents QSR's preprocessing fragility — the problem ISC's unit redefinition dissolves.*

27. **Industry writing on the ingest/query cost ledger:**
    - **The GraphRAG Cost Cliff: How $33,000 Became $33** — Graph Praxis/Medium, 2026 — https://medium.com/graph-praxis/the-graphrag-cost-cliff-how-33-000-became-33-in-eighteen-months-be1b0fbe37e4 — traces GraphRAG indexing economics from prohibitive to viable; frames indexing spend as amortizable.
    - **RAGFlow: From RAG to Context — 2025 year-end review** — https://ragflow.io/blog/rag-review-2025-from-rag-to-context — industry synthesis of the shift toward context engineering and index-time investment.
    - **FalkorDB: Reduce GraphRAG Indexing Costs** — https://www.falkordb.com/blog/reduce-graphrag-indexing-costs/ — practitioner prescriptions for taming ingest-time LLM spend.
    - **CloudZero / NeuralTrust / Glasp 2026 pieces** (https://www.cloudzero.com/blog/inference-cost/, https://neuraltrust.ai/blog/ai-token-optimization-guide, https://glasp.co/articles/context-rot-rag-long-context-hybrid) — the "inference cost paradox": per-token prices fell ~280× while enterprise spend rose ~320% as agentic workflows inflate context 3–5×; standard prescription is retrieval + context discipline, still chunk-level.
    *Relation to ISC-RAG: strong ambient demand for the cost argument; none of it reaches fact-granularity payloads.*

## Shortlist: most citation-worthy for the ISC-RAG paper

1. **arXiv:2603.04814** (Beyond the Context Window) — the only explicit ingest-amortization break-even analysis; ISC-RAG sharpens it to budget-matched fact payloads.
2. **arXiv:2511.03506** (HaluMem) — empirical case that unvalidated fact extraction hallucinates and propagates; directly motivates provenance validation.
3. **arXiv:2603.16415** (IndexRAG) — nearest neighbor on "reason at index time, read compiled facts at query time."
4. **arXiv:2504.19413** (Mem0) — the headline industry-adjacent token-savings claim (90%) from dialogue fact extraction, sans provenance.
5. **arXiv:2602.02007** (xMemory) — compiled-facts-as-default-payload with raw-text fallback; the architecture ISC-RAG's ablations should compare against.
6. **arXiv:2606.20898** (Token Tax) — cost-quality frontier framing at chunk level, the baseline frontier ISC claims to dominate.
7. **arXiv:2601.04932** (GenProve) — quotation-grade provenance is learnable and measurable; legitimizes exact-quote validation.
8. **LazyGraphRAG blog** — the strongest counter-argument (defer semantics to query time); citing and answering it strengthens the thesis.

## Where the ISC-vs-QSR cost argument stands

The argument exists in fragments but has not been made whole, and the specific
ISC-RAG claim appears open. Three partial versions are in print: (a) the
*systems* version — TurboRAG/CacheBlend amortize prefill compute at ingest, but
only lexically; (b) the *agent-memory* version — Mem0 claims 90% token savings
and arXiv:2603.04814 gives an explicit break-even (~10 turns) for fact-based
memory vs long context, but neither controls the read budget nor validates
provenance, and HaluMem shows the extraction step they rely on is a live
hallucination source; (c) the *frontier* version — Token Tax (2606.20898),
CA-RAG, and RAGRouter-Bench measure cost-quality tradeoffs, but every arm is
chunk- or route-level. Meanwhile LazyGraphRAG argues the opposite direction
(indexing spend is often wasted), which makes query volume per corpus the
pivotal variable no one has formalized. Nobody, as of July 2026, has published
the compound claim ISC-RAG makes: that provenance-validated atomic facts *as
the read payload* beat raw chunks at a fraction of the token budget, with the
ingest cost amortized across queries. The budget-matched accuracy-per-token
comparison at fact granularity, with a validation step answering HaluMem, is
unclaimed territory. The nearest works (IndexRAG, xMemory, Knowledge Capsules)
each hold one piece — ingest-time compilation, fact-payload reads, pay-once
framing — but none combines all three with verified provenance, and none works
in the meeting/dialogue domain except unstructured systems like GETALP's.
