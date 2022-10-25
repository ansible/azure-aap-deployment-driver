# export PATH := node_modules/.bin:$(PATH)
export SHELL := /usr/bin/env bash

.PHONY: check-nvm
check-nvm: # will error if nvm is not installed
	@bash -l -c '[ ! -s "${HOME}/.nvm/nvm.sh" ] && { echo "You need nvm!"; exit 1; }'

.PHONY: nvm
nvm: # will run provided command in environment with nvm
	. ${HOME}/.nvm/nvm.sh && nvm use && $(CMD)

.PHONY: build
build:
	make nvm CMD="npm i && npm run build"
