#!/bin/bash
set -x

curl "https://openhab.shosho.ru/rest/spec" -o swagger.json

find . ! -name 'swagger.json' ! -name 'generate.sh' ! -name 'config.yaml' ! -path './templates/*' -type f -exec rm -f {} +
oapi-codegen -generate types,client,spec -o openhab.go -package openhab3 --config config.yaml swagger.json
find . ! -name 'swagger.json' ! -name 'generate.sh' ! -name 'config.yaml' ! -path './templates/*' -type f -exec git add {} +
