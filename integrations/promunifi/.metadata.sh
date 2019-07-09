# Each line must have an export clause.
# This file is parsed and sourced by the Makefile, Docker and Homebrew builds.

# Must match the repo name.
BINARY="unifi-poller"
# github username
GHUSER="davidnewhall"
# docker hub username
DHUSER="golift"
MAINT="David Newhall II <david at sleepers dot pro>"
VENDOR="Go Lift"
DESC="Polls a UniFi controller and stores metrics in InfluxDB"
GOLANGCI_LINT_ARGS="--enable-all -D gochecknoglobals"
# Example must exist at examples/$CONFIG_FILE.example
CONFIG_FILE="up.conf"
LICENSE="MIT"

export BINARY GHUSER DHUSER MAINT VENDOR DESC GOLANGCI_LINT_ARGS CONFIG_FILE LICENSE

# The rest is mostly automatic.

GHREPO="${GHUSER}/${BINARY}"
URL="https://github.com/${GHREPO}"

# This parameter is passed in as -X to go build. Used to override the Version variable in a package.
# This makes a path like github.com/davidnewhall/unifi-poller/unifipoller.Version=1.3.3
# Name the Version-containing library the same as the github repo, without dashes.
VERSION_PATH="github.com/${GHREPO}/$(echo ${BINARY} | tr -d -- -).Version"

# Dynamic. Recommend not changing.
VERSION="$(git tag -l --merged | tail -n1 | tr -d v || echo development)"
# This produces a 0 in some envirnoments (like Homebrew), but it's only used for packages.
ITERATION=$(git rev-list --count --all || echo 0)
DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
COMMIT="$(git rev-parse --short HEAD || echo 0)"

export GHREPO URL VERSION_PATH VERSION ITERATION DATE COMMIT
