package initialize

import (
	"context"
	"errors"
	"fmt"
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
		Level:                     logger.LogLevel(conf.Logger.GormLevel),
		SlowThreshold:             conf.Logger.GormSlowThreshold,
		IgnoreRecordNotFoundError: conf.Logger.GormIgnoreRecordNotFoundError,
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
			logEntry.Errorw(traceMessage, position, lineNum, mistaken, err, execution, cost, sql, sqlRaw)
		} else {
			logEntry.Errorw(traceMessage, position, lineNum, mistaken, err, execution, cost, rows, affected, sql, sqlRaw)
		}
	case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.Level >= logger.Warn:
		sqlRaw, affected := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", g.SlowThreshold)
		if affected == -1 {
			g.WithContext(ctx).Warnw(traceMessage, position, lineNum, slow, slowLog, execution, cost, sql, sqlRaw)
		} else {
			g.WithContext(ctx).Warnw(traceMessage, position, lineNum, slow, slowLog, execution, cost, rows, affected, sql, sqlRaw)
		}
	case g.Level == logger.Info:
		sqlRaw, affected := fc()
		if affected == -1 {
			g.WithContext(ctx).Infow(traceMessage, position, lineNum, execution, cost, sql, sqlRaw)
		} else {
			g.WithContext(ctx).Infow(traceMessage, position, lineNum, execution, cost, rows, affected, sql, sqlRaw)
		}
	}
}

func NewGORM(conf *conf.Config, l logger.Interface) *gorm.DB {
	dsn := conf.Mysql.DSN
	logLevel := conf.Logger.GormLevel
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
		},
		Logger: l.LogMode(logger.LogLevel(logLevel)),
	})
	if err != nil {
		panic(err)
	}
	link, err := db.DB()
	if err != nil {
		panic(err)
	}
	maxIdleConn := conf.Mysql.MaxIdleConn
	maxOpenConn := conf.Mysql.MaxOpenConn
	maxLifetime := conf.Mysql.MaxLifetime
	link.SetMaxIdleConns(maxIdleConn)
	link.SetMaxOpenConns(maxOpenConn)
	link.SetConnMaxLifetime(time.Minute * maxLifetime)
	return db
}
