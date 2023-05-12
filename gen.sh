#!/usr/bin/env bash
set -e

# install gen tools
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0
go install github.com/lyft/protoc-gen-star/protoc-gen-debug@v0.6.2
go install github.com/vektra/mockery/v2@v2.26.1

echo "generating plugin.proto"
protoc -I. --go_out=. --go_opt module=github.com/travix/protoc-gen-go-tf plugin.proto

echo "generating mocks"
# TODO: remove this after https://github.com/vektra/mockery/discussions/549 is close
mockery 2>&1 | grep -v "ALPHA FEATURE" || true

pushd () {
    command pushd "$@" > /dev/null
}
popd () {
    command popd "$@" > /dev/null
}

pushd plugin/testdata/

test_proto_dirs=(
  minimum-valid
)

for dir in "${test_proto_dirs[@]}"; do
  pushd "${dir}"
  echo "generating code_generator_request.pb.bin for testdata/${dir}/*.proto"
  protoc -I. -I../../../ --plugin=protoc-gen-debug="$(go env GOPATH)/bin/protoc-gen-debug" --debug_out=".:." ./*.proto
  popd 1
done
