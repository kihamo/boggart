#!/bin/bash
set -x

find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' -exec rm -f {} +
swagger generate client -f swagger.yml --skip-validation
find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' -exec git add {} +
