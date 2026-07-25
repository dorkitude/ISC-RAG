// Package study2 implements the isolated production-inspired development study.
package study2

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/dorkitude/isc-rag/internal/fireworks"
	"github.com/jackc/pgx/v5/pgxpool"
)

const DirectPrompt = `Return only JSON: {"facts":[{"text":"standalone atomic claim","evidence_spans":["exact source quote"]}]}. Extract explicit claims from the dialogue. Do not infer, summarize, or include questions. Each fact must be standalone, atomic, and retain attribution, hedges, negation, quantities, comparisons, and time limits. Every evidence span must be an exact contiguous source substring.`

// Study2TransferPrompt deliberately transfers mechanics (examples, attribution and
// qualification rules), not its sales ontology or customer/account assumptions.
const Study2TransferPrompt = `You are a precise extractor of explicit facts from a human interview transcript. Return ONLY JSON shaped as {"facts":[{"text":"...","evidence_spans":["..."]}]}.

Extract all salient, explicit claims. A fact is a single subject--relation--target claim that makes sense without the transcript. Speaker labels are provenance, not wording: do not write "the guest said" unless the reported attribution is itself the claim.

Rules:
1. State only what an evidence span explicitly supports. Never use programme knowledge, a summary, or an implication.
2. Fully resolve names, organizations, places, and referents only when the evidence spans establish them. Preserve who believes/reported/said something.
3. Preserve every hedge, negation, condition, comparison, number, date, and time limitation. Do not turn a possibility, prediction, or allegation into a fact.
4. One predicate per card. Split conjunctions only when each resulting card has its own support.
5. Questions, greetings, framing, and claims about the conversation are not facts. A factual answer may be extracted.
6. Evidence spans must be exact contiguous text from SOURCE TURNS. Use one span when sufficient; use up to three only when a second span is necessary to establish an antecedent, reported source, or condition. Do not fabricate connecting text.

Examples:
BAD: "The interview discusses a vaccine." (meta-statement)
BAD: "The policy will reduce costs." when source says "might reduce costs" (lost hedge)
BAD: "She supports it." (unresolved referent)
GOOD: "Dr. Lee said the policy might reduce costs." when a source span names Dr. Lee and says this.
GOOD: "The museum will close on Monday." when that complete commitment is explicit.

Return no duplicate or near-duplicate cards. Prefer coverage across substantive turns before extracting secondary details from the same turn.`

// BundlePrompt is a qualitatively different claim construction strategy: the
// model must build a closed evidence bundle before wording the card.  It is not
// a cosmetic shortening of the direct contract.
const BundlePrompt = `Return ONLY JSON shaped as {"facts":[{"text":"standalone atomic claim","evidence_spans":["exact span 1","exact span 2"]}]}. Work in this order for every card: (1) locate one or more exact source spans, (2) check whether those spans alone establish every name, referent, condition, hedge, number, comparison, date, and reported attribution used in the fact, then (3) write the fact. If any word in the fact needs context that is absent from the evidence bundle, either add the exact antecedent/context span or omit the card.

The evidence bundle is a closed proof, not an illustrative quote. A single span is preferred. Use two or three spans only for a necessary antecedent, named referent, reported source, or condition. Do not turn a question into a fact. Do not use unresolved it/they/he/she/this/that in fact_text. Do not add background knowledge. Preserve modality exactly: might, could, plans to, said, alleged, and according to are meaningful. One card contains one claim. Return no duplicates. Cover distinct substantive turns before secondary details.`

// AttributionPrompt jumps to a different representation: dialogue claims are
// deliberately represented as attributed reports, avoiding the invalid move
// from an interviewee's assertion to an unqualified world fact.
const AttributionPrompt = `Return ONLY JSON: {"facts":[{"text":"Speaker Name said/reported/believed ...","evidence_spans":["exact source quote"]}]}. Extract explicit, substantive declarative claims from SOURCE TURNS. Every fact MUST begin with the exact source speaker label and an attribution verb that matches the source (said, reported, recalled, predicted, believes, or asked). Do not convert a person's statement, opinion, anecdote, allegation, or plan into an unqualified world fact. Preserve all hedges, negation, conditions, dates, numbers, comparisons, and uncertainty. Do not extract questions, introductions, navigation, promotions, or claims merely made by a host while asking a question. One claim per fact; never use unresolved pronouns in the fact. Evidence must be exact contiguous text. Cover different substantive turns before details. Return no duplicates.`

const AtomicAttributionPrompt = `Return ONLY JSON: {"facts":[{"text":"Exact Speaker Label said that a fully-resolved single claim.","evidence_spans":["exact source quote"]}]}. You extract explicit interview claims as attributed reports.

Before emitting each card, perform this hard checklist: (a) the fact has exactly one proposition after "said/reported/believed"; (b) every person, thing, place, and relationship is named, never it/this/that/they/he/she/there; (c) every number, date, comparison, negation, condition, and hedge in the evidence is retained; (d) the selected evidence alone plus the immediately preceding source turn establish the wording; and (e) the source turn is declarative, not a question, introduction, promotion, or host framing.

If a source sentence contains two claims, emit two cards only if each card has separate direct support. If an antecedent cannot be fully resolved without speculation, omit it. Do not compress an anecdote, plan, allegation, opinion, or reported assertion into a bare world fact: state that its exact source speaker said/reported/recalled/believed it. Bad: "Robert Wilkie will stabilize the VA and improve care." Good: "Robert Wilkie said he would work to stabilize the VA." Bad: "The recipe uses it." Good: omit unless the exact source names the ingredient. Do not use conjunctions joining independent claims. Return no duplicates and cover distinct substantive turns first.`

const TurnwiseExtractPrompt = `Return ONLY JSON: {"facts":[{"text":"candidate claim","evidence_spans":["exact quote from TARGET_TURN"]}]}. Extract every explicit substantive declarative claim in TARGET_TURN, up to three claims. Do not extract claims from CONTEXT_ONLY. Preserve exact wording enough to retain hedges, negation, numbers, dates, comparisons, and attribution. Do not extract questions, greetings, promotions, or host framing. Evidence spans must be exact substrings of TARGET_TURN. This is a candidate-generation stage; do not omit a clear claim merely because it needs a previous turn to resolve a name.`

const TurnwiseRewritePrompt = `Return ONLY JSON: {"facts":[{"text":"Exact Speaker Label said/reported/recalled one fully resolved atomic claim.","evidence_spans":["one or more exact source quotes"]}]}. You are the conservative second stage of dialogue fact extraction. Rewrite or reject the supplied candidates using SOURCE_TURNS as the only authority.

Keep a card only if all of these are true: it is one atomic proposition after the attribution verb; every referent is fully named or the card is rejected; every hedge, negation, date, number, comparison, condition, and reported attribution is preserved; the exact evidence spans support the final wording; and it originates in a declarative speaker turn rather than a question or host framing. The fact must begin with the exact source speaker label and a matching attribution verb. Split a candidate only when each new card has direct evidence. Omit rather than infer. Remove duplicates. Return as many defensible cards as possible.`

const ClauseAnchoredPrompt = `Return ONLY JSON: {"facts":[{"text":"one resolved atomic proposition without speaker attribution","evidence_spans":["one exact SOURCE_CLAUSE"]}]}. SOURCE_CLAUSES are a deterministic exhaustive segmentation of one declarative target turn. Consider every supplied clause. For each clause, return one rewritten card only if it can be made standalone, atomic, and fully resolved using CONTEXT_ONLY plus that clause. Otherwise reject it. Never combine clauses. Never use we, us, our, it, this, that, they, them, he, she, there, or any unresolved noun phrase in text. Preserve all negation, hedge, number, date, comparison, condition, and reported attribution in the proposition. Evidence must be the exact clause, not a paraphrase and not context. Do not extract host framing, questions, greetings, or promotion. The system will attach the source speaker deterministically after validation.`

const ProofObligationPrompt = `Return ONLY JSON: {"facts":[{"text":"one fully resolved atomic proposition without speaker attribution","evidence_spans":["one or more exact complete SOURCE_CLAUSES"]}]}. You are a proof-obligation checker, not a free rephraser. SOURCE_CLAUSES are the only clauses you may select as evidence. A final card is allowed only if every resolved entity/referent, predicate, hedge, negation, condition, number, date, comparison, and reported attribution in its text is explicitly supported by one or more selected complete source clauses. Select the clause containing the claim plus every antecedent/qualifier clause needed for the wording; otherwise reject.

Never use we, us, our, it, this, that, they, them, he, she, there, generic you, or an unresolved noun phrase. Never use a synonym or inferred relation that is absent from selected clauses. Do not combine independent claims. Do not emit host framing, questions, greetings, or promotion. Preserve wording conservatively; omission is preferable to an unsupported standalone rewrite. Return no duplicate cards. The system attaches the source speaker deterministically after exact validation.`

const BundleVerifierPrompt = `Return ONLY JSON: {"facts":[{"text":"one fully resolved atomic proposition without speaker attribution","evidence_spans":["one or more exact source quotes"]}]}. You are an independent verifier for high-recall dialogue fact candidates. SOURCE_TURNS are authoritative. For every candidate, either reject it, or emit a repaired card only when its complete final text is entailed by its selected exact evidence bundle.

Requirements for every retained card: exactly one atomic proposition; resolve every referent; preserve speaker attribution, hedge, negation, condition, number, date, comparison, and reported-source language; use no we/us/our/it/this/that/they/them/he/she/there/generic you; and cite every exact span needed to prove the rewritten text. Do not make a bare world fact out of an opinion, anecdote, plan, allegation, or answer. Do not accept a host question or framing as a fact. Deduplicate.

Coverage budget: inspect every candidate turn range. Retain at least one card from each substantive turn that has an explicit defensible claim; rejection is allowed only when no standalone, fully supported card can be made. Do not silently prefer a small set of easy cards. The system assigns the source speaker deterministically after validation.`

const NLIPrompt = `Return ONLY JSON: {"decisions":[{"index":0,"action":"KEEP|DROP|REWRITE","text":"required only for REWRITE","evidence_spans":["required only for REWRITE"]}]}. You are a closed-world claim--evidence adjudicator. Evaluate each immutable high-recall candidate independently against SOURCE_TURNS. You may only KEEP its text unchanged, DROP it, or minimally REWRITE it.

KEEP only if every phrase in the candidate is explicitly entailed by its exact source evidence. DROP if unsupported, non-atomic, a host question/framing claim, or irreparably context-dependent. REWRITE only to remove an unsupported phrase or resolve a reference, and only if every added or changed entity, relation, qualifier, hedge, number, date, comparison, condition, and attribution has an exact evidence span. Never introduce semantic substitutions, causal links, category changes, or inferred relations. Return the complete exact evidence bundle for every REWRITE. Preserve candidate coverage: do not drop the sole defensible candidate from a substantive source turn merely because it is less interesting.`

const DecompositionPrompt = `Return ONLY JSON: {"facts":[{"text":"one atomic standalone proposition without speaker attribution","evidence_spans":["one exact PROPOSITION_UNIT"]}]}. Each PROPOSITION_UNIT is deterministic source text from one eligible declarative turn. Consider every unit separately; never merge units. Copy every MANDATORY_TOKEN exactly into the output card: those tokens encode modality, negation, scope, time, quantity, comparison, and attribution. Resolve a reference only when CONTEXT_ONLY explicitly establishes it, otherwise omit the unit. Evidence must equal a complete exact proposition unit. The system attaches the source speaker after validation.`
const KeepDropPrompt = `Return ONLY JSON: {"decisions":[{"index":0,"action":"KEEP|DROP"}]}. You are a closed-world safety selector. Candidates were derived from one source proposition unit and contain mandatory scope tokens. You may KEEP a card byte-for-byte or DROP it; never rewrite, merge, split, or delete a qualifier. KEEP only fully source-supported, standalone atomic cards. Retain one defensible card per substantive source turn when available.`

const FramePrompt = `Return ONLY JSON: {"frames":[{"clause_id":0,"subject":"","relation":"","object":"","modality":"","negation":"","temporal":"","condition":"","comparison":"","evidence_clause_ids":[0]}]}. Build frames only from ELIGIBLE_SOURCE_UNITS. Every nonempty field must be an exact substring of a selected evidence clause. Do not infer, normalize, or merge units. Put modality, negation, temporal, condition, comparison, and reported attribution into their own fields. Omit ambiguous references, questions, framing, and promotion. Source speaker is deterministic metadata.`

const QuoteRepairPrompt = `Return ONLY JSON: {"decisions":[{"index":0,"defects":["UNRESOLVED_ANTECEDENT|OMITTED_MODIFIER|CONJUNCTION|UNSUPPORTED|NONE"],"action":"KEEP|DROP|REPAIR","text":"required only for REPAIR","evidence_spans":["required only for REPAIR"]}]}. You are a quote-first targeted repair validator. First identify concrete defects in each immutable candidate against exact SOURCE_TURNS. KEEP only defect-free cards exactly as written. DROP unsupported cards and cards whose source modality/scope/negation/time/number/comparison cannot be retained. REPAIR only one unresolved antecedent or conjunction: retain source-clause wording and modifiers, select exact prior span(s) for the antecedent, and return one atomic replacement. Never paraphrase predicates, substitute synonyms, add causal links, or change speaker/agent. Every changed phrase must occur in selected evidence.`

