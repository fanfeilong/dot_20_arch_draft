#!/usr/bin/env sh

set -eu

REPO="${D2A_REPO:-fanfeilong/dot_20_arch_draft}"
VERSION="${D2A_VERSION:-latest}"
INSTALL_DIR="${D2A_INSTALL_DIR:-/usr/local/bin}"

need_cmd() {
	command -v "$1" >/dev/null 2>&1 || {
		echo "missing required command: $1" >&2
		exit 1
	}
}

detect_os() {
	case "$(uname -s)" in
		Darwin) echo "darwin" ;;
		Linux) echo "linux" ;;
		*)
			echo "unsupported operating system: $(uname -s)" >&2
			exit 1
			;;
	esac
}

detect_arch() {
	case "$(uname -m)" in
		x86_64|amd64) echo "amd64" ;;
		arm64|aarch64) echo "arm64" ;;
		*)
			echo "unsupported architecture: $(uname -m)" >&2
			exit 1
			;;
	esac
}

release_base() {
	if [ "$VERSION" = "latest" ]; then
		echo "https://github.com/$REPO/releases/latest/download"
	else
		echo "https://github.com/$REPO/releases/download/$VERSION"
	fi
}

need_cmd curl
need_cmd tar
need_cmd mktemp

OS="$(detect_os)"
ARCH="$(detect_arch)"
ASSET="d2a_${OS}_${ARCH}.tar.gz"
BASE_URL="$(release_base)"
TMP_DIR="$(mktemp -d)"
ARCHIVE_PATH="$TMP_DIR/$ASSET"

cleanup() {
	rm -rf "$TMP_DIR"
}

trap cleanup EXIT INT TERM

echo "Downloading $ASSET from $REPO ($VERSION)..."
curl -fsSL "$BASE_URL/$ASSET" -o "$ARCHIVE_PATH"

mkdir -p "$INSTALL_DIR"
tar -xzf "$ARCHIVE_PATH" -C "$TMP_DIR"
install -m 0755 "$TMP_DIR/d2a" "$INSTALL_DIR/d2a"

echo "Installed d2a to $INSTALL_DIR/d2a"
echo "Run: d2a help"
