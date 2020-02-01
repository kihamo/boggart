#!/bin/bash
set -x

# curl "http://localhost/rest/swagger.json" -o swagger.json

find . ! -name 'swagger.json' ! -name 'generate.sh' ! -name 'client.go' -type f -exec rm -f {} +
go  run ../../vendor/github.com/go-swagger/go-swagger/cmd/swagger/swagger.go generate client -f swagger.json --skip-validation
find . ! -name 'swagger.json' ! -name 'generate.sh' ! -name 'client.go' -type f -exec git add {} +
