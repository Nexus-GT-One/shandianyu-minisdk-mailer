package commandutil

import (
	"os/exec"
	loggerFactory "shandianyu-minisdk-monitor/provider/logger"
)

var logger = loggerFactory.GetLogger()

func Run(command string, args ...string) string {
	cmd := exec.Command(command, args...)

	// 获取输出对象，可以从该对象中读取输出结果
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("%v", err)
		return ""
	}

	logger.Info(string(bytes))
	return string(bytes)
}
