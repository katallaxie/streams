package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ionos-cloud/streams"
	"github.com/ionos-cloud/streams/codec"
	pb "github.com/ionos-cloud/streams/examples/producer/proto"
	"github.com/ionos-cloud/streams/nats/source"
	"github.com/ionos-cloud/streams/noop"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/spf13/cobra"
)

var protoDecoder codec.Decoder[*pb.Demo] = func(b []byte) (*pb.Demo, error) {
	msg := new(pb.Demo)
	if err := proto.Unmarshal(b, msg); err != nil {
		return msg, err
	}

	return msg, nil
}

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

	nc, err := nats.Connect(nats.DefaultURL, nats.Timeout(1*time.Second))
	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			m := &pb.Demo{Name: "foo"}
			b, err := proto.Marshal(m)
			if err != nil {
				return
			}

			err = nc.Publish("foo", b)
			if err != nil {
				return
			}
		}
	}()

	sub, err := nc.SubscribeSync("foo")
	if err != nil {
		return err
	}

	src := source.WithContext(ctx, sub, codec.StringDecoder, protoDecoder)

	err = streams.DefaultRegisterer.Register(streams.DefaultMetrics)
	if err != nil {
		return err
	}

	mon := streams.NewMonitor(streams.DefaultMetrics)

	s := streams.NewStream[string, *pb.Demo](src, streams.WithMonitor(mon))
	s.Log("logging-foo").Sink("mock-sink-1", noop.NewSink[string, *pb.Demo]())

	if s.Error(); err != nil {
		return err
	}

	return nil
}
