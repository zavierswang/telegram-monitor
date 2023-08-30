package jobs

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"strconv"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/models"
	"time"
)

type UserValidity struct {
	Ticker int
	Ctx    context.Context
}

func (u *UserValidity) Run() {
	u.Ticker++
	var users []models.User
	global.App.DB.Find(&users, "is_admin = ?", false)
	now := time.Now()
	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](1).Row(
		tg.NewInlineKeyboardButtonURL("联系客服", fmt.Sprintf("https://t.me/%s", global.App.Config.App.Support)),
	)
	ikb := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)
	//logger.Info("[cron user] %d times", u.Ticker)
	for _, user := range users {
		expiredAt := user.ExpiredAt.UnixMilli()
		chatId, _ := strconv.ParseInt(user.UserID, 10, 64)
		if now.UnixMilli() > expiredAt && u.Ticker%100 == 99 {
			//该用户已经过期
			_ = global.App.Client.SendMessage(tg.ChatID(chatId), "您的帐户已经过期，请联系客服~").
				ParseMode(tg.HTML).
				ReplyMarkup(ikb).
				DoVoid(u.Ctx)
		} else if now.AddDate(0, 0, 2).UnixMilli() > expiredAt &&
			now.UnixMilli() < expiredAt &&
			u.Ticker%30 == 1 {
			//该用户即将过期
			_ = global.App.Client.SendMessage(tg.ChatID(chatId), "<b>🚨预警通知</b>\n\n您的帐户即将过期，请即时联系客服~").
				ParseMode(tg.HTML).
				ReplyMarkup(ikb).
				DoVoid(u.Ctx)
		}
	}
	return
}
