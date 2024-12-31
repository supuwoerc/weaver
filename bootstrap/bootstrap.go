package bootstrap

import (
	"fmt"
	"gin-web/initialize"
	"gin-web/pkg/global"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
)

func Start() {
	initialize.InitConfig()
	global.Logger = initialize.InitZapLogger()
	global.DB = initialize.InitGORM()
	global.RedisClient = initialize.InitRedis()
	global.Localizer = initialize.InitI18N()
	writePid()
	initialize.InitServer(initialize.InitEngine(initialize.LoggerSyncer))
}

func Clean() {
	fmt.Println("关闭服务后的清理...")
}

func writePid() {
	path := viper.GetString("system.pid")
	if strings.TrimSpace(path) == "" {
		panic("system.pid is empty!")
	}
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(file)
	pid := strconv.Itoa(os.Getpid())
	_, err = file.WriteString(pid)
	if err != nil {
		panic(err)
	}
}
