package initialize

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/supuwoerc/weaver/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewZapLogger(conf *conf.Config, sync zapcore.WriteSyncer) *zap.SugaredLogger {
	core := zapcore.NewCore(getEncoder(), sync, zapcore.Level(conf.Logger.Level))
	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Fields(
			zap.String("env", conf.Env),
			zap.Int("pid", os.Getpid()),
		),
	)
	return logger.Sugar()
}
func NewWriterSyncer(conf *conf.Config) zapcore.WriteSyncer {
	projectDir, err := os.Getwd()
	write2Stdout := conf.Logger.Stdout
	targetDir := conf.Logger.Dir
	maxSize := conf.Logger.MaxSize
	maxBackups := conf.Logger.MaxBackups
	maxAge := conf.Logger.MaxAge
	if err != nil {
		panic(err)
	}
	logFileNameWithoutSuffix := filepath.Join(targetDir, time.Now().Format(time.DateOnly))
	if strings.TrimSpace(targetDir) == "" {
		logFileNameWithoutSuffix = filepath.Join(projectDir, "log", time.Now().Format(time.DateOnly))
	}
	logFileName := fmt.Sprintf("%s.txt", logFileNameWithoutSuffix)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    maxSize,    // 文件切割大小
		MaxBackups: maxBackups, // 至多保留多少文件
		MaxAge:     maxAge,     // 保留文件最大天数
		Compress:   true,
	}
	var ws = make([]zapcore.WriteSyncer, 0)
	ws = append(ws, zapcore.AddSync(lumberjackLogger))
	if write2Stdout {
		ws = append(ws, zapcore.AddSync(os.Stdout))
	}
	if len(ws) == 0 {
		panic("missing log output configuration")
	}
	return zapcore.NewMultiWriteSyncer(ws...)
}

func getEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "time"
	config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Local().Format(time.DateTime))
	}
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(config)
}
