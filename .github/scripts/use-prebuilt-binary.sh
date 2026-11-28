#!/bin/sh
# Stand-in for the go tool in GoReleaser builds (see `tool` in .goreleaser.yml).
#
# The release workflow builds each target in its own job and collects the
# binaries in $PREBUILT_DIR, keeping the layout of GoReleaser's dist directory.
# GoReleaser then runs this script as `<tool> build ... -o <output> .`, and it
# copies the matching prebuilt binary to <output> instead of compiling.
set -eu

: "${PREBUILT_DIR:?PREBUILT_DIR is not set}"

if [ "${1:-}" != "build" ]; then
	echo "use-prebuilt-binary: unsupported command: $*" >&2
	exit 1
fi

out=""
while [ $# -gt 0 ]; do
	case "$1" in
	-o)
		out="${2:-}"
		break
		;;
	-o=*)
		out="${1#-o=}"
		break
		;;
	esac
	shift
done
if [ -z "$out" ]; then
	echo "use-prebuilt-binary: -o is not given" >&2
	exit 1
fi

# <dist>/<project>_<os>_<arch>_<variant>/<binary>
src="$PREBUILT_DIR/$(basename "$(dirname "$out")")/$(basename "$out")"
if [ ! -f "$src" ]; then
	echo "use-prebuilt-binary: prebuilt binary not found: $src" >&2
	exit 1
fi

mkdir -p "$(dirname "$out")"
cp "$src" "$out"
# actions/upload-artifact does not preserve file permissions
chmod 0755 "$out"
