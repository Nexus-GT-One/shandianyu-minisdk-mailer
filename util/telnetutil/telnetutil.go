package telnetutil

import (
	"fmt"
	"net"
	"os/exec"
	loggerFactory "shandianyu-minisdk-monitor/provider/logger"
	"strconv"
	"strings"
	"time"
)

var logger = loggerFactory.GetLogger()

type TelnetResult struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

func TelnetTCP(ip string, port int) TelnetResult {
	timeout := time.Second * 5
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		logger.Error("connect %s:%d fail: %s", ip, port, err.Error())
		return TelnetResult{Result: false, Message: err.Error()}
	}
	logger.Info("connect %s:%d success!", ip, port)
	return TelnetResult{Result: true}
}

func TelnetUDP(ip string, port int) TelnetResult {
	cmd := exec.Command("nc", "-zv", "-u", ip, strconv.Itoa(port))

	// 获取命令输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("connect %s:%d fail: %s", ip, port, err.Error())
		return TelnetResult{Result: false, Message: err.Error()}
	}

	// 判断是否连接成功
	if !strings.Contains(string(output), "open") {
		logger.Error("connect %s:%d fail: %s", ip, port, output)
		return TelnetResult{Result: false, Message: string(output)}
	}

	logger.Info("connect %s:%d success!", ip, port)
	return TelnetResult{Result: true}
}
