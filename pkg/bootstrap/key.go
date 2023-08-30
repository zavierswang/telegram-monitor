package bootstrap

import (
	"os"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
	"telegram-monitor/pkg/services/grid"
)

func StorePrivateKey() {
	flag := false
	accounts := grid.DescribeLocalAccounts()
	for key, ks := range accounts {
		for _, k := range ks {
			logger.Info("ks: %+v", k)
			if k.Address.String() == global.App.Config.Telegram.SendAddress {
				flag = true
				logger.Info("%s => %s has been stored", key, k.Address.String())
				return
			}
		}
	}
	if !flag {
		err := grid.ImportPrivateKey(global.App.Config.Telegram.PrivateKey, global.App.Config.Telegram.AliasKey)
		if err != nil {
			flag = false
			logger.Error("import private key failed %v", err)
		}
		return
	}
	if !flag {
		os.Exit(-1)
	}
}
