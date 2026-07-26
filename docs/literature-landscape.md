# ISC-RAG literature landscape

**Scope.** ISC-RAG asks whether ingest-time extraction of atomic, source-linked
facts is a better retrieval representation than fixed-window, turn-aware, and
semantic chunks for real human dialogue. This landscape is deliberately about
the *retrieval unit* and its verifiable provenance, not a general RAG survey.
All links below were opened on 2026-07-24. ACL Anthology, NeurIPS, PMLR, and
arXiv links are direct paper landing pages (and expose a PDF); arXiv items are
marked as preprints.

## Topic map

| Axis | Established line | Open issue for ISC-RAG |
|---|---|---|
| Retrieval architecture | Dense passage retrieval and RAG make a text unit retrievable and then condition a generator on it. | Hold embedding, index, ranking, and answer budget fixed so the unit—not the rest of the stack—is the treatment. |
| Granularity | Passage, sentence, proposition, structure-aware, and contextualized units trade recall, self-containment, and context cost. | Compare three realistic dialogue chunkers with propositions/facts at equal answer-context budgets. |
| Dialogue | Dialogue retrieval must deal with speakers, turn boundaries, coreference, topic shifts, and long context. | Measure on human broadcast interviews, where speaker and exact-turn provenance are native rather than invented labels. |
| Provenance and faithfulness | Evidence extraction, citation attribution, and claim checking evaluate whether text is actually supported. | Require every indexed fact to retain an exact verified quote plus turn/speaker offsets, then separately measure retrieval and answer support. |
| Evaluation | RAG evaluation separates context retrieval from generated-answer quality, but automatic faithfulness metrics are imperfect. | Use deterministic answer/evidence labels as the primary measure, with audits and automated faithfulness checks as secondary diagnostics. |

**ISC-RAG's position.** It sits at the intersection of proposition-level dense
retrieval and evidence-grounded dialogue: an *ingestion-time*, provenance-gated
fact representation for a corpus of real interviews, compared fairly against
three chunking choices under the same retrieval and generation conditions.

## Research world-line and concrete gap

**World-line statement.** If a dialogue corpus can be compiled into small,
source-verifiable propositions without losing meaning, RAG systems could retrieve
the relevant claim with less distracting context while allowing a reader to jump
straight back to the words, speaker, and turn that support it.

**Gap.** Dense X Retrieval directly establishes proposition units as a promising
granularity, but studies Wikipedia/open-domain QA rather than spoken interview
transcripts, and it does not make exact source-quote validation the experimental
contract. Dialogue benchmarks establish that turn dynamics, coreference, and
long context matter, while attribution work shows that a citation is not itself
proof of support. No identified work provides the four-way, controlled comparison
of fixed, turn-aware, semantic, and exact-provenance fact units on a large corpus
of real human interviews with retrieval-first and deterministic answer measures.
This is a gap statement to test, not a claim of priority until a full systematic
review is complete.

## Related-work positioning paragraph

ISC-RAG builds on retrieval-augmented generation and dense passage retrieval,
but treats the indexed representation as the experimental variable. It is closest
to proposition-level retrieval, which indexes self-contained atomic factoids and
reports gains over passages on open-domain QA. ISC-RAG differs by targeting
human dialogue, retaining speaker/turn-level provenance, rejecting facts whose
supporting quotation cannot be located exactly, and comparing fact cards with
three dialogue-appropriate chunking baselines under a fixed context budget. It
therefore complements document-grounded dialogue and attribution benchmarks:
those evaluate grounding or generation, whereas ISC-RAG asks whether a
provenance-preserving *ingestion representation* improves evidence retrieval and
deterministically evaluable answers.

## Verified reference shortlist

### Retrieval and RAG foundations