const PredicatePrompt = `Return ONLY JSON: {"facts":[{"text":"one standalone atomic claim without speaker attribution","evidence_spans":["one exact PREDICATE_UNIT"]}]}. PREDICATE_UNITS are an exhaustive deterministic list of source predicate/coordination items from the target declarative turn. Process every unit; return one candidate per explicit defensible predicate unit. Do not merge units or omit a unit because another unit is more salient. Preserve all source modality, negation, scope, number, date, comparison, condition, and attribution. Use CONTEXT_ONLY only to resolve a referent with exact support; otherwise return the original unit as a candidate. Evidence must equal the exact predicate unit.`

const SubstantivePrompt = `Return ONLY JSON: {"facts":[{"text":"one standalone atomic claim without speaker attribution","evidence_spans":["one exact CLAIM_UNIT plus exact antecedent span when needed"]}]}. CLAIM_UNITS are deterministic source units prefiltered to remove greetings, acknowledgements, commands, promotion, questions, and framing. Render each explicit declarative claim as a standalone sentence only when every referent can be anchored by exact source evidence. Preserve every source hedge, negation, number, date, comparison, condition, and attribution. Omit rather than emit dialogue fragments. Do not merge claim units. Evidence must include every antecedent used for a replacement.`

// FableC4GeneratorPrompt is the frozen natural-language antecedent-window
// generator used by the C4 structural replication.  The output is JSON only
// because this implementation invokes the provider API rather than a file
// writing subagent; its semantic contract is otherwise fixed before gate runs.
const FableC4GeneratorPrompt = `You are compiling provenance-exact fact cards from a broadcast-dialogue transcript. Return ONLY JSON: {"facts":[{"text":"...","primary_quote":"...","antecedent_quote":"... or -","primary_turn":0}]}.

Sweep SOURCE_TURNS in order and produce one fact card for every explicit substantive claim: world facts and events, a speaker's stated position, prediction, estimate, explanation, or first-person experience. Skip greetings, thanks, host housekeeping, bare questions, and filler.

Rules for fact text: (1) standalone: resolve every pronoun and deictic to an explicit named referent from the transcript; convert relative time to an explicit anchor or omit it; never leave a dangling comparative. (2) atomic: exactly one predicate-level claim per card; split compounds. (3) preserve every negation, hedge, quantity, condition, temporal scope, and reported source exactly. (4) if the claim is a stance, prediction, estimate, or first-person report, name the speaker in the fact text. (5) use explicit source content only; no inference, outside knowledge, meta-claims, or question-only claims.

Evidence rules: primary_quote is a verbatim contiguous copy under 240 characters from the one source turn identified by primary_turn, and contains the core assertion including its operator. antecedent_quote is - unless an earlier span is needed to resolve a referent or context; when used it is a verbatim contiguous copy under 240 characters from primary_turn or one of the four preceding turns. Every element of the fact text must be entailed by these two quotes together with their speaker labels. If you cannot quote the evidence, do not emit the card. Return at most 60 cards and no duplicates.`

// FableC4AuditPrompt deliberately has no transcript.  The quoted spans and
// canonical speaker labels are its whole evidence universe.
const FableC4AuditPrompt = `Return ONLY JSON: {"audits":[{"index":0,"status":"PASS|FAIL","missing_elements":["..."]}]}. You audit fact cards for SUPPORT ONLY and do not have the source transcript. For each card, every entity, predicate, number, date, place, descriptor, possessive, and causal link in text must be stated or strictly entailed by primary_quote plus antecedent_quote (if any), read with the supplied speaker labels. Outside knowledge and plausibility never count. Paraphrase is allowed, but if genuinely uncertain, FAIL and list the ungrounded elements. Judge neither style, atomicity, nor duplication. Return exactly one audit for every input index.`

// FableC4RepairPrompt restricts repair to the C4 moves.  It is intentionally
// narrower than prior Study #2 repair prompts: no predicate rewrite, merge,
// split, or new claim is permitted.
const FableC4RepairPrompt = `Return ONLY JSON: {"repairs":[{"index":0,"action":"REANCHOR|RECUT|NAME_FORM|DROP","text":"...","primary_quote":"...","antecedent_quote":"... or -"}]}. You repair only fact cards that failed a strict span-blind support audit. The primary turn is frozen. Apply exactly one allowed move: REANCHOR keeps fact text and primary quote unchanged, replacing only antecedent_quote with an exact span from the primary turn or one of its four predecessors; RECUT keeps fact text and antecedent_quote unchanged but extends or shifts primary_quote within the same primary turn without losing qualifiers; NAME_FORM keeps both quotes unchanged and minimally rewrites only a quoted name or name form so the quote bundle entails it. DROP if none works. Never change a predicate, scope, hedge, negation, quantity, condition, attribution, or causal link. Never merge or split cards. The result must remain standalone and atomic. Verify every output quote is an exact substring of the supplied transcript before emitting it. Return one repair per input index.`

// DialogueActPrompt operates over deterministic question--answer source units,
// not a free transcript window. A response is represented as the answering
// speaker's dialogue act when the question supplies its indispensable scope.
const DialogueActPrompt = `You extract provenance-exact fact cards from deterministic dialogue source units. Return ONLY JSON: {"facts":[{"unit_id":0,"text":"...","answer_quote":"...","question_quote":"... or -"}]}.

Each SOURCE_UNIT is either QUESTION_ANSWER (one interviewer question plus its immediate declarative answer) or STANDALONE (one unpaired declarative turn). Extract every explicit substantive answer/assertion, skipping greetings, housekeeping, bare questions, acknowledgements, and promotion.

For QUESTION_ANSWER units, make each card an attributed dialogue act: name the ANSWER speaker and say that they answered, said, estimated, believed, recalled, or explained a single claim. Do not turn an answer into an unqualified world claim. If the question supplies an essential referent, comparison, timeframe, or requested topic, use the exact question wording only as needed to make the answer standalone and cite it in question_quote. The question itself is not asserted as true. For STANDALONE units, attribute a stance, prediction, estimate, anecdote, or reported assertion to its source speaker; a neutral directly asserted world fact may remain neutral.

Every card is exactly one predicate-level claim. Preserve every answer-side hedge, negation, number, condition, time scope, comparison, and reported-source framing. Never add a relation, cause, identity, category, or universal claim absent from the quotes. answer_quote must be a verbatim contiguous span from the unit's ANSWER/TURN and contain the claim operator. question_quote must be - for STANDALONE, and otherwise either - or a verbatim contiguous span from the unit's QUESTION. The two quoted spans plus their speaker labels must entail every word in the final card. Return no duplicates and at most 60 cards.`

// DialoguePredicatePrompt supplies one deterministic answer predicate unit at
// a time. This changes source representation to address the remaining
// compound-card/fragments failure without relaxing Q/A provenance.
const DialoguePredicatePrompt = `You extract provenance-exact fact cards from deterministic dialogue predicate units. Return ONLY JSON: {"facts":[{"unit_id":0,"text":"...","answer_quote":"...","question_quote":"... or -"}]}.

Each SOURCE_UNIT has exactly one ANSWER_CLAIM_UNIT: a deterministic lexical predicate/coordination item from a speaker turn. Consider every unit and emit at most one card for it. Do not merge units. Skip an item only if it contains no explicit substantive claim.

Render one fully standalone, atomic proposition. For a QUESTION_ANSWER unit, state the answer speaker's dialogue act: the speaker answered, explained, estimated, believed, recalled, or said one claim. Never turn an answer into an unqualified world claim. Use the QUESTION only to resolve an indispensable topic, referent, timeframe, or comparison; when you do, cite that exact span in question_quote. A question is context, not an asserted premise. For a STANDALONE unit, preserve source attribution for a stance, prediction, estimate, anecdote, or reported assertion.

The card may express ONE predicate-level claim only. Do not join independent claims with and/or/but; do not leave a pronoun, deictic, relative time, or dangling comparative. Preserve every hedge, negation, quantity, condition, temporal scope, comparison, and reported-source framing from ANSWER_CLAIM_UNIT exactly. answer_quote must be a verbatim contiguous span inside ANSWER_CLAIM_UNIT, never from another answer clause. question_quote must be - for STANDALONE, and otherwise either - or a verbatim contiguous span from QUESTION. The quoted bundle plus speaker labels must entail every word. Never infer or complete a truncated phrase. Return no duplicates and at most 80 cards.`

const QuoteCapsulePrompt = `Return ONLY JSON: {"cards":[{"unit_id":0,"answer_start":0,"answer_end":0,"question_start":-1,"question_end":-1}]}. You select fact capsules from deterministic dialogue source units. Do not write fact text, quotes, paraphrases, explanations, or any other free-form text.

For each SOURCE_UNIT, select zero or more semantically atomic ANSWER spans using zero-based UTF-8 byte offsets [start,end) into answer_text. Select one span per distinct explicit substantive proposition; split distinct propositions into separate spans. Do not select greetings, acknowledgements, questions, promotion, fragments, or a span whose meaning depends on omitted answer text. Preserve the exact source wording including hedge, negation, quantity, condition, time, comparison, and attribution.

For QUESTION_ANSWER units, optionally select one UTF-8 byte range [question_start,question_end) from question_text only when it supplies indispensable context for interpreting the selected answer. The question is context, never an asserted premise. For STANDALONE units question_start/question_end must be -1. Every boundary must fall exactly between UTF-8 characters. Return at most 80 cards and no duplicate answer spans.`

const QuoteCapsuleExactPrompt = `Return ONLY JSON: {"cards":[{"unit_id":0,"answer_quote":"exact answer span","question_quote":"exact question span or -"}]}. You select fact capsules from deterministic dialogue source units. Do not write fact text, paraphrases, explanations, or any other free-form text.

For each SOURCE_UNIT, select zero or more semantically atomic, verbatim contiguous ANSWER spans. One answer quote is one card; split distinct propositions into separate quotes. Do not select greetings, acknowledgements, questions, promotion, fragments, or a span whose meaning depends on omitted answer text. Preserve exact source wording and every hedge, negation, quantity, condition, time, comparison, and attribution.

For QUESTION_ANSWER units, question_quote is - unless an exact contiguous question span supplies indispensable context for the selected answer. The question is context, not an asserted premise. For STANDALONE units question_quote must be -. Copy quotes character-for-character from the unit text. Return at most 80 cards and no duplicate answer quotes.`

type Config struct {
	StudyVersion, Candidate, ExternalID, Model string
	Replicate, MaxFacts, ChunkChars, Carryover int
}
type Result struct {
	RunID                                                     string
	Raw, Validated, PromptTokens, CompletionTokens, LatencyMS int
	PromptHash                                                string
}
type turn struct {
	ordinal, start, end int
	speaker, text       string
}
type candidate struct {
	Text     string   `json:"text"`
	Evidence []string `json:"evidence_spans"`
}
type response struct {
	Facts []candidate `json:"facts"`
}
type c4Candidate struct {
	Text            string `json:"text"`
	PrimaryQuote    string `json:"primary_quote"`
	AntecedentQuote string `json:"antecedent_quote"`
	PrimaryTurn     int    `json:"primary_turn"`
}
type c4Response struct {
	Facts []c4Candidate `json:"facts"`
}
type c4Card struct {
	Index     int
	Candidate c4Candidate
	Fact      fact
}
type c4Audit struct {
	Index           int      `json:"index"`
	Status          string   `json:"status"`
	MissingElements []string `json:"missing_elements"`
}
type c4AuditResponse struct {
	Audits []c4Audit `json:"audits"`
}
type c4Repair struct {
	Index           int    `json:"index"`
	Action          string `json:"action"`
	Text            string `json:"text"`
	PrimaryQuote    string `json:"primary_quote"`
	AntecedentQuote string `json:"antecedent_quote"`
}
type c4RepairResponse struct {
	Repairs []c4Repair `json:"repairs"`
}
type dialogueUnit struct {
	ID       int
	Question *turn
	Answer   turn
	Paired   bool
}
type dialogueCandidate struct {
	UnitID        int    `json:"unit_id"`
	Text          string `json:"text"`
	AnswerQuote   string `json:"answer_quote"`
	QuestionQuote string `json:"question_quote"`
}
type dialogueResponse struct {
	Facts []dialogueCandidate `json:"facts"`
}
type dialoguePredicateUnit struct {
	ID        int
	Dialogue  dialogueUnit
	ClaimUnit string
}
type capsuleCard struct {
	UnitID        int `json:"unit_id"`
	AnswerStart   int `json:"answer_start"`
	AnswerEnd     int `json:"answer_end"`
	QuestionStart int `json:"question_start"`
	QuestionEnd   int `json:"question_end"`
}
type capsuleResponse struct {
	Cards []capsuleCard `json:"cards"`
}
type capsuleQuoteCard struct {
	UnitID        int    `json:"unit_id"`
	AnswerQuote   string `json:"answer_quote"`
	QuestionQuote string `json:"question_quote"`
}
type capsuleQuoteResponse struct {
	Cards []capsuleQuoteCard `json:"cards"`
}
type nliDecision struct {
	Index    int      `json:"index"`
	Action   string   `json:"action"`
	Text     string   `json:"text"`
	Evidence []string `json:"evidence_spans"`
	Defects  []string `json:"defects"`
}
type nliResponse struct {
	Decisions []nliDecision `json:"decisions"`
}
type frame struct {
	ClauseID   int    `json:"clause_id"`
	Subject    string `json:"subject"`
	Relation   string `json:"relation"`
	Object     string `json:"object"`
	Modality   string `json:"modality"`
	Negation   string `json:"negation"`
	Temporal   string `json:"temporal"`
	Condition  string `json:"condition"`
	Comparison string `json:"comparison"`
	Evidence   []int  `json:"evidence_clause_ids"`
}
type frameResponse struct {
	Frames []frame `json:"frames"`
}
type fact struct {
	text       string
	evidence   []evidence
	start, end int
	speaker    string
}
type evidence struct {
	quote                   string
	start, end, first, last int
	speaker                 string
}

