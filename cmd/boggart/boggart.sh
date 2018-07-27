#!/usr/bin/env bash
set -x

set +e
echo "17" > /sys/class/gpio/unexport
set -e

# export $(grep -v '^#' boggart.env | xargs)
# printenv
./boggart