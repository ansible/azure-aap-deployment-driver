#!/usr/bin/env bash

CONFIG_FILE="${1}/sonar-project.properties"
WORKDIR="$PWD/.sonarqube"
SONAR_SCANNER_VER=4.7.0.2747
SCANNER_PATH="${WORKDIR}/sonar-scanner-${SONAR_SCANNER_VER}-linux"

if [ -z "${1}" ]; then
  echo "Missing parameter specifying directory containing sonar-project.properties file."
  exit 1
fi

if [ ! -f "${CONFIG_FILE}" ]; then
  echo "Can not find sonar-project.properties file: ${CONFIG_FILE}"
  exit 1
fi

mkdir -p "${WORKDIR}"

if [ ! -d "$SCANNER_PATH" ]; then
  echo "Scanner not found, will install it now..."
  "tools/install_sonarqube.sh" "${WORKDIR}" "${SONAR_SCANNER_VER}"
fi

if [ -z "${SONAR_PROJECT_TOKEN}" ]; then
  echo "Environment variable SONAR_PROJECT_TOKEN not set, will attempt to load from .env file if it exists."
  if [ -f "tools/.env" ]; then
    echo "Loading tools/.env..."
    source ./tools/.env
  else
    echo "File tools/.env does not exist, exitting."
    exit 1
  fi
fi

PATH="${PATH}:${SCANNER_PATH}/bin" sonar-scanner "-Dsonar.login=${SONAR_PROJECT_TOKEN}" "-Dproject.settings=${CONFIG_FILE}"
