package main

import (
	"strconv"
	"time"

	"github.com/katallaxie/streams"
	"github.com/katallaxie/streams/sinks"
	"github.com/katallaxie/streams/sources"
)

type message struct {
	msg string
}

func (msg *message) String() string {
	return msg.msg
}

func main() {
	err := sources.NewChanSource(tickerChan(time.Second)).
		Pipe(streams.PassThrough()).
		Pipe(streams.Timeout(5 * time.Second)).
		To(sinks.DefaultStdout)
	if err != nil {
		panic(err)
	}
}

func tickerChan(interval time.Duration) chan any {
	outChan := make(chan any)

	go func() {
		ticker := time.NewTicker(interval)
		for t := range ticker.C {
			outChan <- &message{msg: strconv.FormatInt(t.UnixMilli(), 10)}
		}
	}()

	return outChan
}
