#!/usr/bin/env bash

echo "Updating dependencies"

deps=(
harmony-one/bls
harmony-one/go-sdk
harmony-one/harmony
)

for dep in "${deps[@]}"; do
  echo "Updating dependency ${dep} to latest master"
  go get -u github.com/${dep}@master
done
