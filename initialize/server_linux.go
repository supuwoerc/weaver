package initialize

// 创建http服务器
func InitServer(handle http.Handler) {
	port := viper.GetInt("server.port")
	if port == 0 {
		port = defaultPort
	}
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: handle,
	}
	err := gracehttp.Serve(server)
	if err != nil {
		global.Logger.Error(fmt.Sprintf("服务启动失败：%s\n", err.Error()))
		os.Exit(1)
	}
	global.Logger.Info(fmt.Sprintf("服务关闭..."))
}
