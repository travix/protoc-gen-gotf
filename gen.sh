#!/usr/bin/env bash
set -e
echo "generating plugin.proto"
protoc -I. --go_out=. --go_opt module=github.com/travix/protoc-gen-go-tf plugin.proto

go install github.com/lyft/protoc-gen-star/protoc-gen-debug@latest
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
