//go:build wireinject

package startup

import (
	"webook/interactive/events"
	"webook/interactive/grpc"
	"webook/interactive/repository"
	"webook/interactive/repository/cache"
	"webook/interactive/repository/dao"
	"webook/interactive/service"
	"webook/ioc"

	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(InitRedis, InitDB, InitLogger, ioc.InitKafkaClient, ioc.InitKafkaProducer)
var interactiveSvcProvider = wire.NewSet(
	dao.NewInteractiveDAO,
	cache.NewInteractiveCache,
	repository.NewInteractiveRepository,
	service.NewInteractiveService)

// InitConsumers 初始化消费者列表
func InitConsumers(consumer events.Consumer) []events.Consumer {
	return []events.Consumer{consumer}
}

func InitInteractiveService() service.InteractiveService {
	wire.Build(
		thirdProvider,
		interactiveSvcProvider,
	)
	return service.NewInteractiveService(nil)
}

func InitInteractiveGRPCServer() *grpc.InteractiveServiceServer {
	wire.Build(
		thirdProvider,
		interactiveSvcProvider,
		grpc.NewInteractiveServiceServer,
	)
	return new(grpc.InteractiveServiceServer)
}
