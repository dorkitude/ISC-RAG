#!/usr/bin/env bash
# ISC-RAG raw-data audit gate (host Postgres, no pgvector required).
#
# The audit only reads answers.usage and questions.external_id. The vector type
# appears solely on the embeddings table (migration 002), so a plain PostgreSQL
# 17 install is sufficient: embeddings simply fails to restore and is ignored.
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

DUMP="${1:-$HOME/Downloads/20260819-1620_isc_rag_full.dump}"
REPO="${REPO:-$HOME/firebase/ISC-RAG}"
OUT="${OUT:-$HOME/Desktop/isc_rag_audit}"
DB="${DB:-isc_rag_audit}"
EXPECT_SHA="8efe7f1d0a4693571593a48f35e01d1cce616e822a7119951506fc216cbff04b"

# Homebrew versioned formulae are keg-only; add their bin directories.
for p in /opt/homebrew/opt/postgresql@17/bin /usr/local/opt/postgresql@17/bin \
         /opt/homebrew/opt/libpq/bin /usr/local/opt/libpq/bin; do
  [ -d "$p" ] && PATH="$p:$PATH"
done
export PATH

if ! command -v pg_restore >/dev/null 2>&1 || ! command -v createdb >/dev/null 2>&1; then
  cat >&2 <<'MSG'
FAIL: PostgreSQL client tools not found.

Install them, then rerun this script:

    brew install postgresql@17
    brew services start postgresql@17

pgvector is NOT required: the audit reads only answers and questions, and the
vector type is used solely by the embeddings table.
MSG
  exit 1
fi
echo "using: $(command -v pg_restore)  ($(pg_restore --version))"

mkdir -p "$OUT"

echo "== 1. SHA-256 =="
ACTUAL="$(shasum -a 256 "$DUMP" | awk '{print $1}')"
echo "expected: $EXPECT_SHA"; echo "actual  : $ACTUAL"
[ "$ACTUAL" = "$EXPECT_SHA" ] || { echo "FAIL: checksum mismatch. Stop." >&2; exit 1; }
echo "OK"

echo "== 2. restore into a scratch database =="
dropdb --if-exists "$DB"
createdb "$DB"
# Errors are expected on the embeddings table and the vector extension.
pg_restore -d "$DB" --no-owner -j 8 "$DUMP" 2> "$OUT/00_restore_stderr.txt" || true

EXPECTED_ERR='vector|embeddings'
TOTAL_ERR=$(grep -c 'error' "$OUT/00_restore_stderr.txt" || true)
OTHER_ERR=$(grep 'error' "$OUT/00_restore_stderr.txt" | grep -Ev "$EXPECTED_ERR" | wc -l | tr -d ' ')
echo "restore errors: $TOTAL_ERR total, $OTHER_ERR unrelated to vector/embeddings"
if [ "${OTHER_ERR:-0}" -gt 0 ]; then
  echo "--- unexpected restore errors (first 20) ---"
  grep 'error' "$OUT/00_restore_stderr.txt" | grep -Ev "$EXPECTED_ERR" | head -20
  echo "-------------------------------------------"
fi

echo "== 3. the two tables the audit depends on =="
psql -d "$DB" -At -c "SELECT 'answers', count(*) FROM answers
                      UNION ALL SELECT 'questions', count(*) FROM questions
                      UNION ALL SELECT 'runs', count(DISTINCT experiment_run_id) FROM answers;" \
  | tee "$OUT/00_counts.txt"

echo "== 4. every run_id in the ledgers must exist in the dump =="
cd "$REPO"
cat experiments/*.tsv | awk -F'\t' 'NF==3 {print $3}' | sort -u > "$OUT/00_ledger_runids.txt"
psql -d "$DB" -At -c "SELECT DISTINCT experiment_run_id FROM answers;" | sort -u > "$OUT/00_db_runids.txt"
comm -23 "$OUT/00_ledger_runids.txt" "$OUT/00_db_runids.txt" > "$OUT/00_missing_runids.txt"
if [ -s "$OUT/00_missing_runids.txt" ]; then
  echo "WARNING: run IDs in the ledgers but absent from the dump:"; cat "$OUT/00_missing_runids.txt"
else
  echo "OK: all $(wc -l < "$OUT/00_ledger_runids.txt" | tr -d ' ') ledger run IDs present"
fi

export DATABASE_URL="postgres:///$DB"
if command -v uv >/dev/null 2>&1; then
  RUN() { uv run --with scipy --with psycopg2-binary scripts/stats-mcnemar.py "$@"; }
else
  echo "uv not found; using python3 (scipy and psycopg2 must be installed)"
  RUN() { python3 scripts/stats-mcnemar.py "$@"; }
fi
FK=judge_kimi_frontier_xl_v1;  FS=source_support_kimi_frontier_xl_v1
HK=judge_kimi_heldout500_v1;   HS=source_support_kimi_heldout500_v1

echo "== 5. frontier-xl (frozen 100) =="
{ RUN --runs-tsv experiments/frontier-xl-runs.tsv --correctness-key $FK --support-key $FS --metric correct
  echo
  RUN --runs-tsv experiments/frontier-xl-runs.tsv --correctness-key $FK --support-key $FS --metric corr_supp
} > "$OUT/frontier-xl-mcnemar_REGEN.md"

echo "== 6. held-out 500 (same three sections, same order as the committed file) =="
{ RUN --runs-tsv experiments/heldout-matched-runs.tsv --correctness-key $HK --support-key $HS --metric correct
  echo
  RUN --runs-tsv experiments/heldout-xl-runs.tsv      --correctness-key $HK --support-key $HS --metric correct
  echo
  RUN --runs-tsv experiments/heldout-cr-runs.tsv      --correctness-key $HK --support-key $HS --metric correct
} > "$OUT/heldout-500-mcnemar_REGEN.md"

echo "== 7. contextual-retrieval baseline on frozen 100 (extra) =="
RUN --runs-tsv experiments/cr-comparison-runs.tsv --correctness-key $FK --support-key $FS --metric correct \
  > "$OUT/cr-comparison-mcnemar_REGEN.md"

echo "== 8. compare against the committed tables =="
set +e
python3 "$HERE/audit_compare.py" \
  "$REPO/experiments/frontier-xl-mcnemar-20260726.md=$OUT/frontier-xl-mcnemar_REGEN.md" \
  "$REPO/experiments/heldout-500-mcnemar-20260726.md=$OUT/heldout-500-mcnemar_REGEN.md" \
  | tee "$OUT/AUDIT_REPORT.txt"
RC=${PIPESTATUS[0]}
set -e

echo
echo "Database '$DB' is left in place for follow-up queries: psql -d $DB"
echo "Send back: $OUT/AUDIT_REPORT.txt and $OUT/00_counts.txt"
exit "$RC"
