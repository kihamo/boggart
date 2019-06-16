#!/bin/bash

find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' -type f -exec rm -f {} +
swagger generate client -f swagger.yml
find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' -type f -exec git add {} +