package bootstrap

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"os"
	"os/signal"
	"syscall"
	"telegram-monitor/pkg/controllers"
	"telegram-monitor/pkg/core/cst"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
	"telegram-monitor/pkg/routes"
)

func Telegram() {
	opts := []tg.ClientOption{tg.WithClientServerURL(cst.TelegramApi)}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()
	token := cst.TelegramToken
	if global.App.Config.App.Env == "ttkj" || global.App.Config.App.Env == "prod" {
		token = global.App.Config.Telegram.Token
	}
	global.App.Client = tg.New(token, opts...)
	me, err := global.App.Client.Me(ctx)
	if err != nil {
		logger.Error("authorized failed %v", err)
		os.Exit(2)
	}
	logger.Info("authorized %s successfully.", me.Username.Link())
	//telegram认证成功，启动cron任务
	StartCron()
	controllers.Update(token)
	r := tgb.NewRouter()
	routes.Telegram(r)
	err = tgb.NewPoller(r, global.App.Client).Run(ctx)
	if err != nil {
		os.Exit(2)
	}
}
