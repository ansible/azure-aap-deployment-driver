#!/bin/bash


echo "Server starting."
# Start the api server
/server /dev/fd/1 2>&1 &

echo "MODM - API Server starting."
# Start the api server
/modm/apiserver /dev/fd/1 2>&1 &

echo "MODM - Operator starting."
# Start the operator server
/modm/operator /dev/fd/1 2>&1 &

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?