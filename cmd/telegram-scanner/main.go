package main

import "telegram-monitor/pkg/bootstrap"

const (
	AppName = "telegram-scanner"
)

func main() {
	// 加载配置文件
	bootstrap.LoadConfig(AppName)
	// 初始化MQ
	bootstrap.ConnectMQ()
	// MQ Consume
	bootstrap.StartConsume()
}
