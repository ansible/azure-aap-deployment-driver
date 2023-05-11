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

run_docker() {
   ./scripts/build.sh docker
   start_ngrok_background
   docker compose -f ./tools/docker-compose.yml up  
}

function start_ngrok_background() {
  # start up ngrok and get address
  ngrok http 8080 > /dev/null &
  ngrok_start_result=$?
  export NGROK_ID=$!

  if [ $ngrok_start_result -gt 0 ]; then
    echo "NGROK failed to start."
    echo "exiting."
    exit 1
  fi

  echo "NGROK started: $NGROK_ID"
  sleep 2
  export MODM_PUBLIC_BASE_URL=$(curl -s localhost:4040/api/tunnels | jq '.tunnels[0].public_url' -r)
  echo "NGROK URL:  $MODM_PUBLIC_BASE_URL"
}


TARGET=$1

case $TARGET in

  server)
    run_server
    ;;
  ui)
    run_ui
    ;;
  docker)
    run_docker
    ;;
  *)
    echo -n "unknown"
    ;;
esac
