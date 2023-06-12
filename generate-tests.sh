#!/bin/bash

rm -rf tests/generated || true
mkdir -p tests/generated || true

protoc --plugin=zig-out/bin/protoc-gen-zig \
  --zig_out=tests/generated \
  -Itests/protos_for_test \
  tests/protos_for_test/all.proto

zig fmt tests/generated/tests.pb.zig

echo 'generate-tests.sh finished'