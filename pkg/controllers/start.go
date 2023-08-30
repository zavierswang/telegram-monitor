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
)

func Start(ctx context.Context, update *tgb.MessageUpdate) error {
	bot := NewBot()
	err := update.Client.SetMyCommands(bot.Cmd).DoVoid(ctx)
	if err != nil {
		logger.Error("set command failed %v", err)
		return err
	}
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [start] controller", userId, username)

	buf := new(buffer.Buffer)
	tmpl, err := template.ParseFiles(cst.StartTemplateFile)
	if err != nil {
		logger.Error("[%s %s] template parse file %s, failed %v", userId, username, cst.StartTemplateFile, err)
		return err
	}
	err = tmpl.Execute(buf, username)
	if err != nil {
		logger.Error("[%s %s] template execute file %s, failed %v", userId, username, cst.StartTemplateFile, err)
		return err
	}
	return update.Answer(buf.String()).ParseMode(tg.HTML).ReplyMarkup(bot.ReplayMarkup).DisableWebPagePreview(true).DoVoid(ctx)
}
