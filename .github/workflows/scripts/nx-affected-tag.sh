#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 2 ]]; then
  echo "Usage: nx-affected-tag.sh <tag> <base> [<head>]" >&2
  exit 2
fi

tag="$1"
base="$2"
head="${3:-HEAD}"

NX_CMD="${NX_CMD:-./node_modules/.bin/nx}"

changed_files="$(git diff --name-only "$base" "$head" | paste -sd, -)"

if [[ -z "$changed_files" ]]; then
  printf 'false\n'
  exit 0
fi

projects="$($NX_CMD show projects \
  --affected \
  --files "$changed_files" \
  --projects "tag:$tag" \
  --sep ' ')"

if [[ -n "$projects" && "$projects" != "[]" ]]; then
  printf '%s\n' "$projects" >&2
  printf 'true\n'
else
  printf 'false\n'
fi
