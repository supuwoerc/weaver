package initialize

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/supuwoerc/weaver/conf"
	weaverLogger "github.com/supuwoerc/weaver/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils"
	"gorm.io/plugin/opentelemetry/tracing"
)

const (
	tablePrefix  = "sys_"
	traceMessage = "gorm trace log message"
	position     = "position"
	mistaken     = "error"
	execution    = "execution"
	slow         = "slow-log"
	rows         = "rows"
	sql          = "sql"
)

type GormLogger struct {
	*weaverLogger.Logger
	Level                     logger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

func NewGormLogger(l *weaverLogger.Logger, conf *conf.Config) *GormLogger {
	return &GormLogger{
		Logger:                    l,
		Level:                     logger.LogLevel(conf.GORM.LogLevel),
		SlowThreshold:             conf.GORM.SlowThreshold,
		IgnoreRecordNotFoundError: conf.GORM.IgnoreRecordNotFoundError,
	}
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *g
	newLogger.Level = level
	return &newLogger
}

func (g *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if g.Level >= logger.Info {
		lineNum := utils.FileWithLineNum()
		g.WithContext(ctx).Infof(msg, append([]interface{}{lineNum}, data...)...)
	}
}

func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.Level >= logger.Warn {
		lineNum := utils.FileWithLineNum()
		g.WithContext(ctx).Warnf(msg, append([]interface{}{lineNum}, data...)...)
	}
}

func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.Level >= logger.Error {
		lineNum := utils.FileWithLineNum()
		g.WithContext(ctx).Errorf(msg, append([]interface{}{lineNum}, data...)...)
	}
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.Level <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	lineNum := utils.FileWithLineNum()
	cost := fmt.Sprintf("%.3f ms", float64(elapsed.Nanoseconds())/1e6)
	switch {
	case err != nil && g.Level >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !g.IgnoreRecordNotFoundError):
		sqlRaw, affected := fc()
		stackLevel := zapcore.ErrorLevel
		if errors.Is(err, gorm.ErrRecordNotFound) {
			stackLevel = zapcore.PanicLevel
		}
		logEntry := g.WithContext(ctx).WithOptions(zap.AddStacktrace(stackLevel))
		if affected == -1 {
			logEntry.Errorw(traceMessage, position, lineNum, mistaken, err, execution, cost, sql, g.cleanSQL(sqlRaw))
		} else {
			logEntry.Errorw(traceMessage, position, lineNum, mistaken, err, execution, cost, rows, affected, sql, g.cleanSQL(sqlRaw))
		}
	case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.Level >= logger.Warn:
		sqlRaw, affected := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", g.SlowThreshold)
		if affected == -1 {
			g.WithContext(ctx).Warnw(traceMessage, position, lineNum, slow, slowLog, execution, cost, sql, g.cleanSQL(sqlRaw))
		} else {
			g.WithContext(ctx).Warnw(traceMessage, position, lineNum, slow, slowLog, execution, cost, rows, affected, sql, g.cleanSQL(sqlRaw))
		}
	case g.Level == logger.Info:
		sqlRaw, affected := fc()
		if affected == -1 {
			g.WithContext(ctx).Infow(traceMessage, position, lineNum, execution, cost, sql, g.cleanSQL(sqlRaw))
		} else {
			g.WithContext(ctx).Infow(traceMessage, position, lineNum, execution, cost, rows, affected, sql, g.cleanSQL(sqlRaw))
		}
	}
}

func (g *GormLogger) cleanSQL(sql string) string {
	if sql == "" {
		return sql
	}
	// 替换换行符和多余空格
	sql = strings.ReplaceAll(sql, "\n", " ")
	sql = strings.ReplaceAll(sql, "\t", " ")
	// 清理多余的空格
	re := regexp.MustCompile(`\s+`)
	sql = re.ReplaceAllString(sql, " ")
	// 清理首尾空格
	sql = strings.TrimSpace(sql)
	return sql
}

func NewGORM(conf *conf.Config, l logger.Interface) *gorm.DB {
	logLevel := conf.GORM.LogLevel
	db, err := gorm.Open(mysql.Open(conf.GORM.DSN()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
		},
		Logger: l.LogMode(logger.LogLevel(logLevel)),
	})
	if err != nil {
		panic(err)
	}
	// tracing
	err = db.Use(tracing.NewPlugin(tracing.WithoutMetrics()))
	if err != nil {
		panic(err)
	}
	link, err := db.DB()
	if err != nil {
		panic(err)
	}
	maxIdleConn := conf.GORM.MaxIdleConn
	maxOpenConn := conf.GORM.MaxOpenConn
	maxLifetime := conf.GORM.MaxLifetime
	link.SetMaxIdleConns(maxIdleConn)
	link.SetMaxOpenConns(maxOpenConn)
	link.SetConnMaxLifetime(time.Minute * maxLifetime)
	return db
}
