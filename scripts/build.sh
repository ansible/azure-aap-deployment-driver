#!/bin/bash


function docker_build() {
    docker build . -t aap-deployment-driver -f ./package/Dockerfile
}


function local_build() {
    ./build-modm.sh

    echo "- Building deployment driver"
    echo "  - server"
    make build-server BUILD_DIR=./bin &> /dev/null
    echo "  - web ui"
    make build-web-ui BUILD_DIR=./bin &> /dev/null
}


TARGET=$1

case $TARGET in

  docker)
    docker_build
    ;;
  *)
    echo -n "local"
    local_build
    ;;
esac
