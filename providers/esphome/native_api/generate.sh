#!/bin/bash
set -x

curl "https://raw.githubusercontent.com/esphome/esphome/v1.19.4/esphome/components/api/api.proto" -o ./api.proto
curl "https://raw.githubusercontent.com/esphome/esphome/v1.19.4/esphome/components/api/api_options.proto" -o ./api_options.proto

sed -i '' '2i\
package nativeapi;
' api.proto

sed -i '' '2i\
package nativeapi;
' api_options.proto

protoc --proto_path=./ --go_out=plugins=grpc:./ *.proto
