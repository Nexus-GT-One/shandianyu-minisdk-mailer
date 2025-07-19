package task

import (
	"reflect"
	"runtime/debug"
	"shandianyu-minisdk-mailer/provider/config"
	loggerFactory "shandianyu-minisdk-mailer/provider/logger"
	"shandianyu-minisdk-mailer/provider/mongodb"
	"strings"
	"sync"
)

var taskMap sync.Map
var logger = loggerFactory.GetLogger()
var gameDb = mongodb.GetInstance()
var logDb = mongodb.GetLoggingInstance()
var isProd = reflect.DeepEqual("prod", config.GetString("env"))

type Task struct {
	Crontab  string
	TaskFunc func(param map[string]any)
}

func registerTask(crontab string, taskFunc func(param map[string]any)) {
	programKey := strings.ReplaceAll(getProgramFileName(), ".go", "")
	taskMap.Store(programKey, Task{Crontab: crontab, TaskFunc: taskFunc})
}

func getProgramFileName() string {
	programStack := strings.Split(string(debug.Stack()), "\n")
	programPath := strings.TrimSpace(programStack[len(programStack)-2])
	programPath = strings.TrimSpace(strings.Split(programPath, " ")[0])
	programPath = strings.TrimSpace(programPath[strings.LastIndex(programPath, "/")+1 : strings.LastIndex(programPath, ":")])
	return programPath
}

func GetAllTask() *sync.Map {
	return &taskMap
}
