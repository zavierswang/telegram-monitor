package main

import (
	"telegram-monitor/pkg/bootstrap"
	"telegram-monitor/pkg/core/cst"
)

func main() {
	// 初始配置文件
	bootstrap.LoadConfig(cst.AppName)
	// 初始化DB
	bootstrap.ConnectDB()
	// 初始化缓存
	bootstrap.NewCache()
	// 初始化MQ
	bootstrap.ConnectMQ()
	// Telegram Bot
	go bootstrap.Telegram()
	// 启动HTTP服务
	bootstrap.RunHTTP()
}
