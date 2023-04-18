package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/codec"
	pb "github.com/katallaxie/streams/examples/producer/proto"
	"github.com/katallaxie/streams/kafka/reader"
	"github.com/katallaxie/streams/kafka/source"
	"github.com/katallaxie/streams/msg"
	"google.golang.org/protobuf/proto"

	v8 "github.com/katallaxie/v8go"
	"github.com/katallaxie/v8go-polyfills/console"
	kgo "github.com/segmentio/kafka-go"
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

	dialer := &kgo.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	r := reader.NewReader(
		reader.WithDialer(dialer),
		reader.WithBrokers("localhost:9092"),
		reader.WithGroupID("demo12345"),
		reader.WithTopic("demo12345"),
		reader.WithLogger(nil),
	)

	src := source.WithContext(ctx, r, codec.StringDecoder, protoDecoder, codec.StringEncoder)

	err := streams.DefaultRegisterer.Register(streams.DefaultMetrics)
	if err != nil {
		return err
	}

	m := streams.NewMonitor(streams.DefaultMetrics)

	fn := func(msg msg.Message[string, *pb.Demo]) {
		iso := v8.NewIsolate()
		global := v8.NewObjectTemplate(iso)

		defer iso.Dispose()

		ctx := v8.NewContext(iso, global)
		defer ctx.Close()

		err := console.AddTo(ctx)
		if err != nil {
			return
		}

		v, err := ctx.RunScript("console.log(\"ping pong\"); true", "worker.js")
		if err != nil {
			return
		}
		defer v.Release()
	}

	s := streams.NewStream[string, *pb.Demo](src, streams.WithMonitor(m))
	s.Do("worker", fn).Mark("mark").Drain()

	if s.Error(); err != nil {
		return err
	}

	return nil
}
