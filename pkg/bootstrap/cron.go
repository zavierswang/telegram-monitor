package bootstrap

import (
	"context"
	"github.com/robfig/cron/v3"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
	"telegram-monitor/pkg/jobs"
)

func StartCron() {
	global.App.Cron = cron.New(cron.WithSeconds())
	go func() {
		listenTRC20 := &jobs.ListenTRC20{
			Ctx:    context.Background(),
			Ticker: make(map[string]int),
		}
		_, err := global.App.Cron.AddJob("*/30 * * * * *", listenTRC20)
		if err != nil {
			logger.Error("[cron] start ListenTRC20 cron job failed %v", err)
		}

		userValidity := &jobs.UserValidity{
			Ctx: context.Background(),
		}
		_, err = global.App.Cron.AddJob("*/30 * * * * *", userValidity)
		if err != nil {
			logger.Error("[cron] start UserValidity cron job failed %v", err)
		}
		global.App.Cron.Start()
		defer global.App.Cron.Stop()
		select {}
	}()
}
