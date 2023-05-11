#!/usr/bin/env bash
set -e

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/lyft/protoc-gen-star/protoc-gen-debug@latest

echo "generating plugin.proto"
protoc -I. --go_out=. --go_opt module=github.com/travix/protoc-gen-go-tf plugin.proto

test_proto_dirs=(
  minimum-valid
)
pushd plugin/testdata/
for dir in "${test_proto_dirs[@]}"; do
  pushd "${dir}"
  echo "generating testdata/${dir}/*.proto"
  protoc -I. -I../../../ --plugin=protoc-gen-debug="$(go env GOPATH)/bin/protoc-gen-debug" --debug_out=".:." ./*.proto
  popd
done
