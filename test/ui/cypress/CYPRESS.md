# Cypress Testing

Cypress is being used for both end-to-end tests and component tests.

## Prerequisites

1. You have an AAP instance being deployed on Azure.
2. The Deployment Engine UI is accessible.

## Run

1. Initialize your local environment.
   You may need to create a virtualenv and install the needed packages from ui/requirements.txt if you don't already have them.

    ```shell
    # Create a python venv to install dependencies within an venv. 
    # You can skip it if you want to install dependencies at the OS level.
    cd <repo>/test/ui
    python -m venv venv
    source venv/bin/activate    

    # Clean install a project
    npm ci

    # Verify Cypress installation
    npx cypress verify

    # Install the following dependencies to see useful messages when the engine is dead.
    # This is optional, not manadatory.
    cd <repo>/ui
    pip install -r requirements.txt
    ```

2. Configure tests input.

   Option 1: Set environment variables at OS level
   ```shell
   export CYPRESS_baseUrl=<Deployment Engine UI Url>
   export CYPRESS_DEPLOYMENT_ENGINE_UI_PASSWORD=<Admin password to login Deployment Engine UI>
   export CYPRESS_RH_SSO_URL=https://sso.redhat.com
   export CYPRESS_RH_ACCOUNT_USERNAME=<User to login https://sso.redhat.com>
   export CYPRESS_RH_ACCOUNT_PASSWORD=<Password to login https://sso.redhat.com>
   ```
   Option 2: Set environment variables in `test/ui/cypress.env.json
   ```shell
   cd <repo>/test/ui
   
   Refer to the following example cypress.env.json to create your cypress.env.json file. 
   {
    "baseUrl": "Deployment Engine UI Url",
    "DEPLOYMENT_ENGINE_UI_PASSWORD": "Admin password to login Deployment Engine UI",
    "RH_SSO_URL": "https://sso.redhat.com",
    "RH_ACCOUNT_USERNAME": "User to login https://sso.redhat.com",
    "RH_ACCOUNT_PASSWORD": "Password to login https://sso.redhat.com"
   }
   ```

3. Run Cypress tests.

   ```shell
   cd <repo>/test/ui
   
   Option 1: Run tests from all the specs under `cypress/e2e`
   npx cypress run

   Option 2: Launch Cypress UI to run tests from a specific spec under `cypress/e2e`
   npx cypress open
   ```
