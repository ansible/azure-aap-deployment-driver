# aoc-azure-aap-installer

## Overview

TBA

## Mock server
As of now to run this UI locally we also require python mock server which contains the mock data and API's.

### Installing prereqs for python mock server

Following commands will install dependencies for the mock server script. These are needed to be run only once.
```sh
cd ui
pip install -r requirements.txt
```

### Start mock server

To start the mock server, start another terminal and run following commands:

```sh
cd ui
./run.py
```

The server will listen on port 9090 and provide following APIs:

* GET /step - returns list of steps

* POST /execution/id/restart - to mark failed step with given id for a restart


## UI

The installer UI is in `/ui` folder.