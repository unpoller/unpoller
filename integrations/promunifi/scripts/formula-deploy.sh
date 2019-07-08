#!/bin/bash -x

# Deploys a new homebrew formula file to golift/homebrew-tap.
# Requires SSH credentials in ssh-agent to work.
# Run by Travis-CI when a new release is created on GitHub.
APP=unpacker-poller

if [ -z "$VERSION" ]; then
  VERSION=$TRAVIS_TAG
fi
VERSION=$(echo $VERSION|tr -d v)

make ${APP}.rb VERSION=$VERSION

if [ -z "$VERSION" ]; then
  VERSION=$(grep -E '^\s+url\s+"' ${APP}.rb | cut -d/ -f7 | cut -d. -f1,2,3)
fi

rm -rf homebrew-mugs
git config --global user.email "${APP}@auto.releaser"
git config --global user.name "${APP}-auto-releaser"
git clone git@github.com:golift/homebrew-mugs.git

cp ${APP}.rb homebrew-mugs/Formula
pushd homebrew-mugs
git commit -m "Update ${APP} on Release: ${VERSION}" Formula/${APP}.rb
git push
popd