1. Patrick Lewis et al. (2020). [*Retrieval-Augmented Generation for Knowledge-Intensive NLP Tasks*](https://proceedings.neurips.cc/paper_files/paper/2020/hash/6b493230205f780e1bc26945df7481e5-Abstract.html). NeurIPS.
   Combines a generator with a retrieved dense index and makes provenance and
   updating explicit motivations. It is the architectural starting point, not a
   study of what a retrieval unit should be.

2. Vladimir Karpukhin et al. (2020). [*Dense Passage Retrieval for Open-Domain Question Answering*](https://aclanthology.org/2020.emnlp-main.550/). EMNLP.
   Establishes a strong dense-retrieval passage baseline. ISC-RAG holds this
   family of retrieval machinery constant while changing the indexed unit.

3. Kelvin Guu et al. (2020). [*Retrieval Augmented Language Model Pre-Training*](https://proceedings.mlr.press/v119/guu20a.html). ICML.
   An early retrieval-augmented language-model formulation. Useful background
   for separating parametric memory from retrieved text evidence.

4. Xiaohua Wang et al. (2024). [*Searching for Best Practices in Retrieval-Augmented Generation*](https://aclanthology.org/2024.emnlp-main.981/). EMNLP.
   Systematically treats a RAG pipeline as a collection of interacting choices.
   It supports freezing non-representation choices rather than optimizing each
   condition differently.

### Retrieval unit granularity and chunking

5. Tong Chen et al. (2024). [*Dense X Retrieval: What Retrieval Granularity Should We Use?*](https://aclanthology.org/2024.emnlp-main.845/). EMNLP.
   The nearest conceptual predecessor. It defines propositions as atomic,
   self-contained factoids and finds proposition indexing stronger than passage
   indexing for its open-domain retrieval/QA settings and computation budget.

6. Antonio Jimeno Yepes et al. (2024). [*Financial Report Chunking for Effective Retrieval Augmented Generation*](https://arxiv.org/abs/2402.05131). arXiv preprint.
   Tests document-structure-aware chunks in financial reports. It motivates a
   structural baseline, but its source genre and annotations differ sharply from
   conversational turns.

7. Michael Günther et al. (2024). [*Late Chunking: Contextual Chunk Embeddings Using Long-Context Embedding Models*](https://arxiv.org/abs/2409.04701). arXiv preprint.
   Argues that embedding isolated chunks can discard surrounding context, and
   applies segmentation after contextual token encoding. It is a future
   embedding-side ablation, rather than one of ISC-RAG's primary four units.

8. Zijie Zhong et al. (2024). [*Mix-of-Granularity: Optimize the Chunking Granularity for Retrieval-Augmented Generation*](https://arxiv.org/abs/2406.00456). arXiv preprint.
   Explores selecting different granularities for different knowledge sources.
   ISC-RAG instead asks the more controlled question: for one dialogue corpus,
   does a provenance-gated fact unit beat three fixed representation policies?

### Rephrasing and decontextualization

9. Eunsol Choi et al. (2021). [*Decontextualization: Making Sentences Stand-Alone*](https://aclanthology.org/2021.tacl-1.27/). TACL.
   Defines rewriting a sentence with its context so it can be interpreted alone
   while preserving its meaning, and evaluates it as preprocessing for retrieval.
   ISC-RAG adopts this need for standalone units but adds atomicity, dialogue
   attribution, and exact-quote provenance as a hard contract.

10. Benjamin Newman et al. (2023). [*A Question Answering Framework for Decontextualizing User-facing Snippets from Scientific Documents*](https://aclanthology.org/2023.emnlp-main.193/). EMNLP.
   Frames decontextualization as supplying missing context through question
   generation, answering, and rewriting. It motivates auditing which contextual
   additions a rewrite introduces; ISC-RAG must check those additions against a
   transcript quotation rather than treating fluent rewriting as faithful.

11. Zhenyun Deng, Yulong Chen, and Andreas Vlachos (2025). [*Improving Zero-shot Sentence Decontextualisation with Content Selection and Planning*](https://aclanthology.org/2025.emnlp-main.1258/). EMNLP.
   Explicitly selects context needed to resolve ambiguity before rewriting a
   sentence. It is a relevant future extraction baseline or ablation for
   ISC-RAG, especially for pronouns, ellipsis, and speaker-dependent claims.

### Dialogue, interviews, and conversational retrieval

12. Chenguang Zhu et al. (2021). [*MediaSum: A Large-scale Media Interview Dataset for Dialogue Summarization*](https://aclanthology.org/2021.naacl-main.474/). NAACL.
   Introduces the 463.6K-transcript NPR/CNN interview collection used by
   ISC-RAG. Its scale and multiparty broadcasts make it a useful real-dialogue
   retrieval corpus, though its original task is summarization.

13. Ming Zhong et al. (2021). [*QMSum: A New Benchmark for Query-based Multi-domain Meeting Summarization*](https://aclanthology.org/2021.naacl-main.472/). NAACL.
    Provides query-to-relevant-span annotations over long meetings and a
    locate-then-summarize framing. It is especially relevant to the future
    multi-turn/multi-fact ISC-RAG evaluation slice.

14. Song Feng et al. (2020). [*doc2dial: A Goal-Oriented Document-Grounded Dialogue Dataset*](https://aclanthology.org/2020.emnlp-main.652/). EMNLP.
    Grounds information-seeking dialogues in external documents. It differs from
    ISC-RAG because MediaSum itself is the human dialogue corpus being searched,
    rather than a dialogue about a separate source document.

15. Prakhar Gupta et al. (2022). [*DialFact: A Benchmark for Fact-Checking in Dialogue*](https://aclanthology.org/2022.acl-long.263/). ACL.
    Pairs conversational claims with evidence retrieval and verification tasks,
    identifying dialogue-specific issues such as colloquial wording and
    coreference. Its evidence is Wikipedia; ISC-RAG's evidence is the source
    transcript itself.

16. Yohan Lee, Yongwoo Song, and Sangyeop Kim (2025). [*Finding Diamonds in Conversation Haystacks: A Benchmark for Conversational Data Retrieval*](https://aclanthology.org/2025.emnlp-industry.162/). EMNLP Industry Track.
    A dedicated conversational retrieval benchmark with query and conversation
    labels. It supports the premise that conversational retrieval remains hard,
    but studies product-insight tasks rather than source-grounded fact cards.

17. Nouha Dziri et al. (2022). [*Evaluating Attribution in Dialogue Systems: The BEGIN Benchmark*](https://aclanthology.org/2022.tacl-1.62/). TACL.
    Human annotations show that common automatic grounding metrics can fail,
    especially with longer knowledge sources. This is a reason to retain exact
    evidence labels and manual audits in ISC-RAG.

### Atomic facts, evidence, and provenance

18. Sewon Min et al. (2023). [*FActScore: Fine-grained Atomic Evaluation of Factual Precision in Long Form Text Generation*](https://aclanthology.org/2023.emnlp-main.741/). EMNLP.
    Evaluates generated text by decomposing it into atomic facts and checking
    support. ISC-RAG applies the atomic-unit idea earlier, at corpus ingestion,
    and still needs to report fact recall as well as precision/fidelity.

19. Pranav Mani, Davis Liang, and Zachary Chase Lipton (2024). [*Fast Evidence Extraction for Grounded Language Model Outputs*](https://aclanthology.org/2024.fever-1.24/). FEVER Workshop.
    Frames locating evidence spans for generated lines as a scalable task. It is
    relevant to validating the source span behind every extracted fact.

20. Liyan Tang, Philippe Laban, and Greg Durrett (2024). [*MiniCheck: Efficient Fact-Checking of LLMs on Grounding Documents*](https://aclanthology.org/2024.emnlp-main.499/). EMNLP.
    Studies efficient checking of whether claims are supported by grounding
    documents. It is a useful secondary audit tool, but cannot replace exact
    quote-location validation in the primary pipeline.

21. Joe Stacey et al. (2024). [*Atomic Inference for NLI with Generated Facts as Atoms*](https://aclanthology.org/2024.emnlp-main.569/). EMNLP.
    Shows that generated facts need a multi-stage treatment; direct use can hurt.
    This strengthens the case for separating extraction, provenance validation,
    and optional enrichment rather than treating all generated fact text as safe.

22. James Thorne et al. (2018). [*The Fact Extraction and VERification (FEVER) Shared Task*](https://aclanthology.org/W18-5501/). FEVER Workshop.
    A foundational claim-and-evidence benchmark. ISC-RAG borrows the principle of
    evaluating evidence retrieval separately from downstream decisions, while
    using transcript spans instead of Wikipedia evidence pages.

### RAG evaluation and answer attribution

23. Shahul Es et al. (2024). [*RAGAs: Automated Evaluation of Retrieval Augmented Generation*](https://aclanthology.org/2024.eacl-demo.16/). EACL Demonstrations.
    Separates context relevance, faithfulness, and answer quality without human
    reference labels. It is useful as a diagnostic complement, not a replacement
    for ISC-RAG's evidence-grounded deterministic labels.

24. Jirui Qi et al. (2024). [*Model Internals-based Answer Attribution for Trustworthy Retrieval-Augmented Generation*](https://aclanthology.org/2024.emnlp-main.347/). EMNLP.
    Examines whether RAG answer tokens can be attributed to retrieved documents.
    It motivates reporting citation/evidence correctness rather than accepting a
    model's requested citation format at face value.

25. Weijia Zhang et al. (2024). [*Towards Fine-Grained Citation Evaluation in Generated Text: A Comparative Analysis of Faithfulness Metrics*](https://aclanthology.org/2024.inlg-main.35/). INLG.
    Finds no single automatic support metric reliable in every setting. ISC-RAG
    should therefore make exact evidence overlap and deterministic answer match
    primary, and use LLM judging as sensitivity analysis only.

### Contextual chunking and proposition successors (added 2026-07-25)

26. Anthropic (2024). [*Introducing Contextual Retrieval*](https://www.anthropic.com/engineering/contextual-retrieval). Engineering blog post.
    Prepends 50–100 tokens of LLM-generated situating context to each chunk
    before embedding and BM25 indexing, cutting top-20 retrieval failures by
    49%. Ingest-time semantic work, and therefore a partial form of ISC — but
    the retrieval unit remains a raw chunk that grows larger, and the generated
    context is not validated against the source. It compiles the pointer, not
    the read-time payload.

27. Nan Zhang et al. (2024). [*SiReRAG: Indexing Similar and Related Information for Multihop Reasoning*](https://arxiv.org/abs/2412.06206). ICLR 2025.
    Extracts propositions and entities, grouping propositions by shared
    entities into a relatedness tree alongside a similarity tree. Propositions
    organize multihop retrieval structure; there is no source-span validation
    and no read-budget-controlled answer evaluation.

28. Jingjin Wang et al. (2025). [*PropRAG: Guiding Retrieval with Beam Search over Proposition Paths*](https://aclanthology.org/2025.emnlp-main.317/). EMNLP.
    Uses contextually rich propositions as nodes and beam search over
    proposition paths to build multi-hop reasoning chains. Like ISC-RAG it uses
    a prompted large model for extraction rather than a distilled one, but the
    propositions serve reasoning-path discovery, not compiled, provenance-gated
    read-time context.

29. Yanning Hou et al. (2026). [*AtomicRAG: Atom–Entity Graphs for Retrieval-Augmented Generation*](https://arxiv.org/abs/2604.20844). arXiv preprint.
    Replaces text chunks with "knowledge atoms" — self-contained factual
    units — in an atom–entity graph, arguing that chunks rigidly group multiple
    atomic facts together. The closest recent descendant of Dense X; still
    Wikipedia-style benchmarks, no exact-quote validation, no dialogue.

30. Krittin Pachtrachai et al. (2026). [*From Transcripts to AI Agents: Knowledge Extraction, RAG Integration, and Robust Evaluation of Conversational AI Assistants*](https://arxiv.org/abs/2602.15859). arXiv preprint.
    The closest domain neighbor: extracts structured knowledge from real call
    transcripts as the sole RAG grounding source, evaluated with a
    transcript-grounded user simulator. An industrial end-to-end system paper;
    it does not run a controlled representation comparison or enforce
    quote-level provenance.

### Ingest-time compilation and token economics (added 2026-07-25)

These entries come from the broader
[literature deep-dive](literature-deep-dive-isc-cost-20260725.md), which holds
24 further verified candidates.

31. Natchanon Pollertlam and Witchayut Kornsuwannawit (2026). [*Beyond the Context Window: A Cost-Performance Analysis of Fact-Based Memory vs. Long-Context LLMs for Persistent Agents*](https://arxiv.org/abs/2603.04814). arXiv preprint.
    The only published explicit break-even analysis for ingest-time fact
    memory: fact-based memory becomes cheaper than long-context after roughly
    ten interaction turns at 100k context, under a cost model that includes
    prompt caching. It compares memory against long context rather than
    compiled facts against retrieved chunks at matched read budgets, and it
    validates no provenance — the sharpened, budget-matched version of its
    amortization argument is the one ISC-RAG makes.

32. Ding Chen et al. (2025). [*HaluMem: Evaluating Hallucinations in Memory Systems of Agents*](https://arxiv.org/abs/2511.03506). arXiv preprint.
    An operation-level benchmark (~15k memory points over million-token
    conversations) showing that agent memory systems hallucinate during fact
    extraction and updating, and that these errors propagate into downstream
    QA. This is the strongest empirical motivation for ISC-RAG's exact-quote
    validation gate: extraction without validation is a measured failure mode,
    not a hypothetical one.

33. Zhenghua Bao and Yi Shi (2026). [*IndexRAG: Bridging Facts for Cross-Document Reasoning at Index Time*](https://arxiv.org/abs/2603.16415). arXiv preprint.
    Moves multi-hop reasoning offline by pre-generating cross-document
    "bridging facts" as retrievable units, leaving query time a single-pass
    retrieval and one LLM call. The nearest structural neighbor to ingest-time
    semantic compilation for QA; it has no provenance validation, no dialogue
    domain, and frames the benefit as latency rather than token economics.

## Directly relevant people and groups to follow

This is a relevance list, **not a citation-count or leadership ranking**. The
paper landing pages above verify the authorship; current institutional affiliations
change and should be checked on official profile pages before using them in a
proposal or paper.

- **Tong Chen, Hongwei Wang, and Dong Yu** — proposition-level retrieval
  ([Dense X Retrieval](https://aclanthology.org/2024.emnlp-main.845/)).
- **Patrick Lewis and Douwe Kiela** — retrieval-augmented generation
  ([RAG](https://proceedings.neurips.cc/paper_files/paper/2020/hash/6b493230205f780e1bc26945df7481e5-Abstract.html)).
- **Chenguang Zhu and Yang Liu** — large-scale interview dialogue data
  ([MediaSum](https://aclanthology.org/2021.naacl-main.474/)).
- **Ming Zhong and Yang Liu** — query-focused long meeting spans
  ([QMSum](https://aclanthology.org/2021.naacl-main.472/)).
- **Greg Durrett** — factual support evaluation and efficient checking
  ([MiniCheck](https://aclanthology.org/2024.emnlp-main.499/)).
- **Evangelos Kanoulas** — citation/faithfulness evaluation
  ([fine-grained citation evaluation](https://aclanthology.org/2024.inlg-main.35/)).

## Implications for the experimental design

1. Treat `retrieval unit` as the only primary independent variable. Fix corpus
   documents, embedding model, distance metric, index configuration, ranking,
   answer model/prompt, and *answer-context token budget*.
2. Make evidence-retrieval metrics primary: gold-turn overlap recall@k, MRR, and
   nDCG, plus evidence density/context tokens. This avoids confusing a better
   generator with a better representation.
3. Require exact-quote and offset validation before a fact enters the primary
   condition. Report extraction yield, quote-validation failure rate, duplicate
   rate, and a blind human audit of fact meaning and attribution.
4. Score answers with strict normalized exact match and gold-evidence overlap;
   report automated faithfulness only as a labelled secondary measure. The
   canonical questions currently cover one-hop, single-turn claims, so state
   clearly that multi-turn synthesis and coreference are follow-up strata.
5. Include `fact-only` versus `fact + parent turns` as a named ablation. This
   distinguishes a compact index advantage from a context-restoration advantage.

## Quality report and open uncertainties

- **Verified references:** 33 total; **32/33 (97%)** are 2020 or later;
  **8/33 (24%)** are explicitly marked arXiv preprints and **1/33** is an
  industry engineering blog post; **0** have an unverified URL in this
  document. Entries 26–33 were landing-page verified on 2026-07-25, including
  author lists.
- **Strongest nearest work:** Dense X Retrieval. The paper should cite it early,
  explain the common proposition idea, and make the different corpus,
  provenance contract, baselines, and evaluation explicit.
- **Open uncertainty — novelty:** A broader database search (Semantic Scholar,
  ACL Anthology, arXiv, and IR venues) is still needed before claiming that no
  prior work compares exactly these four conditions on dialogue.
- **Open uncertainty — corpus transfer:** MediaSum is edited broadcast dialogue;
  it is real human speech but is not private, spontaneous workplace conversation.
  A second corpus such as meetings or public institutional Q&A is needed for a
  general claim about human dialogue.
- **Open uncertainty — facts:** Exact quote location proves provenance, not that
  a normalized fact preserves scope, negation, temporality, or speaker stance.
  The planned human fidelity audit remains necessary.
