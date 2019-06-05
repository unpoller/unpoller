#!/bin/bash

# This script builds a simple macos Installer pkg. Run by the Makefile.
# Use: `make osxpkg`

OUTPUT=osxpkg
BINARY=unifi-poller
VERSION=$(git tag -l --merged | tail -n1 | tr -d v)

fpm -h > /dev/null 2>&1
if [ "$?" != "0" ]; then
  echo "fpm missing. Install fpm: https://fpm.readthedocs.io/en/latest/installing.html"
  exit 1
fi

echo "Building '${OUTPUT}' package for ${BINARY} version ${VERSION}."

PREFIX=/usr/local
BINFIX=/usr/local

# Make a build environment.
rm -rf package_build
mkdir -p package_build${BINFIX}/bin package_build${PREFIX}/etc/${BINARY} package_build${BINFIX}/share/man/man1
mkdir -p package_build${PREFIX}/var/log

# Copy the binary, config file and man page into the env.
cp ${BINARY}.macos package_build${BINFIX}/bin/${BINARY}
cp *.1.gz package_build${BINFIX}/share/man/man1
cp examples/up.conf.example package_build${PREFIX}/etc/${BINARY}/

# Copy in launch agent.
mkdir -p package_build/Library/LaunchAgents
cp init/launchd/com.github.davidnewhall.unifi-poller.plist package_build/Library/LaunchAgents/

# Make a package.
fpm -s dir -t ${OUTPUT} \
  --name ${BINARY} \
  --version ${VERSION} \
  --iteration $(git rev-list --all --count) \
  --after-install scripts/after-install-osx.sh \
  --osxpkg-identifier-prefix com.github.davidnewhall \
  --license MIT \
  --maintainer 'david at sleepers dot pro' \
  --url 'https://github.com/davidnewhall/unifi-poller' \
  --description 'This daemon polls a Unifi controller at a short interval and stores the collected metric data in an Influx Database.' \
  --chdir package_build
