package job

import (
	"time"
	"webook/pkg/logger"

	"github.com/robfig/cron/v3"
)

type CronJobBuilder struct {
	l logger.LoggerV1
	//这里还可以加入监控告警等
}
type fn func() error

func (f fn) Run() {
	_ = f()
}

func NewCronJobBuilder(l logger.LoggerV1) *CronJobBuilder {
	return &CronJobBuilder{l: l}
}

func (b *CronJobBuilder) Build(j Job) cron.Job {
	name := j.Name()
	return fn(func() error {
		//在这里可以进行监控或者告警以及日志信息的打印
		b.l.Info("cron job",
			logger.String("name", name),
			logger.String("start", time.Now().Format(time.DateTime)),
		)
		defer func() {
			b.l.Info("cron job",
				logger.String("name", name),
				logger.String("end", time.Now().Format(time.DateTime)),
			)
		}()
		err := j.Run()
		if err != nil {
			b.l.Error("cron job",
				logger.String("name", name),
				logger.Error(err),
			)
		}
		return nil
	})

}
