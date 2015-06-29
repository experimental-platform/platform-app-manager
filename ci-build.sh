#!/bin/bash
# THIS ONLY WORK IN OUR CI!

echo "Current path: " $(pwd)

docker run --rm -v /data/jenkins/jobs/${JOB_NAME}/workspace:/usr/src/app-manager -w /usr/src/app-manager golang:1.4 /bin/bash -c 'go get -d && go build -v'
