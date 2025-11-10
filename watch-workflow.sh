#!/usr/bin/env bash

if [ "$#" -lt 2 ] || [[ "$*" =~ '-h' ]]; then
  cat <<EOF
usage:
  $0 <workflow.yml> <unique id>

note:
  all GH_* env vars apply, such as GH_REPO
  https://cli.github.com/manual/gh_help_environment
EOF
  exit 1
fi

set -xeu

workflow="$1"
unique_id="$2"

list_resp="$(mktemp)"
trap "rm $list_resp" EXIT

echo 'waiting 30s for run to start'
for _ in {1..30}; do
  gh run list --workflow "$workflow" \
    --json databaseId,displayTitle,status,url \
    --jq ".[] | select(.displayTitle | contains(\"$unique_id\"))" \
    | tee "$list_resp" | grep . && break
  sleep 1
  printf .
done

echo 'watch the run'
gh run watch --compact --interval 3 --exit-status $(jq -r '.databaseId' < "$list_resp")
