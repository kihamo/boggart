#!/bin/bash

wget -c http://raspberry:7080/rest/swagger.json
swagger generate client -f swagger.json --skip-validation