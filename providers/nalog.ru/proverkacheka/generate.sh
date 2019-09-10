#!/bin/bash
set -x

find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' -type f -exec rm -f {} +
go run ../../../vendor/github.com/go-swagger/go-swagger/cmd/swagger/swagger.go generate client -f swagger.yml
find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' -type f -exec git add {} +