// Some OpenAI-compatible deployments reject reasoning_effort altogether.
// Preserve zero-temperature JSON generation for those models while retaining
// explicit no-reasoning requests for the DeepSeek study condition.
func studyChat(ctx context.Context, client *fireworks.Client, model, system, user string, max int) (fireworks.ChatResult, error) {
	if strings.Contains(model, "gpt-oss") {
		return client.ChatJSON(ctx, model, system, user, max)
	}
	return client.ChatJSONNoReasoning(ctx, model, system, user, max)
}

func Run(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, cfg Config) (Result, error) {
	if cfg.StudyVersion == "" || cfg.Candidate == "" || cfg.ExternalID == "" || cfg.Model == "" || cfg.Replicate < 1 {
		return Result{}, fmt.Errorf("study-version, candidate, external-id, model and positive replicate are required")
	}
	if cfg.Candidate == "fable-c4-sbar" {
		return runFableC4(ctx, pool, client, cfg)
	}
	if cfg.Candidate == "dialogue-act-pairs" {
		return runDialogueActPairs(ctx, pool, client, cfg)
	}
	if cfg.Candidate == "dialogue-act-predicate-units" {
		return runDialogueActPredicateUnits(ctx, pool, client, cfg)
	}
	if cfg.Candidate == "quote-capsules" {
		return runQuoteCapsules(ctx, pool, client, cfg)
	}
	if cfg.Candidate == "quote-capsules-exact" {
		return runQuoteCapsulesExact(ctx, pool, client, cfg)
	}
	if cfg.MaxFacts <= 0 {
		cfg.MaxFacts = 45
	}
	if cfg.ChunkChars <= 0 {
		cfg.ChunkChars = 10000
	}
	var doc string
	if err := pool.QueryRow(ctx, `SELECT id FROM documents WHERE corpus='mediasum' AND corpus_split='validation' AND external_id=$1`, cfg.ExternalID).Scan(&doc); err != nil {
		return Result{}, fmt.Errorf("validation document %q: %w", cfg.ExternalID, err)
	}
	transcript, turns, err := source(ctx, pool, doc)
	if err != nil {
		return Result{}, err
	}
	prompt := DirectPrompt
	if cfg.Candidate == "bundle" || cfg.Candidate == "expanded-bundle" {
		prompt = BundlePrompt
	} else if cfg.Candidate == "attributed-packed" {
		prompt = AttributionPrompt
	} else if cfg.Candidate == "atomic-attributed" {
		prompt = AtomicAttributionPrompt
	} else if cfg.Candidate == "turnwise-staged" || cfg.Candidate == "turnwise-local-staged" {
		prompt = TurnwiseExtractPrompt
	} else if cfg.Candidate == "clause-anchored" {
		prompt = ClauseAnchoredPrompt
	} else if cfg.Candidate == "proof-obligation" {
		prompt = ProofObligationPrompt
	} else if cfg.Candidate == "highrecall-verified" {
		prompt = Study2TransferPrompt
	} else if cfg.Candidate == "highrecall-nli" {
		prompt = Study2TransferPrompt
	} else if cfg.Candidate == "prefilter-nli" {
		prompt = Study2TransferPrompt
	} else if cfg.Candidate == "decomposed-nli" {
		prompt = DecompositionPrompt
	} else if cfg.Candidate == "framed-render" {
		prompt = FramePrompt
	} else if cfg.Candidate == "quote-repair" {
		prompt = Study2TransferPrompt
	} else if cfg.Candidate == "predicate-repair" {
		prompt = PredicatePrompt
	} else if cfg.Candidate == "substantive-render" {
		prompt = SubstantivePrompt
	} else if cfg.Candidate == "explicit" || cfg.Candidate == "carryover" || cfg.Candidate == "speaker-packed" || cfg.Candidate == "expanded-packed" {
		prompt = Study2TransferPrompt
	}
	modelTurns := turns
	if cfg.Candidate == "prefilter-nli" || cfg.Candidate == "decomposed-nli" || cfg.Candidate == "framed-render" || cfg.Candidate == "quote-repair" || cfg.Candidate == "predicate-repair" || cfg.Candidate == "substantive-render" {
		modelTurns = eligibleTurns(turns)
	}
	chunks := [][]turn{modelTurns}
	if cfg.Candidate == "turnwise-staged" || cfg.Candidate == "turnwise-local-staged" || cfg.Candidate == "clause-anchored" || cfg.Candidate == "proof-obligation" {
		chunks = turnWindows(turns)
	}
	if cfg.Candidate == "decomposed-nli" || cfg.Candidate == "framed-render" {
		chunks = turnWindows(modelTurns)
	}
	if cfg.Candidate == "predicate-repair" {
		chunks = turnWindows(modelTurns)
	}
	if cfg.Candidate == "substantive-render" {
		chunks = turnWindows(modelTurns)
	}
	if cfg.Candidate == "speaker-packed" || cfg.Candidate == "carryover" || cfg.Candidate == "bundle" || cfg.Candidate == "highrecall-verified" || cfg.Candidate == "highrecall-nli" || cfg.Candidate == "prefilter-nli" || cfg.Candidate == "quote-repair" {
		chunks = pack(modelTurns, cfg.ChunkChars, cfg.Carryover)
	}
	all := []candidate{}
	raw := []string{}
	var pt, ct int
	started := time.Now()
	type chunkResult struct {
		facts  []candidate
		raw    string
		pt, ct int
		err    error
	}
	results := make([]chunkResult, len(chunks))
	var workers sync.WaitGroup
	// Bound concurrency to avoid a burst that would confound the study through
	// rate-limit retries, while still making turn-wise extraction runnable.
	workerCount := 8
	if cfg.Candidate == "framed-render" || cfg.Candidate == "predicate-repair" {
		// Frame extraction is one compact response per target turn; a wider,
		// bounded fan-out avoids losing the whole document to the local command
		// window while preserving per-turn independent requests.
		workerCount = 20
	}
	limit := make(chan struct{}, workerCount)
	for i, input := range chunks {
		workers.Add(1)
		go func(i int, chunk []turn) {
			defer workers.Done()
			limit <- struct{}{}
			defer func() { <-limit }()
			labeled := label(chunk)
			user := fmt.Sprintf("<SOURCE_METADATA>programme: MediaSum interview; turn range: %d-%d; speaker labels are exact source metadata</SOURCE_METADATA>\n<SOURCE_TURNS>\n%s\n</SOURCE_TURNS>\nReturn at most %d cards.", chunk[0].ordinal, chunk[len(chunk)-1].ordinal, labeled, cfg.MaxFacts)
			if cfg.Candidate == "turnwise-staged" || cfg.Candidate == "turnwise-local-staged" {
				context := chunk[:len(chunk)-1]
				target := chunk[len(chunk)-1]
				user = fmt.Sprintf("<SOURCE_METADATA>programme: MediaSum interview; source speaker labels are exact metadata</SOURCE_METADATA>\n<CONTEXT_ONLY>\n%s\n</CONTEXT_ONLY>\n<TARGET_TURN number=\"%d\" speaker=\"%s\">\n%s\n</TARGET_TURN>\nExtract only claims explicitly stated in TARGET_TURN. CONTEXT_ONLY may resolve a referent but must not supply a new claim. Return at most %d cards.", label(context), target.ordinal, target.speaker, target.text, cfg.MaxFacts)
			}
			if cfg.Candidate == "clause-anchored" || cfg.Candidate == "proof-obligation" {
				context := chunk[:len(chunk)-1]
				target := chunk[len(chunk)-1]
				clauses := append([]string{}, splitContextClauses(context)...)
				clauses = append(clauses, splitClauses(target.text)...)
				clausesJSON, _ := json.Marshal(clauses)
				user = fmt.Sprintf("<CONTEXT_ONLY>\n%s\n</CONTEXT_ONLY>\n<TARGET_TURN number=\"%d\" speaker=\"%s\">\n%s\n</TARGET_TURN>\n<SOURCE_CLAUSES>\n%s\n</SOURCE_CLAUSES>\nEvaluate every target source clause. Return at most %d cards.", label(context), target.ordinal, target.speaker, target.text, clausesJSON, cfg.MaxFacts)
			}
			if cfg.Candidate == "decomposed-nli" {
				context := chunk[:len(chunk)-1]
				target := chunk[len(chunk)-1]
				units := propositionUnits(target.text)
				payload := make([]map[string]any, len(units))
				for j, unit := range units {
					payload[j] = map[string]any{"text": unit, "mandatory_tokens": mandatoryTokens(unit)}
				}
				encoded, _ := json.Marshal(payload)
				user = fmt.Sprintf("<CONTEXT_ONLY>\n%s\n</CONTEXT_ONLY>\n<TARGET_TURN number=\"%d\" speaker=\"%s\">\n%s\n</TARGET_TURN>\n<PROPOSITION_UNITS>\n%s\n</PROPOSITION_UNITS>\nEvaluate every unit; return at most %d cards.", label(context), target.ordinal, target.speaker, target.text, encoded, cfg.MaxFacts)
			}
			if cfg.Candidate == "framed-render" {
				context := chunk[:len(chunk)-1]
				target := chunk[len(chunk)-1]
				units := propositionUnits(target.text)
				payload := make([]map[string]any, len(units))
				for j, unit := range units {
					payload[j] = map[string]any{"clause_id": j, "text": unit}
				}
				encoded, _ := json.Marshal(payload)
				user = fmt.Sprintf("<CONTEXT_ONLY>\n%s\n</CONTEXT_ONLY>\n<ELIGIBLE_SOURCE_UNITS>\n%s\n</ELIGIBLE_SOURCE_UNITS>\nReturn frames only.", label(context), encoded)
			}
			if cfg.Candidate == "predicate-repair" {
				context := chunk[:len(chunk)-1]
				target := chunk[len(chunk)-1]
				units, _ := json.Marshal(predicateUnits(target.text))
				user = fmt.Sprintf("<CONTEXT_ONLY>\n%s\n</CONTEXT_ONLY>\n<TARGET_TURN>\n%s\n</TARGET_TURN>\n<PREDICATE_UNITS>\n%s\n</PREDICATE_UNITS>\nReturn one candidate per unit; do not apply a global cap.", label(context), target.text, units)
			}
			if cfg.Candidate == "substantive-render" {
				context := chunk[:len(chunk)-1]
				target := chunk[len(chunk)-1]
				units, _ := json.Marshal(substantiveUnits(target.text))
				user = fmt.Sprintf("<CONTEXT_ONLY>\n%s\n</CONTEXT_ONLY>\n<CLAIM_UNITS>\n%s\n</CLAIM_UNITS>\nReturn only standalone claim cards.", label(context), units)
			}
			r, e := studyChat(ctx, client, cfg.Model, prompt, user, 6000)
			if e != nil {
				results[i].err = e
				return
			}
			if cfg.Candidate == "framed-render" {
				frames, frameErr := decodeFrames(r.Text)
				if frameErr != nil {
					results[i].err = frameErr
					return
				}
				results[i] = chunkResult{facts: renderFrames(frames, propositionUnits(chunk[len(chunk)-1].text)), raw: r.Text, pt: r.PromptTokens, ct: r.CompletionTokens}
				return
			}
			p, e := decode(r.Text)
			if e != nil {
				results[i].err = e
				return
			}
			if cfg.Candidate == "turnwise-local-staged" {
				payload, _ := json.Marshal(map[string]any{"candidate_cards": p.Facts})
				second, secondErr := studyChat(ctx, client, cfg.Model, TurnwiseRewritePrompt, user+"\n<CANDIDATE_CARDS>\n"+string(payload)+"\n</CANDIDATE_CARDS>\nRewrite or reject only these target-turn candidates.", 3000)
				if secondErr != nil {
					results[i].err = secondErr
					return
				}
				secondParsed, secondErr := decode(second.Text)
				if secondErr != nil {
					results[i].err = secondErr
					return
				}
				results[i] = chunkResult{facts: secondParsed.Facts, raw: r.Text + "\n---LOCAL-REWRITE---\n" + second.Text, pt: r.PromptTokens + second.PromptTokens, ct: r.CompletionTokens + second.CompletionTokens}
				return
			}
			results[i] = chunkResult{facts: p.Facts, raw: r.Text, pt: r.PromptTokens, ct: r.CompletionTokens}
		}(i, input)
	}
	workers.Wait()
	for _, r := range results {
		if r.err != nil {
			return Result{}, r.err
		}
		pt += r.pt
		ct += r.ct
		raw = append(raw, r.raw)
		all = append(all, r.facts...)
	}
	if cfg.Candidate == "turnwise-staged" {
		payload, _ := json.Marshal(map[string]any{"candidate_cards": all})
		r, e := studyChat(ctx, client, cfg.Model, TurnwiseRewritePrompt, "<SOURCE_TURNS>\n"+label(turns)+"\n</SOURCE_TURNS>\n<CANDIDATE_CARDS>\n"+string(payload)+"\n</CANDIDATE_CARDS>\nReturn only revised defensible cards.", 6000)
		if e != nil {
			return Result{}, e
		}
		pt += r.PromptTokens
		ct += r.CompletionTokens
		raw = append(raw, "---REWRITE---\n"+r.Text)
		p, e := decode(r.Text)
		if e != nil {
			return Result{}, e
		}
		all = p.Facts
	}
	if cfg.Candidate == "highrecall-verified" {
		payload, _ := json.Marshal(map[string]any{"high_recall_candidates": all})
		r, e := studyChat(ctx, client, cfg.Model, BundleVerifierPrompt, "<SOURCE_TURNS>\n"+label(turns)+"\n</SOURCE_TURNS>\n<HIGH_RECALL_CANDIDATES>\n"+string(payload)+"\n</HIGH_RECALL_CANDIDATES>\nVerify every candidate turn range; return only defensible repaired cards.", 2500)
		if e != nil {
			return Result{}, e
		}
		pt += r.PromptTokens
		ct += r.CompletionTokens
		raw = append(raw, "---BUNDLE-VERIFIER---\n"+r.Text)
		p, e := decode(r.Text)
		if e != nil {
			return Result{}, e
		}
		all = p.Facts
	}
	if cfg.Candidate == "highrecall-nli" {
		items := make([]map[string]any, len(all))
		for i, c := range all {
			items[i] = map[string]any{"index": i, "text": c.Text, "evidence_spans": c.Evidence}
		}
		payload, _ := json.Marshal(map[string]any{"candidates": items})
		r, e := studyChat(ctx, client, cfg.Model, NLIPrompt, "<SOURCE_TURNS>\n"+label(turns)+"\n</SOURCE_TURNS>\n<CANDIDATES>\n"+string(payload)+"\n</CANDIDATES>\nReturn one decision per candidate index.", 2500)
		if e != nil {
			return Result{}, e
		}
		pt += r.PromptTokens
		ct += r.CompletionTokens
		raw = append(raw, "---NLI-ADJUDICATION---\n"+r.Text)
		decisions, e := decodeNLI(r.Text)
		if e != nil {
			return Result{}, e
		}
		all = applyNLI(all, decisions)
	}
	if cfg.Candidate == "prefilter-nli" {
		items := make([]map[string]any, len(all))
		for i, c := range all {
			items[i] = map[string]any{"index": i, "text": c.Text, "evidence_spans": c.Evidence}
		}
		payload, _ := json.Marshal(map[string]any{"candidates": items})
		r, e := studyChat(ctx, client, cfg.Model, NLIPrompt, "<ELIGIBLE_DECLARATIVE_SOURCE_TURNS>\n"+label(modelTurns)+"\n</ELIGIBLE_DECLARATIVE_SOURCE_TURNS>\n<CANDIDATES>\n"+string(payload)+"\n</CANDIDATES>\nQuestion/framing/promotion turns have already been removed. Return one decision per candidate index.", 2500)
		if e != nil {
			return Result{}, e
		}
		pt += r.PromptTokens
		ct += r.CompletionTokens
		raw = append(raw, "---PREFILTER-NLI---\n"+r.Text)
		decisions, e := decodeNLI(r.Text)
		if e != nil {
			return Result{}, e
		}
		all = applyNLI(all, decisions)
	}
	if cfg.Candidate == "quote-repair" {
		items := make([]map[string]any, len(all))
		for i, c := range all {
			items[i] = map[string]any{"index": i, "text": c.Text, "evidence_spans": c.Evidence}
		}
		payload, _ := json.Marshal(map[string]any{"candidates": items})
		r, e := studyChat(ctx, client, cfg.Model, QuoteRepairPrompt, "<ELIGIBLE_SOURCE_TURNS>\n"+label(modelTurns)+"\n</ELIGIBLE_SOURCE_TURNS>\n<CANDIDATES>\n"+string(payload)+"\n</CANDIDATES>\nReturn one defect decision per candidate index.", 2500)
		if e != nil {
			return Result{}, e
		}
		pt += r.PromptTokens
		ct += r.CompletionTokens
		raw = append(raw, "---QUOTE-REPAIR---\n"+r.Text)
		decisions, e := decodeNLI(r.Text)
		if e != nil {
			return Result{}, e
		}
		all = applyNLI(all, decisions)
	}
	if cfg.Candidate == "predicate-repair" {
		var repaired []candidate
		for start := 0; start < len(all); start += 20 {
			end := start + 20
			if end > len(all) {
				end = len(all)
			}
			batch := all[start:end]
			payload, _ := json.Marshal(map[string]any{"candidates": batch})
			ordinal := start/20 + 1
			var checkpointCandidates, checkpointRepair *string
			checkpointErr := pool.QueryRow(ctx, `SELECT candidate_json::text,repair_json::text FROM study_two_batch_checkpoints WHERE study_version=$1 AND candidate=$2 AND replicate=$3 AND document_id=$4 AND batch_ordinal=$5`, cfg.StudyVersion, cfg.Candidate, cfg.Replicate, doc, ordinal).Scan(&checkpointCandidates, &checkpointRepair)
			if checkpointErr == nil && checkpointRepair != nil && *checkpointRepair != "" {
				var saved struct {
					Candidates []candidate `json:"candidates"`
				}
				if checkpointCandidates != nil {
					_ = json.Unmarshal([]byte(*checkpointCandidates), &saved)
				}
				decisions, decodeErr := decodeNLI(*checkpointRepair)
				if decodeErr != nil {
					return Result{}, decodeErr
				}
				repaired = append(repaired, applyNLI(saved.Candidates, decisions)...)
				raw = append(raw, "---PREDICATE-REPAIR-RESUMED---\n"+*checkpointRepair)
				continue
			}
			if _, e := pool.Exec(ctx, `INSERT INTO study_two_batch_checkpoints(study_version,candidate,replicate,document_id,batch_ordinal,config,candidate_json) VALUES($1,$2,$3,$4,$5,jsonb_build_object('batch_size',20),$6) ON CONFLICT(study_version,candidate,replicate,document_id,batch_ordinal) DO UPDATE SET candidate_json=EXCLUDED.candidate_json,repair_json=NULL,created_at=now()`, cfg.StudyVersion, cfg.Candidate, cfg.Replicate, doc, ordinal, string(payload)); e != nil {
				return Result{}, e
			}
			r, e := studyChat(ctx, client, cfg.Model, QuoteRepairPrompt, "<ELIGIBLE_SOURCE_TURNS>\n"+label(modelTurns)+"\n</ELIGIBLE_SOURCE_TURNS>\n<CANDIDATES>\n"+string(payload)+"\n</CANDIDATES>\nReturn one defect decision per candidate index.", 2500)
			if e != nil {
				return Result{}, e
			}
			pt += r.PromptTokens
			ct += r.CompletionTokens
			raw = append(raw, "---PREDICATE-REPAIR---\n"+r.Text)
			decisions, e := decodeNLI(r.Text)
			if e != nil {
				return Result{}, e
			}
			kept := applyNLI(batch, decisions)
			repaired = append(repaired, kept...)
			if _, e = pool.Exec(ctx, `UPDATE study_two_batch_checkpoints SET repair_json=$6 WHERE study_version=$1 AND candidate=$2 AND replicate=$3 AND document_id=$4 AND batch_ordinal=$5`, cfg.StudyVersion, cfg.Candidate, cfg.Replicate, doc, ordinal, r.Text); e != nil {
				return Result{}, e
			}
		}
		all = repaired
	}
	if cfg.Candidate == "decomposed-nli" {
		items := make([]map[string]any, len(all))
		for i, c := range all {
			items[i] = map[string]any{"index": i, "text": c.Text, "evidence_spans": c.Evidence}
		}
		payload, _ := json.Marshal(map[string]any{"candidates": items})
		r, e := studyChat(ctx, client, cfg.Model, KeepDropPrompt, "<ELIGIBLE_DECLARATIVE_SOURCE_TURNS>\n"+label(modelTurns)+"\n</ELIGIBLE_DECLARATIVE_SOURCE_TURNS>\n<CANDIDATES>\n"+string(payload)+"\n</CANDIDATES>\nReturn one decision per candidate index.", 2500)
		if e != nil {
			return Result{}, e
		}
		pt += r.PromptTokens
		ct += r.CompletionTokens
		raw = append(raw, "---KEEP-DROP-NLI---\n"+r.Text)
		decisions, e := decodeNLI(r.Text)
		if e != nil {
			return Result{}, e
		}
		all = applyNLI(all, decisions)
	}
	validated := validate(all, transcript, turns, cfg.Candidate == "expanded-bundle" || cfg.Candidate == "expanded-packed" || cfg.Candidate == "attributed-packed" || cfg.Candidate == "atomic-attributed" || cfg.Candidate == "turnwise-staged" || cfg.Candidate == "turnwise-local-staged" || cfg.Candidate == "clause-anchored" || cfg.Candidate == "proof-obligation" || cfg.Candidate == "highrecall-verified" || cfg.Candidate == "highrecall-nli" || cfg.Candidate == "prefilter-nli" || cfg.Candidate == "decomposed-nli" || cfg.Candidate == "framed-render" || cfg.Candidate == "quote-repair" || cfg.Candidate == "predicate-repair", cfg.Candidate == "clause-anchored" || cfg.Candidate == "proof-obligation" || cfg.Candidate == "highrecall-verified" || cfg.Candidate == "highrecall-nli" || cfg.Candidate == "prefilter-nli" || cfg.Candidate == "decomposed-nli" || cfg.Candidate == "framed-render" || cfg.Candidate == "quote-repair" || cfg.Candidate == "predicate-repair")
	if cfg.Candidate == "proof-obligation" {
		validated = bindProofClauses(validated, turns)
		validated = proofFilter(validated, turns)
	}
	latency := int(time.Since(started).Milliseconds())
	policy := "all evidence spans exact contiguous source substrings"
	if cfg.Candidate == "expanded-bundle" || cfg.Candidate == "expanded-packed" || cfg.Candidate == "attributed-packed" || cfg.Candidate == "atomic-attributed" || cfg.Candidate == "turnwise-staged" || cfg.Candidate == "turnwise-local-staged" || cfg.Candidate == "clause-anchored" || cfg.Candidate == "proof-obligation" || cfg.Candidate == "highrecall-verified" || cfg.Candidate == "highrecall-nli" || cfg.Candidate == "prefilter-nli" || cfg.Candidate == "decomposed-nli" {
		policy += "; deterministic proof expansion: evidence turn plus two preceding source turns"
	}
	if cfg.Candidate == "turnwise-staged" || cfg.Candidate == "turnwise-local-staged" || cfg.Candidate == "clause-anchored" || cfg.Candidate == "proof-obligation" {
		policy += "; stage one targets declarative turns only; stage two rewrites or rejects all candidates"
	}
	if cfg.Candidate == "highrecall-verified" {
		policy += "; high-recall speaker-packed stage followed by independent full-source bundle verifier with per-turn coverage budget"
	}
	if cfg.Candidate == "highrecall-nli" {
		policy += "; high-recall speaker-packed stage followed by closed-world KEEP/DROP/minimal-REWRITE NLI adjudication with a per-turn coverage budget"
	}
	if cfg.Candidate == "prefilter-nli" {
		policy += "; deterministic interrogative/framing/promotion prefilter, high-recall eligible-turn extraction, then closed-world NLI adjudication with deterministic source speaker assignment"
	}
	if cfg.Candidate == "decomposed-nli" {
		policy += "; deterministic eligible-turn proposition decomposition with mandatory scope tokens, then KEEP/DROP-only NLI selection and deterministic source speaker assignment"
	}
	reasoning := "none"
	if strings.Contains(cfg.Model, "gpt-oss") {
		reasoning = "provider does not accept reasoning_effort; omitted"
	}
	material := prompt + "\n---CONFIG---\n" + fmt.Sprintf("candidate=%s; chunk_chars=%d; carryover_turns=%d; max_facts=%d; validation=%s; model=%s; temperature=0; reasoning_effort=%s", cfg.Candidate, cfg.ChunkChars, cfg.Carryover, cfg.MaxFacts, policy, cfg.Model, reasoning)
	ph := hash(material)
	rh := hash(strings.Join(raw, "\n---CHUNK---\n"))
	tx, e := pool.Begin(ctx)
	if e != nil {
		return Result{}, e
	}
	defer tx.Rollback(ctx)
	var run string
	e = tx.QueryRow(ctx, `INSERT INTO study_two_runs(study_version,candidate,replicate,document_id,model,effective_prompt,prompt_hash,config,raw_output,raw_output_hash,prompt_tokens,completion_tokens,latency_ms) VALUES($1,$2,$3,$4,$5,$6,$7,jsonb_build_object('chunk_chars',$8::int,'carryover_turns',$9::int,'max_facts',$10::int,'validation_policy',$11::text,'temperature',0,'reasoning_effort','none'),$12,$13,$14,$15,$16) ON CONFLICT(study_version,candidate,replicate,document_id) DO UPDATE SET model=EXCLUDED.model,effective_prompt=EXCLUDED.effective_prompt,prompt_hash=EXCLUDED.prompt_hash,config=EXCLUDED.config,raw_output=EXCLUDED.raw_output,raw_output_hash=EXCLUDED.raw_output_hash,prompt_tokens=EXCLUDED.prompt_tokens,completion_tokens=EXCLUDED.completion_tokens,latency_ms=EXCLUDED.latency_ms,created_at=now() RETURNING id`, cfg.StudyVersion, cfg.Candidate, cfg.Replicate, doc, cfg.Model, material, ph, cfg.ChunkChars, cfg.Carryover, cfg.MaxFacts, policy, strings.Join(raw, "\n---CHUNK---\n"), rh, pt, ct, latency).Scan(&run)
	if e != nil {
		return Result{}, e
	}
	if _, e = tx.Exec(ctx, `DELETE FROM study_two_facts WHERE run_id=$1`, run); e != nil {
		return Result{}, e
	}
	for i, f := range validated {
		var fid string
		e = tx.QueryRow(ctx, `INSERT INTO study_two_facts(run_id,ordinal,text,content_hash,turn_start,turn_end,speaker_label) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id`, run, i+1, f.text, hash(f.text), f.start, f.end, f.speaker).Scan(&fid)
		if e != nil {
			return Result{}, e
		}
		for j, x := range f.evidence {
			if _, e = tx.Exec(ctx, `INSERT INTO study_two_evidence(fact_id,ordinal,quote,quote_char_start,quote_char_end,turn_start,turn_end,speaker_label) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`, fid, j+1, x.quote, x.start, x.end, x.first, x.last, x.speaker); e != nil {
				return Result{}, e
			}
		}
	}
	if e = tx.Commit(ctx); e != nil {
		return Result{}, e
	}
	return Result{run, len(all), len(validated), pt, ct, latency, ph}, nil
}

