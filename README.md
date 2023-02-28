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

## SonarQube

Sonar analysis is performed by Github Actions on the code repository for this
project.  Results available at [sonarcloud.io](https://sonarcloud.io/project/overview?id=aoc-aap-test-installer)