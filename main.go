package main

import (
	"context"
	_ "embed"
	"shandianyu-minisdk-mailer/crontab"
	"shandianyu-minisdk-mailer/provider/config"
	loggerFactory "shandianyu-minisdk-mailer/provider/logger"
)

//go:embed banner
var banner string

func main() {
	logger := loggerFactory.GetLogger()
	logger.Info(banner)
	crontab.NewCrontab().Start()
	logger.Info("%s started", config.GetString("application.name"))
	<-context.Background().Done()
}
