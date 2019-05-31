#!/bin/bash

# This script builds a simple macos Installer pkg. Run by the Makefile.
# Use: `make osx`

OUTPUT=$1
BINARY=unifi-poller
VERSION=$(git tag -l --merged | tail -n1 | tr -d v)

fpm -h > /dev/null 2>&1
if [ "$?" != "0" ]; then
  echo "fpm missing. Install fpm: https://fpm.readthedocs.io/en/latest/installing.html"
  exit 1
fi

echo "Building 'osxpkg' package."

PREFIX=/usr/local
BINFIX=/usr/local

# Make a build environment.
rm -rf package_build
mkdir -p package_build${BINFIX}/bin package_build${PREFIX}/etc/${BINARY} package_build${BINFIX}/share/man/man1
mkdir -p package_build${PREFIX}/var/log

# Copy the binary, config file and man page into the env.
cp ${BINARY} package_build${BINFIX}/bin
cp *.1.gz package_build${BINFIX}/share/man/man1
cp examples/up.conf.example package_build${PREFIX}/etc/${BINARY}/

# Copy in launch agent.
mkdir -p package_build/Library/LaunchAgents
cp init/launchd/com.github.davidnewhall.unifi-poller.plist package_build/Library/LaunchAgents/

# Make a package.
fpm -s dir -t osxpkg \
  --name ${BINARY} \
  --version ${VERSION} \
  --after-install scripts/after-install-osx.sh \
  --osxpkg-identifier-prefix com.github.davidnewhall \
  --chdir package_build
