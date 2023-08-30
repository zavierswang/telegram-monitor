package bootstrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"net/http"
	"strconv"
	mq "telegram-monitor/pkg/common/rabbit"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
	"telegram-monitor/pkg/models"
	"telegram-monitor/pkg/services/tron"
	"time"
)

func ConnectMQ() {
	var err error
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d%s", global.App.Config.MQ.Username, global.App.Config.MQ.Password, global.App.Config.MQ.Host, global.App.Config.MQ.Port, global.App.Config.MQ.VHost)
	global.App.MQ, err = mq.NewRabbitMQ("telegram", dsn, "monitor", "usdt")
	if err != nil {
		logger.Error("new rabbitmq queue failed %v", err)
	}
	logger.Info("connect rabbitmq %s successfully", dsn)
}

func StartConsume() {
	err := global.App.MQ.Consume(func(delivery amqp.Delivery) {
		var addr models.Address
		_ = json.Unmarshal(delivery.Body, &addr)
		key := fmt.Sprintf("%s_%s", addr.Username, addr.Address)
		now := time.Now()
		params := map[string]string{
			"limit":            "50",
			"start":            "0",
			"sort":             "-timestamp",
			"count":            "true",
			"filterTokenValue": "0",
			"start_timestamp":  strconv.FormatInt(now.Add(time.Second*-30).UnixMilli(), 10),
			"end_timestamp":    strconv.FormatInt(now.UnixMilli(), 10),
			"relatedAddress":   addr.Address,
		}
		headers := map[string]string{
			"TRON-PRO-API-KEY": global.App.Config.Telegram.TronScanApiKey,
		}
		transfers, err := tron.TRC20Transfer(params, headers, false, false)
		if err != nil {
			logger.Error("[consume] tron.TRC20Transfer => %s %s failed %v", addr.Username, addr.Address, err)
			return
		}
		if len(transfers) == 0 {
			logger.Info("[consume] %s not found latest txIds", key)
			return
		}
		notifier := tron.Notifier{
			ID:        addr.ID,
			UserId:    addr.UserID,
			Address:   addr.Address,
			Transfers: transfers,
		}
		logger.Info("[consume] %s found latest txIds and length equal %d", key, len(notifier.Transfers))
		buf, _ := json.Marshal(&notifier)
		resp, err := http.Post(global.App.Config.Telegram.Callback, "application/json", bytes.NewBuffer(buf))
		if err != nil {
			logger.Error("[consume] %s callback trc20 api failed %v", key, err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			logger.Error("[consume] %s callback notifier errCode %d", key, resp.StatusCode)
			return
		}
	})
	if err != nil {
		logger.Error("[consume] delivery failed %v", err)
	}
	select {}
}
