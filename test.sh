#!/bin/bash

cd generator

rm protoc-gen-main || true 
go build  -o protoc-gen-main main.go

cd ..

rm -rf tmp
mkdir -p tmp || true

export GO_IMPORT_PATH=tests/protos_for_test
protoc --plugin=zig-out/bin/protoc-gen-zig \
  --zig_out=tmp \
  /usr/local/lib/protobuf/include/google/protobuf/compiler/plugin.proto \
  /usr/local/lib/protobuf/include/google/protobuf/descriptor.proto

zig fmt tmp/google/protobuf/compiler/plugin.proto.zig
zig fmt tmp/google/protobuf/descriptor.proto.zig


echo 'test.sh finished'