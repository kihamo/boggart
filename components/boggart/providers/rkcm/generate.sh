#!/bin/bash
set -x

find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' ! -path './docs/*' -type f -exec rm -f {} +
go run $GOPATH/src/github.com/go-swagger/go-swagger/cmd/swagger/swagger.go generate client -f swagger.yml
find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' ! -path './docs/*' -type f -exec git add {} +
