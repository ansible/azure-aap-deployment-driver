#!/usr/bin/env bash

WORKDIR="$PWD/.sonarqube"
SONAR_SCANNER_VER=4.7.0.2747
SCANNER_PATH="${WORKDIR}/sonar-scanner-${SONAR_SCANNER_VER}-linux"

mkdir -p ${WORKDIR}

if [ ! -d "$SCANNER_PATH" ]; then
  "tools/install_sonarqube.sh" "${WORKDIR}" "${SONAR_SCANNER_VER}"
fi

# Execute unit tests with code coverage
cd server
go test -cover -coverprofile=../coverage.txt -count=1 ./...
cd ..

if [ -z "${SONAR_PROJECT_TOKEN}" ]; then
  echo "Environment variable SONAR_PROJECT_TOKEN not set, will attempt to load from .env file if it exists."
  if [ -f "tools/.env" ]; then
    source ./tools/.env
  else
    echo "File tools/.env does not exist, exitting."
    exit 1
  fi
fi

PATH="${PATH}:${SCANNER_PATH}/bin" sonar-scanner -Dsonar.login=${SONAR_PROJECT_TOKEN}
