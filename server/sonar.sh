#!/usr/bin/env bash


WORKDIR="$PWD/.sonarqube"
SONAR_SCANNER_VER=4.7.0.2747
SCANNER_PATH="${WORKDIR}/sonar-scanner-${SONAR_SCANNER_VER}-linux"
SONAR_PROJECT_TOKEN=***REMOVED***

mkdir -p ${WORKDIR}

if [ ! -d "$SCANNER_PATH" ]; then
  "tools/install_sonarqube.sh" "${WORKDIR}" "${SONAR_SCANNER_VER}"
fi

# Execute unit tests with code coverage
go test ./... -coverprofile=coverage.txt

PATH="${PATH}:${SCANNER_PATH}/bin" sonar-scanner -Dsonar.login=${SONAR_PROJECT_TOKEN}
