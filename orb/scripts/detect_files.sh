#!/bin/bash
set -e

if [[ -n "${PLUTO_DIRECTORY}" ]]; then
    PLUTO_ARGS="$PLUTO_ARGS --directory ${PLUTO_DIRECTORY}"
fi

export PLUTO_ARGS

pluto detect-files "$PLUTO_ARGS"
