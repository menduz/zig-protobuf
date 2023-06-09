#!/bin/bash

cd generator

rm protoc-gen-main || true 
go build  -o protoc-gen-main main.go

cd ..

rm -rf tests/generated_output
mkdir -p tests/generated_output || true

export GO_IMPORT_PATH=tests/protos_for_test
protoc --plugin=generator/protoc-gen-main \
  --main_out=tests/generated_output \
  /usr/local/lib/protobuf/include/google/protobuf/compiler/plugin.proto \
  /usr/local/lib/protobuf/include/google/protobuf/descriptor.proto

zig fmt tests/generated_output/google/protobuf/compiler/plugin.pb.zig
zig fmt tests/generated_output/google/protobuf/descriptor.pb.zig

echo 'generation finished'