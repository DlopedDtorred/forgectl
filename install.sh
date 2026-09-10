#!/usr/bin/env bash
#
# forgectl - Installer
# ------------------------------------------------
# Installs the forgectl CLI by building from source.
#
# Usage:
#   ./install.sh                 # Build & install to /usr/local/bin
#   ./install.sh --prefix ~/.local  # Install to a custom prefix
#   ./install.sh --bin-dir ~/bin    # Install to a custom bin dir
#   ./install.sh --remove           # Remove forgectl
#   ./install.sh --help             # Show this help
#
# Defaults:
#   BIN_DIR = /usr/local/bin
#
# Requirements:
#   - Go 1.22+ (only when building from source)
#   - git   (to fetch dependencies)
#   - make  (optional, falls back to go build)
# ------------------------------------------------

set -euo pipefail

# ---------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------

VERSION="0.2.0"
REPO_URL="https://github.com/forgectl/forgectl"
PROJECT_NAME="forgectl"
PREFIX="${PREFIX:-/usr/local}"
BIN_DIR="${BIN_DIR:-${PREFIX}/bin}"
SHARE_DIR="${PREFIX}/share/${PROJECT_NAME}"

GREEN=$'\033[32m'
YELLOW=$'\033[33m'
RED=$'\033[31m'
CYAN=$'\033[36m'
BOLD=$'\033[1m'
RESET=$'\033[0m'

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

info()  { printf '%s  %s%s\n' "$CYAN" "$*" "$RESET"; }
ok()    { printf '%s  ✓ %s%s\n' "$GREEN" "$*" "$RESET"; }
warn()  { printf '%s  ⚠ %s%s\n' "$YELLOW" "$*" "$RESET"; }
err()   { printf '%s  ✗ %s%s\n' "$RED" "$*" "$RESET" >&2; }
die()   { err "$1"; exit "${2:-1}"; }

usage() {
    sed -n '2,20p' "$0" | sed 's/^# \{0,1\}//'
    exit 0
}

require() {
    if ! command -v "$1" >/dev/null 2>&1; then
        die "'$1' is required but was not found in PATH."
    fi
}

# ---------------------------------------------------------------------------
# Detect source directory
# ---------------------------------------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
IS_REPO_ROOT=false

if [ -f "${SCRIPT_DIR}/go.mod" ]; then
    IS_REPO_ROOT=true
    SRC_DIR="${SCRIPT_DIR}"
fi

# ---------------------------------------------------------------------------
# Actions
# ---------------------------------------------------------------------------

action_build() {
    echo
    info "${BOLD}forgectl ${VERSION} installer${RESET}"
    info "--------------------------------------------"
    info "Source:  ${SRC_DIR}"
    info "Bin dir: ${BIN_DIR}"
    echo

    require "go"
    require "git"

    local build_dir
    build_dir="$(mktemp -d)"

    if [ "$IS_REPO_ROOT" = "true" ]; then
        ok "Developer install: building from repository clone"
        build_dir="$SRC_DIR"
    else
        info "Downloading source from ${REPO_URL}..."
        git clone --depth 1 --branch "v${VERSION}" "$REPO_URL" "$build_dir" 2>/dev/null \
            || git clone --depth 1 "$REPO_URL" "$build_dir"
        ok "Source downloaded"
    fi

    # Build
    info "Building ${PROJECT_NAME}..."
    (
        cd "$build_dir"
        if [ -f Makefile ] && command -v make >/dev/null 2>&1; then
            make build || go build -ldflags "-s -w" -o "$PROJECT_NAME" .
        else
            go build -ldflags "-s -w" -o "$PROJECT_NAME" .
        fi
    )
    ok "Binary built"

    # Install
    info "Installing to ${BIN_DIR}..."
    install -d "$BIN_DIR"
    install -m 0755 "${build_dir}/${PROJECT_NAME}" "${BIN_DIR}/${PROJECT_NAME}"
    ok "Installed ${BIN_DIR}/${PROJECT_NAME}"

    # Install shared files (completions, etc.) — best effort only
    if [ -d "${build_dir}/.github" ]; then
        if [ -d "$PREFIX" ] && [ -w "$PREFIX" ]; then
            install -d "$SHARE_DIR"
            cp -r "${build_dir}/.github" "${SHARE_DIR}/" 2>/dev/null || true
        else
            warn "Skipping shared files (${SHARE_DIR} not writable)"
        fi
    fi

    echo
    ok "${BOLD}forgectl ${VERSION} installed successfully!${RESET}"
    info "  Run 'forgectl --help' to get started"
    info "  Run 'forgectl doctor' to verify your setup"
    echo
}

action_remove() {
    echo
    warn "Removing forgectl..."
    if [ -f "${BIN_DIR}/${PROJECT_NAME}" ]; then
        rm -f "${BIN_DIR}/${PROJECT_NAME}"
        ok "Removed ${BIN_DIR}/${PROJECT_NAME}"
    else
        warn "No forgectl binary found in ${BIN_DIR}"
    fi
    if [ -d "$SHARE_DIR" ]; then
        rm -rf "$SHARE_DIR"
        ok "Removed ${SHARE_DIR}"
    fi
    rm -f "${HOME}/.local/share/bash-completion/completions/forgectl" \
          "${HOME}/.config/fish/completions/forgectl.fish" 2>/dev/null || true
    echo
    info "forgectl was removed from your system. Config in ~/.forgectl was kept."
    echo
}

action_version() {
    echo "${VERSION}"
    exit 0
}

# ---------------------------------------------------------------------------
# Parse arguments
# ---------------------------------------------------------------------------

while [ $# -gt 0 ]; do
    case "$1" in
        --prefix)        PREFIX="$2"; BIN_DIR="${PREFIX}/bin"; shift 2 ;;
        --bin-dir)       BIN_DIR="$2"; shift 2 ;;
        --remove|uninstall) DO_REMOVE=true; shift ;;
        --version)       action_version ;;
        --help|-h)       usage ;;
        *)               err "Unknown option: $1" ; usage ;;
    esac
done

if [ "${DO_REMOVE:-}" = "true" ]; then
    action_remove
else
    action_build
fi