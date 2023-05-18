# protoc-gen-gotf

protoc plugin that generates go code for terraform provider using protobuf messages and services

## After

```shell
go install github.com/srikrsna/protoc-gen-gotag@latest
protoc -I. --gotag_out=auto="tfsdk":. --gotag_opt module=<GO_MOD_NAME> example.proto
```
