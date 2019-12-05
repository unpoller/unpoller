#!/bin/bash

# Deploy our built packages to packagecloud.

REPO=unstable
[ "$TRAVIS_BRANCH" != "$TRAVIS_TAG" ] || REPO=stable
echo "deploying packages from branch: $TRAVIS_BRANCH, tag: $TRAVIS_TAG to repo: $REPO"

source .metadata.sh
# deb
cmd="package_cloud push golift/${REPO}/debian/trusty"
$cmd release/unifi-poller_${VERSION}-${ITERATION}_arm64.deb
$cmd release/unifi-poller_${VERSION}-${ITERATION}_amd64.deb
$cmd release/unifi-poller_${VERSION}-${ITERATION}_armhf.deb
$cmd release/unifi-poller_${VERSION}-${ITERATION}_i386.deb
# rpm
cmd="package_cloud push golift/${REPO}/el/6"
$cmd release/unifi-poller-${VERSION}-${ITERATION}.arm64.rpm
$cmd release/unifi-poller-${VERSION}-${ITERATION}.x86_64.rpm
$cmd release/unifi-poller-${VERSION}-${ITERATION}.armhf.rpm
$cmd release/unifi-poller-${VERSION}-${ITERATION}.i386.rpm
