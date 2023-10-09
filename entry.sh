#!/usr/bin/env bash

# Start up and tee console output to log
./start.sh 2>&1 | tee -a ${BASE_PATH}/console.log
