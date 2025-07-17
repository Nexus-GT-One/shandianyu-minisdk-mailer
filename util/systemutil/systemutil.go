package systemutil

import (
	"net"
	"os"
	"runtime/debug"
	"shandianyu-minisdk-mailer/provider/logger"
	"strings"
)

/**
 * 模拟三元表达式，简化代码
 * @param condition 条件
 * @param trueVal 条件为真时的值
 * @param falseVal 条件为假时的值
 */
func If(condition bool, trueVal any, falseVal any) any {
	if condition {
		return trueVal
	}
	return falseVal
}

/**
 * 封装一下协程，保证协程在出现异常的情况下程序不崩溃
 * @param businessFunc 处理业务的方法 (不可为空)
 * @param catchFunc 异常处理的方法 (可为空)
 * @param finallyFunc 无论异常与否都会最终执行的方法 (可为空)
 */
func Goroutine(funs ...func()) {
	if len(funs) <= 0 || len(funs) > 3 {
		panic("参数错误")
		return
	}

	// 未传方法进来的，给一个空方法
	for i := len(funs); i < 3; i++ {
		funs = append(funs, func() {})
	}

	businessFunc := funs[0]
	catchFunc := funs[1]
	finallyFunc := funs[2]
	go func(f func()) {
		defer func() {
			if e := recover(); e != nil {
				logger.GetLogger().Info("协程捕获异常: %v \n 协程堆栈信息：\n %v", e, string(debug.Stack()))
				catchFunc()
				finallyFunc()
			}
		}()
		businessFunc()
	}(businessFunc)
	finallyFunc()
}

// 获取本机IP
func GetLocalIp() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		logger.GetLogger().Info("net.Interfaces failed, err: %v", err.Error())
		return ""
	}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}

	return ""
}

// 获取主机名
func GetHostName() string {
	hostname, _ := os.Hostname()
	return strings.TrimSpace(hostname)
}

// 判断当前条件为false的时候抛出异常
//
// condition-判断条件; err-要抛出的异常
func AssertAndThrowError(condition bool, err error) {
	if !condition {
		return
	}

	panic(err)
}
