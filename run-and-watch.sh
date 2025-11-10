#!/usr/bin/env bash

set -euo pipefail

workflow="$1"
uid="$2"
shift; shift

gh workflow run "$workflow" "$@"
./watch-workflow.sh "$workflow" "$uid"

