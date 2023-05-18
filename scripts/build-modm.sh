#!/bin/bash

mkdir -p ./bin

# ---------------------------------------------------------------------------
# Microsoft's marketplace offer deployment manager
# NOTES:
#   The version must match what's depended on in ./server/go.mod, otherwise
#   a version mismatch will occur.
#

# this can be set as an incoming flag value
modm_version=v1.0.0

if [ "$MODM_VERSION" != "" ]; then
    modm_version=$MODM_VERSION
fi

repository_name=commercial-marketplace-offer-deploy
modm_build_dir=./build/$repository_name
modm_repository_url=https://github.com/microsoft/$repository_name.git

echo ""
echo "building MODM (Marketplace Offer Deployment Manager)"
echo "----------------------------------------------------"
echo "version:   $modm_version"
echo "source:    $modm_repository_url"
echo "output:    $modm_build_dir"
echo ""

echo "- cleaning build directory"
rm -rf $modm_build_dir/
mkdir -p $modm_build_dir

echo "- cloning repository"
git clone --filter=blob:none --depth=1 -b $modm_version --single-branch --quiet $modm_repository_url $modm_build_dir &> /dev/null
rm -rf $modm_build_dir/.git

# build the binaries then copy output to ./bin
echo "- building binaries..."
cd $modm_build_dir
make build &> /dev/null

echo "- copying binaries to ./bin"
mkdir -p ../../bin
cp ./bin/apiserver ./bin/operator ../../bin
cd ../../

echo ""
echo "done."
echo ""
