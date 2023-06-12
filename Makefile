BUILD_DIR := build
INSTALLER_SERVER_DIR := server
INSTALLER_WEBUI_DIR := ui
CONTAINER_REGISTRY_DEFAULT_NAMESPACE ?= aoc-${USER}
DRIVER_RELEASE_TAG ?=$(shell git rev-parse --short HEAD)
IMAGE_NAME ?= installer
IMAGE_TAG ?= latest
MODM_REPOSITORY_NAME := commercial-marketplace-offer-deploy
MODM_BUILD_DIR := build/${MODM_REPOSITORY_NAME}
MODM_REPOSITORY_URL := https://github.com/microsoft/${MODM_REPOSITORY_NAME}.git
MODM_VERSION := v1.3.4
MODM_IMAGE_NAME ?= modm
MODM_IMAGE_TAG ?= ${MODM_VERSION}

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

resolve-registry:
ifndef CONTAINER_REGISTRY_DEFAULT_SERVER
	$(error Environment variable CONTAINER_REGISTRY_DEFAULT_SERVER is not set)
endif
ifndef CONTAINER_REGISTRY_NAMESPACE
CONTAINER_REGISTRY_NAMESPACE := ${CONTAINER_REGISTRY_DEFAULT_NAMESPACE}
endif
ifndef CONTAINER_REGISTRY
CONTAINER_REGISTRY := ${CONTAINER_REGISTRY_DEFAULT_SERVER}/${CONTAINER_REGISTRY_NAMESPACE}
endif

assemble: clean resolve-registry build-server build-web-ui build-modm assemble-deployment-driver assemble-modm

assemble-deployment-driver: clean resolve-registry build-server build-web-ui
	@echo "Building docker image: ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
	docker build -f ./Dockerfile.deployment-driver --build-arg DRIVER_RELEASE_TAG=${DRIVER_RELEASE_TAG} -t ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} .
	docker tag ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} ${CONTAINER_REGISTRY}/${IMAGE_NAME}:latest

assemble-modm: clean resolve-registry build-modm
	@echo "Building docker image: ${CONTAINER_REGISTRY}/${MODM_IMAGE_NAME}:${MODM_IMAGE_TAG}"
	docker build -f ./Dockerfile.modm -t ${CONTAINER_REGISTRY}/${MODM_IMAGE_NAME}:${MODM_IMAGE_TAG} .
	docker tag ${CONTAINER_REGISTRY}/${MODM_IMAGE_NAME}:${MODM_IMAGE_TAG} ${CONTAINER_REGISTRY}/${MODM_IMAGE_NAME}:latest

save-image: assemble
	@echo "Saving docker image: ${IMAGE_NAME}:${IMAGE_TAG} to tar.gz archive..."
	docker image save ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} | gzip --best --stdout > build/${IMAGE_NAME}_${IMAGE_TAG}.tar.gz

push-authenticated: check-credentials login-to-registry push-image logout-from-registry

login-to-registry:
	@echo "Logging in to container registry: ${CONTAINER_REGISTRY}"
	echo $${CONTAINER_REGISTRY_PASSWORD} | docker login --username ${CONTAINER_REGISTRY_USERNAME} --password-stdin ${CONTAINER_REGISTRY}

logout-from-registry:
	docker logout ${CONTAINER_REGISTRY}

push-image: assemble
	@echo "Pushing deployment driver image to registry: ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
	docker push --all-tags ${CONTAINER_REGISTRY}/${IMAGE_NAME}
	@echo "Pushing MODM image to registry: ${CONTAINER_REGISTRY}/${MODM_IMAGE_NAME}:${MODM_IMAGE_TAG}"
	docker push --all-tags ${CONTAINER_REGISTRY}/${MODM_IMAGE_NAME}

build-server:
	@echo "Building installer server"
	make build -C ${INSTALLER_SERVER_DIR}
	cp ${INSTALLER_SERVER_DIR}/build/server ${BUILD_DIR}

build-modm:
	@echo "Building MODM"
	@echo "Cleaning build dir"
	rm -rf ${MODM_BUILD_DIR}
	@echo "Cloning repository"
	git clone --filter=blob:none --depth=1 -b ${MODM_VERSION} --single-branch --quiet -c advice.detachedHead=false ${MODM_REPOSITORY_URL} ${MODM_BUILD_DIR}
	rm -rf ${MODM_BUILD_DIR}/.git
	@echo "Building binaries..."
	cd ${MODM_BUILD_DIR} && make build
	@echo "Copying binaries to ./build"
	cp ${MODM_BUILD_DIR}/bin/apiserver ${MODM_BUILD_DIR}/bin/operator ${BUILD_DIR}
	@echo "Removing MODM build dir"
	rm -rf ${MODM_BUILD_DIR}

build-web-ui:
	@echo "Building installer web UI"
	make build -C ${INSTALLER_WEBUI_DIR}
	cp -ap ${INSTALLER_WEBUI_DIR}/build/. ${BUILD_DIR}/public
