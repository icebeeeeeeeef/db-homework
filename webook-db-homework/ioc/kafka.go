package ioc

import (
	"webook/interactive/events"
	"webook/interactive/repository"
	"webook/internal/config"
	events2 "webook/internal/events"
	"webook/pkg/logger"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

// InitKafkaProducer 初始化 Kafka 生产者
func InitKafkaProducer(l logger.LoggerV1) sarama.SyncProducer {
	cfg := getKafkaConfig()
	saramaConfig := cfg.ToSaramaConfig()

	producer, err := sarama.NewSyncProducer(cfg.Addrs, saramaConfig)
	if err != nil {
		l.Error("初始化 Kafka 生产者失败", logger.Field{
			Key:   "error",
			Value: err,
		})
		panic(err)
	}

	l.Info("Kafka 生产者初始化成功", logger.Field{
		Key:   "addrs",
		Value: cfg.Addrs,
	})

	return producer
}

// InitKafkaClient 初始化 Kafka 客户端
func InitKafkaClient(l logger.LoggerV1) sarama.Client {
	cfg := getKafkaConfig()
	saramaConfig := cfg.ToSaramaConfig()

	client, err := sarama.NewClient(cfg.Addrs, saramaConfig)
	if err != nil {
		l.Error("初始化 Kafka 客户端失败", logger.Field{
			Key:   "error",
			Value: err,
		})
		panic(err)
	}

	l.Info("Kafka 客户端初始化成功", logger.Field{
		Key:   "addrs",
		Value: cfg.Addrs,
	})

	return client
}

// InitKafkaConsumer 初始化 Kafka 消费者
func InitKafkaConsumer(l logger.LoggerV1) sarama.ConsumerGroup {
	cfg := getKafkaConfig()
	saramaConfig := cfg.ToSaramaConfig()

	consumer, err := sarama.NewConsumerGroup(cfg.Addrs, cfg.Consumer.GroupID, saramaConfig)
	if err != nil {
		l.Error("初始化 Kafka 消费者失败", logger.Field{
			Key:   "error",
			Value: err,
		})
		panic(err)
	}

	l.Info("Kafka 消费者初始化成功", logger.Field{
		Key:   "addrs",
		Value: cfg.Addrs,
	}, logger.Field{
		Key:   "group_id",
		Value: cfg.Consumer.GroupID,
	})

	return consumer
}

// InitArticleEventProducer 初始化文章事件生产者
func InitArticleEventProducer(producer sarama.SyncProducer) events2.Producer {
	return events2.NewKafkaProducer(producer)
}

// InitArticleEventConsumer 初始化文章事件消费者
func InitArticleEventConsumer(client sarama.Client, l logger.LoggerV1, repo repository.InteractiveRepository) events.Consumer {
	return events.NewKafkaBatchConsumer(client, l, repo)
}

// getKafkaConfig 从配置中获取 Kafka 配置
func getKafkaConfig() *config.KafkaConfig {
	cfg := config.DefaultKafkaConfig()

	// 从 viper 读取配置
	if addrs := viper.GetStringSlice("kafka.addrs"); len(addrs) > 0 {
		cfg.Addrs = addrs
	}

	if groupID := viper.GetString("kafka.consumer.group_id"); groupID != "" {
		cfg.Consumer.GroupID = groupID
	}

	if offsetInitial := viper.GetString("kafka.consumer.offset_initial"); offsetInitial != "" {
		cfg.Consumer.OffsetInitial = offsetInitial
	}

	if returnSuccesses := viper.GetBool("kafka.producer.return_successes"); viper.IsSet("kafka.producer.return_successes") {
		cfg.Producer.ReturnSuccesses = returnSuccesses
	}

	if retryMax := viper.GetInt("kafka.producer.retry_max"); retryMax > 0 {
		cfg.Producer.RetryMax = retryMax
	}

	if retryBackoff := viper.GetString("kafka.producer.retry_backoff"); retryBackoff != "" {
		cfg.Producer.RetryBackoff = retryBackoff
	}

	if dialTimeout := viper.GetString("kafka.net.dial_timeout"); dialTimeout != "" {
		cfg.Net.DialTimeout = dialTimeout
	}

	if readTimeout := viper.GetString("kafka.net.read_timeout"); readTimeout != "" {
		cfg.Net.ReadTimeout = readTimeout
	}

	if writeTimeout := viper.GetString("kafka.net.write_timeout"); writeTimeout != "" {
		cfg.Net.WriteTimeout = writeTimeout
	}

	// 读取 topics 配置
	if topics := viper.GetStringMapString("kafka.topics"); len(topics) > 0 {
		cfg.Topics = topics
	}

	return cfg
}
