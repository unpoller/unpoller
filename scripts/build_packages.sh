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

# Make a build environment.
mkdir -p package_build/usr/bin package_build/etc/${BINARY} package_build/lib/systemd/system package_build/usr/share/man/man1

# Copy the binary, config file and man page into the env.
cp ${BINARY} package_build/usr/bin
cp *.1.gz package_build/usr/share/man/man1
cp examples/up.conf.example package_build/etc/${BINARY}/up.conf

# Fix the paths in the systemd unit file before copying it into the emv.
sed "s#ExecStart.*#ExecStart=/usr/bin/${BINARY} --config=/etc/${BINARY}/up.conf#" \
  init/systemd/unifi-poller.service > package_build/lib/systemd/system/${BINARY}.service

fpm -s dir -t ${OUTPUT} \
  -n ${BINARY} \
  -v ${VERSION} \
  --after-install scripts/after-install.sh \
  -C package_build
