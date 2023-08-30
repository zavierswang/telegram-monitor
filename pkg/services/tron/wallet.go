package tron

import (
	"fmt"
	"telegram-monitor/pkg/common/httpclient"
	"telegram-monitor/pkg/core/cst"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
)

type AccountWalletResp struct {
	Data []AccountWallet `json:"data"`
}

type AccountWallet struct {
	TokenAbbr       string `json:"token_abbr"`
	TokenValueInUsd string `json:"token_value_in_usd"`
	Balance         string `json:"balance"`
}

func GetAddressWallet(params map[string]string) (*AccountWallet, error) {
	url := fmt.Sprintf("%s%s", cst.TronScanApi, cst.TronWallet)
	headers := map[string]string{
		"TRON-PRO-API-KEY": global.App.Config.Telegram.TronScanApiKey,
	}
	var wallets AccountWalletResp
	err := httpclient.GetJson(url, params, headers, &wallets)
	if err != nil {
		logger.Error("[tron] GetAddressWallet failed %w", url, err)
		return nil, err
	}
	for _, wallet := range wallets.Data {
		if wallet.TokenAbbr == "USDT" {
			return &wallet, nil
		}
	}
	return nil, err
}
