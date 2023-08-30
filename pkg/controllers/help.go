package controllers

import (
	"context"
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

func Help(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [help] controller", userId, username)
	tmpl, err := template.ParseFiles(cst.HelpTemplateFile)
	if err != nil {
		logger.Error("[%s %s] template parse file %s, failed %v", userId, username, cst.HelpTemplateFile, err)
		return err
	}
	var user models.User
	err = global.App.DB.First(&user, "user_id = ?", userId).Error
	if err != nil {
		logger.Error("[%s %s] illegl user", userId, username)
		return err
	}
	buf := new(buffer.Buffer)
	err = tmpl.Execute(buf, user)
	if err != nil {
		logger.Error("[%s %s] template execute file %s, failed %v", userId, username, cst.HelpTemplateFile, err)
		return err
	}
	err = update.Answer(buf.String()).ParseMode(tg.HTML).DisableWebPagePreview(true).DoVoid(ctx)
	if err != nil {
		logger.Error("[%s %s] update answer failed %v", userId, username, err)
	}
	return err
}
