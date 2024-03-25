# Overview

This folder contains end-to-end tests for deployment driver engine (at this time only web UI actually).

## Requirements

* Node version manager `nvm`
* Go for back-end testing

### Running UI E2E

Requires deployment driver UI running and its URL needs to be passed to the test.

#### 1. Running UI E2E for locally running UI with development backend

The web UI te needs to have a back-end running for fetching the data. Easiest is to use the python script in the UI folder.

You may need to create a virtualenv and install the needed packages from ui/requirements.txt if you don't already have them.

```sh
cd ui
virtualenv venv
source venv/bin/activate
pip install -r requirements.txt
```

In one terminal, run following commands to get a fake back-end started:

```sh
cd ui
./run.py
```

In another terminal, run the following commands to get the web ui development web server started:

```sh
cd ui
nvm use
npm i
BROWSER=none npm start
```

And finally, in the another terminal, run the following commands to get cypress run the tests:

```sh
cd test/ui
nvm use
npm i
npx cypress run
```

The last command above will run the tests in a headless mode with "electron" browser. To run it with cypress UI where you can choose your browser, use command:

```sh
npx cypress open --e2e
```

#### 2. Running UI E2E with an AAP instance on Azure

##### Prerequisites

1. You have an AAP instance is being deployed on Azure.
2. The Deployment Engine UI is accessible.

##### Run

1. You may need to create a virtualenv and install the needed packages if you don't want to install them at OS level.

```shell
cd <repo>/test/ui

# Create a virtualenv and activate it
virtualenv venv
source venv/bin/activate

# Clean install the project
npm ci
```

2. You are encouraged to install the packages from `<repo>/ui/requirements.txt` if you want to see useful messages when the deployment engine is dead.

```shell
cd <repo>/ui
pip install -r requirements.txt
```

3. Configure environment variables for Cypress automation to run.

   Option 1: Set environment variables at OS level
   ```shell
   export CYPRESS_DEPLOYMENT_DRIVER_URL=<Deployment Engine UI Url>
   export CYPRESS_DEPLOYMENT_ENGINE_UI_PASSWORD=<Admin password to login Deployment Engine UI>
   export CYPRESS_RH_SSO_URL=https://sso.redhat.com
   export CYPRESS_RH_ACCOUNT_USERNAME=<User to login https://sso.redhat.com>
   export CYPRESS_RH_ACCOUNT_PASSWORD=<Password to login https://sso.redhat.com>
   ```
   Option 2: Set environment variables in `test/ui/cypress.env.json`

   ```shell
   cd <repo>/test/ui
   
   Refer to the following example to create your cypress.env.json file if you don't have it.
   {
    "DEPLOYMENT_DRIVER_URL": "Deployment Engine UI Url",
    "DEPLOYMENT_ENGINE_UI_PASSWORD": "Admin password to login Deployment Engine UI",
    "RH_SSO_URL": "https://sso.redhat.com",
    "RH_ACCOUNT_USERNAME": "User to login https://sso.redhat.com",
    "RH_ACCOUNT_PASSWORD": "Password to login https://sso.redhat.com"
   }
   ```

4. Run Cypress tests.

   ```shell
   cd <repo>/test/ui
   
   Option 1: Run tests from all the specs under `cypress/e2e`
   npx cypress run

   Option 2: Launch Cypress UI to run tests from a specific spec under `cypress/e2e`
   npx cypress open
   ```
