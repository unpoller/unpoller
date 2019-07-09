#!/bin/bash -x

# Deploys a new homebrew formula file to golift/homebrew-tap.
# Requires SSH credentials in ssh-agent to work.
# Run by Travis-CI when a new release is created on GitHub.

source .metadata.sh

make ${BINARY}.rb

git config --global user.email "${BINARY}@auto.releaser"
git config --global user.name "${BINARY}-auto-releaser"

rm -rf homebrew-mugs
git clone git@github.com:golift/homebrew-mugs.git

cp ${BINARY}.rb homebrew-mugs/Formula
pushd homebrew-mugs
git commit -m "Update ${BINARY} on Release: v${VERSION}-${ITERATION}" Formula/${BINARY}.rb
git push
popd
