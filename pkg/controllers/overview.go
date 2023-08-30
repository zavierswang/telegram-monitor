package controllers

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"strings"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
	"telegram-monitor/pkg/models"
)

func Overview(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [overview] controller", userId, username)

	var addrs []models.Address
	global.App.DB.Find(&addrs, "user_id = ?", userId)
	if len(addrs) == 0 {
		return update.Answer("<b>您还未添加钱包地址, 请输入TRC20地址👇</b>").ParseMode(tg.HTML).DoVoid(ctx)
	}
	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](1).Row()
	for _, addr := range addrs {
		a := []rune(addr.Address)
		b := fmt.Sprintf("%s...%s", string(a[0:6]), string(a[30:]))
		display := fmt.Sprintf("【%s】%s 查看流水👈", addr.Remark, b)
		layout.Insert(
			tg.NewInlineKeyboardButtonCallback(display, fmt.Sprintf("/statistics %s", addr.Address)),
		)
	}
	layout.Insert(tg.NewInlineKeyboardButtonCallback("关闭", Callback.Close))
	inlineKeyboard := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)

	return update.Answer("<b>您的钱包地址列表</b>").ParseMode(tg.HTML).ReplyMarkup(inlineKeyboard).DoVoid(ctx)
}
