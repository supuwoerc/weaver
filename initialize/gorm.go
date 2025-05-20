package initialize

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/supuwoerc/weaver/conf"
	weaverLogger "github.com/supuwoerc/weaver/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils"
)

const TablePrefix = "sys_"

const (
	infoStr      = "%s\n[info] "
	warnStr      = "%s\n[warn] "
	errStr       = "%s\n[error] "
	traceStr     = "%s\n[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
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
		g.WithContext(ctx).Infof(infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.Level >= logger.Warn {
		g.WithContext(ctx).Warnf(warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.Level >= logger.Error {
		g.WithContext(ctx).Warnf(errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.Level <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && g.Level >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !g.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			g.WithContext(ctx).Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.WithContext(ctx).Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.Level >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", g.SlowThreshold)
		if rows == -1 {
			g.WithContext(ctx).Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.WithContext(ctx).Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case g.Level == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			g.WithContext(ctx).Infof(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.WithContext(ctx).Infof(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func NewGORM(conf *conf.Config, l logger.Interface) *gorm.DB {
	dsn := conf.Mysql.DSN
	logLevel := conf.Logger.GormLevel
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   TablePrefix,
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
