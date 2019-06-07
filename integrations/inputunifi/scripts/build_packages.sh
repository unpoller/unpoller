#!/bin/bash

# This script builds a deb, rpm or osx package. Run by the Makefile.
# Use: `make rpm`, `make deb`, or `make osxpkg`

set -e -o pipefail

BINARY=unifi-poller
OUTPUT=$1
VERSION=$2
ITERATION=$3
[ "$VERSION" != "" ] || VERSION=$(git tag -l --merged | tail -n1 | tr -d v)
[ "$ITERATION" != "" ] || ITERATION=$(git rev-list --all --count)

if [ "$OUTPUT" != "deb" ] && [ "$OUTPUT" != "rpm" ] && [ "$OUTPUT" != "osxpkg" ]; then
    echo "first argument must be 'deb' or 'rpm' or 'osxpkg'"
    exit 1
fi

fpm -h > /dev/null 2>&1
if [ "$?" != "0" ]; then
    echo "Package Build Failure. FPM missing. Install FPM: https://fpm.readthedocs.io/en/latest/installing.html"
    exit 1
fi

echo "Building '${OUTPUT}' package for ${BINARY} version '${VERSION}-${ITERATION}'."

# These paths work on Linux. Suggest not changing.
PREFIX=
BINFIX=/usr
UNAME=linux
AFTER=scripts/after-install.sh
if [ "$OUTPUT" = "osxpkg" ]; then
  # These paths work on OSX. Do not change.
  PREFIX=/usr/local
  BINFIX=/usr/local
  UNAME=macos
  AFTER=scripts/after-install-osx.sh
fi

# Make a build environment.
rm -rf package_build
mkdir -p package_build${BINFIX}/bin package_build${PREFIX}/etc/${BINARY}
mkdir -p package_build${BINFIX}/share/man/man1 package_build${BINFIX}/share/doc/unifi-poller

# Copy the binary, config file and man page into the env.
cp ${BINARY}.${UNAME} package_build${BINFIX}/bin/${BINARY}
cp *.1.gz package_build${BINFIX}/share/man/man1
cp examples/*.conf.example package_build${PREFIX}/etc/${BINARY}/
cp examples/* package_build${BINFIX}/share/doc/unifi-poller

# Copy startup file. Different for osx vs linux.
if [ "$UNAME" = "linux" ]; then
    cp examples/up.conf.example package_build${PREFIX}/etc/${BINARY}/up.conf
    # Fix the paths in the systemd unit file before copying it into the emv.
    mkdir -p package_build/lib/systemd/system
    sed "s#ExecStart.*#ExecStart=${BINFIX}/bin/${BINARY} --config=${PREFIX}/etc/${BINARY}/up.conf#" \
      init/systemd/unifi-poller.service > package_build/lib/systemd/system/${BINARY}.service

else # macos
    # Sometimes the log folder is missing on osx. Create it.
    mkdir -p package_build${PREFIX}/var/log
    mkdir -p package_build/Library/LaunchAgents
    cp init/launchd/com.github.davidnewhall.unifi-poller.plist package_build/Library/LaunchAgents/
fi

# Make a package.
fpm -s dir -t ${OUTPUT} \
  --name ${BINARY} \
  --version ${VERSION} \
  --iteration ${ITERATION} \
  --after-install ${AFTER} \
  --before-remove scripts/before-remove.sh \
  --osxpkg-identifier-prefix com.github.davidnewhall \
  --license MIT \
  --url 'https://github.com/davidnewhall/unifi-poller' \
  --maintainer 'david at sleepers dot pro' \
  --description 'This daemon polls a Unifi controller at a short interval and stores the collected metric data in an Influx Database.' \
  --chdir package_build
