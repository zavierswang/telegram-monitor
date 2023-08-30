package jobs

import (
	"context"
	"encoding/json"
	mq "telegram-monitor/pkg/common/rabbit"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
	"telegram-monitor/pkg/models"
	"time"
)

type ListenTRC20 struct {
	Ticker     map[string]int
	TRC20Queue map[string][]string
	Ctx        context.Context
}

func (l *ListenTRC20) Run() {
	now := time.Now()
	var users []models.User
	err := global.App.DB.Find(&users, "expired_at > ? OR is_admin = ?", now, true).Error
	if err != nil {
		logger.Error("[scheduler] mysql query users failed %v", err)
		return
	}
	for _, user := range users {
		l.Ticker[user.UserID]++
		var addrs []models.Address
		global.App.DB.Find(&addrs, "user_id = ?", user.UserID)
		for _, addr := range addrs {
			buff, _ := json.Marshal(addr)
			//logger.Info("[scheduler] %s %s => %d times send message to rabbitmq", user.Username, addr.Address, l.Ticker[user.UserID])
			err = global.App.MQ.SendMessage(mq.Message{Body: buff})
			if err != nil {
				logger.Error("rabbit send delay message failed %v", err)
				continue
			}
		}
	}
	return
}

func (l *ListenTRC20) Consume() {

}
