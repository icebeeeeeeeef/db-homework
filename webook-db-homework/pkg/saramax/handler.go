package saramax

import (
	"encoding/json"
	"webook/pkg/logger"

	"github.com/IBM/sarama"
)

type Handler interface {
	Setup(session sarama.ConsumerGroupSession) error
	Cleanup(session sarama.ConsumerGroupSession) error
	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}

type Handler_[T any] struct {
	fn func(msg *sarama.ConsumerMessage, t T) error
	l  logger.LoggerV1
}

func NewHandler_[T any](fn func(msg *sarama.ConsumerMessage, t T) error, l logger.LoggerV1) Handler {
	return &Handler_[T]{fn: fn, l: l}
}

func (h *Handler_[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler_[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

//把处理数据之前的重复逻辑全部封装起来

func (h *Handler_[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
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
		err = h.fn(msg, t)
		//在这里执行重试,比如简单的回调多次
		//在这里回调真正的处理逻辑，也就是我们只要写好具体的消费逻辑就好
		if err != nil {
			h.l.Error("处理消息失败", logger.Error(err),
				logger.String("topic", msg.Topic),
				logger.Int("partition", int(msg.Partition)),
				logger.Int("offset", int(msg.Offset)),
			)
			return err
		} else {
			session.MarkMessage(msg, "")
		}
	}

	return nil
}
