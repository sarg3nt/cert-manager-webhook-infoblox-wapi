#!/bin/bash

#cspell:ignore diffoci

set -euo pipefail
IFS=$'\n\t'

main() {
  if [ -z "${TAG_OLD:-}" ]; then
    echo "Error: TAG_OLD is not set."
    exit 1
  fi

  if [ -z "${TAG_PATCH:-}" ]; then
    echo "Error: TAG_PATCH is not set."
    exit 1
  fi

  echo "Downloading the diffoci binary."
  latest_release_url=$(gh release view -R reproducible-containers/diffoci --json assets -q '.assets[] | select(.name | test("linux-amd64")) | .url')
  echo "Using the latest release URL: $latest_release_url"
  curl -L -o diffoci "$latest_release_url"
  chmod +x diffoci
  # Github runner does not print empty echos. :(
  echo "-"

  echo "Pulling the previous Docker image to compare."
  docker pull "${TAG_OLD}"
  echo "-"

  echo "Checking if the images are different with diffoci."
  OLD_IMAGE="docker://${TAG_OLD}"
  NEW_IMAGE="docker://${TAG_PATCH}"
  set +e
  ./diffoci diff --semantic "$OLD_IMAGE" "$NEW_IMAGE"
  DIFFOCI_EXIT_CODE=$?
  set -e
  echo "-"

  # Check the exit code of diffoci. If it is zero then there are no changes, otherwise there are.
  if [ $DIFFOCI_EXIT_CODE -eq 0 ]; then
    echo "The images appear to be the same, exiting."
    echo "continue=false" >>"$GITHUB_OUTPUT"
    exit 0
  fi

  echo "The images appear to be different.  Continuing."
  echo "continue=true" >>"$GITHUB_OUTPUT"
}

if ! (return 0 2>/dev/null); then
  (main "$@")
fi
