#!/bin/sh -l

ARGS="--driver ${INPUT_DRIVER} --region ${INPUT_REGION} --bucket ${INPUT_BUCKET}"
if [ "$INPUT_EXCLUDE" ]; then
  ARGS+=" --exclude $INPUT_EXCLUDE"
fi
ARGS+=" $INPUT_DIST"

# shellcheck disable=SC2086
uptoc ${ARGS}