// runFableC4 is a structural reproduction of the Fable study's C4-SBAR
// architecture on Study #2's already-fixed documents. It does not consult a
// claim ledger at any stage; the ledger is evaluation-only.
func runFableC4(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, cfg Config) (Result, error) {
	if cfg.MaxFacts <= 0 {
		cfg.MaxFacts = 60
	}
	var doc string
	if err := pool.QueryRow(ctx, `SELECT id FROM documents WHERE corpus='mediasum' AND corpus_split='validation' AND external_id=$1`, cfg.ExternalID).Scan(&doc); err != nil {
		return Result{}, fmt.Errorf("validation document %q: %w", cfg.ExternalID, err)
	}
	transcript, turns, err := source(ctx, pool, doc)
	if err != nil {
		return Result{}, err
	}
	started := time.Now()
	var raw []string
	var pt, ct int
	generation, err := studyChat(ctx, client, cfg.Model, FableC4GeneratorPrompt, "<SOURCE_TURNS>\n"+label(turns)+"\n</SOURCE_TURNS>\nReturn at most "+fmt.Sprint(cfg.MaxFacts)+" cards.", 7000)
	if err != nil {
		return Result{}, err
	}
	pt += generation.PromptTokens
	ct += generation.CompletionTokens
	raw = append(raw, "---S1-GENERATE---\n"+generation.Text)
	generated, err := decodeC4(generation.Text)
	if err != nil {
		return Result{}, fmt.Errorf("C4 generator JSON: %w", err)
	}
	stage1 := validateC4(generated.Facts, transcript, turns)
	stage1JSON, _ := json.Marshal(c4AuditCards(stage1))
	raw = append(raw, "---S2-MECHANICAL-VALIDATE---\n"+string(stage1JSON))
	firstAudit, err := studyChat(ctx, client, cfg.Model, FableC4AuditPrompt, "<CARDS>\n"+string(stage1JSON)+"\n</CARDS>\nThe cards include all canonical speaker labels needed for this audit. Return one audit per card.", 5000)
	if err != nil {
		return Result{}, err
	}
	pt += firstAudit.PromptTokens
	ct += firstAudit.CompletionTokens
	raw = append(raw, "---S3-SPAN-BLIND-AUDIT---\n"+firstAudit.Text)
	audits, err := decodeC4Audit(firstAudit.Text)
	if err != nil {
		return Result{}, fmt.Errorf("C4 span-blind audit JSON: %w", err)
	}
	passed, failed := partitionC4(stage1, audits)
	final := append([]c4Card(nil), passed...)
	if len(failed) > 0 {
		failedJSON, _ := json.Marshal(c4RepairCards(failed, audits))
		repair, repairErr := studyChat(ctx, client, cfg.Model, FableC4RepairPrompt, "<SOURCE_TURNS>\n"+label(turns)+"\n</SOURCE_TURNS>\n<FAILED_CARDS>\n"+string(failedJSON)+"\n</FAILED_CARDS>\nReturn one repair per failed index.", 6000)
		if repairErr != nil {
			return Result{}, repairErr
		}
		pt += repair.PromptTokens
		ct += repair.CompletionTokens
		raw = append(raw, "---S4-TARGETED-REPAIR---\n"+repair.Text)
		repairs, decodeErr := decodeC4Repair(repair.Text)
		if decodeErr != nil {
			return Result{}, fmt.Errorf("C4 repair JSON: %w", decodeErr)
		}
		repaired := validateC4Repairs(failed, repairs, transcript, turns)
		repairedJSON, _ := json.Marshal(c4AuditCards(repaired))
		raw = append(raw, "---S4B-MECHANICAL-REVALIDATE---\n"+string(repairedJSON))
		if len(repaired) > 0 {
			secondAudit, auditErr := studyChat(ctx, client, cfg.Model, FableC4AuditPrompt, "<CARDS>\n"+string(repairedJSON)+"\n</CARDS>\nThe cards include all canonical speaker labels needed for this audit. Return one audit per card.", 5000)
			if auditErr != nil {
				return Result{}, auditErr
			}
			pt += secondAudit.PromptTokens
			ct += secondAudit.CompletionTokens
			raw = append(raw, "---S5-FRESH-SPAN-BLIND-REAUDIT---\n"+secondAudit.Text)
			second, decodeErr := decodeC4Audit(secondAudit.Text)
			if decodeErr != nil {
				return Result{}, fmt.Errorf("C4 span-blind re-audit JSON: %w", decodeErr)
			}
			kept, _ := partitionC4(repaired, second)
			final = append(final, kept...)
		}
	}
	// S6 is mechanical assembly: audit-pass originals plus audit-pass repaired
	// cards, deduplicated by text and ordered by original candidate index.
	sort.Slice(final, func(i, j int) bool { return final[i].Index < final[j].Index })
	finalFacts := dedupeC4(final)
	finalJSON, _ := json.Marshal(c4AuditCards(final))
	raw = append(raw, "---S6-ASSEMBLE---\n"+string(finalJSON))
	policy := "C4-SBAR: primary quote exact in cited turn (<=240 chars); optional antecedent exact in primary turn or four preceding turns (<=240 chars); span-blind support audit; repair limited to REANCHOR, RECUT, NAME_FORM, or DROP; repaired cards exact-revalidated and freshly re-audited; assembly keeps PASS only"
	material := FableC4GeneratorPrompt + "\n---S3-AUDIT---\n" + FableC4AuditPrompt + "\n---S4-REPAIR---\n" + FableC4RepairPrompt + "\n---CONFIG---\n" + fmt.Sprintf("candidate=%s; max_facts=%d; full_transcript=true; antecedent_window=4; model=%s; temperature=0; reasoning_effort=none; validation=%s", cfg.Candidate, cfg.MaxFacts, cfg.Model, policy)
	return persistStudy2Run(ctx, pool, cfg, doc, finalFacts, len(generated.Facts), material, policy, strings.Join(raw, "\n---STAGE---\n"), pt, ct, int(time.Since(started).Milliseconds()))
}

