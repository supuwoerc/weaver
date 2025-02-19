package initialize

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var LoggerSyncer zapcore.WriteSyncer

func InitZapLogger() *zap.SugaredLogger {
	level := viper.Get("logger.level")
	logMode := level.(int)
	LoggerSyncer = getWriterSyncer()
	core := zapcore.NewCore(getEncoder(), LoggerSyncer, zapcore.Level(logMode))
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
	write2Stdout := viper.GetBool("logger.stdout")
	targetDir := viper.GetString("logger.dir")
	maxSize := viper.GetInt("logger.maxSize")
	maxBackups := viper.GetInt("logger.maxBackups")
	maxAge := viper.GetInt("logger.maxAge")
	colorful := viper.GetBool("logger.colorful")
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
	if !colorful {
		ws = append(ws, zapcore.AddSync(lumberjackLogger))
	}
	if write2Stdout {
		ws = append(ws, zapcore.AddSync(os.Stdout))
	}
	if len(ws) == 0 {
		panic("缺少日志输出配置信息")
	}
	return zapcore.NewMultiWriteSyncer(ws...)
}
