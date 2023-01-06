#!/bin/bash

set -exo pipefail

export GITTAG=$1
export PACKAGECLOUD_TOKEN=$2
export PACKAGE_NAME=$3
export ARTIFACT_PATH=$4
export ARTIFACT_OS=$5

if [ ! -n "${GITTAG}" ]; then
  echo "GITTAG must be set"
  exit 1
fi

if [ ! -n "${PACKAGECLOUD_TOKEN}" ]; then
  echo "PACKAGECLOUD_TOKEN must be set"
  exit 1
fi

if [ ! -n "${PACKAGE_NAME}" ]; then
  echo "PACKAGE_NAME must be set"
  exit 1
fi

export PACKAGE_VERSION=${GITTAG}
export PACKAGE_DIR="./dist/${PACKAGE_NAME}_linux_amd64"

# NOTE: compatibility with goreleaser 1.8.3 and later
# See more: https://github.com/goreleaser/goreleaser/commit/63436392db6ac0557513535fc3ee4223a44810ed
if [[ -d "${PACKAGE_DIR}_v1" ]]; then
  export PACKAGE_DIR="${PACKAGE_DIR}_v1"
fi

if [[ ! -d "${PACKAGE_DIR}" ]]; then
  export PACKAGE_DIR="./dist/unpoller_linux_amd64"

  if [[ ! -d ${PACKAGE_DIR} ]]; then
    export PACKAGE_DIR="${PACKAGE_DIR}_v1"
  fi
fi

export PACKAGE_CLOUD_REPO="golift/pkgs"
if [[ ${PACKAGE_VERSION} =~ .+-rc ]]; then
  export PACKAGE_CLOUD_REPO="golift/unstable"
fi

export SUPPORTED_UBUNTU_VERSIONS="focal"
export SUPPORTED_REDHAT_VERSIONS="6"

if [[ $ARTIFACT_PATH == *termux* ]]; then
  # skip termux builds
  exit 0
fi

for ubuntu_version in ${SUPPORTED_UBUNTU_VERSIONS}
do
  if [[ $ARTIFACT_PATH == *.deb ]]; then
    package_cloud push ${PACKAGE_CLOUD_REPO}/ubuntu/${ubuntu_version} $ARTIFACT_PATH --skip-errors
  fi
done

for redhat_version in ${SUPPORTED_REDHAT_VERSIONS}
do
  if [[ $ARTIFACT_PATH == *.rpm ]]; then
    package_cloud push ${PACKAGE_CLOUD_REPO}/el/${redhat_version} $ARTIFACT_PATH --skip-errors
  fi
done
