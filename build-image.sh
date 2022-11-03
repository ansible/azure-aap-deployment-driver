#!/usr/bin/env bash
#
echo
echo "Rebuilding deployment template"
echo "========================"
az bicep build -f ../installer.bicep --outfile installer.template.json
./prepare-install-templates.sh

echo
echo "Building server"
echo "========================"
cd server
# TZ: Commented out because it was failing.
# Also, calling 'go build main.go' did not leave main binary in place but calling 'go build' did, although named 'server'
#go build -ldflags="-extldflags=-static" main.go
go build
cd ..

echo
echo "Removing old docker image tag"
echo "========================"
docker rmi quay.io/aoc/installer:latest

echo
echo "Building docker image"
echo "========================"
docker build -t quay.io/aoc/installer:latest .

echo
echo "Pushing new docker image"
echo "========================"
docker push quay.io/aoc/installer:latest
