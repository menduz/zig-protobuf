#!/bin/bash

cd generator

# rm protoc-gen-main || true 
go build  -o protoc-gen-main main.go

cd ..

rm -rf bootstrapped-generator/google
mkdir -p bootstrapped-generator/google || true

export GO_IMPORT_PATH=tests/protos_for_test
protoc --plugin=generator/protoc-gen-main \
  --main_out=bootstrapped-generator \
  /usr/local/lib/protobuf/include/google/protobuf/compiler/plugin.proto \
  /usr/local/lib/protobuf/include/google/protobuf/descriptor.proto

zig fmt bootstrapped-generator/google/protobuf/compiler/plugin.pb.zig
zig fmt bootstrapped-generator/google/protobuf/descriptor.pb.zig

cp -r bootstrapped-generator/google tests/generated_output

echo 'generation finished'