# Overview

This folder contains end-to-end tests for deployment driver engine (at this time only web UI actually).

## Requirements

* Node version manager `nvm`
* Go for back-end testing

### Running UI E2E

Requires deployment driver UI running and its URL needs to be passed to the test.

#### 1. Running UI E2E for locally running UI

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

#### 2. Running E2E for UI on a URL

**WARNING:** The tests wont work just yet in this scenario because the deployment driver forces users to log-in with Red Hat SSO and the dialog that does that prevents any other testing/verification until user has logged in. In the near future we will have a way to disable SSO login requirement for testing.

To run the E2E tests against a deployment driver running on another URL, all you need to do is run:

```sh
CYPRESS_DEPLOYMENT_DRIVER_URL=https://somehost.somewhere.com npx cypress run
```
