#!/usr/bin/env bash

server_local_development_env () {
    # export / set all environment variables from the local development env defaults in /configs
    echo "Copying environment .env from configs/.env.development.local"
    cp ./configs/.env.development.local ./build/.env

    # Azure environment variables
    if command -v az &> /dev/null
    then
        echo "Azure CLI found. Setting Azure environment variables using 'az account show'."
        export AZURE_TENANT_ID=$(az account show --query tenantId -o tsv)
        export AZURE_SUBSCRIPTION_ID=$(az account show --query id -o tsv)
    fi

    # ensure we have a resource group (we are assuming az login has occurred)
    local $(cat ./configs/.env.development.local | grep 'RESOURCE_GROUP_NAME' | xargs)

    echo "Ensuring resource group '${RESOURCE_GROUP_NAME}'"
    az group create -n $RESOURCE_GROUP_NAME -l eastus2 -o none

    # local sqlite database so it has a db to run
    touch "build/${DB_PATH}"
    echo ""
}

run_server () {
    # setup local environment variables and start the server
    server_local_development_env
    echo -e "Starting Server.\n" 
    pushd ./build
        ./server
    popd
}

run_ui () {
    echo -e "Starting UI.\n" 
    pushd ./ui
        npm start
    popd
}

TARGET=$1

case $TARGET in

  server)
    run_server
    ;;
  ui)
    run_ui
    ;;
  *)
    echo -n "unknown"
    ;;
esac
