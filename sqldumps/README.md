# Private database backups

Full ISC-RAG database dumps contain MediaSum transcripts, supporting
quotations, extracted facts, questions, generated answers, and raw model
artifacts. They are disaster-recovery material, not public release artifacts.

The `sqldumps/` directory ignores every dump. Store encrypted backups only in a
private location with access limited to the research team. Do not attach a full
dump to a public GitHub release.

Create a private backup:

```sh
pg_dump "$DATABASE_URL" -Fc -f "sqldumps/$(date +%Y%m%d-%H%M)_isc_rag_full.dump"
shasum -a 256 sqldumps/*.dump
```

For public verification, regenerate `results/public-results.sql`. That export
contains configuration and aggregate results but no source or model payloads;
see [`results/README.md`](../results/README.md).
