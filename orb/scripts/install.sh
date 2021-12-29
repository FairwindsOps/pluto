#!/bin/bash

if [[ -z "${VERSION}" ]]; then
  VERSION=latest
fi

curl -s https://api.github.com/repos/FairwindsOps/pluto/releases/${VERSION} \
| grep "browser_download_url.*linux_amd64.tar.gz" \
| cut -d '"' -f 4 \
| wget -qi -

tarball="$(find . -name "*linux_amd64.tar.gz")"
tar -xzf $tarball

sudo mv pluto /bin

location="$(which pluto)"
echo "Pluto binary location: $location"

version="$(pluto version)"
echo "Pluto binary version: $version"
