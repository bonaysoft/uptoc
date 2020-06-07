#!/bin/sh -l

ARGS="--driver ${INPUT_DRIVER} --region ${INPUT_REGION} --bucket ${INPUT_BUCKET}"
if [ "$INPUT_EXCLUDE" ]; then
  ARGS+=" --exclude $INPUT_EXCLUDE"
fi

# shellcheck disable=SC2086
uptoc ${ARGS} ${INPUT_DIST}
