#!/bin/bash

SOURCE_TEMPLATE="installer.template.json"
TEMPLATES_DIR="templates"

# get the list of modules from the source template
deployModulesArr=( $(jq -r '.resources[].name' ${SOURCE_TEMPLATE}) )

echo "Cleaning ${TEMPLATES_DIR}"
rm -rf ${TEMPLATES_DIR}

for module in "${deployModulesArr[@]}"
do
  echo " -> Module: ${module}"
  mkdir -p "${TEMPLATES_DIR}/${module}"
  # extract parameters
  jq --arg MODULE "${module}" '.resources[] | select(.name==$MODULE) |.properties.parameters' > "${TEMPLATES_DIR}/${module}/${module}.parameters.json" ${SOURCE_TEMPLATE}
  # extract dependencies
  cat "${SOURCE_TEMPLATE}" | \
  jq --arg MODULE "${module}" '.resources[] | select(.name==$MODULE) | .dependsOn[] | split(",")[1] | capture("[^a-zA-Z]+(?<name>[a-zA-Z]+)").name' | \
  jq -s '.' > "${TEMPLATES_DIR}/${module}/${module}.dependencies.json"
  # extract template
  jq --arg MODULE "${module}" '.resources[] | select(.name==$MODULE) |.properties.template' > "${TEMPLATES_DIR}/${module}/${module}.json" ${SOURCE_TEMPLATE}
done

