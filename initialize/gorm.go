package initialize

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

const TablePrefix = "sys_"

func NewGORM(v *viper.Viper) *gorm.DB {
	dsn := v.GetString("mysql.dsn")
	logLevel := v.Get("gorm.logLevel")
	level := logLevel.(int)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   TablePrefix,
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.LogLevel(level)),
	})
	if err != nil {
		panic(err)
	}
	link, err := db.DB()
	if err != nil {
		panic(err)
	}
	maxIdleConn := v.GetInt("mysql.maxIdleConn")
	maxOpenConn := v.GetInt("mysql.maxOpenConn")
	maxLifetime := v.GetDuration("mysql.maxLifetime")
	link.SetMaxIdleConns(maxIdleConn)
	link.SetMaxOpenConns(maxOpenConn)
	link.SetConnMaxLifetime(time.Minute * maxLifetime)
	return db
}
