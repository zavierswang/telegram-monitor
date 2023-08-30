package controllers

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"strings"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
	"telegram-monitor/pkg/middleware"
	"telegram-monitor/pkg/models"
)

func Remark(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [remark] controller", userId, username)

	sess := middleware.SessionManager.Get(ctx)
	remarkAddress := sess.RemarkAddress
	label := strings.TrimSpace(update.Text)
	logger.Info("[%s %s] %s => %s", userId, username, remarkAddress, label)
	var addr models.Address
	err := global.App.DB.First(&addr, "user_id = ? AND address = ?", userId, remarkAddress).Error
	if err != nil {
		logger.Warn("[%s %s] not found address with %s", userId, username, remarkAddress)
		return update.Answer("💣您输入的地址未找到，请确认备注已添加的地址~").DoVoid(ctx)
	}
	addr.Remark = label
	err = global.App.DB.Model(&models.Address{}).Where("user_id = ? AND address = ?", userId, remarkAddress).Update("remark", label).Error
	if err != nil {
		logger.Error("[%s %s] remark address failed %v", userId, username, err)
		return err
	}
	middleware.SessionManager.Reset(sess)
	return update.Answer("<b>添加备注成功🎉</b>").ParseMode(tg.HTML).DisableWebPagePreview(true).DoVoid(ctx)
}
