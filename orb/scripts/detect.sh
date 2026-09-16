#!/bin/bash
set -e

if [[ -z "${PLUTO_FILE}" ]]; then
    echo "Error: requires a file argument"
    exit 1
fi

pluto detect "$PLUTO_FILE"
