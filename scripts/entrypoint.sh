#!/bin/bash -l

ARGS="--driver ${INPUT_DRIVER} --region ${INPUT_REGION} --bucket ${INPUT_BUCKET}"
if [ "$INPUT_EXCLUDE" ]; then
  ARGS+=" --exclude ${INPUT_EXCLUDE}"
fi
if [ "$INPUT_SAVEROOT" ]; then
  ARGS+=" --save_root ${INPUT_SAVEROOT}"
fi
echo "${ARGS} ${INPUT_DIST}"

# shellcheck disable=SC2086
uptoc ${ARGS} ${INPUT_DIST}
