#!/bin/bash
set -e

SRC_PATH=$(pwd)
PROJECT_NAME="github.com/experimental-platform/platform-app-manager"

docker run -v "${SRC_PATH}:/go/src/$PROJECT_NAME" -w "/go/src/$PROJECT_NAME" golang:1.4 /bin/bash -c 'go get -d && go build -v'
