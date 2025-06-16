# minohm-go

A Go implementation of the [minohm][] interface, for using Ohm grammars from Go.

A _grammar blob_ is an [Ohm][] grammar that has been compiled to .wasm via the `@ohm-js/wasm` NPM package. To use a grammar blob to match some input, you need a miniohm implementation for your host language of choice (JavaScript, Go, Python, etc.) This package provides a miniohm implementation for the [Go Programming Language][go].

[minohm]: https://github.com/ohmjs/ohm/blob/main/doc/design/miniohm.md
[Ohm]: https://ohmjs.org
[go]: https://go.dev/

## Overview

The implementation consists of two main components:

1. **matcher.go**: A Go implementation of the JavaScript `WasmMatcher` class from the Ohm `wasm` package
2. **testmain.go**: A command-line program that demonstrates how to use the WasmMatcher

## WasmMatcher

The `WasmMatcher` class provides a high-level API for working with Ohm grammar blobs:

```go
matcher := NewWasmMatcher(ctx)
err := matcher.LoadModule("path/to/grammar.wasm")
matcher.SetInput("text to match")
success, err := matcher.Match()
success, err := matcher.MatchRule("specificRule")
cstRoot, err := matcher.GetCstRoot()
```

## Walking the CST

A full implementation of semantics, operations, etc. is not part of the miniohm interface. Instead, you can walk the CST (concrete syntax tree) directly using the CstNode interface. See testmain.go for an example.

## Developing

Useful commands:

```sh
make # Build
make test # Run tests
```
