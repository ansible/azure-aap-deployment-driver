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
  echo "Unexpected value in MODM_ROLE environment variable. Expected either 'apiserver' or 'operator', got: ${MODM_ROLE}"
  exit 1
fi

# start the executable in background so its PID can be stored in a file
echo "Starting ${EXECUTABLE}..."
${EXECUTABLE} &
MODM_PID=$!

# store PID in a file that's not in persistent volume and wait for the process to end
echo ${MODM_PID} > /tmp/modm_process_pid
wait -n ${MODM_PID}
