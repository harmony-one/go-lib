#!/usr/bin/env bash

echo "Updating dependencies"

deps=(
harmony-one/bls
harmony-one/go-sdk
harmony-one/harmony
)

for dep in "${deps[@]}"; do
  commit=$(curl -Ls https://api.github.com/repos/${dep}/commits/master | jq '.sha' | tr -d '"')
  echo "Updating dependency ${dep} to latest master commit ${commit}"
  go get -u github.com/${dep}@${commit}
done

echo "Tidying up go.mod ..."
go mod tidy