func decodeC4(s string) (c4Response, error) {
	if i, j := strings.Index(s, "{"), strings.LastIndex(s, "}"); i >= 0 && j > i {
		s = s[i : j+1]
	}
	var r c4Response
	return r, json.Unmarshal([]byte(s), &r)
}
func decodeC4Audit(s string) ([]c4Audit, error) {
	if i, j := strings.Index(s, "{"), strings.LastIndex(s, "}"); i >= 0 && j > i {
		s = s[i : j+1]
	}
	var r c4AuditResponse
	if err := json.Unmarshal([]byte(s), &r); err != nil {
		return nil, err
	}
	return r.Audits, nil
}
func decodeC4Repair(s string) ([]c4Repair, error) {
	if i, j := strings.Index(s, "{"), strings.LastIndex(s, "}"); i >= 0 && j > i {
		s = s[i : j+1]
	}
	var r c4RepairResponse
	if err := json.Unmarshal([]byte(s), &r); err != nil {
		return nil, err
	}
	return r.Repairs, nil
}

func validateC4(in []c4Candidate, raw string, ts []turn) []c4Card {
	seen := map[string]bool{}
	var out []c4Card
	for i, c := range in {
		if card, ok := validateC4One(i, c, raw, ts); ok && !seen[hash(card.Fact.text)] {
			seen[hash(card.Fact.text)] = true
			out = append(out, card)
		}
	}
	return out
}
func validateC4One(index int, c c4Candidate, raw string, ts []turn) (c4Card, bool) {
	c.Text, c.PrimaryQuote, c.AntecedentQuote = strings.TrimSpace(c.Text), strings.TrimSpace(c.PrimaryQuote), strings.TrimSpace(c.AntecedentQuote)
	if c.Text == "" || c.PrimaryQuote == "" || len(c.PrimaryQuote) > 240 {
		return c4Card{}, false
	}
	var primary *turn
	for i := range ts {
		if ts[i].ordinal == c.PrimaryTurn {
			primary = &ts[i]
			break
		}
	}
	if primary == nil || !strings.Contains(primary.text, c.PrimaryQuote) {
		return c4Card{}, false
	}
	pa := primary.start + strings.Index(primary.text, c.PrimaryQuote)
	f := fact{text: c.Text, start: primary.ordinal, end: primary.ordinal, speaker: primary.speaker, evidence: []evidence{{quote: c.PrimaryQuote, start: pa, end: pa + len(c.PrimaryQuote), first: primary.ordinal, last: primary.ordinal, speaker: primary.speaker}}}
	if c.AntecedentQuote != "" && c.AntecedentQuote != "-" {
		if len(c.AntecedentQuote) > 240 {
			return c4Card{}, false
		}
		var antecedent *turn
		for i := range ts {
			if ts[i].ordinal < primary.ordinal-4 || ts[i].ordinal > primary.ordinal || !strings.Contains(ts[i].text, c.AntecedentQuote) {
				continue
			}
			antecedent = &ts[i]
			break
		}
		if antecedent == nil {
			return c4Card{}, false
		}
		aa := antecedent.start + strings.Index(antecedent.text, c.AntecedentQuote)
		f.evidence = append(f.evidence, evidence{quote: c.AntecedentQuote, start: aa, end: aa + len(c.AntecedentQuote), first: antecedent.ordinal, last: antecedent.ordinal, speaker: antecedent.speaker})
		if antecedent.ordinal < f.start {
			f.start = antecedent.ordinal
		}
	}
	return c4Card{Index: index, Candidate: c, Fact: f}, true
}

