
ARG MODM_VERSION=v0.2.2
ARG NODE_VERSION=18.16.0

FROM node:18 as uibuilder
WORKDIR /builder
# ---------------------------------------------------------------
# Build the web

# UI / frontend - will be hosted via proxy server
RUN mkdir ./ui && mkdir -p ./build/public
COPY ./ui/ ./ui
RUN cd ./ui && npm install
RUN cd ./ui && npm run build
RUN cp -r ./ui/build ./build/public


# Use the offical golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.19-buster as builder
WORKDIR /builder


# git install
RUN apt update && apt-get install git
RUN mkdir ./scripts && mkdir ./bin


COPY ./package/entrypoint.sh ./bin
COPY ./scripts/build-modm.sh ./scripts

# ---------------------------------------------------------------
# Build the server
COPY ./server/go.sum ./
COPY ./server/go.mod ./
RUN go mod download

COPY ./server/ ./

RUN go build -o ./bin/ ./

# ---------------------------------------------------------------
# Build MODM using script which clones a specific version and outs to bin folder
# this script will copy to bin

RUN export MODM_VERSION=${MODM_VERSION}
RUN ./scripts/build-modm.sh


# /
# ---------------------------------------------------------------
# Final image

FROM debian:buster-slim

USER root

# install the azure CLI so that CLI context on the host, locally, can be used when making calls to azure
RUN apt-get clean && apt-get update
RUN apt install curl -y
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash

# Copy the binary to the production image from the builder stage.
COPY --from=uibuilder /builder/build/public /public

COPY --from=builder /builder/bin/server /server
COPY --from=builder /builder/bin/operator /operator
COPY --from=builder /builder/bin/apiserver /apiserver
COPY --from=builder /builder/bin/entrypoint.sh /entrypoint.sh

CMD ["/entrypoint.sh"]