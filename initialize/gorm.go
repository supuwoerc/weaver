package initialize

import (
	"fmt"
	"gin-web/pkg/global"
	"gin-web/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

const TablePrefix = "sys_"

func InitGORM() *gorm.DB {
	dsn := viper.GetString("mysql.dsn")
	logLevel := viper.Get("gorm.logLevel")
	level := logLevel.(int)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   TablePrefix,
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.LogLevel(level)),
	})
	if err != nil {
		global.Logger.Error(fmt.Sprintf("GORM初始化失败：%s\n", err.Error()))
		panic(err)
	}
	link, err := db.DB()
	if err != nil {
		global.Logger.Error(fmt.Sprintf("DB初始化失败：%s\n", err.Error()))
		panic(err)
	}
	maxIdleConn := viper.GetInt("mysql.maxIdleConn")
	maxOpenConn := viper.GetInt("mysql.maxOpenConn")
	maxLifetime := viper.GetDuration("mysql.maxLifetime")
	link.SetMaxIdleConns(maxIdleConn)
	link.SetMaxOpenConns(maxOpenConn)
	link.SetConnMaxLifetime(time.Minute * maxLifetime)
	autoMigrate(db)
	return db
}

func autoMigrate(db *gorm.DB) {
	_ = db.AutoMigrate(&dao.User{})
}
