# Protoc plugin Cli Options

<details open>
<summary>Setting options in protoc command</summary>

```shell
protoc -I. \
  --gotf_out=. --gotf_opt module=${module} --gotf_opt log_level=debug \
  <YOUR_PROTO_FILES>
```

</details>

<details>
<summary>Setting options in buf.gen.yaml</summary>

```yaml
version: v1
plugins:
  - plugin: gotf
    out: .
    opt:
    - module=github.com/travix/gotf-example
    - log_level=debug
```

</details>

## `module=`

_Required_ option in the plugin argument.

The value of this option will be used as a go module name in the generated code.
It is required to create import statements in generated terraform go code.

Current code generation is tested with this option set. Generated code is not guaranteed to work without it.

## `log_level=`

_Optional_ option in the plugin argument.
The default value is `warn`.

Available values are `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`.
