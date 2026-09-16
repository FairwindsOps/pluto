#!/bin/bash

function version_gt() {
    test "$(printf '%s\n' "$@" | sort -V | head -n 1)" != "$1";
}

echo "Starting release."

command -v circleci >/dev/null 2>&1 || { echo >&2 "I require circleci but it's not installed.  Aborting."; exit 1; }
echo "Found circleci command."

cli_version=$(circleci version | cut -d+ -f1)
required_version=0.1.5705

if version_gt "$required_version" "$cli_version"; then
     echo "This script requires circleci version greater than or equal to 0.1.5705!"
     exit 1
fi

commit=$(git log -n1 --pretty='%h')
tag=$(git describe --exact-match --tags "$commit")

retVal=$?
echo "retVal = $retVal"
if [ $retVal -ne 0 ]; then
    echo "You need to checkout a valid tag for this to work."
    exit $retVal
fi

echo "Release: $commit - $tag"

echo "Validating..."
make orb-validate || { echo 'Orb failed to validate.' ; exit 1; }

echo "Releasing..."
circleci orb publish orb.yml "fairwinds/pluto@${tag:1}"
