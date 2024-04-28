package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func InitZapLogger() *zap.SugaredLogger {
	logMode := zapcore.DebugLevel
	// TODO:根据环境来切换日志级别
	core := zapcore.NewCore(getEncoder(), getWriterSyncer(), logMode)
	return zap.New(core).Sugar()
}

func getEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "time"
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Local().Format(time.DateTime))
	}
	return zapcore.NewJSONEncoder(config)
}

func getWriterSyncer() zapcore.WriteSyncer {
	projectDir, err := os.Getwd()
	// TODO:根据配置读取日志存储位置
	targetDir := "log"
	if err != nil {
		panic(err)
	}
	separator := filepath.Separator
	logFileNameWithoutSuffix := strings.Join([]string{projectDir, targetDir, time.Now().Format(time.DateOnly)}, string(separator))
	logFileName := fmt.Sprintf("%s.txt", logFileNameWithoutSuffix)
	// TODO:根据配置修改日志分割配置
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    1,  // 文件切割大小
		MaxBackups: 10, // 至多保留多少文件
		MaxAge:     90, // 保留文件最大天数
		Compress:   true,
	}
	defer func() {
		_ = lumberjackLogger.Close()
	}()
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberjackLogger), zapcore.AddSync(os.Stdout))
}
