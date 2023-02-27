BUILD_DIR := build
INSTALLER_SERVER_DIR := server
INSTALLER_WEBUI_DIR := ui
CONTAINER_REGISTRY_DEFAULT_SERVER ?= aocinstallerdev.azurecr.io
CONTAINER_REGISTRY_DEFAULT_NAMESPACE ?= aoc-${USER}
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

resolve-registry:
ifndef CONTAINER_REGISTRY_NAMESPACE
CONTAINER_REGISTRY_NAMESPACE := ${CONTAINER_REGISTRY_DEFAULT_NAMESPACE}
endif
ifndef CONTAINER_REGISTRY
CONTAINER_REGISTRY := ${CONTAINER_REGISTRY_DEFAULT_SERVER}/${CONTAINER_REGISTRY_NAMESPACE}
endif

assemble: clean resolve-registry build-server build-web-ui
	@echo "Building docker image: ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
	docker rmi ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}
	docker build -t ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} .
	docker tag ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} ${CONTAINER_REGISTRY}/${IMAGE_NAME}:latest

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
	@echo "Pushing image to registry: ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
	docker push --all-tags ${CONTAINER_REGISTRY}/${IMAGE_NAME}

build-server:
	@echo "Building installer server"
	make build -C ${INSTALLER_SERVER_DIR}
	cp ${INSTALLER_SERVER_DIR}/build/server ${BUILD_DIR}

build-web-ui:
	@echo "Building installer web UI"
	make build -C ${INSTALLER_WEBUI_DIR}
	cp -ap ${INSTALLER_WEBUI_DIR}/build/. ${BUILD_DIR}/public
run:
	@echo "\n*** Starting mock API service ***\n"
	./ui/run.py &

	@echo "\n*** Starting the server ***\n"
	./build/server

	@echo "\n*** Starting the UI and launching browser... ***\n"
	cd /ui
	npm start
