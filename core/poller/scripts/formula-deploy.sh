#!/bin/bash

# Deploys a new unifi-poller.rb formula file to golift/homebrew-tap.
# Requires SSH credentials in ssh-agent to work.
# Run by Travis-CI when a new release is created on GitHub.

if [ -z "$VERSION" ]; then
  VERSION=$TRAVIS_TAG
fi
VERSION=$(echo $VERSION|tr -d v)

make unifi-poller.rb VERSION=$VERSION

if [ -z "$VERSION" ]; then
  VERSION=$(grep -E '^\s+url\s+"' unifi-poller.rb | cut -d/ -f7 | cut -d. -f1,2,3)
fi

rm -rf homebrew-mugs
git config --global user.email "unifi@auto.releaser"
git config --global user.name "unifi-auto-releaser"
git clone git@github.com:golift/homebrew-mugs.git

cp unifi-poller.rb homebrew-mugs/Formula
pushd homebrew-mugs/Formula
git commit -m "Update unifi-poller on Release: ${VERSION}" unifi-poller.rb
git push
popd
