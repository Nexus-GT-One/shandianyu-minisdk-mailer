package crontab

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"shandianyu-minisdk-mailer/crontab/task"
	"shandianyu-minisdk-mailer/provider/config"
	loggerFactory "shandianyu-minisdk-mailer/provider/logger"
	"shandianyu-minisdk-mailer/util/systemutil"
)

var logger = loggerFactory.GetLogger()

type Crontab struct{ cron *cron.Cron }

func NewCrontab() *Crontab {
	return &Crontab{cron: cron.New(cron.WithSeconds())}
}

func (c *Crontab) Start() {
	// only模式，方便开发，通过配置文件控制，只执行一次
	if len(config.GetString("crontab.only")) > 0 {
		crawler, _ := task.GetAllTask().Load(fmt.Sprintf("%v", config.GetString("crontab.only")))
		if crawler == nil {
			logger.Fatal("error only program name: %s", config.GetString("crontab.only"))
			return
		}

		(crawler.(task.Task)).TaskFunc(make(map[string]any))
		return
	}

	task.GetAllTask().Range(func(key, value any) bool {
		t := value.(task.Task)
		c.register(t.Crontab, func() { systemutil.Goroutine(func() { t.TaskFunc(make(map[string]any)) }) })
		logger.Info("[注册定时任务]\t%v\t%v", t.Crontab, key)
		return true
	})
	c.cron.Start()
}

func (c *Crontab) register(spec string, cmd func()) {
	_, err := c.cron.AddFunc(spec, cmd)
	if err != nil {
		logger.Fatal("%v", err)
	}
}
