# protoc-gen-gotf

Protoc plugin that generates go code for terraform provider using protobuf messages and services.

### Motive

Writing terraform provider is a nontrivial task,
updating attributes and schema of resources and datasources specially is an error-prone activity.
For Resources and Datasources, a lot of boilerplate code is required to be written,
such as handling authentication with provider service,
marshalling and unmarshalling of models and if gRPC is used initializing service clients, etc.

This plugin's aim is
to generate the boilerplate code,
synchronization of protobuf messages and their fields with terraform block schema and attributes using protobuf options.

## 💻 Install

Pre-build binaries are available on [releases page].

Or install via go cli

```shell
go install github.com/travix/protoc-gen-gotf@latest
```

## ✏️ Example

> - [gotf-example] repository for a working example.
> - [gotf] for a go interfaces used by generated code.

## [🧑‍💻 Usage][Usage]

## [📋 Documentation][Documentation]

[gotf-example]: https://github.com/travix/gotf-example
[gotf]: https://github.com/travix/gotf
[releases page]: https://github.com/travix/protoc-gen-gotf/releases/latest
[Documentation]: ./docs/README.md
[Usage]: ./docs/usage.md
