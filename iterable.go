package streams

import (
	"context"
	"sync"

	"github.com/ionos-cloud/streams/msg"
)

type channelIterable struct {
	next        <-chan *msg.Message
	opts        []Option
	subscribers []chan *msg.Message
	producing   bool
	sync.RWMutex
}

func newChannelIterable(next <-chan *msg.Message, opts ...Option) Iterable {
	return &channelIterable{
		next:        next,
		subscribers: make([]chan *msg.Message, 0),
		opts:        opts,
	}
}

// Observe ...
func (c *channelIterable) Observe(opts ...Option) <-chan *msg.Message {
	out := make(chan *msg.Message)
	c.Lock()
	defer c.Unlock()

	c.subscribers = append(c.subscribers, out)

	return out
}

func (c *channelIterable) connect(ctx context.Context) {
	c.Lock()
	if !c.producing {
		go c.produce(ctx)
		c.producing = true
	}
	c.Unlock()
}

func (c *channelIterable) produce(ctx context.Context) {
	defer c.close()

	for {
		select {
		case <-ctx.Done():
			return
		case m, ok := <-c.next:
			if !ok {
				return
			}

			c.RLock()
			for _, sub := range c.subscribers {
				sub <- m
			}
			c.RUnlock()
		}
	}
}

func (c *channelIterable) close() func() {
	return func() {
		c.RLock()

		for _, sub := range c.subscribers {
			close(sub)
		}

		c.RUnlock()
	}
}
