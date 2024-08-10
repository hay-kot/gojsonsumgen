# Go JSON Sum Type Generator

This is an experimental code generator for generating sum types in go that are JSON marshable

## Installation

```bash
go install github.com/hay-kot/gojsonsumgen
```

## Usage

```bash
gojsonsumgen gen .
```

The `gen` command will walk through the directory and generate sum types for all the files that have the matching directive comment

```go
// gosumtype: <type name>
```

### Example Comment

```go
// ActionType is a sum type discriminator for ActionDef
//
// gosumtype: ActionDef
//
//	opt:tag:  type
//	opt:name: Type
//
//	'action-shell'  : ActionShell
//	'action-python' : ActionPython
//	'action-http'   : ActionHTTP
type ActionType string
```

See [example](https://github.com/hay-kot/gojsonsumgen/tree/main/examples) for more examples of input and output files.
