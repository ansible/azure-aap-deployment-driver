# aoc-azure-aap-installer

## Overview

This repository contains the Ansible on Clouds managed application installer.

Installer consists of following components:

- Installer engine driving deployments of the ARM templates
- Installer web UI providing user way of interacting with installer
- Nginx web server and reverse proxy serving the installer web UI and proxy-ing API requests to installer engine

## Git Submodules

This repository includes git submodules to reference other repositories:

- *aap-azurerm* - the Azure Bicep files that are converted to Azure ARM templates used by the AAP installer

### Working with submodules

#### Updating

After cloning the repository, run `git submodule update --remote`

#### Status

To see simple status showing latest commit on each submodule: `git submodule status`

To see more about the commits: `git submodule summary` or `git submodule summary --files`
