#!/bin/bash

# Opens a PR that records a monthly release in the Release Please manifest and
# the Helm chart. The monthly workflow tags releases on its own, so without this
# Release Please keeps computing the next version from a stale manifest and
# proposes versions that already exist.

set -euo pipefail
IFS=$'\n\t'

MANIFEST_FILE=".release-please-manifest.json"
CHART_FILE="charts/cert-manager-webhook-infoblox-wapi/Chart.yaml"

main() {
  if [[ -z "${VERSION:-}" ]]; then
    echo "Error: VERSION is not set."
    exit 1
  fi

  if [[ -z "${GH_TOKEN:-}" ]]; then
    echo "Error: GH_TOKEN is not set."
    exit 1
  fi

  local version="${VERSION#v}"
  if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: VERSION '$VERSION' is not a semantic version."
    exit 1
  fi

  local branch="chore/sync-release-please-${version}"
  echo "Syncing Release Please manifest and chart to $version on branch $branch"

  if git ls-remote --exit-code --heads origin "$branch" >/dev/null; then
    echo "Branch $branch already exists, nothing to do."
    return 0
  fi

  git switch --create "$branch"

  jq --arg v "$version" '.["."] = $v' "$MANIFEST_FILE" >"${MANIFEST_FILE}.tmp"
  mv "${MANIFEST_FILE}.tmp" "$MANIFEST_FILE"

  sed -i -E "s/^(version|appVersion): [^ ]+ # x-release-please-version$/\1: ${version} # x-release-please-version/" "$CHART_FILE"

  if git diff --quiet; then
    echo "Manifest and chart already at $version, nothing to do."
    return 0
  fi
  git diff

  git config user.name "github-actions[bot]"
  git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
  git commit --all --message "chore(release): sync Release Please manifest and chart to ${version}"
  git push origin "$branch"

  gh pr create \
    --base main \
    --head "$branch" \
    --title "chore(release): sync Release Please manifest and chart to ${version}" \
    --body "The monthly release published v${version} outside Release Please. This updates \`${MANIFEST_FILE}\` and the chart \`version\`/\`appVersion\` to ${version} so the next Release Please PR proposes the correct version and the chart points at the new image."
}

if ! (return 0 2>/dev/null); then
  (main "$@")
fi
