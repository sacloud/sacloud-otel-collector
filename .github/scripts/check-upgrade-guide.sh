#!/bin/bash
# Check that the upgrade guide (docs/UPGRADE_vX.Y_to_vX.Z.md) exists when the
# version is bumped to a new minor (or major) version.
#
# Usage: check-upgrade-guide.sh <current version> <next version>
#   e.g. check-upgrade-guide.sh v0.8.2 v0.9.0
set -euo pipefail

CURRENT="${1:?usage: $0 <current version> <next version>}"
NEXT="${2:?usage: $0 <current version> <next version>}"

minor_of() {
  if [[ ! "$1" =~ ^v([0-9]+)\.([0-9]+)\.[0-9]+$ ]]; then
    echo "::error::$1 is not a vX.Y.Z version" >&2
    exit 1
  fi
  echo "v${BASH_REMATCH[1]}.${BASH_REMATCH[2]}"
}

FROM=$(minor_of "$CURRENT")
TO=$(minor_of "$NEXT")

if [ "$FROM" = "$TO" ]; then
  echo "$CURRENT -> $NEXT is not a minor version bump, no upgrade guide required"
  exit 0
fi

GUIDE="docs/UPGRADE_${FROM}_to_${TO}.md"
if [ ! -f "$GUIDE" ]; then
  echo "::error::$CURRENT -> $NEXT is a minor version bump, but $GUIDE does not exist"
  exit 1
fi
echo "$CURRENT -> $NEXT is a minor version bump, $GUIDE exists"
