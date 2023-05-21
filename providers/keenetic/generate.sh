#!/bin/bash
set -x

# $GOPATH/src/github.com/go-swagger/go-swagger/vendor/github.com/go-openapi/swag/util.go удалить ненужные вхождения HTTP

find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' ! -name 'error.go' ! -name 'round_tripper.go' ! -name 'runtime.go' ! -path './docs/*' ! -path './static/*' -type f -exec rm -f {} +
swagger generate client -f swagger.yml
find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' ! -name 'error.go' ! -name 'round_tripper.go' ! -name 'runtime.go' ! -path './docs/*' ! -path './static/*' -type f -exec git add {} +
