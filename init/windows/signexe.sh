#!/usr/bin/env bash

set -e -o pipefail

# Authenticode-sign a Windows PE via golift/codesign (YubiKey-backed signerd).
# GoReleaser calls this from builds.hooks.post on windows binaries, after the
# CLI is installed into a temp GOBIN (never /usr/bin/codesign).
#
# Local snapshots skip when CODESIGN_URL is unset. GitHub Actions must fail
# closed — the release contract is an Authenticode-signed Windows binary.
#
# Prefer CODESIGN_BIN, then GOBIN/codesign, then GOPATH/bin, then PATH
# entries that are not Apple's /usr/bin/codesign.

function pick_codesign() {
  if [ -n "${CODESIGN_BIN:-}" ]; then
    echo "${CODESIGN_BIN}"
    return
  fi
  if [ -n "${GOBIN:-}" ] && [ -x "${GOBIN}/codesign" ]; then
    echo "${GOBIN}/codesign"
    return
  fi
  gopath="$(go env GOPATH 2>/dev/null || true)"
  if [ -n "${gopath}" ] && [ -x "${gopath}/bin/codesign" ]; then
    echo "${gopath}/bin/codesign"
    return
  fi
  while IFS= read -r p; do
    case "$p" in
      /usr/bin/codesign|/bin/codesign) continue ;;
    esac
    echo "$p"
    return
  done < <(type -a -p codesign 2>/dev/null || true)
  return 1
}

function sign() {
  if [ -z "${CODESIGN_URL:-}" ]; then
    if [ -n "${GITHUB_ACTIONS:-}" ]; then
      echo "CODESIGN_URL unset; refusing to ship an unsigned Windows binary in CI" >&2
      exit 1
    fi
    echo "Skipped signing ${FILE} (CODESIGN_URL unset) .." >&2
    exit 0
  fi

  bin="$(pick_codesign)" || {
    echo "CODESIGN_URL is set but golift codesign CLI not found (set CODESIGN_BIN)" >&2
    exit 1
  }

  CODESIGN_NAME="${CODESIGN_NAME:-unpoller}" \
  CODESIGN_WEBSITE="${CODESIGN_WEBSITE:-https://unpoller.com}" \
  "${bin}" -- "${FILE}"
  echo "Signed ${FILE} .." >&2
}

[ -z "$1" ] || FILE="$1" sign
