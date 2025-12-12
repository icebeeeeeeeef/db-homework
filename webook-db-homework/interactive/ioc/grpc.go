package ioc

import (
	"webook/pkg/grpcx"

	grpc2 "webook/interactive/grpc"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCServer(intrSvc *grpc2.InteractiveServiceServer) *grpcx.Server {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	intrSvc.Register(server)
	return &grpcx.Server{
		Server: server,
		Addr:   cfg.Addr}
}
