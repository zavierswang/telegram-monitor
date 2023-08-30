package controllers

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"go.uber.org/zap/buffer"
	"html/template"
	"os"
	"strings"
	"telegram-monitor/pkg/core/cst"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
	"telegram-monitor/pkg/models"
	"telegram-monitor/pkg/utils"
)

func List(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [list] controller", userId, username)

	var user models.User
	err := global.App.DB.First(&user, "user_id = ?", userId).Error
	if err != nil {
		logger.Error("[%s %s] not found user", userId, username)
		return err
	}
	if !user.IsAdmin {
		logger.Warn("[%s %s] is not administrator", userId, username)
		return update.Answer("您不是管理员，不要乱点~").DoVoid(ctx)
	}
	var users []models.User
	global.App.DB.Find(&users)
	fs, err := os.ReadFile(cst.ListTemplateFile)
	if err != nil {
		logger.Error("[%s %s] open template file %s, failed %v", userId, username, cst.ListTemplateFile, err)
		return err
	}
	tmpl, err := template.New("list").Funcs(template.FuncMap{"format": utils.DateTime, "replace": utils.Replace, "bool": utils.Boolen}).Parse(string(fs))
	if err != nil {
		logger.Error("[%s %s] template parse file %s, failed %v", userId, username, cst.ListTemplateFile, err)
		return err
	}
	buf := new(buffer.Buffer)
	err = tmpl.Execute(buf, users)
	if err != nil {
		logger.Error("[%s %s] template execute file %s, failed %v", userId, username, cst.ListTemplateFile, err)
		return err
	}
	return update.Answer(buf.String()).ParseMode(tg.HTML).DoVoid(ctx)
}
