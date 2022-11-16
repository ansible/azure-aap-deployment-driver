#!/usr/bin/env bash

SOURCE_TEMPLATES_DIR=$1
SOURCE_TEMPLATE="$1/installer.bicep"
INSTALLER_TEMPLATE="$2/installer.template.json"
OUTPUT_TEMPLATES_DIR="$2/templates"

if [ -f "${SOURCE_TEMPLATE}" ]; then
  az bicep build -f "${SOURCE_TEMPLATE}" --outfile "${INSTALLER_TEMPLATE}"
else
  echo "Can not find source template ${SOURCE_TEMPLATE}"
  exit 1
fi

# get the list of modules from the source template
deployModulesArr=( $(jq -r '.resources[].name' ${INSTALLER_TEMPLATE}) )

echo "Cleaning ${OUTPUT_TEMPLATES_DIR}"
rm -rf ${OUTPUT_TEMPLATES_DIR}

for module in "${deployModulesArr[@]}"
do
  echo " -> Module: ${module}"
  mkdir -p "${OUTPUT_TEMPLATES_DIR}/${module}"
  # extract parameters
  jq --arg MODULE "${module}" '.resources[] | select(.name==$MODULE) |.properties.parameters' > "${OUTPUT_TEMPLATES_DIR}/${module}/${module}.parameters.json" "${INSTALLER_TEMPLATE}"
  # extract dependencies
  cat "${INSTALLER_TEMPLATE}" | \
  jq --arg MODULE "${module}" '.resources[] | select(.name==$MODULE) | .dependsOn[] | split(",")[1] | capture("[^a-zA-Z]+(?<name>[a-zA-Z]+)").name' | \
  jq -s '.' > "${OUTPUT_TEMPLATES_DIR}/${module}/${module}.dependencies.json"
  # extract template
  jq --arg MODULE "${module}" '.resources[] | select(.name==$MODULE) |.properties.template' > "${OUTPUT_TEMPLATES_DIR}/${module}/${module}.json" "${INSTALLER_TEMPLATE}"
done

rm -f "${INSTALLER_TEMPLATE}"
