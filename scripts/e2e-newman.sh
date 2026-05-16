#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
collection="${repo_root}/tests/e2e/postman/opendart-smoke.postman_collection.json"
report="${repo_root}/test-results/newman.xml"
base_url="${OPENDART_BASE_URL:-https://opendart.fss.or.kr}"

if [[ -z "${OPENDART_API_KEY:-}" ]]; then
  echo "error: OPENDART_API_KEY is required for live OpenDART e2e smoke tests." >&2
  echo "hint: OPENDART_API_KEY=... scripts/e2e-newman.sh" >&2
  exit 1
fi

mkdir -p "$(dirname "${report}")"

if command -v newman >/dev/null 2>&1; then
  newman_cmd=(newman)
elif command -v npx >/dev/null 2>&1; then
  newman_cmd=(npx --yes newman)
else
  echo "error: newman is not installed and npx is not available." >&2
  echo "hint: install newman or run with a Node.js environment that provides npx." >&2
  exit 1
fi

echo "Running OpenDART live e2e smoke against ${base_url}"
echo "Writing Newman JUnit report to ${report}"

"${newman_cmd[@]}" run "${collection}" \
  --env-var "base_url=${base_url}" \
  --env-var "crtfc_key=${OPENDART_API_KEY}" \
  --reporters junit \
  --reporter-junit-export "${report}"

echo "OpenDART live e2e smoke passed."
