package initialize

import (
	"gin-web/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

const TablePrefix = "sys_"

func NewGORM(conf *conf.Config) *gorm.DB {
	dsn := conf.Mysql.DSN
	logLevel := conf.Logger.GormLevel
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   TablePrefix,
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.LogLevel(logLevel)),
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
