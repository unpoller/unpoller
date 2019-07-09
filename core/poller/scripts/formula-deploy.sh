#!/bin/bash -x

# Deploys a new homebrew formula file to golift/homebrew-tap.
# Requires SSH credentials in ssh-agent to work.
# Run by Travis-CI when a new release is created on GitHub.

source .metadata.sh

make ${BINARY}.rb

git config --global user.email "${BINARY}@auto.releaser"
git config --global user.name "${BINARY}-auto-releaser"

rm -rf homebrew_release_repo
git clone git@github.com:${HBREPO}.git homebrew_release_repo

cp ${BINARY}.rb homebrew_release_repo/Formula
pushd homebrew_release_repo
git commit -m "Update ${BINARY} on Release: v${VERSION}-${ITERATION}" Formula/${BINARY}.rb
git push
popd
