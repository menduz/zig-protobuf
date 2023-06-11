#!/bin/bash

rm -rf tmp || true
mkdir -p tmp || true

protoc --plugin=zig-out/bin/protoc-gen-zig \
  --zig_out=tmp \
  -Itests/protos_for_test \
  /usr/local/lib/protobuf/include/google/protobuf/compiler/plugin.proto \
  /usr/local/lib/protobuf/include/google/protobuf/descriptor.proto \
  tests/protos_for_test/custom-import.proto

zig fmt tmp/tests.pb.zig
zig fmt tmp/google/protobuf.pb.zig
zig fmt tmp/google/protobuf/compiler.pb.zig

rm -rf dcl || true
mkdir -p dcl || true
protoc --plugin=zig-out/bin/protoc-gen-zig \
  --zig_out=dcl \
  -I=protocol/proto \
  -I=protocol/public \
  protocol/public/sdk-components.proto \
  protocol/public/sdk-apis.proto

echo 'test.sh finished'