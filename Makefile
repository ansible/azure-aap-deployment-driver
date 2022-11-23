export SHELL := /usr/bin/bash

BUILD_DIR := build
INSTALLER_SERVER_DIR := server
INSTALLER_WEBUI_DIR := ui
CONTAINER_REGISTRY ?= quay.io/aoc
IMAGE_NAME ?= installer
IMAGE_TAG ?= latest

.PHONY: clean assemble save-image push-image check-credentials build-server build-web-ui

all: assemble

clean:
	rm -rf build
	mkdir -p build/public

check-credentials:
ifndef CONTAINER_REGISTRY_USERNAME
	$(error Environment variable CONTAINER_REGISTRY_USERNAME is not set)
endif
ifndef CONTAINER_REGISTRY_PASSWORD
	$(error Environment variable CONTAINER_REGISTRY_PASSWORD is not set)
endif

.ONESHELL:
assemble: clean build-server build-web-ui
	@echo "Building docker image: ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
	docker rmi ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}
	docker build -t ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} .
	docker tag ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} ${CONTAINER_REGISTRY}/${IMAGE_NAME}:latest

save-image: assemble
	@echo "Saving docker image: ${IMAGE_NAME}:${IMAGE_TAG} to tar.gz archive..."
	docker image save ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} | gzip --best --stdout > build/${IMAGE_NAME}_${IMAGE_TAG}.tar.gz

.ONESHELL:
push-image: check-credentials assemble
	@echo "Pushing image to registry: ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
	echo $${CONTAINER_REGISTRY_PASSWORD} | docker login --username ${CONTAINER_REGISTRY_USERNAME} --password-stdin ${CONTAINER_REGISTRY}
	docker push --all-tags ${CONTAINER_REGISTRY}/${IMAGE_NAME}
	docker logout ${CONTAINER_REGISTRY}

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
