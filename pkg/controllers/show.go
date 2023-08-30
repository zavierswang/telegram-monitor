package controllers

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"go.uber.org/zap/buffer"
	"html/template"
	"os"
	"regexp"
	"strings"
	"telegram-monitor/pkg/core/cst"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
	"telegram-monitor/pkg/models"
	"telegram-monitor/pkg/utils"
)

func Show(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [show] controller", userId, username)
	var user models.User
	fs, err := os.ReadFile(cst.ShowTemplateFile)
	if err != nil {
		logger.Error("[%s %s] read template file %s, failed %v", userId, username, cst.ShowTemplateFile, err)
		return err
	}
	tmpl, err := template.New("show").Funcs(template.FuncMap{"format": utils.DateTime, "replace": utils.Replace}).Parse(string(fs))
	if err != nil {
		logger.Error("[%s %s] template parse file %s, failed %v", userId, username, cst.ShowTemplateFile, err)
		return err
	}
	buf := new(buffer.Buffer)
	text := strings.TrimSpace(update.Text)
	compile, err := regexp.Compile(`show(?P<space>\s+)?(?P<username>\w+)?`)
	if err != nil {
		logger.Error("[%s %s] regexp compile failed %v", userId, username, err)
		return update.Answer("您的指令格式错误~").DoVoid(ctx)
	}
	groups := utils.FindGroups(compile, text)
	if groups["username"] == "" || strings.ToLower(groups["username"]) == username {
		// 查询用户自己的个人信息
		err = global.App.DB.First(&user, "username = ?", username).Error
		if err != nil {
			logger.Error("[%s %s] not found user %s", userId, username, groups["username"])
			return update.Answer("非法用户，禁止操作~").DoVoid(ctx)
		}
	} else {
		// 管理员查询其它用户
		var u models.User
		err = global.App.DB.First(&u, "user_id = ?", userId).Error
		if err != nil {
			logger.Error("[%s %s] not found user %s", userId, username, groups["username"])
			return update.Answer("非法用户，禁止操作~").DoVoid(ctx)
		}
		if !u.IsAdmin {
			logger.Error("[%s %s] is not administrator", userId, username)
			return update.Answer("您不是管理员，不要乱点~").DoVoid(ctx)
		}
		err = global.App.DB.First(&user, "username = ?", groups["username"]).Error
		if err != nil {
			logger.Error("[%s %s] not found username %s", userId, username, groups["username"])
			return update.Answer("您输入的用户名不存在~").DoVoid(ctx)
		}
	}
	err = tmpl.Execute(buf, user)
	if err != nil {
		logger.Error("[%s %s] template execute file %s, failed %v", userId, username, cst.ShowTemplateFile, err)
		return err
	}
	return update.Answer(buf.String()).ParseMode(tg.HTML).DoVoid(ctx)
}
