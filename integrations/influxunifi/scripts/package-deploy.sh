#!/bin/bash

# Deploy our built packages to jfrog bintray.

COMPONENT=unstable
if [ "$TRAVIS_BRANCH" == "$TRAVIS_TAG" ] && [ "$TRAVIS_BRANCH" != "" ]; then
  COMPONENT=main
fi
echo "deploying packages from branch: $TRAVIS_BRANCH, tag: $TRAVIS_TAG to repo: $COMPONENT"

source .metadata.sh

for os in el centos; do
  for arch in arm64 armhf x86_64 i386; do
    file="unifi-poller-${VERSION}-${ITERATION}.${arch}.rpm"
    opts="publish=1;override=1"
    url="https://api.bintray.com/content/golift/${os}/unifi-poller/${VERSION}-${ITERATION}/${COMPONENT}/${arch}/${file}"
    echo curl -T "release/${file}" "${url};${opts}"
    curl -T "release/${file}" -u "${JFROG_USER_API_KEY}" "${url};${opts}"
    echo
  done
done

for os in ubuntu debian; do
  for arch in arm64 armhf amd64 i386; do
    file="unifi-poller_${VERSION}-${ITERATION}_${arch}.deb"
    opts="deb_distribution=xenial,bionic,focal,jesse,stretch,buster,bullseye;deb_component=${COMPONENT};deb_architecture=${arch};publish=1;override=1"
    url="https://api.bintray.com/content/golift/${os}/unifi-poller/${VERSION}-${ITERATION}/${file}"
    echo curl -T "release/${file}" "${url};${opts}"
    curl -T "release/${file}" -u "${JFROG_USER_API_KEY}" "${url};${opts}"
    echo
  done
done
