# Schemagen [![GoDoc](https://godoc.org/github.com/calmera/schemagen?status.svg)](https://godoc.org/github.com/calmera/schemagen)

This is a tool that fetches Avro schemas from [Confluent Schema Registry](https://github.com/confluentinc/schema-registry) and compiles them to Go code.

Code generation is entirely based on [gogen-avro](https://github.com/alanctgardner/gogen-avro).
Schema retrieval is based on [srclient](https://github.com/riferrei/srclient)

This project started out as a fork of [schemagen](https://github.com/burdiyan/schemagen), but has been
modified to include support for secured registries (like the confluent cloud one).

## Installation

Right now the only way to install `schemagen` is to build it from source:

```
go install github.com/calmera/schemagen/cmd/schemagen
```

## Getting Started

1. Create a file named `.schemagen.yaml` in the root of your project.
2. Specify Schema Registry URL, subjects and versions of the schema you want to download and compile.
3. Run `schemagen` to download the schemas from Schema Registry and compile them.

### Config Example

```
registry: http://confluent-schema-registry.default.svc.cluster.local:8081
target_dir: ./out
package: mypkg
schemas:
  - subject: my-topic-value
    version: latest
  - subject: another-topic-value
    version: "2"
```
