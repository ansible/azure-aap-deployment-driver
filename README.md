# azure-aap-deployment-driver

## Overview

This repository contains the Ansible on Clouds managed application Deployment
Driver.

The Deployment Driver consists of following components:

- The engine driving deployments of the ARM templates
- The web UI providing the user a way to interact with the deployment engine/server
- Nginx web server and reverse proxy serving the installer web UI and proxy-ing API requests to installer engine

## Development Team

This component is primarily developed by the Ansible Automation Platform on Azure team in Red Hat.

[Github Issues](https://github.com/ansible/azure-aap-deployment-driver/issues) can be used to file tickets for help, bugs, vulnerabilities or other security issues.

Contributions and suggestions are welcome!  Please see below for getting started.

## Development flow

**NOTE:** The following sections are just outlines, more information needs to be added.

### Running the engine locally

In this flow the main (starting) part of the installation is done by creating and deploying a managed application from aap-azurerm repo and rest of the deployment is done locally (with code in this repo).

First, you will need to create a modified container:

1. Modify the container so it does not run the engine that deploys the templates. Just comment out the line `./server &` in the `start.sh` script
2. Deploy the container to the registry. Login to Azure and Azure Container registry and run `make push-image`. Running the make file locally will push the container to registry `aocinstallerdev.azurecr.io` and name space `aoc-${USER}` where the `${USER}` will be your current user name

Next, the code in the `aap-azurerm` repo needs to be pointed to the modified container:

1. In the file `main.bicep` modify parameter `containerRegistry` to point to the `aocinstallerdev.azurecr.io` registry.
2. In the file `modules/containerInstance.bicep` modify parameter `image:` to point to your container.
3. Run the `create.sh` script to create managed app definition and deploy it.

Finally, run the deployment engine locally:

1. Generate templates in the `app-azurerm` repo. After running `./create.sh ...` they will be in `build` and `dist` folders
2. Copy the `templates` folder into folder `installerstore` in the root of this project
3. Create a `server/.env` file from the template and put values and outputs from step 3. above into it. For the value of "MAIN_OUTPUTS", assuming you are deploying via Azure Portal, go to the managed resource group the managed app was deployed to, under Deployments click on "containerGroupDeploy" and then on Template. On the "Parameters" tab, find the "mainOutputs" parameter and grab the JSON string from it's value field. Make sure there are no new lines in that JSON string when putting it to the .env file. The JSON string must start with { and with }.
4. Run the server and it should start deploying.

### Logging into Azure Container registry

- log in to Azure (make sure to pick the right tenant with `--tenant ...` parameter)
- log in to ACR: `az acr login --name REGISTRY_NAME`  (use only the registry name, not its URL)

## SonarQube

Sonar analysis is performed by Github Actions on the code repository for this
project.  Results available at [sonarcloud.io](https://sonarcloud.io/project/overview?id=aoc-aap-test-installer)

## Developer Setup

### Prequisties

- Linux or Mac OS
- [VS Code](https://code.visualstudio.com/download)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [NPM](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)
- [NVM](https://github.com/nvm-sh/nvm)
- flask (`pip install flask` or `pip3 install flask`)
- flask_cors (`pip install flask_cors` or `pip3 install flask_cors`)

### Setup

- Clone the repo locally and open in VS Code
- Open a terminal and navigate to the `ui` folder

### Start the UI
- Run `make` to run the `./Makefile` installation script
    - On the "Watch Usage" prompt, press `q` to continue the script
- Run `npm start` 
    - This will start the front-end and open it in the browser
- Press `login` to continue in your browser
- Press `CTRL + C` in the terminal window to quit the npm start
- Run ` ./run.py &` to standup a mock API request from Azure

### Start the server

- In the terminal, navigate to the root folder of the project
- Run `make` to run the main `./Makefile` installation script
    - Info: The script will build and tag a docker image; you can run `docker images` to validate.
- Run `./build/server`
    - BUG: This will generate a fatal error: `AZURE_SUBSCRIPTION_ID environment variable must be set.`
- Run `cd server` to change to the server directory
- Run `make test` to run a test



