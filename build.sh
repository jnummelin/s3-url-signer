#!/bin/sh

# To easily cross-compile binaries
go get github.com/mitchellh/gox

VERSION=$(git describe --all --exact-match `git rev-parse HEAD` | grep tags | sed 's/tags\///')
GIT_COMMIT=$(git rev-list -1 HEAD)

gox -output="output/{{.Dir}}_{{.OS}}_{{.Arch}}" \
  -osarch="darwin/amd64 linux/amd64 linux/arm64 linux/arm win/amd64" \
  -ldflags "-s -w  -X github.com/jnummelin/s3-url-signer/version.GitCommit=${GIT_COMMIT} -X github.com/jnummelin/s3-url-signer/version.Version=${VERSION}"
