package main

import (
	"context"
	"log"
	"os"

	"github.com/ionos-cloud/streams"
	"github.com/ionos-cloud/streams/codec"
	pb "github.com/ionos-cloud/streams/examples/producer/proto"
	"github.com/ionos-cloud/streams/msg"
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

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}

	m := &pb.Demo{Name: "foo"}
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	err = nc.Publish("foo", b)
	if err != nil {
		return err
	}

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

	ss := s.Branch("branch", func(m msg.Message[string, *pb.Demo]) (bool, error) {
		return m.Key() == "foo", nil
	}, func(m msg.Message[string, *pb.Demo]) (bool, error) {
		return m.Key() == "bar", nil
	})
	ss[0].Log("logging-foo").Sink("mock-sink-1", noop.NewSink[string, *pb.Demo]())
	ss[1].Log("logging-bar").Sink("mock-sink-2", noop.NewSink[string, *pb.Demo]())

	if s.Error(); err != nil {
		return err
	}

	return nil
}
