package main

import (
	"context"
	"fmt"
	"time"

	"github.com/katallaxie/pkg/conv"
	"github.com/katallaxie/pkg/errorx"
	"github.com/katallaxie/streams"
	natsx "github.com/katallaxie/streams/nats"
	"github.com/katallaxie/streams/sinks"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func main() {
	opts := &server.Options{
		JetStream: true,
		StoreDir:  "./tmp",
		Debug:     true,
	}
	ns, err := server.NewServer(opts)
	errorx.Panic(err)

	go ns.Start()

	if !ns.ReadyForConnections(4 * time.Second) {
		panic("not ready for connection")
	}

	nc, err := nats.Connect(ns.ClientURL())
	errorx.Panic(err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	js, err := nc.JetStream()
	errorx.Panic(err)

	streamName := "EVENS"

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"events.>"},
	})
	errorx.Panic(err)

	ack, err := js.Publish("events.1", []byte("hello"))
	errorx.Panic(err)
	fmt.Println("Published message with ack:", ack)

	sub, err := js.PullSubscribe("events.1", "consumer", nats.AckExplicit())
	if err != nil {
		panic(err)
	}

	s, err := natsx.NewJetStreamSource(ctx, &natsx.JetStreamSourceConfig{
		Conn:           nc,
		FetchBatchSize: 2,
		JetStreamCtx:   js,
		Sub:            sub,
	})
	errorx.Panic(err)

	s.Pipe(streams.DefaultPassThrough).Pipe(streams.NewMap(conv.String)).To(sinks.DefaultStdout)
}
