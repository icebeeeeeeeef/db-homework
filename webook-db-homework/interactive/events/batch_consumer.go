package events

import (
	"context"
	"time"
	"webook/interactive/repository"
	"webook/pkg/logger"
	"webook/pkg/saramax"

	"github.com/IBM/sarama"
)

type KafkaBatchConsumer struct {
	client    sarama.Client
	l         logger.LoggerV1
	repo      repository.InteractiveRepository
	batchSize int
	timeout   time.Duration
}

func NewKafkaBatchConsumer(client sarama.Client, l logger.LoggerV1, repo repository.InteractiveRepository) saramax.Consumer {
	return &KafkaBatchConsumer{client: client, l: l, repo: repo, batchSize: 100, timeout: time.Second * 10}
}

func (c *KafkaBatchConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", c.client)
	if err != nil {
		return err
	}
	c.l.Info("启动 batch consumer", logger.Field{
		Key:   "topics",
		Value: []string{"read_event"},
	})
	go func() {
		err := cg.Consume(context.Background(), []string{"read_event"},
			saramax.NewBatchHandler_[ReadEvent](c.Consume, c.l, c.batchSize, c.timeout))
		if err != nil {
			c.l.Error("消费读事件失败", logger.Error(err))
		}
	}()
	return nil
}

func (c *KafkaBatchConsumer) Consume(msgs []*sarama.ConsumerMessage, evts []ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	biz := make([]string, 0, len(msgs))
	ids := make([]int64, 0, len(evts))
	for _, evt := range evts {
		biz = append(biz, "article")
		ids = append(ids, evt.Aid)
	}
	err := c.repo.BatchIncRead(ctx, biz, ids)
	if err != nil {
		c.l.Error("批量增加阅读量失败", logger.Error(err))
	} else {
		c.l.Info("批量增加阅读量成功", logger.Field{
			Key:   "count",
			Value: len(ids),
		})
	}
	return nil
}
