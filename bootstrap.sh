#!/usr/bin/env bash
# Bootstrap ISC-RAG's pinned source repository and MediaSum resource.
set -euo pipefail

project_dir="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
source_dir="$project_dir/sources/MediaSum"
data_dir="$project_dir/resources/mediasum"
revision="bcbcef7de3b4e702e526352d825a4ff06de2becb"
source_repo="https://github.com/zcgzcgzcg1/MediaSum.git"
accept_research_only=false
skip_data=false

usage() {
  cat <<'USAGE'
Usage: ./bootstrap.sh --accept-mediasum-research-only [--skip-data]

Initializes the pinned official MediaSum source repository and downloads the
pinned MediaSum train, validation, and test files into resources/mediasum/.

Options:
  --accept-mediasum-research-only  Confirm MediaSum will be used for research only.
  --skip-data                      Initialize only the upstream source submodule.
  -h, --help                       Show this help text.
USAGE
}

die() {
  printf 'bootstrap: %s\n' "$*" >&2
  exit 1
}

need() {
  command -v "$1" >/dev/null 2>&1 || die "required command not found: $1"
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --accept-mediasum-research-only) accept_research_only=true ;;
    --skip-data) skip_data=true ;;
    -h|--help) usage; exit 0 ;;
    *) die "unknown option: $1" ;;
  esac
  shift
done

need git

if git -C "$project_dir" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  (
    cd "$project_dir"
    git submodule update --init --recursive -- sources/MediaSum
  )
else
  if [ ! -d "$source_dir/.git" ]; then
    mkdir -p "$(dirname -- "$source_dir")"
    git clone "$source_repo" "$source_dir"
  fi
  git -C "$source_dir" checkout 7c92589bb6e69e678d043a9cf8c8c7285cd88055
fi

printf 'Upstream source ready: %s\n' "$source_dir"

if [ "$skip_data" = true ]; then
  exit 0
fi

[ "$accept_research_only" = true ] || die \
  'MediaSum is research-only. Re-run with --accept-mediasum-research-only to download it.'

need curl
need shasum
mkdir -p "$data_dir"

verify_file() {
  expected_hash="$1"
  target="$2"
  [ -f "$target" ] || return 1
  actual_hash="$(shasum -a 256 "$target" | awk '{print $1}')"
  [ "$actual_hash" = "$expected_hash" ]
}

fetch_file() {
  name="$1"
  expected_hash="$2"
  target="$data_dir/$name"
  url="https://huggingface.co/datasets/nbroad/mediasum/resolve/$revision/$name?download=true"

  if verify_file "$expected_hash" "$target"; then
    printf 'Verified existing %s\n' "$name"
    return
  fi

  temp="$data_dir/.${name}.download.$$"
  trap 'rm -f "$temp"' EXIT HUP INT TERM
  printf 'Downloading %s\n' "$name"
  curl --fail --location --retry 4 --retry-all-errors --output "$temp" "$url"
  verify_file "$expected_hash" "$temp" || die "checksum mismatch for $name"
  mv "$temp" "$target"
  trap - EXIT HUP INT TERM
  printf 'Verified %s\n' "$name"
}

fetch_file train.json cd713e926637b34930c1c5595c9e97cec6d514f4d95921a73e3da766b5ba75af
fetch_file validation.json 313dfd682b25a6151a760632c06ba2bbb6fbb585a7df5301b524bdb07cb499a0
fetch_file test.json 136c241cdb4fba03395df259d1153d60a61fe274215e78ee1fc96211feff6cde

printf 'MediaSum resource ready: %s\n' "$data_dir"
