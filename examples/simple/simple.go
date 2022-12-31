package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ionos-cloud/streams"
	"github.com/ionos-cloud/streams/kafka"
	"github.com/ionos-cloud/streams/kafka/reader"
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
	rand.Seed(time.Now().UnixNano())

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

	src := kafka.WithContext[string, string](ctx, r, kafka.StringDecoder{}, kafka.StringDecoder{}, kafka.StringEncoder{})

	s := streams.NewStream[string, string](src)
	s.Log().Sink(noop.NewSink[string, string]())

	<-ctx.Done()

	return nil
}
