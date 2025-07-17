package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path"
	"runtime/debug"
	"shandianyu-minisdk-mailer/provider/config"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"strings"
	"sync"
	"time"
)

var callerMap sync.Map

type Logger struct {
	zapLogger *zap.Logger
}

func init() {
	cmd, _ := os.Getwd()
	log.SetFlags(log.Flags() | log.Lshortfile | log.Lmicroseconds)
	logrotate, _ := rotatelogs.New(
		path.Join(cmd, config.GetString("logger.path")+"/"+config.GetString("application.name")+"-%Y-%m-%d.log"),
		rotatelogs.WithMaxAge(30*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "file",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller: func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			stack := strings.Split(string(debug.Stack()), "\n")
			prefix := strings.TrimSpace(strings.ReplaceAll(stack[16], strings.ReplaceAll(cmd, "\\", "/"), ""))
			prefix = strings.TrimSpace(arrayutil.First(strings.Split(prefix, " ")))
			prefix = strings.ReplaceAll(prefix[1:], "/", ".")
			enc.AppendString(strings.TrimSpace(strings.ReplaceAll(prefix, "srv.shandianyu.shandianyu-minisdk-monitor.", "")))
		},
		EncodeName: zapcore.FullNameEncoder,
	}

	// 设置日志级别
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(logrotate), zapcore.AddSync(os.Stdout)}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writes...),
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)

	zap.ReplaceGlobals(zap.New(core, zap.AddCaller()))
}

func GetLogger() *Logger {
	return &Logger{zapLogger: zap.L()}
}

func (o *Logger) Debug(msg string, args ...any) {
	o.zapLogger.Debug(fmt.Sprintf(msg, args...))
}

func (o *Logger) Info(msg string, args ...any) {
	o.zapLogger.Info(fmt.Sprintf(msg, args...))
}

func (o *Logger) Warn(msg string, args ...any) {
	o.zapLogger.Warn(fmt.Sprintf(msg, args...))
}

func (o *Logger) Error(msg string, args ...any) {
	o.zapLogger.Error(fmt.Sprintf(msg, args...))
}

func (o *Logger) Panic(msg string, args ...any) {
	o.zapLogger.Panic(fmt.Sprintf(msg, args...))
}

func (o *Logger) Fatal(msg string, args ...any) {
	o.zapLogger.Fatal(fmt.Sprintf(msg, args...))
	os.Exit(0)
}
