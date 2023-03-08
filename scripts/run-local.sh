#!/usr/bin/env bash

run_server () {
    # local sqlite database so it has a db to run
    touch "build/${DB_PATH}"
    echo -e "\nStarting Server.\n" 
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
