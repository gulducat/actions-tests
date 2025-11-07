#!/usr/bin/env bash

set -euo pipefail

function run() {
  local workflow uid
  workflow="$1" ; shift
  uid="$1" ; shift
  echo "workflow: $workflow"
  echo "unique id: '$uid'"
  echo "run args: $@"

  gh workflow run "$workflow" "$@"
  
  echo 'wait for run to start'
  for _ in {1..20}; do
    gh run list --workflow "$workflow" \
      --json databaseId,displayTitle,status,url \
      --jq ".[] | select(.displayTitle | contains(\"$uid\"))" \
      | tee run.json | grep . && break
    sleep 1
    printf .
  done
  
  echo 'watch the run'
  gh run watch --compact --interval 3 --exit-status $(jq -r '.databaseId' < run.json)
}

run "$@"

