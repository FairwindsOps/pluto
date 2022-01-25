#!/bin/bash
set -e

if [[ -z "${PLUTO_FILE}" ]]; then
    echo "Error: requires a file argument"
    exit 1
fi

if [[ "${PLUTO_IGNORE_DEPRECATIONS}" = true ]]; then
    PLUTO_ARGS="$PLUTO_ARGS --ignore-deprecations"
fi

if [[ "${PLUTO_IGNORE_REMOVALS}" = true ]]; then
    PLUTO_ARGS="$PLUTO_ARGS --ignore-removals"
fi

if [[ -n "${PLUTO_TARGET_VERSIONS}" ]]; then
    PLUTO_ARGS="$PLUTO_ARGS --target-versions k8s=${PLUTO_TARGET_VERSIONS}"
fi

export PLUTO_ARGS
export PLUTO_FILE_PATH

pluto detect $PLUTO_FILE $PLUTO_ARGS
