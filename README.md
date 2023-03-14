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

Contributions and suggestions are welcome!

### Running in VSCode Dev Container

For [VS Code](https://code.visualstudio.com/) users, there is the ability to run your local development environment inside a development container.  This allows you to ensure that you have the needed prerequisites and avoid any portability issues. 

Working with VS Code development containers requires you to follow the installation instructions at [https://code.visualstudio.com/docs/devcontainers/containers#_installation](https://code.visualstudio.com/docs/devcontainers/containers#_installation).  Once you have completed the installation istructions, you can open
this cloned repo's folder in VS Code, or clone the repo to a development container.  

For more details on working with development containers, please see [https://code.visualstudio.com/docs/devcontainers/containers](https://code.visualstudio.com/docs/devcontainers/containers).

## SonarQube

Sonar analysis is performed by Github Actions on the code repository for this
project.  Results available at [sonarcloud.io](https://sonarcloud.io/project/overview?id=aoc-aap-test-installer)

## Developer setup

Instructions to clone the repo and run the UI + server on your development machine are [here](./docs/developer-setup.md).