package saramax

import (
	"context"
	"encoding/json"
	"time"
	"webook/pkg/logger"

	"github.com/IBM/sarama"
)

type BatchHandler interface {
	Setup(session sarama.ConsumerGroupSession) error
	Cleanup(session sarama.ConsumerGroupSession) error
	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}

type BatchHandler_[T any] struct {
	fn        func(msg []*sarama.ConsumerMessage, ts []T) error
	l         logger.LoggerV1
	batchSize int
	timeout   time.Duration
}

func NewBatchHandler_[T any](fn func(msg []*sarama.ConsumerMessage, ts []T) error, l logger.LoggerV1, batchSize int, timeout time.Duration) BatchHandler {
	return &BatchHandler_[T]{fn: fn, l: l, batchSize: batchSize, timeout: timeout}
}

func (h *BatchHandler_[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *BatchHandler_[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

//把处理数据之前的重复逻辑全部封装起来

func (h *BatchHandler_[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgsCh := claim.Messages()
	for {
		ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
		defer cancel()
		var msgs = make([]*sarama.ConsumerMessage, 0, h.batchSize)
		var ts = make([]T, 0, h.batchSize)
		var done = false
		for i := 0; i < h.batchSize && !done; i++ {
			select {
			case msg, ok := <-msgsCh:
				if !ok {
					//这里代表通道关闭了，没必要继续
					cancel()
					return nil
				}
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					h.l.Error("反序列化消息失败", logger.Error(err),
						logger.String("topic", msg.Topic),
						logger.Int("partition", int(msg.Partition)),
						logger.Int("offset", int(msg.Offset)),
					)
					continue
				}
				//把他们两个放在一起，确保长度相同
				msgs = append(msgs, msg)
				ts = append(ts, t)
			case <-ctx.Done(): //避免等待一批时间过长
				done = true
			}
		}
		cancel()
		if len(msgs) == 0 {
			continue
		}
		err := h.fn(msgs, ts)
		if err != nil {
			h.l.Error("处理消息失败", logger.Error(err)) //这里其实实际上应该记录整个批次的信息才算完整
			continue
		}
		for _, msg := range msgs {
			session.MarkMessage(msg, "")
		}
	}
}
