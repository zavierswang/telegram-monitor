package tron

import (
	"errors"
	"fmt"
	"strconv"
	"telegram-monitor/pkg/common/httpclient"
	"telegram-monitor/pkg/core/cst"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
	"time"
)

type AddressInfo struct {
	Address             string                   `json:"address"`
	Balance             int64                    `json:"balance"`
	TransactionsIn      int64                    `json:"transactions_in"`       //转入笔数
	TransactionsOut     int64                    `json:"transactions_out"`      //转出笔数
	Transactions        int64                    `json:"transactions"`          //总笔数
	CreateTime          int64                    `json:"date_created"`          //首次时间
	LatestOperationTime int64                    `json:"latest_operation_time"` //最新时间
	WithPriceTokens     []map[string]interface{} `json:"withPriceTokens,omitempty"`
}

type AddressDetailInfo struct {
	Address    string
	TokenType  string
	Balance    float64
	TotalCount int64
	InCount    int64
	OutCount   int64
	CreateTime string
	LatestTime string
}

func GetAddressAccountInfo(params map[string]string) (*AddressDetailInfo, error) {
	var body AddressInfo
	uri := fmt.Sprintf("%s%s", cst.TronScanApi, cst.TronAccountV2)
	headers := map[string]string{
		"TRON-PRO-API-KEY": global.App.Config.Telegram.TronScanApiKey,
	}
	err := httpclient.GetJson(uri, params, headers, &body)
	if err != nil {
		logger.Error("[tron] GetAddressAccountInfo request failed %v", uri, err)
		return nil, err
	}

	var data AddressDetailInfo
	data.Address = body.Address
	data.TotalCount = body.Transactions
	data.InCount = body.TransactionsIn
	data.OutCount = body.TransactionsOut
	data.CreateTime = time.UnixMilli(body.CreateTime).Format(cst.DateTimeFormatter)
	data.LatestTime = time.UnixMilli(body.LatestOperationTime).Format(cst.DateTimeFormatter)
	data.TokenType = "USDT"
	walletParam := map[string]string{
		"address":    params["address"],
		"asset_type": "0",
	}
	wallet, err := GetAddressWallet(walletParam)
	if err != nil {
		logger.Error("[tron] GetAccountWallet failed %v", err)
		return nil, err
	}
	//logger.Info("wallet: %+v", wallet)
	if wallet == nil {
		logger.Warn("[tron] wallet is nil")
		return &data, errors.New("wallet is nil")
	}
	balanceFloat, err := strconv.ParseFloat(wallet.Balance, 64)
	if err != nil {
		return nil, err
	}
	data.Balance, err = strconv.ParseFloat(fmt.Sprintf("%.3f", balanceFloat), 64)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
