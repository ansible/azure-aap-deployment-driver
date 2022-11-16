# Getting Started

This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

## Install dependencies

### Mock server

As of now to run this UI locally we also require python mock server which contains the mock data and API's.

Following commands will install dependencies for the mock server script. These are needed to be run only once.

```sh
pip install -r requirements.txt
```

To start the mock server, start another terminal and run following commands:

```sh
./run.py
```

The server will listen on port 9090 and provide following APIs:

* GET /step - returns list of steps

* POST /execution/id/restart - to mark failed step with given id for a restart

### NodeJS

1. Install nvm from <https://github.com/nvm-sh/nvm>
2. Install required NodeJS with command: `nvm install` (run the command from the /ui folder)
3. Use the required NodeJS: `nvm use`

### NPM modules

Run command `npm i` to install all dependencies required by this project.\
Some will be reported as deprecated and some will report vulnerabilities. This is normal as they are all from development dependencies.

## Available Scripts

In the project directory, you can run:

### `npm start`

Runs the app in the development mode.\
It should open <http://localhost:9999/> in the default browser. If that does not work, try <http://127.0.0.1:9999/> instead.

The page will reload if you make edits.\
You will also see any lint errors in the console.

### `npm test`

Launches the test runner in the interactive watch mode.\
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.
For configuring Jest, the test runner, see [Test Configuration](https://create-react-app.dev/docs/running-tests#configuration).

### `npm run build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

### `npm run eject`

**Note: this is a one-way operation. Once you `eject`, you can’t go back!**

If you aren’t satisfied with the build tool and configuration choices, you can `eject` at any time. This command will remove the single build dependency from your project.

Instead, it will copy all the configuration files and the transitive dependencies (webpack, Babel, ESLint, etc) right into your project so you have full control over them. All of the commands except `eject` will still work, but they will point to the copied scripts so you can tweak them. At this point you’re on your own.

You don’t have to ever use `eject`. The curated feature set is suitable for small and middle deployments, and you shouldn’t feel obligated to use this feature. However we understand that this tool wouldn’t be useful if you couldn’t customize it when you are ready for it.

## Learn More

You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).

To learn React, check out the [React documentation](https://reactjs.org/).
