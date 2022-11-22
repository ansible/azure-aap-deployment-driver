export SHELL := /usr/bin/bash

BUILD_DIR := build
INSTALLER_SERVER_DIR := server
INSTALLER_WEBUI_DIR := ui
DOCKER_IMAGE_REGISTRY ?= quay.io/aoc
DOCKER_IMAGE_NAME ?= installer
DOCKER_IMAGE_TAG ?= latest

.PHONY: clean assemble build-server build-web-ui

all: assemble

clean:
	rm -rf build
	mkdir -p build/public

.ONESHELL:
assemble: clean build-server build-web-ui
	@echo "Building docker image: ${DOCKER_IMAGE_REGISTRY}/${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}"
	docker rmi ${DOCKER_IMAGE_REGISTRY}/${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}
	docker build -t ${DOCKER_IMAGE_REGISTRY}/${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} .

save-image: assemble
	@echo "Saving docker image: ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} to tar.gz archive..."
	docker image save ${DOCKER_IMAGE_REGISTRY}/${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} | gzip --best --stdout > build/${DOCKER_IMAGE_NAME}_${DOCKER_IMAGE_TAG}.tar.gz

.ONESHELL:
build-server:
	@echo "Building installer server"
	cd ${INSTALLER_SERVER_DIR}
	make build
	cd ..
	cp ${INSTALLER_SERVER_DIR}/build/server ${BUILD_DIR}

.ONESHELL:
build-web-ui:
	@echo "Building installer web UI"
	cd ${INSTALLER_WEBUI_DIR}
	make build
	cd ..
	cp -ap ${INSTALLER_WEBUI_DIR}/build/. ${BUILD_DIR}/public
