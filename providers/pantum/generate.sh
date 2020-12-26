#!/bin/bash
set -x

find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' ! -name 'product_id_enumer.go' ! -name 'cartridge_status_enumer.go' ! -name 'drum_status_enumer.go' -type f -exec rm -f {} +
swagger generate client -f swagger.yml
find . ! -name 'swagger.yml' ! -name 'generate.sh' ! -name 'client.go' ! -name 'product_id_enumer.go' ! -name 'cartridge_status_enumer.go' ! -name 'drum_status_enumer.go' -type f -exec git add {} +
