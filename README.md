# :surfing_woman: Streams

[![Release](https://github.com/katallaxie/streams/actions/workflows/main.yml/badge.svg)](https://github.com/katallaxie/streams/actions/workflows/main.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/katallaxie/streams.svg)](https://pkg.go.dev/github.com/katallaxie/streams)
[![Go Report Card](https://goreportcard.com/badge/github.com/katallaxie/streams)](https://goreportcard.com/report/github.com/katallaxie/streams)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)

A teeny-tiny package to create stream processing workloads. 

## Docs

You can find the documentation hosted on [godoc.org](https://godoc.org/github.com/katallaxie/streams).

## Examples

See the [examples](/examples) directory for more.

## Operators

* `Do`: Execute a function for each element in the stream.
* `Filter`: Filter elements from the stream.
* `Map`: Transform elements in the stream.
* `Reduce`: Reduce elements in the stream.
* `Split`: Split the stream into multiple streams.
* `Merge`: Merge multiple streams into one.
* `FlatMap`: Transform elements in the stream into multiple elements.
* `Skip`: Skip elements in the stream.

## Source 

* `Channel`: Takes a channel as an input

## Sink

* `Channel`: Takes a channel as an output
* `FSM`: Takes a finite state machine as an output
* `Stdout`: Takes the standard output as an output
* `Ignore`: Ignores the output

## License

[Apache 2.0](/LICENSE)
