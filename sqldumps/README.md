# sqldumps

Disaster-recovery dumps of the ISC-RAG PostgreSQL database (the system of
record for the frozen-100 and heldout-500 protocol data).

Dumps are pg_dump custom format, named `YYYYMMDD-HHMM_isc_rag_full.dump`.
They are too large for regular git objects (~800MB), so `*.dump` is gitignored
here and each dump is attached to a GitHub release tagged `sqldump-YYYYMMDD`
on this repository instead.

Create:

    pg_dump "$DATABASE_URL" -Fc -f "sqldumps/$(date +%Y%m%d-%H%M)_isc_rag_full.dump"
    gh release create "sqldump-$(date +%Y%m%d)" sqldumps/*.dump --notes "Full pg_dump -Fc of isc_rag" --title "SQL dump $(date +%Y%m%d)"

Restore (into a fresh database):

    pg_restore -d "$DATABASE_URL" --no-owner -j 8 <file>.dump

Inspect contents without restoring:

    pg_restore -l <file>.dump
