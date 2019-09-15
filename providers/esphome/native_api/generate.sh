#!/bin/bash
set -x

curl "https://raw.githubusercontent.com/esphome/esphome/dev/esphome/components/api/api.proto" -o ./api.proto
curl "https://raw.githubusercontent.com/esphome/esphome/dev/esphome/components/api/api_options.proto" -o ./api_options.proto

sed -i '' '2i\
package native_api;
' api.proto

sed -i '' '2i\
package native_api;
' api_options.proto

protoc --proto_path=./ --go_out=plugins=grpc:./ *.proto
