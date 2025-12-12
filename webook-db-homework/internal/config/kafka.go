package config

import (
	"time"

	"github.com/IBM/sarama"
)

// KafkaConfig Kafka 配置
type KafkaConfig struct {
	Addrs    []string          `yaml:"addrs"`    // Kafka 服务器地址列表
	Topics   map[string]string `yaml:"topics"`   // Topic 名称映射
	Consumer ConsumerConfig    `yaml:"consumer"` // 消费者配置
	Producer ProducerConfig    `yaml:"producer"` // 生产者配置
	Net      NetConfig         `yaml:"net"`      // 网络配置
}

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
	GroupID       string `yaml:"group_id"`       // 消费者组ID
	OffsetInitial string `yaml:"offset_initial"` // 初始偏移量: oldest, newest
}

// ProducerConfig 生产者配置
type ProducerConfig struct {
	ReturnSuccesses bool   `yaml:"return_successes"` // 是否返回成功确认
	RetryMax        int    `yaml:"retry_max"`        // 最大重试次数
	RetryBackoff    string `yaml:"retry_backoff"`    // 重试间隔
}

// NetConfig 网络配置
type NetConfig struct {
	DialTimeout  string `yaml:"dial_timeout"`  // 连接超时
	ReadTimeout  string `yaml:"read_timeout"`  // 读取超时
	WriteTimeout string `yaml:"write_timeout"` // 写入超时
}

// DefaultKafkaConfig 默认 Kafka 配置
func DefaultKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		Addrs: []string{"localhost:9094"},
		Topics: map[string]string{
			"read_event":    "read_event",
			"article_event": "article_event",
		},
		Consumer: ConsumerConfig{
			GroupID:       "webook_consumer_group",
			OffsetInitial: "oldest",
		},
		Producer: ProducerConfig{
			ReturnSuccesses: true,
			RetryMax:        3,
			RetryBackoff:    "100ms",
		},
		Net: NetConfig{
			DialTimeout:  "5s",
			ReadTimeout:  "5s",
			WriteTimeout: "5s",
		},
	}
}

// GetTopic 获取指定类型的 Topic 名称
func (c *KafkaConfig) GetTopic(topicType string) string {
	if topic, exists := c.Topics[topicType]; exists {
		return topic
	}
	return topicType // 如果没找到，返回类型名作为默认值
}

// ToSaramaConfig 转换为 Sarama 配置
func (c *KafkaConfig) ToSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()

	// 生产者配置
	config.Producer.Return.Successes = c.Producer.ReturnSuccesses
	config.Producer.Retry.Max = c.Producer.RetryMax

	// 解析重试间隔
	if backoff, err := time.ParseDuration(c.Producer.RetryBackoff); err == nil {
		config.Producer.Retry.Backoff = backoff
	}

	// 消费者配置
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	if c.Consumer.OffsetInitial == "oldest" {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	// 网络配置
	if dialTimeout, err := time.ParseDuration(c.Net.DialTimeout); err == nil {
		config.Net.DialTimeout = dialTimeout
	}
	if readTimeout, err := time.ParseDuration(c.Net.ReadTimeout); err == nil {
		config.Net.ReadTimeout = readTimeout
	}
	if writeTimeout, err := time.ParseDuration(c.Net.WriteTimeout); err == nil {
		config.Net.WriteTimeout = writeTimeout
	}

	return config
}
