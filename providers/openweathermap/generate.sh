#!/bin/bash
set -x

# curl "http://localhost/rest/swagger.json" -o swagger.json

find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' -type f -exec rm -f {} +
swagger generate client -f swagger.yml --skip-validation
find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' -type f -exec git add {} +