func c4AuditCards(cards []c4Card) []map[string]any {
	out := make([]map[string]any, 0, len(cards))
	for _, c := range cards {
		antecedent, antecedentTurn, antecedentSpeaker := "-", -1, ""
		if len(c.Fact.evidence) > 1 {
			antecedent, antecedentTurn, antecedentSpeaker = c.Fact.evidence[1].quote, c.Fact.evidence[1].first, c.Fact.evidence[1].speaker
		}
		out = append(out, map[string]any{"index": c.Index, "text": c.Fact.text, "primary_quote": c.Fact.evidence[0].quote, "antecedent_quote": antecedent, "primary_turn": c.Candidate.PrimaryTurn, "antecedent_turn": antecedentTurn, "primary_speaker": c.Fact.speaker, "antecedent_speaker": antecedentSpeaker})
	}
	return out
}
func partitionC4(cards []c4Card, audits []c4Audit) ([]c4Card, []c4Card) {
	byIndex := map[int]c4Audit{}
	for _, a := range audits {
		byIndex[a.Index] = a
	}
	var pass, fail []c4Card
	for _, c := range cards {
		if a, ok := byIndex[c.Index]; ok && strings.EqualFold(strings.TrimSpace(a.Status), "PASS") {
			pass = append(pass, c)
		} else {
			fail = append(fail, c)
		}
	}
	return pass, fail
}
func c4RepairCards(cards []c4Card, audits []c4Audit) []map[string]any {
	byIndex := map[int]c4Audit{}
	for _, a := range audits {
		byIndex[a.Index] = a
	}
	out := c4AuditCards(cards)
	for _, row := range out {
		row["missing_elements"] = byIndex[row["index"].(int)].MissingElements
	}
	return out
}
func validateC4Repairs(failed []c4Card, repairs []c4Repair, raw string, ts []turn) []c4Card {
	byIndex := map[int]c4Repair{}
	for _, r := range repairs {
		byIndex[r.Index] = r
	}
	var out []c4Card
	for _, old := range failed {
		r, ok := byIndex[old.Index]
		if !ok || strings.EqualFold(strings.TrimSpace(r.Action), "DROP") {
			continue
		}
		action := strings.ToUpper(strings.TrimSpace(r.Action))
		if action != "REANCHOR" && action != "RECUT" && action != "NAME_FORM" {
			continue
		}
		c := c4Candidate{Text: r.Text, PrimaryQuote: r.PrimaryQuote, AntecedentQuote: r.AntecedentQuote, PrimaryTurn: old.Candidate.PrimaryTurn}
		switch action {
		case "REANCHOR":
			if c.Text != old.Candidate.Text || c.PrimaryQuote != old.Candidate.PrimaryQuote {
				continue
			}
		case "RECUT":
			if c.Text != old.Candidate.Text || c.AntecedentQuote != old.Candidate.AntecedentQuote {
				continue
			}
		case "NAME_FORM":
			if c.PrimaryQuote != old.Candidate.PrimaryQuote || c.AntecedentQuote != old.Candidate.AntecedentQuote {
				continue
			}
		}
		if card, valid := validateC4One(old.Index, c, raw, ts); valid {
			out = append(out, card)
		}
	}
	return out
}
func dedupeC4(cards []c4Card) []fact {
	seen := map[string]bool{}
	var out []fact
	for _, c := range cards {
		key := hash(c.Fact.text)
		if !seen[key] {
			seen[key] = true
			out = append(out, c.Fact)
		}
	}
	return out
}

