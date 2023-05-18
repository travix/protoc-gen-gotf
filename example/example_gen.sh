#!/usr/bin/env bash
set -euo pipefail

require() {
  if ! command -v "$1" &>/dev/null && [[ -n "$2" ]]; then
    go install "$2"
  fi
  if ! command -v "$1" >/dev/null; then
    echo >&2 "$1 not found"
    if [[ -n "$2" ]]; then
      echo >&2 "is \${GOPATH}/bin in your \$PATH?"
    fi
    exit 1
  fi
}

require protoc-gen-gotag github.com/srikrsna/protoc-gen-gotag@v0.6.2
require protoc-gen-gotf

module="github.com/travix/protoc-gen-gotf/example"

echo "generating protobuf and terraform code"
protoc -I. -I../ \
  --go_out=. --go_opt module=${module} \
  --go-grpc_out=. --go-grpc_opt module=${module} \
  --gotf_out=. --gotf_opt=log_level=debug --gotf_opt module=${module} \
  example.proto

echo "adding tfsdk struct tag to protobuf"
protoc -I. -I../ --gotag_out=auto="tfsdk":. --gotag_opt module=${module} example.proto
