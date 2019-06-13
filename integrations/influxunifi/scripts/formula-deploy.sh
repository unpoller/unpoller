#!/bin/bash

# Deploys a new unifi-poller.rb  formula file to golift/homebrew-tap.

make unifi-poller.rb
VERSION=$(grep -E '^\s*version\s*"' unifi-poller.rb | cut -d\" -f 2)

rm -rf homebrew-mugs
git config --global user.email "unifi@auto.releaser"
git config --global user.name "unifi-auto-releaser"
git clone git@github.com:golift/homebrew-mugs.git

cp unifi-poller.rb homebrew-mugs/Formula
pushd homebrew-mugs
echo "Showing diff:"
git diff
git commit -m "Update unifi-poller on Release: v${VERSION}" Formula/unifi-poller.rb
git push
popd
