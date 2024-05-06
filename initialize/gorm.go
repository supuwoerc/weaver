package initialize

import (
	"fmt"
	"gin-web/pkg/global"
	"gin-web/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func InitGORM() *gorm.DB {
	// TODO:从配置文件读取DSN相关配置信息
	dsn := "gin_web:gin_web@tcp(127.0.0.1:3306)/gin_web?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "sys_",
			SingularTable: true,
		},
		// TODO:区分环境设置日志级别
		Logger: logger.Default.LogMode(logger.Info),
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
	// TODO: 从配置文件中读取
	link.SetMaxIdleConns(3)
	link.SetMaxOpenConns(3)
	link.SetConnMaxLifetime(time.Hour)
	autoMigrate(db)
	return db
}

func autoMigrate(db *gorm.DB) {
	_ = db.AutoMigrate(&dao.User{})
}
