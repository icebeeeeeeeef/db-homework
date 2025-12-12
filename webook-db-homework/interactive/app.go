package main

import (
	"webook/pkg/grpcx"
	"webook/pkg/saramax"
)

type App struct {
	Server    *grpcx.Server
	consumers []saramax.Consumer
}
