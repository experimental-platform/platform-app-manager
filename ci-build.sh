#!/bin/bash
set -e

SRC_PATH=$(pwd)
PROJECT_NAME="github.com/experimental-platform/platform-app-manager"

./install-glide-v1.sh && glide up
docker run -v "${SRC_PATH}:/go/src/$PROJECT_NAME" -w "/go/src/$PROJECT_NAME" -e GO15VENDOREXPERIMENT=1 golang:1.5 /bin/bash -c "go build -v"
