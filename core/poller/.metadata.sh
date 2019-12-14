# Each line must have an export clause.
# This file is parsed and sourced by the Makefile, Docker and Homebrew builds.
# Powered by Application Builder: https://github.com/golift/application-builder

# Must match the repo name.
BINARY="unifi-poller"
# github username
GHUSER="davidnewhall"
# Github repo containing homebrew formula repo.
HBREPO="golift/homebrew-mugs"
MAINT="David Newhall II <david at sleepers dot pro>"
VENDOR="Go Lift <code at golift dot io>"
DESC="Polls a UniFi controller, exports metrics to InfluxDB and Prometheus"
GOLANGCI_LINT_ARGS="--enable-all -D gochecknoglobals -D dupl -D lll -D funlen -D wsl -e G402"
# Example must exist at examples/$CONFIG_FILE.example
CONFIG_FILE="up.conf"
LICENSE="MIT"
# FORMULA is either 'service' or 'tool'. Services run as a daemon, tools do not.
# This affects the homebrew formula (launchd) and linux packages (systemd).
FORMULA="service"

export BINARY GHUSER HBREPO MAINT VENDOR DESC GOLANGCI_LINT_ARGS CONFIG_FILE LICENSE FORMULA

# The rest is mostly automatic.
# Fix the repo if it doesn't match the binary name.
# Provide a better URL if one exists.

# Used as go import path in docker and homebrew builds.
IMPORT_PATH="github.com/${GHUSER}/${BINARY}"
# Used for source links and wiki links.
SOURCE_URL="https://${IMPORT_PATH}"
# Used for documentation links.
URL="${SOURCE_URL}"

# This parameter is passed in as -X to go build. Used to override the Version variable in a package.
# This makes a path like github.com/user/hello-world/helloworld.Version=1.3.3
# Name the Version-containing library the same as the github repo, without dashes.
VERSION_PATH="${IMPORT_PATH}/pkg/poller.Version"

# Dynamic. Recommend not changing.
VVERSION=$(git describe --abbrev=0 --tags $(git rev-list --tags --max-count=1))
VERSION="$(echo $VVERSION | tr -d v | grep -E '^\S+$' || echo development)"
# This produces a 0 in some envirnoments (like Homebrew), but it's only used for packages.
ITERATION=$(git rev-list --count --all || echo 0)
DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
BRANCH="$(git rev-parse --abbrev-ref HEAD || echo unknown)"
COMMIT="$(git rev-parse --short HEAD || echo 0)"

# Used by homebrew downloads.
#SOURCE_PATH=https://codeload.${IMPORT_PATH}/tar.gz/v${VERSION}
# This is a custom download path for homebrew formula.
SOURCE_PATH=https://golift.io/${BINARY}/archive/v${VERSION}.tar.gz

export IMPORT_PATH SOURCE_URL URL VERSION_PATH VVERSION VERSION ITERATION DATE BRANCH COMMIT SOURCE_PATH
