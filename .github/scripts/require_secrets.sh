#!/usr/bin/env bash
# Fail closed when a publishing secret is empty. GitHub maps missing secrets
# to "" so `if: secrets.FOO != ''` is not a substitute.
set -euo pipefail

missing=0
need() {
  local name=$1
  local val=${!name-}
  if [ -z "${val}" ]; then
    echo "missing secret: ${name}" >&2
    missing=1
  fi
}

need GORELEASER_KEY
need GPG_SIGNING_KEY
need DOCKER_USERNAME
need DOCKER_PASSWORD
need CODESIGN_URL
need CODESIGN_CLIENT_CERT
need CODESIGN_CLIENT_KEY
need MACOS_SIGN_P12
need MACOS_SIGN_PASSWORD
need MACOS_NOTARY_KEY
need MACOS_NOTARY_KEY_ID
need MACOS_NOTARY_ISSUER_ID
need PACKAGECLOUD_TOKEN
need HOMEBREW_TAP_GITHUB_TOKEN

if [ "${missing}" -ne 0 ]; then
  echo "refusing to publish with empty secrets" >&2
  exit 1
fi

echo "required secrets are present"
