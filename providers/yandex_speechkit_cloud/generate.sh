#!/bin/bash
set -x

# $GOPATH/src/github.com/go-swagger/go-swagger/vendor/github.com/go-openapi/swag/util.go удалить ненужные вхождения HTTP

find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' -type f -exec rm -f {} +
swagger generate client -f swagger.yml
find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' -type f -exec git add {} +
