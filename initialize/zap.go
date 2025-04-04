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

func NewZapLogger(v *viper.Viper, sync zapcore.WriteSyncer) *zap.SugaredLogger {
	level := v.Get("logger.level")
	logMode := level.(int)
	core := zapcore.NewCore(getEncoder(), sync, zapcore.Level(logMode))
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
}
func NewWriterSyncer(v *viper.Viper) zapcore.WriteSyncer {
	projectDir, err := os.Getwd()
	write2Stdout := v.GetBool("logger.stdout")
	targetDir := v.GetString("logger.dir")
	maxSize := v.GetInt("logger.maxSize")
	maxBackups := v.GetInt("logger.maxBackups")
	maxAge := v.GetInt("logger.maxAge")
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
		panic("缺少日志输出配置信息")
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
