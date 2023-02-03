package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// ResetFirstOffsetConsumerGroup resets the first offset of a consumer group.
func ResetFirstOffsetConsumerGroup(ctx context.Context, brokers []string, topic string, gid string) error {
	conn, err := kafka.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	parts, err := conn.ReadPartitions(topic)
	if err != nil {
		return err
	}

	offsets := make(map[int]int64, len(parts))

	for _, p := range parts {
		addr := fmt.Sprintf("%s:%d", p.Leader.Host, p.Leader.Port)
		c, err := kafka.DialLeader(ctx, "tcp", addr, topic, p.ID)
		if err != nil {
			return err
		}
		defer c.Close()

		offset, err := c.ReadFirstOffset()
		if err != nil {
			return err
		}
		offsets[p.ID] = offset
	}

	group, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		Brokers: brokers,
		Topics:  []string{topic},
		ID:      gid,
	})
	if err != nil {
		return err
	}
	defer group.Close()
	gen, err := group.Next(ctx)
	if err != nil {
		return err
	}

	err = gen.CommitOffsets(map[string]map[int]int64{topic: offsets})
	if err != nil {
		return err
	}
	return nil
}

// ResetTimestampOffsetConsumerGroup resets the first offset of a consumer group.
func ResetTimestampOffsetConsumerGroup(ctx context.Context, brokers []string, topic string, gid string, timestamp time.Time) error {
	conn, err := kafka.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	parts, err := conn.ReadPartitions(topic)
	if err != nil {
		return err
	}

	offsets := make(map[int]int64, len(parts))

	for _, p := range parts {
		addr := fmt.Sprintf("%s:%d", p.Leader.Host, p.Leader.Port)
		c, err := kafka.DialLeader(ctx, "tcp", addr, topic, p.ID)
		if err != nil {
			return err
		}
		defer c.Close()

		offset, err := c.ReadOffset(timestamp)
		if err != nil {
			return err
		}
		offsets[p.ID] = offset
	}

	group, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		Brokers: brokers,
		Topics:  []string{topic},
		ID:      gid,
	})
	if err != nil {
		return err
	}
	defer group.Close()
	gen, err := group.Next(ctx)
	if err != nil {
		return err
	}

	err = gen.CommitOffsets(map[string]map[int]int64{topic: offsets})
	if err != nil {
		return err
	}
	return nil
}
