#!/bin/bash
set -e

SRC_PATH=$(pwd)

docker run -v ${SRC_PATH}:/usr/src/app-manager -w /usr/src/app-manager golang:1.4 /bin/bash -c 'go get -d && go build -v'
