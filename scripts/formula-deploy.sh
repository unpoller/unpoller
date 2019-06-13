#!/bin/bash

# Deploys a new unifi-poller.rb  formula file to golift/homebrew-tap.

make unifi-poller.rb
VERSION=$(grep -E '^\s*version\s*"' unifi-poller.rb | cut -d\" -f 2)

rm -rf homebrew-repo
git clone https://$GITHUB_API_KEY@github.com/golift/homebrew-repo.git

cp unifi-poller.rb homebrew-repo/Formula
pushd homebrew-repo
echo "Showing diff:"
git diff
git config user.name "unifi-poller-bot"
git config user.email "unifi@poller.bot"
git commit -m "Update unifi-poller on Release: v${VERSION}" Formula/unifi-poller.rb
git push
popd
