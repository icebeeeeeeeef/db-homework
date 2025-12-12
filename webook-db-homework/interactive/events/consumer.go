package events

import (
	"context"
	"time"
	"webook/interactive/repository"
	"webook/pkg/logger"
	"webook/pkg/saramax"

	"github.com/IBM/sarama"
)

var _ saramax.Consumer = (*KafkaConsumer)(nil)

type KafkaConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
	repo   repository.InteractiveRepository
}

func NewKafkaConsumer(client sarama.Client, l logger.LoggerV1, repo repository.InteractiveRepository) *KafkaConsumer {
	return &KafkaConsumer{client: client, l: l, repo: repo}
}

func (c *KafkaConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", c.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(), []string{"read_event"},
			saramax.NewHandler_[ReadEvent](c.Consume, c.l))
		if err != nil {
			c.l.Error("消费读事件失败", logger.Error(err))
		}
	}()
	return nil
}

func (c *KafkaConsumer) Consume(msg *sarama.ConsumerMessage, evt ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.repo.IncRead(ctx, "article", evt.Aid)
}
