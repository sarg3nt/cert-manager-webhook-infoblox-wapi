#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

main() {
  if [[ -z "${REGISTRY:-}" ]]; then
    echo "Error: REGISTRY is not set."
    exit 1
  fi

  if [[ -z "${IMAGE_NAME:-}" ]]; then
    echo "Error: IMAGE_NAME is not set."
    exit 1
  fi

  echo "Getting latest tag using gh and the GitHub API"
  latest_version=$(gh api /repos/sarg3nt/cert-manager-webhook-infoblox-wapi/tags |
    jq -r '[.[] | select(.name | startswith("v"))] | sort_by(.name) | reverse | .[0].name')
  echo "Latest Version: $latest_version"
  # Github runner does not print empty echos. :(
  echo "-"

  echo "Finding the latest tag version and setting major, minor, patch and new_patch."
  if [[ $latest_version =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
    major=${BASH_REMATCH[1]}
    echo "Major: $major"
    minor=${BASH_REMATCH[2]}
    echo "Minor: $minor"
    patch=${BASH_REMATCH[3]}
    echo "Old Patch: $patch"
    new_patch=$((patch + 1))
    echo "New Patch: $new_patch"

    tag_major="${REGISTRY}/${IMAGE_NAME}:${major}"
    tag_minor="${REGISTRY}/${IMAGE_NAME}:${major}.${minor}"
    tag_patch="${REGISTRY}/${IMAGE_NAME}:${major}.${minor}.${new_patch}"
    tag_old="${REGISTRY}/${IMAGE_NAME}:${major}.${minor}.${patch}"
    tag_latest="${REGISTRY}/${IMAGE_NAME}:latest"

    echo "TAG_MAJOR: $tag_major"
    echo "TAG_MAJOR=$tag_major" >>"$GITHUB_ENV"
    echo "TAG_MINOR: $tag_minor"
    echo "TAG_MINOR=$tag_minor" >>"$GITHUB_ENV"
    echo "TAG_PATCH: $tag_patch"
    echo "TAG_PATCH=$tag_patch" >>"$GITHUB_ENV"
    echo "TAG_LATEST: $tag_latest"
    echo "TAG_LATEST=$tag_latest" >>"$GITHUB_ENV"
    echo "TAG_OLD: $tag_old"
    echo "TAG_OLD=$tag_old" >>"$GITHUB_ENV"

    new_version="v${major}.${minor}.${new_patch}"
    echo "New Version: $new_version"
    echo "VERSION=$new_version" >>"$GITHUB_ENV"
  else
    echo "Could not determine the latest tag version."
    exit 1
  fi
}

if ! (return 0 2>/dev/null); then
  (main "$@")
fi
