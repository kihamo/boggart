#!/bin/bash
set -x

# GOROOT="/usr/local/Cellar/go@1.12/1.12.10/libexec"
# $GOPATH/src/github.com/go-swagger/go-swagger/vendor/github.com/go-openapi/swag/util.go удалить ненужные вхождения HTTP

find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' -type f -exec rm -f {} +
#/usr/local/Cellar/go@1.12/1.12.10/bin/go run ../../../vendor/github.com/go-swagger/go-swagger/cmd/swagger/swagger.go generate client -f swagger.yml
go run ../../../vendor/github.com/go-swagger/go-swagger/cmd/swagger/swagger.go generate client -f swagger.yml
find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' -type f -exec git add {} +
