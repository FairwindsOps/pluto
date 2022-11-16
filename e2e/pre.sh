#!/bin/bash

set -e

go get github.com/markbates/pkger/cmd/pkger
make build
docker cp ./ e2e-command-runner:/pluto
