# Developer Setup

## Prequisties

- Linux or MacOS
- [VS Code](https://code.visualstudio.com/download)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [NPM](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)

## Setup

- Clone the repo locally and open in VS Code
- Open a terminal and navigate to the `ui` folder

## Start the UI
- Run `make` to run the `./Makefile` installation script
    - On the "Watch Usage" prompt, press `q` to continue the script
- Run `npm start` 
    - This will start the front-end and open it in the browser
- Press `login` to continue in your browser
- Press `CTRL + C` in the terminal window to quit the npm start
- Run ` ./run.py &` to standup a mock API request from Azure

## Start the server

- In the terminal, navigate to the root folder of the project
- Run `make` to run the main `./Makefile` installation script
    - Info: The script will build and tag a docker image; you can run `docker images` to validate.
- Run `./build/server`
    - BUG: This will generate a fatal error: `AZURE_SUBSCRIPTION_ID environment variable must be set.`
- Run `cd server` to change to the server directory
- Run `make test` to run a test



