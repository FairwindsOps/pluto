#!/bin/bash
if [[ -n "${PLUTO_DIRECTORY}" ]]; then
    PLUTO_ARGS="$PLUTO_ARGS --directory ${PLUTO_DIRECTORY}"
fi

if [[ "${PLUTO_IGNORE_DEPRECATIONS}" ]]; then
    PLUTO_ARGS="$PLUTO_ARGS --ignore-deprecations"
fi

if [[ "${PLUTO_IGNORE_REMOVALS}" ]]; then
    PLUTO_ARGS="$PLUTO_ARGS --ignore-removals"
fi

if [[ -n "${PLUTO_TARGET_VERSIONS}" ]]; then
    PLUTO_ARGS="$PLUTO_ARGS --target-versions k8s=${PLUTO_TARGET_VERSIONS}"
fi

export PLUTO_ARGS

pluto detect-files $PLUTO_ARGS
