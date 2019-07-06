#!/bin/bash

# This script creates multi-architecture manifests from images that were built in Docker Cloud.
# Wait for all the images to build then run this to update them.

# TODO: Get someone at Docker to add this post-build config to Docker Cloud Auto Builds.

set -e -o pipefail

CMD=$1
REPO=golift/unifi-poller
VERSION=$(git tag -l --merged | tail -n1 | tr -d v)

USER=$(echo $REPO | cut -d/ -f1)
/bin/echo -n "Docker User '$USER' "
docker login --username=$USER

docker manifest create --amend ${REPO}:latest ${REPO}:latest_linux_amd64 ${REPO}:latest_linux_arm
docker manifest annotate ${REPO}:latest ${REPO}:latest_linux_arm --os linux --arch arm
docker manifest push ${REPO}:latest

if [ "$CMD" = "release" ]; then
  # stable tag. stable
  docker manifest create ${REPO}:stable ${REPO}:${VERSION}_linux_amd64 ${REPO}:${VERSION}_linux_arm
  docker manifest annotate ${REPO}:stable ${REPO}:${VERSION}_linux_arm --os linux --arch arm
  docker manifest push ${REPO}:stable
  # version tag. 1.2.3
  docker manifest create ${REPO}:${VERSION} ${REPO}:${VERSION}_linux_amd64 ${REPO}:${VERSION}_linux_arm
  docker manifest annotate ${REPO}:${VERSION} ${REPO}:${VERSION}_linux_arm --os linux --arch arm
  docker manifest push ${REPO}:${VERSION}
  # short version tag. 1.2
  VER=$(echo $VERSION | cut -d. -f1,2)
  docker manifest create ${REPO}:${VER} ${REPO}:${VERSION}_linux_amd64 ${REPO}:${VERSION}_linux_arm
  docker manifest annotate ${REPO}:${VER} ${REPO}:${VERSION}_linux_arm --os linux --arch arm
  docker manifest push ${REPO}:${VER}
fi
