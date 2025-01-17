#!/usr/bin/env bash

gitroot=$(git rev-parse --show-toplevel)

pushd "$gitroot" >/dev/null || exit 1

# list all directories that contain .go files except for vendor directories
dirs=$(grep --include="*.go" --exclude-dir="vendor" --exclude-dir="plugins" -r "package main" -l | xargs dirname | sort -u)

# use parallel to execute "cd {} && go build" in for each directory in $dirs
parallel -j4 --no-notice --bar --eta "cd {} && go build" ::: "$dirs"

popd >/dev/null || exit 1
