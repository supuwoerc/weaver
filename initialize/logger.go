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

func InitZapLogger() *zap.SugaredLogger {
	level := viper.Get("logger.level")
	logMode := level.(int)
	core := zapcore.NewCore(getEncoder(), getWriterSyncer(), zapcore.Level(logMode))
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
	targetDir := viper.GetString("logger.dir")
	maxSize := viper.GetInt("logger.maxSize")
	maxBackups := viper.GetInt("logger.maxBackups")
	maxAge := viper.GetInt("logger.maxAge")
	if err != nil {
		panic(err)
	}
	separator := filepath.Separator
	logFileNameWithoutSuffix := strings.Join([]string{projectDir, targetDir, time.Now().Format(time.DateOnly)}, string(separator))
	logFileName := fmt.Sprintf("%s.txt", logFileNameWithoutSuffix)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    maxSize,    // 文件切割大小
		MaxBackups: maxBackups, // 至多保留多少文件
		MaxAge:     maxAge,     // 保留文件最大天数
		Compress:   true,
	}
	defer func() {
		_ = lumberjackLogger.Close()
	}()
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberjackLogger), zapcore.AddSync(os.Stdout))
}
