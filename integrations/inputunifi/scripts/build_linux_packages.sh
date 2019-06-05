#!/bin/bash

# This script builds a deb or rpm package. Run by the Makefile.
# Use: `make rpm` or `make deb`

OUTPUT=$1
BINARY=unifi-poller
VERSION=$(git tag -l --merged | tail -n1 | tr -d v)

if [ "$OUTPUT" != "deb" ] && [ "$OUTPUT" != "rpm" ]; then
  echo "first argument must be 'deb' or 'rpm'"
  exit 1
fi

fpm -h > /dev/null 2>&1
if [ "$?" != "0" ]; then
  echo "fpm missing. Install fpm: https://fpm.readthedocs.io/en/latest/installing.html"
  exit 1
fi

echo "Building '${OUTPUT}' package."

# eh, don't change these.
PREFIX=
BINFIX=/usr

# Make a build environment.
rm -rf package_build
mkdir -p package_build${BINFIX}/bin package_build${PREFIX}/etc/${BINARY} package_build${BINFIX}/share/man/man1

# Copy the binary, config file and man page into the env.
cp ${BINARY}.linux package_build${BINFIX}/bin/${BINARY}
cp *.1.gz package_build${BINFIX}/share/man/man1
cp examples/up.conf.example package_build${PREFIX}/etc/${BINARY}/up.conf

# Fix the paths in the systemd unit file before copying it into the emv.
mkdir -p package_build/lib/systemd/system
sed "s#ExecStart.*#ExecStart=${BINFIX}/bin/${BINARY} --config=${PREFIX}/etc/${BINARY}/up.conf#" \
  init/systemd/unifi-poller.service > package_build/lib/systemd/system/${BINARY}.service

# Make a package.
fpm -s dir -t ${OUTPUT} \
  --name ${BINARY} \
  --version ${VERSION} \
  --iteration $(git rev-list --all --count) \
  --after-install scripts/after-install.sh \
  --before-remove scripts/before-remove.sh \
  --license MIT \
  --url 'https://github.com/davidnewhall/unifi-poller' \
  --maintainer 'david at sleepers dot pro' \
  --description 'This daemon polls a Unifi controller at a short interval and stores the collected metric data in an Influx Database.' \
  --chdir package_build
