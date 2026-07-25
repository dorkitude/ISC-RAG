# Study #4 materials — deterministic speaker metadata

This directory holds the reproducible, non-source-derived materials for the
role-metadata prompt study. `claim-ledger.md` was locked before extraction;
`prompts/` contains five independently seeded contracts and the two sequential
refinements; and `run.sh` makes a temperature-zero Fireworks request against the
fixed development document `NPR-1040`.

Run from the ISC-RAG project root after placing `FIREWORKS_API_KEY` and
`DATABASE_URL` in the local `.env`:

```sh
experiments/fact-prompt-study-4-role-metadata/run.sh seed-a-role-aware-direct
```

The runner reads only the deterministic `turns.speaker_name` and
`turns.speaker_role` fields added in migration 011. It does not infer roles from
dialogue and it does not write primary facts or study database tables. Transcript
copies, requests containing them, completions, and provider response metadata
are local-only in the ignored `artifacts/fact-prompt-study-4-role-metadata/`.
The audit and outcome are in
`docs/fact-prompt-study-4-role-metadata.md`.
