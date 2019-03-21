#!/bin/bash

export GOPATH=$(pwd)/../../../../

export CGO_ENABLED=0

export GOOS=linux

export GOARCH=amd64

govendor sync
go build -o bin/magneto gitlab.kingbay-tech.com/engine-lottery/magneto/cmd/magneto

docker build -t weiwen/magneto -f ./build/Dockerfile .

docker-compose -f ./deployments/docker-compose.yml up -d