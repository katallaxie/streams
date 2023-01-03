# :surfing_woman: Streams

[![Release](https://github.com/ionos-cloud/streams/actions/workflows/main.yml/badge.svg)](https://github.com/ionos-cloud/streams/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ionos-cloud/streams)](https://goreportcard.com/report/github.com/ionos-cloud/streams)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)

A teeny-tiny package to create stream processing workloads. It is intended to be used with [Apache Kafka](https://kafka.apache.org/).  

:warning: This is experimental. APIs may change. :warning:

## Getting Started

There are only a few packages that help Gophers to create stream processing workloads. This package is one of them. It is intended to be used with [Apache Kafka](https://kafka.apache.org/).

```bash
go get github.com/ionos-cloud/streams
```

It features a channel based API to consume messages from a Kafka topic and a channel based API to produce messages to a Kafka topic. It assumes the use of a [consumer group](https://docs.confluent.io/platform/current/clients/consumer.html#:~:text=A%20consumer%20group%20is%20a,proportional%20share%20of%20the%20partitions.) for the consumption of messages.

The package connects a `source` with a sink via small functional operatios.

* `Map`
* `Filter`
* `Log`
* `FanOut`
* `Do`
* `Merge`
* `Branch`

There is support for [Prometheus](https://prometheus.io/) metrics.

## Docs

You can find the documentation hosted on [godoc.org](https://godoc.org/github.com/ionos-cloud/streams).

## License

[Apache 2.0](/LICENSE)