# Usage

## 1. import gotf.proto

- Add import a statement to your proto files.
    ```proto
    syntax = "proto3";
    
    import "gotf.proto";
    ```
  Proto syntax `proto3` is supported, might work with `proto2` but not guaranteed.
- Source [gotf.proto]
    - Either copy [gotf.proto] to your project
    - Or use [buf.build/travix/gotf] as a dependency in your buf.yaml file.

## 2. Add gotf options

Add gotf options to your proto files. See [Proto options] or [gotf.proto] for more details.

## 2. Generate protobuf files

- Generate protobuf files (`*.pb`)
- If you wish to use gRPC clients in your provider also generate gRPC service, client files (`*_grpc.pb.go`)
<details open>
<summary>Example protoc command</summary>

```shell
module=<YOUR_GO_MOD_NAME>
protoc -I. \
  --go_out=. --go_opt module=${module} \
  --go-grpc_out=. --go-grpc_opt module=${module} \
  <YOUR_PROTO_FILES>
```
you might need to pass `-I<PATH_TO_DIR_WITH_gotf.proto>` if gotf.proto is not in the same directory as your proto files.

</details>

<details>
<summary>Example buf generate</summary>

**buf.yaml**
```yaml
version: v1beta1
deps:
  - buf.build/travix/gotf
breaking:
  use:
    - FILE
lint:
  use:
    - DEFAULT
```
**buf.gen.yaml**
```yaml
version: v1
plugins:
  - plugin: go
    out: .
    opt:
    - module=<YOUR_GO_MOD_NAME>
  - plugin: go-grpc
    out: .
    opt:
    - module=<YOUR_GO_MOD_NAME>
```
**commands**
```shell
buf mod update
buf generate --path <YOUR_PROTO_FILES>
```

</details>

## 3. Generate go terraform code

Run protoc again with gotf plugin to generate go terraform code

<details open>
<summary>Example protoc command</summary>

```shell
protoc -I. \
  --gotf_out=. --gotf_opt module=${module} \
  <YOUR_PROTO_FILES>
```
you might need to add `-I<PATH_TO_DIR_WITH_gotf.proto>` arg to protoc if gotf.proto is not in the same directory as your proto files.

</details>

<details>
<summary>Example buf generate</summary>

**buf.gen.tf.yaml**
```yaml
version: v1
plugins:
  - plugin: gotf
    out: .
    opt:
    - module=<YOUR_GO_MOD_NAME>
```
**commands**
```shell
buf generate --path <YOUR_PROTO_FILES> --template buf.gen.tag.yaml
```

</details>

## 4. Implement executors for provider, resources and datasources

Generated terraform code is a bridge between your service and terraform cli.
It contains interfaces that need to be implemented you.
Check [Executor] for more details on executor.

Generated code and executor implementation are still pieces that need to be put together.
See [gotf-example] for a crude example of how to build and use provider locally.

For authoritative guide follow [https://developer.hashicorp.com/terraform/plugin](https://developer.hashicorp.com/terraform/plugin).

[gotf]: https://github.com/travix/gotf
[gotf.proto]: ../gotf.proto
[buf.build/travix/gotf]: https://buf.build/travix/gotf
[Proto options]: ./proto-options.md
[Executor]: ./executor.md
[gotf-example]: https://github.com/travix/gotf-example
