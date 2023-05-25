#!/usr/bin/env bash
set -euo pipefail

require() {
  if ! command -v "$1" &>/dev/null && [[ -n "$2" ]]; then
    go install "$2"
  fi
  if ! command -v "$1" >/dev/null; then
    >&2 echo "$1 not found"
    if [[ -n "$2" ]]; then
     >&2 echo "is \${GOPATH}/bin in your \$PATH?"
    fi
    exit 1
  fi
}
# install gen tools
require protoc-gen-go google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0
require protoc-gen-debug github.com/lyft/protoc-gen-star/protoc-gen-debug@v0.6.2
require mockery github.com/vektra/mockery/v2@v2.26.1
require protoc

echo "generating gotf.proto"
protoc -I. --go_out=. --go_opt module=github.com/travix/protoc-gen-gotf \
  --go_opt=Mgotf.proto="github.com/travix/protoc-gen-gotf/pb;pb" \
  gotf.proto

echo "generating mocks"
# TODO: remove this after https://github.com/vektra/mockery/discussions/549 is close
mockery 2>&1 | grep -v "ALPHA FEATURE" || true

pushd () {
    command pushd "$@" > /dev/null
}
popd () {
    command popd > /dev/null
}

for dir in testdata/*/; do
  pushd "${dir}"
  echo "generating code_generator_request.pb.bin for testdata/${dir}*.proto"
  protoc -I. -I../../ --plugin=protoc-gen-debug="$(go env GOPATH)/bin/protoc-gen-debug" --debug_out=".:." ./*.proto
  popd
done
