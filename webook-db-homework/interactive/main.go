package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	initViperV1()
	app := InitApp()
	for _, consumer := range app.consumers {
		err := consumer.Start()
		if err != nil {
			panic(err)
		}

	}
	err := app.Server.Serve()
	if err != nil {
		panic(err)
	}

}

func initViperV1() {
	cfile := pflag.String("config", "./config/dev.yaml", "指定文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件发生变化:", e.Name)
	})
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
