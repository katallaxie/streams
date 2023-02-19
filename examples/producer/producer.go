package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/ionos-cloud/streams/examples/producer/proto"
	"google.golang.org/protobuf/proto"

	"github.com/segmentio/kafka-go"
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

	topic := "demo12345"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	m := &pb.Demo{
		Name: "demo",
	}
	msg, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Key: []byte("foo"), Value: msg},
		kafka.Message{Key: []byte("bar"), Value: msg},
		kafka.Message{Value: msg},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	return nil
}
