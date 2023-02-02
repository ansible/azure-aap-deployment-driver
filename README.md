# aoc-azure-aap-installer

## Overview

This repository contains the Ansible on Clouds managed application installer.

Installer consists of following components:

- Installer engine driving deployments of the ARM templates
- Installer web UI providing user way of interacting with installer
- Nginx web server and reverse proxy serving the installer web UI and proxy-ing API requests to installer engine

## Development flow

**NOTE:** Following sections are just outlines, more information needs to be added.

### Running engine locally

In this flow the main (starting) part of the installation is done by creating and deploying managed application from aap-azurerm repo and rest of the deployment is done locally (with code in this repo).

First, you will need to create a modified container:

1. Modify the container so it does not run the engine that deploys the templates. Just comment out the line `./server &` in the `start.sh` script
2. Deploy the container to the registry. Login to Azure and Azure Container registry and run `make push-image`. Running the make file locally will push the container to registry `aocinstallerdev.azurecr.io` and name space `aoc-${USER}` where the `${USER}` will be your current user name

Next, the code in the `aap-azurerm` repo needs to be pointed to the modified container:

1. In the file `main.bicep` modify parameter `containerRegistry` to point to the `aocinstallerdev.azurecr.io` registry.
2. In the file `modules/containerInstance.bicep` modify parameter `image:` to point to your container.
3. Run the `create.sh` script to create managed app definition and deploy it.

Finally, run deployment engine locally:

1. Generate templates in the `app-azurerm` repo. After running `./create.sh ...` they will be in `build` and `dist` folders
2. Copy the `templates` folder into folder `installerstore` in the root of this project
3. Create a `server/.env` file from the template and put values and outputs from step 3. above into it. For the value of "MAIN_OUTPUTS", assuming you are deploying via Azure Portal, go to the managed resource group the managed app was deployed to, under Deployments click on "containerGroupDeploy" and then on Template. On the "Parameters" tab, find the "mainOutputs" parameter and grab the JSON string from it's value field. Make sure there are no new lines in that JSON string when putting it to the .env file. The JSON string must start with { and with }.
4. Run the server and it should start deploying.

### Logging into Azure Container registry

- login to Azure (make sure to pick the right tenant with `--tenant ...` parameter)
- login to ACR: `az acr login --name REGISTRY_NAME`  (use only the registry name, not its URL)
