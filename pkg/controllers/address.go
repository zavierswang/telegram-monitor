package controllers

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"go.uber.org/zap/buffer"
	"html/template"
	"strings"
	"telegram-monitor/pkg/core/cst"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
	"telegram-monitor/pkg/models"
)

func Address(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [address] controller", userId, username)

	var addrs []models.Address
	global.App.DB.Find(&addrs, "user_id = ?", userId)

	if len(addrs) == 0 {
		tmpl, err := template.ParseFiles(cst.AddAddressTemplateFile)
		if err != nil {
			logger.Error("[%s %s] template parse file %s, failed %v", userId, username, cst.AddAddressTemplateFile, err)
			return err
		}
		buf := new(buffer.Buffer)
		err = tmpl.Execute(buf, nil)
		if err != nil {
			logger.Error("[%s %s] template execute file %s, failed %v", userId, username, cst.AddAddressTemplateFile, err)
			return err
		}
		return update.Answer(buf.String()).ParseMode(tg.HTML).DisableWebPagePreview(true).DoVoid(ctx)
	}
	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](1).Row()
	for _, addr := range addrs {
		a := []rune(addr.Address)
		b := fmt.Sprintf("%s...%s", string(a[0:10]), string(a[25:]))
		if !addr.IsMonitor {
			layout.Insert(tg.NewInlineKeyboardButtonCallback(fmt.Sprintf("【%s】 %s  🔔 添加", addr.Remark, b), fmt.Sprintf("/add %s", addr.Address)))
		} else {
			layout.Insert(tg.NewInlineKeyboardButtonCallback(fmt.Sprintf("【%s】 %s  🔕 删除", addr.Remark, b), fmt.Sprintf("/remove %s", addr.Address)))
		}
	}
	layout.Insert(tg.NewInlineKeyboardButtonCallback("关闭", Callback.Close))
	ikb := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)
	tmpl, err := template.ParseFiles(cst.AddressTemplateFile)
	if err != nil {
		logger.Error("[%s %s] template parse file %s, failed %v", userId, username, cst.AddressTemplateFile, err)
		return err
	}
	buf := new(buffer.Buffer)
	err = tmpl.Execute(buf, nil)
	if err != nil {
		logger.Error("[%s %s] template execute file %s, failed %v", userId, username, cst.AddressTemplateFile, err)
		return err
	}
	return update.Answer(buf.String()).ParseMode(tg.HTML).ReplyMarkup(ikb).DisableWebPagePreview(true).DoVoid(ctx)
}
