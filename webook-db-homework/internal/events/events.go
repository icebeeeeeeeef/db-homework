package events

import articleEvents "webook/internal/events/article"

// 导出 article 包中的类型
type Producer = articleEvents.Producer

// Consumer 和 BatchConsumer 实现在 interactive 模块，不在 internal 下
// 保留 ReadEvent 以便主服务产出阅读事件
type ReadEvent = articleEvents.ReadEvent

// 导出构造函数
var NewKafkaProducer = articleEvents.NewKafkaProducer

// 注意：Consumer 和 BatchConsumer 的构造函数在 interactive 模块，不在此处重导出
