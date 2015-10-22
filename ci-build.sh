#!/bin/bash
set -e

SRC_PATH=$(pwd)
PROJECT_NAME="github.com/experimental-platform/platform-app-manager"

docker run -v "${SRC_PATH}:/go/src/$PROJECT_NAME" -w "/go/src/$PROJECT_NAME" -e GO15VENDOREXPERIMENT=1 golang:1.5 /bin/bash -c "./install-glide-v1.sh && glide up && go build -v"
