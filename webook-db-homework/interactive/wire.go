//go:build wireinject

package main

import (
	"webook/interactive/grpc"

	"webook/interactive/events"
	"webook/interactive/ioc"
	"webook/interactive/repository"
	"webook/interactive/repository/cache"
	"webook/interactive/repository/dao"
	"webook/interactive/service"

	"github.com/google/wire"
)

var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	cache.NewInteractiveCache,
	dao.NewInteractiveDAO,
	repository.NewInteractiveRepository,
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitLogger,
	ioc.InitRedis,
	ioc.InitKafka,
)

func InitApp() *App {
	wire.Build(
		interactiveSvcProvider,
		thirdPartySet,
		events.NewKafkaConsumer,
		grpc.NewInteractiveServiceServer,
		ioc.InitGRPCServer,
		ioc.NewConsumers,
		wire.Struct(new(App), "*"), //将App中的字段注入到wire中
	)
	return new(App)
}
