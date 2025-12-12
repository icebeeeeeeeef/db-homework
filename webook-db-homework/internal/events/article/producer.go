package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

type Producer interface {
	ProduceReadEvent(ctx context.Context, evt ReadEvent) error
}

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(producer sarama.SyncProducer) Producer {
	return &KafkaProducer{
		producer: producer,
	}
}

func (p *KafkaProducer) ProduceReadEvent(ctx context.Context, evt ReadEvent) error {

	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "read_event",
		Value: sarama.ByteEncoder(data),
	})
	return err
}

func (evt ReadEvent) String() string {
	return fmt.Sprintf("Uid: %d, Aid: %d", evt.Uid, evt.Aid)
}

type ReadEvent struct {
	Uid int64
	Aid int64
}
