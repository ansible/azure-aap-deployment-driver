#!/usr/bin/env bash

# Verify that required env variables are set, only checking for variables needed by this script
if [ -z ${MODM_ROLE} ]; then
  echo "Environment variable MODM_ROLE not set or is empty."
  exit 1
fi

if [ "${MODM_ROLE}" == "operator" ]; then
  echo "Role specified to run: operator"
  EXECUTABLE="./operator"
elif [ "${MODM_ROLE}" == "apiserver" ]; then
  echo "Role specified to run: apiserver"
  EXECUTABLE="./apiserver"
else
  echo "Role specified does not match neither 'apiserver' nor 'operator'."
  exit 1
fi

echo "Starting ${EXECUTABLE}..."
${EXECUTABLE}
