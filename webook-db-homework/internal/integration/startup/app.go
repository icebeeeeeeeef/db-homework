package startup

import (
	"webook/interactive/events"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

// App 应用结构体
type App struct {
	Server    *gin.Engine
	Consumers []events.Consumer
	Cron      *cron.Cron
}
