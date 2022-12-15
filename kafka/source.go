package kafka

import (
	"context"

	"github.com/ionos-cloud/streams/msg"
	kgo "github.com/segmentio/kafka-go"
)

type kafka struct {
	reader *kgo.Reader
}

// New ...
func New(r *kgo.Reader) *kafka {
	k := new(kafka)
	k.reader = r

	return k
}

// Commit ...
func (k *kafka) Commit(mm ...*msg.Message) error {
	return k.reader.CommitMessages(context.Background())
}

// Message ...
func (k *kafka) Messages() chan *msg.Message {
	out := make(chan *msg.Message)

	go func() {
		for {
			m, err := k.reader.FetchMessage(context.Background())
			if err != nil {
				break
			}

			out <- &msg.Message{
				Name: string(m.Key),
			}
		}
		close(out)
	}()

	return out
}
