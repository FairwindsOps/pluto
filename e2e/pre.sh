#!/bin/bash

set -e

go build -ldflags "-s -w" -o pluto
docker cp ./ e2e-command-runner:/pluto
