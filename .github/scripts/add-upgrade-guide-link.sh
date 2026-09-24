#!/bin/bash
# Prepend a link to the upgrading guide (docs/UPGRADE_*_to_vX.Y.md) to the
# release notes of a vX.Y.0 release.
#
# Usage: add-upgrade-guide-link.sh <tag>
#   DRY_RUN=1 prints the new release notes without updating the release.
set -euo pipefail

TAG="${1:?usage: $0 <tag>}"
REPO="${GITHUB_REPOSITORY:-sacloud/sacloud-otel-collector}"

if [[ ! "$TAG" =~ ^v([0-9]+)\.([0-9]+)\.0$ ]]; then
  echo "$TAG is not a vX.Y.0 release, skipped"
  exit 0
fi
MINOR="v${BASH_REMATCH[1]}.${BASH_REMATCH[2]}"

GUIDE=$(find docs -maxdepth 1 -name "UPGRADE_*_to_${MINOR}.md" | sort | head -n 1)
if [ -z "$GUIDE" ]; then
  echo "no upgrading guide for $MINOR found, skipped"
  exit 0
fi

URL="https://github.com/${REPO}/blob/${TAG}/${GUIDE}"
BODY=$(gh release view "$TAG" --repo "$REPO" --json body --jq .body)
if [[ "$BODY" == *"$URL"* ]]; then
  echo "release notes of $TAG already have the link, skipped"
  exit 0
fi

NOTES=$(mktemp)
trap 'rm -f "$NOTES"' EXIT
{
  echo "## Upgrading"
  echo
  echo "This release contains breaking changes. See [$(basename "$GUIDE")]($URL) before upgrading."
  echo
  echo "$BODY"
} > "$NOTES"

if [ -n "${DRY_RUN:-}" ]; then
  cat "$NOTES"
  exit 0
fi
gh release edit "$TAG" --repo "$REPO" --notes-file "$NOTES"
echo "added the link to $GUIDE to release notes of $TAG"
