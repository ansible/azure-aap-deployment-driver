# Developer Setup

## Prequisties

- Linux or MacOS
- [VS Code](https://code.visualstudio.com/download)
- [Docker](https://www.docker.com/)
- [Node.js](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)
- [Python](https://www.python.org/)

## Setup

- Clone the repo locally and open in VS Code
- In the terminal, navigate to the root directory of the project
- Run `make` to run the `./Makefile` installation script and install dependencies

## Start the server

Note: running the server actually runs the deployment engine and it requires several environment variables in order to connect to Azure.

- In the terminal, navigate to the `server` directory
- Run `make test`

## Start the UI

- If you are not running the server and want to use the mock API server:
    - In the terminal, navigate to the `ui` folder of the project
    - Run ` ./run.py &` to standup a mock API server
- Run `npm start` 
    - This will start the front-end and open it in the browser
- Press `login` to continue in your browser
    - Note: In development mode the password can be any 12 characters