func persistStudy2Run(ctx context.Context, pool *pgxpool.Pool, cfg Config, doc string, facts []fact, rawCount int, material, policy, raw string, pt, ct, latency int) (Result, error) {
	if strings.Contains(cfg.Model, "gpt-oss") {
		policy += "; provider does not accept reasoning_effort; omitted"
	}
	ph, rh := hash(material), hash(raw)
	tx, err := pool.Begin(ctx)
	if err != nil {
		return Result{}, err
	}
	defer tx.Rollback(ctx)
	var run string
	err = tx.QueryRow(ctx, `INSERT INTO study_two_runs(study_version,candidate,replicate,document_id,model,effective_prompt,prompt_hash,config,raw_output,raw_output_hash,prompt_tokens,completion_tokens,latency_ms) VALUES($1,$2,$3,$4,$5,$6,$7,jsonb_build_object('max_facts',$8::int,'validation_policy',$9::text,'temperature',0,'reasoning_effort','none'),$10,$11,$12,$13,$14) ON CONFLICT(study_version,candidate,replicate,document_id) DO UPDATE SET model=EXCLUDED.model,effective_prompt=EXCLUDED.effective_prompt,prompt_hash=EXCLUDED.prompt_hash,config=EXCLUDED.config,raw_output=EXCLUDED.raw_output,raw_output_hash=EXCLUDED.raw_output_hash,prompt_tokens=EXCLUDED.prompt_tokens,completion_tokens=EXCLUDED.completion_tokens,latency_ms=EXCLUDED.latency_ms,created_at=now() RETURNING id`, cfg.StudyVersion, cfg.Candidate, cfg.Replicate, doc, cfg.Model, material, ph, cfg.MaxFacts, policy, raw, rh, pt, ct, latency).Scan(&run)
	if err != nil {
		return Result{}, err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM study_two_facts WHERE run_id=$1`, run); err != nil {
		return Result{}, err
	}
	for i, f := range facts {
		var fid string
		if err = tx.QueryRow(ctx, `INSERT INTO study_two_facts(run_id,ordinal,text,content_hash,turn_start,turn_end,speaker_label) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id`, run, i+1, f.text, hash(f.text), f.start, f.end, f.speaker).Scan(&fid); err != nil {
			return Result{}, err
		}
		for j, x := range f.evidence {
			if _, err = tx.Exec(ctx, `INSERT INTO study_two_evidence(fact_id,ordinal,quote,quote_char_start,quote_char_end,turn_start,turn_end,speaker_label) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`, fid, j+1, x.quote, x.start, x.end, x.first, x.last, x.speaker); err != nil {
				return Result{}, err
			}
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return Result{}, err
	}
	return Result{run, rawCount, len(facts), pt, ct, latency, ph}, nil
}

// runDialogueActPairs makes the question--answer relationship an explicit
// stored evidence bundle. Pairing is deterministic and happens before any LLM
// call, so coverage cannot depend on a model guessing conversational scope.
func runDialogueActPairs(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, cfg Config) (Result, error) {
	if cfg.MaxFacts <= 0 {
		cfg.MaxFacts = 60
	}
	var doc string
	if err := pool.QueryRow(ctx, `SELECT id FROM documents WHERE corpus='mediasum' AND corpus_split='validation' AND external_id=$1`, cfg.ExternalID).Scan(&doc); err != nil {
		return Result{}, fmt.Errorf("validation document %q: %w", cfg.ExternalID, err)
	}
	transcript, turns, err := source(ctx, pool, doc)
	if err != nil {
		return Result{}, err
	}
	units := dialogueUnits(turns)
	payload, _ := json.Marshal(dialogueUnitPayload(units))
	started := time.Now()
	r, err := studyChat(ctx, client, cfg.Model, DialogueActPrompt, "<SOURCE_UNITS>\n"+string(payload)+"\n</SOURCE_UNITS>\nReturn at most "+fmt.Sprint(cfg.MaxFacts)+" cards.", 7000)
	if err != nil {
		return Result{}, err
	}
	decoded, err := decodeDialogue(r.Text)
	if err != nil {
		return Result{}, fmt.Errorf("dialogue-act JSON: %w", err)
	}
	facts := validateDialogue(decoded.Facts, transcript, units)
	policy := "deterministic immediate Q/A pairing: a question turn pairs only with its immediately following different-speaker declarative turn; paired answers store exact answer quote and optional exact question quote; unpaired declarative turns are standalone; answer speaker is attached deterministically; all stored quotes exact contiguous source substrings"
	material := DialogueActPrompt + "\n---CONFIG---\n" + fmt.Sprintf("candidate=%s; max_facts=%d; model=%s; temperature=0; reasoning_effort=none; validation=%s", cfg.Candidate, cfg.MaxFacts, cfg.Model, policy)
	raw := "---SOURCE-UNIT-METADATA---\n" + string(payload) + "\n---DIALOGUE-ACT-GENERATION---\n" + r.Text
	return persistStudy2Run(ctx, pool, cfg, doc, facts, len(decoded.Facts), material, policy, raw, r.PromptTokens, r.CompletionTokens, int(time.Since(started).Milliseconds()))
}

func dialogueUnits(ts []turn) []dialogueUnit {
	pairedAnswer := map[int]bool{}
	var out []dialogueUnit
	for i, t := range ts {
		if isQuestion(t.text) && i+1 < len(ts) && declarative(ts[i+1].text) && t.speaker != ts[i+1].speaker {
			q := t
			out = append(out, dialogueUnit{ID: len(out), Question: &q, Answer: ts[i+1], Paired: true})
			pairedAnswer[i+1] = true
			continue
		}
		if declarative(t.text) && !pairedAnswer[i] {
			out = append(out, dialogueUnit{ID: len(out), Answer: t})
		}
	}
	return out
}
func isQuestion(text string) bool { return strings.Contains(strings.TrimSpace(text), "?") }
func dialogueUnitPayload(units []dialogueUnit) []map[string]any {
	out := make([]map[string]any, 0, len(units))
	for _, u := range units {
		kind := "STANDALONE"
		row := map[string]any{"unit_id": u.ID, "kind": kind, "answer_turn": u.Answer.ordinal, "answer_speaker": u.Answer.speaker, "answer_text": u.Answer.text}
		if u.Paired && u.Question != nil {
			row["kind"] = "QUESTION_ANSWER"
			row["question_turn"] = u.Question.ordinal
			row["question_speaker"] = u.Question.speaker
			row["question_text"] = u.Question.text
		}
		out = append(out, row)
	}
	return out
}
func decodeDialogue(s string) (dialogueResponse, error) {
	if i, j := strings.Index(s, "{"), strings.LastIndex(s, "}"); i >= 0 && j > i {
		s = s[i : j+1]
	}
	var r dialogueResponse
	return r, json.Unmarshal([]byte(s), &r)
}
func validateDialogue(in []dialogueCandidate, raw string, units []dialogueUnit) []fact {
	byID := map[int]dialogueUnit{}
	for _, u := range units {
		byID[u.ID] = u
	}
	seen := map[string]bool{}
	var out []fact
	for _, c := range in {
		u, ok := byID[c.UnitID]
		c.Text, c.AnswerQuote, c.QuestionQuote = strings.TrimSpace(c.Text), strings.TrimSpace(c.AnswerQuote), strings.TrimSpace(c.QuestionQuote)
		if !ok || c.Text == "" || c.AnswerQuote == "" || len(c.AnswerQuote) > 240 || !strings.Contains(u.Answer.text, c.AnswerQuote) || seen[hash(c.Text)] {
			continue
		}
		if !u.Paired && c.QuestionQuote != "" && c.QuestionQuote != "-" {
			continue
		}
		answerAt := u.Answer.start + strings.Index(u.Answer.text, c.AnswerQuote)
		f := fact{text: c.Text, start: u.Answer.ordinal, end: u.Answer.ordinal, speaker: u.Answer.speaker, evidence: []evidence{{quote: c.AnswerQuote, start: answerAt, end: answerAt + len(c.AnswerQuote), first: u.Answer.ordinal, last: u.Answer.ordinal, speaker: u.Answer.speaker}}}
		if c.QuestionQuote != "" && c.QuestionQuote != "-" {
			if !u.Paired || u.Question == nil || len(c.QuestionQuote) > 240 || !strings.Contains(u.Question.text, c.QuestionQuote) {
				continue
			}
			questionAt := u.Question.start + strings.Index(u.Question.text, c.QuestionQuote)
			f.evidence = append(f.evidence, evidence{quote: c.QuestionQuote, start: questionAt, end: questionAt + len(c.QuestionQuote), first: u.Question.ordinal, last: u.Question.ordinal, speaker: u.Question.speaker})
			f.start = u.Question.ordinal
		}
		seen[hash(c.Text)] = true
		out = append(out, f)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].start < out[j].start })
	return out
}

func runDialogueActPredicateUnits(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, cfg Config) (Result, error) {
	if cfg.MaxFacts <= 0 {
		cfg.MaxFacts = 80
	}
	var doc string
	if err := pool.QueryRow(ctx, `SELECT id FROM documents WHERE corpus='mediasum' AND corpus_split='validation' AND external_id=$1`, cfg.ExternalID).Scan(&doc); err != nil {
		return Result{}, fmt.Errorf("validation document %q: %w", cfg.ExternalID, err)
	}
	transcript, turns, err := source(ctx, pool, doc)
	if err != nil {
		return Result{}, err
	}
	units := dialoguePredicateUnits(dialogueUnits(turns))
	payload, _ := json.Marshal(dialoguePredicatePayload(units))
	started := time.Now()
	r, err := studyChat(ctx, client, cfg.Model, DialoguePredicatePrompt, "<SOURCE_UNITS>\n"+string(payload)+"\n</SOURCE_UNITS>\nReturn at most "+fmt.Sprint(cfg.MaxFacts)+" cards.", 8000)
	if err != nil {
		return Result{}, err
	}
	decoded, err := decodeDialogue(r.Text)
	if err != nil {
		return Result{}, fmt.Errorf("dialogue predicate JSON: %w", err)
	}
	facts := validateDialoguePredicates(decoded.Facts, transcript, units)
	policy := "deterministic immediate Q/A pairing followed by deterministic lexical answer predicate/coordination units; one model card per predicate unit; answer quote must be exact within its own predicate unit; optional question quote only from that unit's paired question; answer speaker assigned deterministically"
	material := DialoguePredicatePrompt + "\n---CONFIG---\n" + fmt.Sprintf("candidate=%s; max_facts=%d; model=%s; temperature=0; reasoning_effort=none; validation=%s", cfg.Candidate, cfg.MaxFacts, cfg.Model, policy)
	raw := "---PREDICATE-SOURCE-UNIT-METADATA---\n" + string(payload) + "\n---DIALOGUE-PREDICATE-GENERATION---\n" + r.Text
	return persistStudy2Run(ctx, pool, cfg, doc, facts, len(decoded.Facts), material, policy, raw, r.PromptTokens, r.CompletionTokens, int(time.Since(started).Milliseconds()))
}
func dialoguePredicateUnits(dialogues []dialogueUnit) []dialoguePredicateUnit {
	var out []dialoguePredicateUnit
	for _, d := range dialogues {
		for _, claim := range predicateUnits(d.Answer.text) {
			claim = strings.TrimSpace(claim)
			if claim == "" {
				continue
			}
			out = append(out, dialoguePredicateUnit{ID: len(out), Dialogue: d, ClaimUnit: claim})
		}
	}
	return out
}
func dialoguePredicatePayload(units []dialoguePredicateUnit) []map[string]any {
	out := make([]map[string]any, 0, len(units))
	for _, u := range units {
		kind := "STANDALONE"
		row := map[string]any{"unit_id": u.ID, "kind": kind, "answer_turn": u.Dialogue.Answer.ordinal, "answer_speaker": u.Dialogue.Answer.speaker, "answer_claim_unit": u.ClaimUnit}
		if u.Dialogue.Paired && u.Dialogue.Question != nil {
			row["kind"] = "QUESTION_ANSWER"
			row["question_turn"] = u.Dialogue.Question.ordinal
			row["question_speaker"] = u.Dialogue.Question.speaker
			row["question_text"] = u.Dialogue.Question.text
		}
		out = append(out, row)
	}
	return out
}
func validateDialoguePredicates(in []dialogueCandidate, raw string, units []dialoguePredicateUnit) []fact {
	byID := map[int]dialoguePredicateUnit{}
	for _, u := range units {
		byID[u.ID] = u
	}
	seen := map[string]bool{}
	var out []fact
	for _, c := range in {
		u, ok := byID[c.UnitID]
		c.Text, c.AnswerQuote, c.QuestionQuote = strings.TrimSpace(c.Text), strings.TrimSpace(c.AnswerQuote), strings.TrimSpace(c.QuestionQuote)
		if !ok || c.Text == "" || c.AnswerQuote == "" || len(c.AnswerQuote) > 240 || !strings.Contains(u.ClaimUnit, c.AnswerQuote) || seen[hash(c.Text)] {
			continue
		}
		d := u.Dialogue
		if !d.Paired && c.QuestionQuote != "" && c.QuestionQuote != "-" {
			continue
		}
		answerAt := d.Answer.start + strings.Index(d.Answer.text, c.AnswerQuote)
		f := fact{text: c.Text, start: d.Answer.ordinal, end: d.Answer.ordinal, speaker: d.Answer.speaker, evidence: []evidence{{quote: c.AnswerQuote, start: answerAt, end: answerAt + len(c.AnswerQuote), first: d.Answer.ordinal, last: d.Answer.ordinal, speaker: d.Answer.speaker}}}
		if c.QuestionQuote != "" && c.QuestionQuote != "-" {
			if !d.Paired || d.Question == nil || len(c.QuestionQuote) > 240 || !strings.Contains(d.Question.text, c.QuestionQuote) {
				continue
			}
			questionAt := d.Question.start + strings.Index(d.Question.text, c.QuestionQuote)
			f.evidence = append(f.evidence, evidence{quote: c.QuestionQuote, start: questionAt, end: questionAt + len(c.QuestionQuote), first: d.Question.ordinal, last: d.Question.ordinal, speaker: d.Question.speaker})
			f.start = d.Question.ordinal
		}
		seen[hash(c.Text)] = true
		out = append(out, f)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].start < out[j].start })
	return out
}

func runQuoteCapsules(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, cfg Config) (Result, error) {
	if cfg.MaxFacts <= 0 {
		cfg.MaxFacts = 80
	}
	var doc string
	if err := pool.QueryRow(ctx, `SELECT id FROM documents WHERE corpus='mediasum' AND corpus_split='validation' AND external_id=$1`, cfg.ExternalID).Scan(&doc); err != nil {
		return Result{}, err
	}
	_, turns, err := source(ctx, pool, doc)
	if err != nil {
		return Result{}, err
	}
	units := dialogueUnits(turns)
	payload, _ := json.Marshal(dialogueUnitPayload(units))
	started := time.Now()
	r, err := studyChat(ctx, client, cfg.Model, QuoteCapsulePrompt, "<SOURCE_UNITS>\n"+string(payload)+"\n</SOURCE_UNITS>\nReturn at most "+fmt.Sprint(cfg.MaxFacts)+" cards.", 5000)
	if err != nil {
		return Result{}, err
	}
	decoded, err := decodeCapsules(r.Text)
	if err != nil {
		return Result{}, err
	}
	facts := validateCapsules(decoded.Cards, units)
	policy := "deterministic immediate Q/A pairing; model emits only unit IDs and UTF-8 byte offsets; exact boundary validation; deterministic attributed quote-capsule renderer; one answer span/card and optional paired question context span"
	material := QuoteCapsulePrompt + "\n---CONFIG---\n" + fmt.Sprintf("candidate=%s; max_facts=%d; model=%s; temperature=0; reasoning_effort=none; validation=%s", cfg.Candidate, cfg.MaxFacts, cfg.Model, policy)
	raw := "---CAPSULE-SOURCE-UNITS---\n" + string(payload) + "\n---CAPSULE-BOUNDARIES---\n" + r.Text
	return persistStudy2Run(ctx, pool, cfg, doc, facts, len(decoded.Cards), material, policy, raw, r.PromptTokens, r.CompletionTokens, int(time.Since(started).Milliseconds()))
}
func decodeCapsules(s string) (capsuleResponse, error) {
	if i, j := strings.Index(s, "{"), strings.LastIndex(s, "}"); i >= 0 && j > i {
		s = s[i : j+1]
	}
	var r capsuleResponse
	return r, json.Unmarshal([]byte(s), &r)
}
func validUTF8Boundary(s string, i int) bool {
	return i >= 0 && i <= len(s) && (i == 0 || i == len(s) || (s[i]&0xc0) != 0x80)
}
func validateCapsules(in []capsuleCard, units []dialogueUnit) []fact {
	byID := map[int]dialogueUnit{}
	for _, u := range units {
		byID[u.ID] = u
	}
	seen := map[string]bool{}
	var out []fact
	for _, c := range in {
		u, ok := byID[c.UnitID]
		if !ok || c.AnswerStart < 0 || c.AnswerEnd <= c.AnswerStart || c.AnswerEnd > len(u.Answer.text) || !validUTF8Boundary(u.Answer.text, c.AnswerStart) || !validUTF8Boundary(u.Answer.text, c.AnswerEnd) {
			continue
		}
		a := strings.TrimSpace(u.Answer.text[c.AnswerStart:c.AnswerEnd])
		low := strings.ToLower(strings.Trim(a, ". !?"))
		if len(a) < 8 || strings.Contains(a, "?") || low == "yes" || low == "no" || low == "okay" || low == "right" {
			continue
		}
		q := ""
		var qe *evidence
		// JSON decoders give omitted optional integer fields their zero value.
		// Treat the impossible empty [0,0) range as the documented no-question
		// sentinel; a real selected question span must be nonempty.
		if !(c.QuestionStart == 0 && c.QuestionEnd == 0) && (c.QuestionStart >= 0 || c.QuestionEnd >= 0) {
			if !u.Paired || u.Question == nil || c.QuestionStart < 0 || c.QuestionEnd <= c.QuestionStart || c.QuestionEnd > len(u.Question.text) || !validUTF8Boundary(u.Question.text, c.QuestionStart) || !validUTF8Boundary(u.Question.text, c.QuestionEnd) {
				continue
			}
			q = strings.TrimSpace(u.Question.text[c.QuestionStart:c.QuestionEnd])
			if q == "" {
				continue
			}
			at := u.Question.start + c.QuestionStart
			x := evidence{quote: q, start: at, end: at + len(q), first: u.Question.ordinal, last: u.Question.ordinal, speaker: u.Question.speaker}
			qe = &x
		}
		at := u.Answer.start + c.AnswerStart
		ae := evidence{quote: a, start: at, end: at + len(a), first: u.Answer.ordinal, last: u.Answer.ordinal, speaker: u.Answer.speaker}
		text := u.Answer.speaker + " said: \"" + a + "\""
		if u.Paired {
			text = u.Answer.speaker + " answered"
			if q != "" {
				text += " the question \"" + q + "\""
			}
			text += ": \"" + a + "\""
		}
		key := hash(text)
		if seen[key] {
			continue
		}
		seen[key] = true
		f := fact{text: text, start: u.Answer.ordinal, end: u.Answer.ordinal, speaker: u.Answer.speaker, evidence: []evidence{ae}}
		if qe != nil {
			f.evidence = append(f.evidence, *qe)
			f.start = qe.first
		}
		out = append(out, f)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].start < out[j].start })
	return out
}

func runQuoteCapsulesExact(ctx context.Context, pool *pgxpool.Pool, client *fireworks.Client, cfg Config) (Result, error) {
	if cfg.MaxFacts <= 0 {
		cfg.MaxFacts = 80
	}
	var doc string
	if err := pool.QueryRow(ctx, `SELECT id FROM documents WHERE corpus='mediasum' AND corpus_split='validation' AND external_id=$1`, cfg.ExternalID).Scan(&doc); err != nil {
		return Result{}, err
	}
	_, turns, err := source(ctx, pool, doc)
	if err != nil {
		return Result{}, err
	}
	units := dialogueUnits(turns)
	payload, _ := json.Marshal(dialogueUnitPayload(units))
	started := time.Now()
	r, err := studyChat(ctx, client, cfg.Model, QuoteCapsuleExactPrompt, "<SOURCE_UNITS>\n"+string(payload)+"\n</SOURCE_UNITS>\nReturn at most "+fmt.Sprint(cfg.MaxFacts)+" cards.", 5000)
	if err != nil {
		return Result{}, err
	}
	decoded, err := decodeCapsuleQuotes(r.Text)
	if err != nil {
		return Result{}, err
	}
	facts := validateCapsuleQuotes(decoded.Cards, units)
	policy := "deterministic immediate Q/A pairing; model emits only exact answer quote and optional paired-question quote strings; validator locates exact contiguous spans inside that unit; deterministic attributed quote-capsule renderer; no model-written fact text"
	material := QuoteCapsuleExactPrompt + "\n---CONFIG---\n" + fmt.Sprintf("candidate=%s; max_facts=%d; model=%s; temperature=0; reasoning_effort=none; validation=%s", cfg.Candidate, cfg.MaxFacts, cfg.Model, policy)
	raw := "---EXACT-CAPSULE-UNITS---\n" + string(payload) + "\n---EXACT-CAPSULE-SELECTION---\n" + r.Text
	return persistStudy2Run(ctx, pool, cfg, doc, facts, len(decoded.Cards), material, policy, raw, r.PromptTokens, r.CompletionTokens, int(time.Since(started).Milliseconds()))
}
func decodeCapsuleQuotes(s string) (capsuleQuoteResponse, error) {
	if i, j := strings.Index(s, "{"), strings.LastIndex(s, "}"); i >= 0 && j > i {
		s = s[i : j+1]
	}
	var r capsuleQuoteResponse
	return r, json.Unmarshal([]byte(s), &r)
}
func validateCapsuleQuotes(in []capsuleQuoteCard, units []dialogueUnit) []fact {
	by := map[int]dialogueUnit{}
	for _, u := range units {
		by[u.ID] = u
	}
	seen := map[string]bool{}
	var out []fact
	for _, c := range in {
		u, ok := by[c.UnitID]
		a := strings.TrimSpace(c.AnswerQuote)
		if !ok || a == "" || len(a) > 240 || !strings.Contains(u.Answer.text, a) {
			continue
		}
		low := strings.ToLower(strings.Trim(a, ". !?"))
		if len(a) < 8 || strings.Contains(a, "?") || low == "yes" || low == "no" || low == "okay" {
			continue
		}
		at := u.Answer.start + strings.Index(u.Answer.text, a)
		ae := evidence{quote: a, start: at, end: at + len(a), first: u.Answer.ordinal, last: u.Answer.ordinal, speaker: u.Answer.speaker}
		q := strings.TrimSpace(c.QuestionQuote)
		var qe *evidence
		if q != "" && q != "-" {
			if !u.Paired || u.Question == nil || len(q) > 240 || !strings.Contains(u.Question.text, q) {
				continue
			}
			qat := u.Question.start + strings.Index(u.Question.text, q)
			x := evidence{quote: q, start: qat, end: qat + len(q), first: u.Question.ordinal, last: u.Question.ordinal, speaker: u.Question.speaker}
			qe = &x
		}
		text := u.Answer.speaker + " said: \"" + a + "\""
		if u.Paired {
			text = u.Answer.speaker + " answered"
			if q != "" && q != "-" {
				text += " the question \"" + q + "\""
			}
			text += ": \"" + a + "\""
		}
		if seen[hash(text)] {
			continue
		}
		seen[hash(text)] = true
		f := fact{text: text, start: u.Answer.ordinal, end: u.Answer.ordinal, speaker: u.Answer.speaker, evidence: []evidence{ae}}
		if qe != nil {
			f.evidence = append(f.evidence, *qe)
			f.start = qe.first
		}
		out = append(out, f)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].start < out[j].start })
	return out
}

// turnWindows supplies the preceding two turns solely as local disambiguating
// context. The final turn is the only extraction target, so this differs from
// both transcript packing and overlapping chunk extraction.
func turnWindows(ts []turn) [][]turn {
	var windows [][]turn
	for i, t := range ts {
		if !declarative(t.text) {
			continue
		}
		start := i - 2
		if start < 0 {
			start = 0
		}
		window := append([]turn(nil), ts[start:i]...)
		window = append(window, t)
		windows = append(windows, window)
	}
	return windows
}

func declarative(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" || strings.HasSuffix(text, "?") {
		return false
	}
	// Interview transcripts include host introductions and acknowledgements;
	// retain mixed declarative turns but exclude trivial acknowledgement-only ones.
	lower := strings.ToLower(text)
	return lower != "yes." && lower != "yeah." && lower != "thank you." && lower != "thanks."
}

// eligibleTurns is a deterministic pre-generation guard. It deliberately
// removes utterances that are syntactically interrogative or clear broadcast
// framing/promotion, the observed source of guest-misattribution errors.
func eligibleTurns(ts []turn) []turn {
	var out []turn
	for _, t := range ts {
		text := strings.TrimSpace(t.text)
		lower := strings.ToLower(text)
		if !declarative(text) || strings.Contains(text, "?") {
			continue
		}
		if strings.HasPrefix(lower, "thank you") || strings.HasPrefix(lower, "thanks") || strings.HasPrefix(lower, "coming up") || strings.HasPrefix(lower, "you're listening") || strings.HasPrefix(lower, "for more") {
			continue
		}
		out = append(out, t)
	}
	return out
}
func pack(ts []turn, n, carry int) [][]turn {
	var out [][]turn
	var cur []turn
	size := 0
	for _, t := range ts {
		if len(cur) > 0 && size+len(t.text) > n {
			base := append([]turn{}, cur...)
			out = append(out, base)
			prev := cur
			if carry > len(prev) {
				carry = len(prev)
			}
			cur = append([]turn{}, prev[len(prev)-carry:]...)
			size = 0
			for _, x := range cur {
				size += len(x.text)
			}
		}
		cur = append(cur, t)
		size += len(t.text)
	}
	if len(cur) > 0 {
		out = append(out, cur)
	}
	return out
}
func source(ctx context.Context, p *pgxpool.Pool, id string) (string, []turn, error) {
	rows, e := p.Query(ctx, `SELECT ordinal,speaker_label,text,char_start,char_end FROM turns WHERE document_id=$1 ORDER BY ordinal`, id)
	if e != nil {
		return "", nil, e
	}
	defer rows.Close()
	var out []turn
	var raw []string
	for rows.Next() {
		var t turn
		if e = rows.Scan(&t.ordinal, &t.speaker, &t.text, &t.start, &t.end); e != nil {
			return "", nil, e
		}
		out = append(out, t)
		raw = append(raw, t.text)
	}
	return strings.Join(raw, "\n"), out, rows.Err()
}
func label(ts []turn) string {
	var p []string
	for _, t := range ts {
		p = append(p, fmt.Sprintf("[TURN %d | %s]\n%s", t.ordinal, t.speaker, t.text))
	}
	return strings.Join(p, "\n\n")
}
func decode(s string) (response, error) {
	if i, j := strings.Index(s, "{"), strings.LastIndex(s, "}"); i >= 0 && j > i {
		s = s[i : j+1]
	}
	var r response
	return r, json.Unmarshal([]byte(s), &r)
}
func decodeNLI(s string) ([]nliDecision, error) {
	if i, j := strings.Index(s, "{"), strings.LastIndex(s, "}"); i >= 0 && j > i {
		s = s[i : j+1]
	}
	var r nliResponse
	if err := json.Unmarshal([]byte(s), &r); err != nil {
		return nil, err
	}
	return r.Decisions, nil
}
func decodeFrames(s string) ([]frame, error) {
	if i, j := strings.Index(s, "{"), strings.LastIndex(s, "}"); i >= 0 && j > i {
		s = s[i : j+1]
	}
	var r frameResponse
	if err := json.Unmarshal([]byte(s), &r); err != nil {
		return nil, err
	}
	return r.Frames, nil
}
func renderFrames(frames []frame, units []string) []candidate {
	var out []candidate
	for _, f := range frames {
		if f.ClauseID < 0 || f.ClauseID >= len(units) || len(f.Evidence) == 0 {
			continue
		}
		evidence := []string{}
		valid := true
		for _, id := range f.Evidence {
			if id < 0 || id >= len(units) {
				valid = false
				break
			}
			evidence = append(evidence, units[id])
		}
		if !valid {
			continue
		}
		text := strings.TrimSpace(strings.Join([]string{f.Subject, f.Modality, f.Negation, f.Relation, f.Object, f.Temporal, f.Condition, f.Comparison}, " "))
		if text == "" {
			continue
		}
		proof := strings.Join(evidence, " ")
		for _, field := range []string{f.Subject, f.Relation, f.Object, f.Modality, f.Negation, f.Temporal, f.Condition, f.Comparison} {
			if strings.TrimSpace(field) != "" && !strings.Contains(proof, field) {
				valid = false
				break
			}
		}
		if valid {
			out = append(out, candidate{Text: text, Evidence: evidence})
		}
	}
	return out
}

func applyNLI(in []candidate, decisions []nliDecision) []candidate {
	byIndex := map[int]nliDecision{}
	for _, d := range decisions {
		if d.Index >= 0 && d.Index < len(in) {
			byIndex[d.Index] = d
		}
	}
	var out []candidate
	for i, c := range in {
		d, ok := byIndex[i]
		if !ok {
			continue
		}
		switch strings.ToUpper(strings.TrimSpace(d.Action)) {
		case "KEEP":
			out = append(out, c)
		case "REWRITE":
			if strings.TrimSpace(d.Text) != "" && len(d.Evidence) > 0 {
				out = append(out, candidate{Text: d.Text, Evidence: d.Evidence})
			}
		}
	}
	return out
}
func validate(cs []candidate, raw string, ts []turn, expandAntecedents, deterministicSpeaker bool) []fact {
	seen := map[string]bool{}
	var out []fact
	for _, c := range cs {
		c.Text = strings.TrimSpace(c.Text)
		if c.Text == "" || seen[hash(c.Text)] || len(c.Evidence) == 0 || len(c.Evidence) > 3 {
			continue
		}
		f := fact{text: c.Text}
		ok := true
		for _, q := range c.Evidence {
			q = strings.TrimSpace(q)
			at := strings.Index(raw, q)
			if q == "" || at < 0 || len(q) > 320 {
				ok = false
				break
			}
			a, b, sp := bounds(ts, at, at+len(q))
			f.evidence = append(f.evidence, evidence{q, at, at + len(q), a, b, sp})
		}
		if !ok {
			continue
		}
		// A distinct validation policy: persist the source turn that contains
		// the selected evidence plus its two predecessors.  This makes the
		// decontextualization proof bundle explicit and auditable instead of
		// silently relying on prompt-only carryover.  These are source spans,
		// not generated context.
		if expandAntecedents {
			ordinal := f.evidence[0].first
			for _, t := range ts {
				if t.ordinal < ordinal-2 || t.ordinal > ordinal {
					continue
				}
				already := false
				for _, x := range f.evidence {
					if x.start == t.start && x.end == t.end {
						already = true
					}
				}
				if !already {
					f.evidence = append(f.evidence, evidence{t.text, t.start, t.end, t.ordinal, t.ordinal, t.speaker})
				}
			}
		}
		for _, x := range f.evidence {
			if f.start == 0 || x.first < f.start {
				f.start = x.first
			}
			if x.last > f.end {
				f.end = x.last
			}
			if f.speaker == "" {
				f.speaker = x.speaker
			}
		}
		if deterministicSpeaker && f.speaker != "" {
			f.text = f.speaker + " said that " + strings.TrimSpace(f.text)
		}
		seen[hash(c.Text)] = true
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].start < out[j].start })
	return out
}

// splitClauses is intentionally lexical and deterministic. It creates the
// source-first coverage budget; an LLM may rewrite or reject a clause, but may
// not invent a candidate that this segmentation did not expose.
func splitClauses(text string) []string {
	var out []string
	start := 0
	for i, r := range text {
		if r != '.' && r != ';' && r != '!' {
			continue
		}
		part := strings.TrimSpace(text[start : i+len(string(r))])
		if part != "" && !strings.HasSuffix(part, "?") {
			out = append(out, part)
		}
		start = i + len(string(r))
	}
	if part := strings.TrimSpace(text[start:]); part != "" && !strings.HasSuffix(part, "?") {
		out = append(out, part)
	}
	return out
}

func propositionUnits(text string) []string {
	var out []string
	for _, clause := range splitClauses(text) {
		parts := strings.Split(clause, ";")
		for _, p := range parts {
			for _, q := range strings.Split(p, ", and ") {
				q = strings.TrimSpace(q)
				if q != "" {
					out = append(out, q)
				}
			}
		}
	}
	return out
}

func predicateUnits(text string) []string {
	var out []string
	for _, sentence := range splitClauses(text) {
		for _, item := range strings.Split(sentence, ";") {
			for _, part := range strings.Split(item, ", and ") {
				part = strings.TrimSpace(part)
				if part != "" {
					out = append(out, part)
				}
			}
		}
	}
	return out
}

func substantiveUnits(text string) []string {
	lower := strings.ToLower(strings.TrimSpace(text))
	if strings.Contains(text, "?") || strings.HasPrefix(lower, "thank") || strings.HasPrefix(lower, "yes") || strings.HasPrefix(lower, "yeah") || strings.HasPrefix(lower, "please") || strings.HasPrefix(lower, "coming up") || strings.HasPrefix(lower, "you're listening") {
		return nil
	}
	var out []string
	for _, u := range predicateUnits(text) {
		l := strings.ToLower(strings.TrimSpace(u))
		if l == "" || l == "okay." || l == "all right." || l == "right." {
			continue
		}
		out = append(out, u)
	}
	return out
}

func mandatoryTokens(text string) []string {
	lower := strings.ToLower(text)
	markers := []string{"not", "never", "no ", "might", "may", "could", "would", "should", "will", "said", "reported", "according to", "today", "yesterday", "currently", "former", "first", "more", "less", "%"}
	var out []string
	for _, marker := range markers {
		if strings.Contains(lower, marker) {
			out = append(out, marker)
		}
	}
	for _, token := range strings.Fields(text) {
		if strings.IndexFunc(token, func(r rune) bool { return r >= '0' && r <= '9' }) >= 0 {
			out = append(out, token)
		}
	}
	return out
}

func splitContextClauses(ts []turn) []string {
	var out []string
	for _, t := range ts {
		out = append(out, splitClauses(t.text)...)
	}
	return out
}

// proofFilter is deterministic: an accepted proof quote must be a complete
// lexical source clause, and the final representation cannot retain dialogue
// deixis. It intentionally rejects dubious output rather than repairing it.
func proofFilter(in []fact, ts []turn) []fact {
	allowed := map[string]bool{}
	for _, t := range ts {
		for _, clause := range splitClauses(t.text) {
			allowed[strings.TrimSpace(clause)] = true
		}
	}
	// "that" is deliberately excluded because deterministic attribution adds
	// the grammatical prefix "<speaker> said that ...".
	bad := map[string]bool{"we": true, "us": true, "our": true, "it": true, "this": true, "they": true, "them": true, "he": true, "she": true, "there": true, "you": true}
	var out []fact
	for _, f := range in {
		ok := true
		// The expansion spans are permitted, but at least one original exact
		// model-selected span must itself be a complete source clause.
		hasClause := false
		for _, e := range f.evidence {
			if allowed[strings.TrimSpace(e.quote)] {
				hasClause = true
				break
			}
		}
		if !hasClause {
			continue
		}
		for _, token := range strings.Fields(strings.ToLower(f.text)) {
			token = strings.Trim(token, ".,;:!?()[]\"'")
			if bad[token] {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, f)
		}
	}
	return out
}

// bindProofClauses deterministically expands a model-selected exact substring
// to its containing lexical source clause. This is provenance normalization,
// not semantic rewriting: the stored replacement remains an exact source span.
func bindProofClauses(in []fact, ts []turn) []fact {
	for i := range in {
		for _, t := range ts {
			for _, clause := range splitClauses(t.text) {
				at := strings.Index(t.text, clause)
				if at < 0 {
					continue
				}
				clauseStart := t.start + at
				clauseEnd := clauseStart + len(clause)
				for _, e := range in[i].evidence {
					if e.start < clauseStart || e.end > clauseEnd {
						continue
					}
					already := false
					for _, prior := range in[i].evidence {
						if prior.start == clauseStart && prior.end == clauseEnd {
							already = true
						}
					}
					if !already {
						in[i].evidence = append(in[i].evidence, evidence{clause, clauseStart, clauseEnd, t.ordinal, t.ordinal, t.speaker})
					}
				}
			}
		}
	}
	return in
}
func bounds(ts []turn, a, b int) (int, int, string) {
	first, last := ts[0], ts[len(ts)-1]
	for _, t := range ts {
		if a >= t.start && a <= t.end {
			first = t
		}
		if b >= t.start && b <= t.end+1 {
			last = t
			break
		}
	}
	return first.ordinal, last.ordinal, first.speaker
}
func hash(s string) string {
	x := sha256.Sum256([]byte(strings.TrimSpace(s)))
	return hex.EncodeToString(x[:])
}
