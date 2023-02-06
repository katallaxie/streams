package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/ionos-cloud/streams/kafka/reader"
	"github.com/ionos-cloud/streams/kafka/writer"
	"github.com/ionos-cloud/streams/view"
	kgo "github.com/segmentio/kafka-go"
)

type table struct {
	reader *kgo.Reader
	writer *kgo.Writer
}

// NewTable ...
func NewTable() view.Table {
	t := new(table)

	dialer := &kgo.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	r := reader.NewReader(
		reader.WithDialer(dialer),
		reader.WithBrokers("localhost:9092"),
		reader.WithTopic("samples_view"),
	)

	w := writer.NewWriter(
		writer.WithBrokers("localhost:9092"),
		writer.WithTopic("samples_view"),
	)

	t.writer = w
	t.reader = r

	return t
}

// Set ...
func (t *table) Set(key string, value []byte) error {
	fmt.Println("Set", key, string(value))

	err := t.writer.WriteMessages(context.Background(), kgo.Message{
		Key:   []byte(key),
		Value: value,
	})

	if err != nil {
		return err
	}

	return nil
}

// Delete ...
func (t *table) Delete(key string) error {
	return nil
}

// Next ...
func (t *table) Next() <-chan view.NextCursor {
	out := make(chan view.NextCursor)

	go func() {
		for {
			m, err := t.reader.ReadMessage(context.Background())
			if err != nil {
				break
			}

			out <- view.NextCursor{Key: string(m.Key), Value: m.Value}
		}

		close(out)
	}()

	return out
}
