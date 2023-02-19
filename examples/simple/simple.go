package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ionos-cloud/streams"
	"github.com/ionos-cloud/streams/codec"
	"github.com/ionos-cloud/streams/kafka"
	"github.com/ionos-cloud/streams/kafka/reader"
	"github.com/ionos-cloud/streams/msg"
	"github.com/ionos-cloud/streams/noop"

	kgo "github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "simple",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd.Context())
	},
}

func init() {
	rootCmd.SilenceUsage = true
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	dialer := &kgo.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	r := reader.NewReader(
		reader.WithDialer(dialer),
		reader.WithBrokers("localhost:9092"),
		reader.WithGroupID("demo12345"),
		reader.WithTopic("demo12345"),
	)

	src := kafka.WithContext(ctx, r, codec.StringDecoder, codec.StringDecoder, codec.StringEncoder)

	err := streams.DefaultRegisterer.Register(streams.DefaultMetrics)
	if err != nil {
		return err
	}

	m := streams.NewMonitor(streams.DefaultMetrics)

	s := streams.NewStream[string, string](src, streams.WithMonitor(m), streams.WithBuffer(1))
	ss := s.Branch("branch1", func(m msg.Message[string, string]) (bool, error) {
		return m.Key() == "foo", nil
	}, func(m msg.Message[string, string]) (bool, error) {
		return m.Key() == "bar", nil
	})
	ss[0].Log("logging-foo").Sink("mock-sink-1", noop.NewSink[string, string]())
	ss[1].Log("logging-bar").Sink("mock-sink-2", noop.NewSink[string, string]())

	if s.Error(); err != nil {
		return err
	}

	return nil
}
