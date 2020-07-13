#!/bin/bash
set -x

find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' ! -path './docs/*' ! -path './static/*' -type f -exec rm -f {} +
swagger generate client -f swagger.yml
find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' ! -path './docs/*' ! -path './static/*' -type f -exec git add {} +
