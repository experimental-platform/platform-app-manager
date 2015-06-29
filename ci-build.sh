#!/bin/bash
# THIS ONLY WORK IN OUR CI!

docker run --rm -v $(pwd):/usr/src/app-manager -w /usr/src/app-manager golang:1.4 /bin/bash -c 'go get -d && go build -v'
