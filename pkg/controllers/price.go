package controllers

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"go.uber.org/zap/buffer"
	"html/template"
	"strings"
	"telegram-monitor/pkg/core/cst"
	"telegram-monitor/pkg/core/logger"
	"telegram-monitor/pkg/services"
)

func Price(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [otc] controller", userId, username)

	resp, err := services.GetOkxTradingOrders()
	if err != nil {
		logger.Error("[%s %s] trading orders failed %v", userId, username, err)
		return err
	}
	//logger.Info("[%s %s] okx tradig orders response: %+v", userId, username, resp)
	data := resp.Data.Sell
	if len(resp.Data.Sell) >= 10 {
		data = resp.Data.Sell[0:9]
	}
	buf := new(buffer.Buffer)
	tmpl, err := template.ParseFiles(cst.PriceTemplateFile)
	if err != nil {
		logger.Error("[%s %s] template parse file %s, failed %v", userId, username, cst.PriceTemplateFile, err)
		return err
	}
	err = tmpl.Execute(buf, &data)
	if err != nil {
		logger.Error("[%s %s] template execute file %s, failed %v", userId, username, cst.PriceTemplateFile, err)
		return err
	}
	layout := tg.NewButtonLayout(1, tg.NewInlineKeyboardButtonCallback("关闭", Callback.Close))
	ikb := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)
	return update.Answer(buf.String()).ParseMode(tg.HTML).ReplyMarkup(ikb).DoVoid(ctx)
}
