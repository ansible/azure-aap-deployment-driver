#!/usr/bin/env bash

if [ $# -ne 2 ]; then echo "Working directory required as first arg, scanner version as second"; fi

WORKDIR=$1
SONAR_SCANNER_VER=$2
SONAR_SCANNER_URL="https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-${SONAR_SCANNER_VER}-linux.zip"
SONAR_SERVER_URL="https://sonarcloud.io"

# Fetch scanner
wget -qO "$WORKDIR/scanner.zip" "$SONAR_SCANNER_URL"

# Unpack scanner
unzip -q "$WORKDIR/scanner.zip" -d "$WORKDIR"

# Inject URI for sonarqube host to config file
echo "sonar.host.url=${SONAR_SERVER_URL}" >> "${WORKDIR}/sonar-scanner-${SONAR_SCANNER_VER}-linux/conf/sonar-scanner.properties"

# Done!
