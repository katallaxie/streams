package table

import (
	"context"
	"errors"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ionos-cloud/streams"
	"github.com/ionos-cloud/streams/kafka/reader"
	"github.com/ionos-cloud/streams/kafka/writer"
	"github.com/ionos-cloud/streams/msg"

	kgo "github.com/segmentio/kafka-go"
)

// Topic is a kafka topic.
type Topic string

// String is returning the topic as string.
func (t Topic) String() string {
	return strings.Join([]string{prefix, string(t)}, sep)
}

type table struct {
	dialer *kgo.Dialer

	reader *kgo.Reader
	writer *kgo.Writer

	brokers []string
	topic   Topic
	ctx     context.Context

	err     error
	errOnce sync.Once

	streams.Table
	streams.Sink[string, []byte]
}

const (
	prefix = "table"
	sep    = "."
)

// Opt ...
type Opt func(t *table)

// WithBrokers ...
func WithBrokers(brokers ...string) Opt {
	return func(t *table) {
		t.brokers = append(t.brokers, brokers...)
	}
}

// WithTopic ...
func WithTopic(topic Topic) Opt {
	return func(t *table) {
		t.topic = topic
	}
}

// WithDialer ...
func WithDialer(dialer *kgo.Dialer) Opt {
	return func(t *table) {
		t.dialer = dialer
	}
}

// WithContext ...
func WithContext(ctx context.Context, opts ...Opt) streams.Table {
	t := new(table)
	t.ctx = ctx

	for _, opt := range opts {
		opt(t)
	}

	if t.dialer == nil {
		t.dialer = &kgo.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
		}
	}

	r := reader.NewReader(
		reader.WithDialer(t.dialer),
		reader.WithBrokers(t.brokers...),
		reader.WithTopic(string(t.topic)),
	)
	t.reader = r

	w := writer.NewWriter(
		writer.WithBrokers(t.brokers...),
		writer.WithTopic(string(t.topic)),
	)
	t.writer = w

	return t
}

// Set ...
func (t *table) Set(key string, value []byte) error {
	err := t.writer.WriteMessages(t.ctx, kgo.Message{
		Key:   []byte(key),
		Value: value,
	})
	if err != nil {
		return err
	}

	return nil
}

// Setup ...
func (t *table) Setup() error {
	topic := string(t.topic)

	conn, err := t.dialer.DialContext(t.ctx, "tcp", t.brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	var controllerConn *kgo.Conn
	controllerConn, err = kgo.DialContext(t.ctx, "tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfig := []kgo.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
			ConfigEntries: []kgo.ConfigEntry{
				{
					ConfigName:  "cleanup.policy",
					ConfigValue: "compact",
				},
			},
		},
	}

	err = controllerConn.CreateTopics(topicConfig...)
	if err != nil && !errors.Is(err, kgo.TopicAlreadyExists) {
		return err
	}

	return nil
}

// Delete ...
func (t *table) Delete(key string) error {
	err := t.writer.WriteMessages(t.ctx, kgo.Message{
		Key:   []byte(key),
		Value: nil, // Tombstone record
	})
	if err != nil {
		return err
	}

	return nil
}

// Next ...
func (t *table) Next() <-chan streams.NextCursor {
	out := make(chan streams.NextCursor)

	go func() {
		for {
			m, err := t.reader.ReadMessage(t.ctx)
			if err != nil {
				t.errOnce.Do(func() {
					t.err = err
				})

				break
			}

			latest := t.reader.Lag() == 0

			out <- streams.NextCursor{Key: string(m.Key), Value: m.Value, Latest: latest}
		}

		close(out)
	}()

	return out
}

// Write confirms to the Sink interface. And allows to persist to tables as a sink.
func (t *table) Write(msgs ...msg.Message[string, []byte]) error {
	for _, m := range msgs {
		err := t.Set(m.Key(), m.Value())
		if err != nil {
			return err
		}
	}

	return nil
}
