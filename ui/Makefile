NVM_SCRIPT ?= ${HOME}/.nvm/nvm.sh

all: test build

# this will print a message if nvm is not installed
${NVM_SCRIPT}:
	@echo NVM was not found, please provide the location of the script via env variable NVM_SCRIPT.
	false

.PHONY: nvm
nvm: ${NVM_SCRIPT} # will run provided command in environment with nvm
	. ${NVM_SCRIPT} --no-use && nvm install && npm i && ${CMD}

.PHONY: build
build: ${NVM_SCRIPT}
	make nvm CMD="npm run build"

.PHONY: test
test: ${NVM_SCRIPT}
	make nvm CMD="npm run test"

.PHONY: audit-prod-only
audit-prod-only: ${NVM_SCRIPT}
	make nvm CMD="npm audit --omit dev"
