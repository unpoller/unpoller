# Each line must have an export clause.
# This file is parsed and sourced by the Makefile, Docker and Homebrew builds.
# Powered by Application Builder: https://github.com/golift/application-builder

# Bring in dynamic repo/pull/source info.
source $(dirname "${BASH_SOURCE[0]}")/init/buildinfo.sh

# Must match the repo name.
BINARY="unpoller"
# Github repo containing homebrew formula repo.
HBREPO="golift/homebrew-mugs"
MAINT="David Newhall II <david at sleepers dot pro>"
VENDOR="Go Lift <code at golift dot io>"
DESC="Polls a UniFi controller, exports metrics to InfluxDB and Prometheus"
GOLANGCI_LINT_ARGS="--enable-all"
# Example must exist at examples/$CONFIG_FILE.example
CONFIG_FILE="up.conf"
LICENSE="MIT"
# FORMULA is either 'service' or 'tool'. Services run as a daemon, tools do not.
# This affects the homebrew formula (launchd) and linux packages (systemd).
FORMULA="service"

# Used for source links and wiki links.
SOURCE_URL="https://github.com/${BINARY}/${BINARY}"

VERSION_PATH="golift.io/version"

# This is a custom download path for homebrew formula.
SOURCE_PATH=https://golift.io/${BINARY}/archive/v${VERSION}.tar.gz

export BINARY HBREPO MAINT VENDOR DESC GOLANGCI_LINT_ARGS CONFIG_FILE
export LICENSE FORMULA SOURCE_URL VERSION_PATH SOURCE_PATH

### Optional ###

# Import this signing key only if it's in the keyring.
gpg --list-keys 2>/dev/null | grep -q B93DD66EF98E54E2EAE025BA0166AD34ABC5A57C
[ "$?" != "0" ] || export SIGNING_KEY=B93DD66EF98E54E2EAE025BA0166AD34ABC5A57C
