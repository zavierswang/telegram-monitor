package controllers

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg/tgb"
	"regexp"
	"strconv"
	"strings"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
	"telegram-monitor/pkg/models"
	"telegram-monitor/pkg/utils"
	"time"
)

func Grant(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [grant] controller", userId, username)
	text := strings.TrimSpace(update.Text)
	var u models.User
	global.App.DB.First(&u, "user_id = ?", userId)
	if !u.IsAdmin {
		logger.Error("[%s %s] is not administrator", userId, username)
		return update.Answer("您不是管理员，不要乱点~").DoVoid(ctx)
	}
	compile, err := regexp.Compile(`^grant\s+(?P<username>\w+)(?P<space>\s+)?(?P<days>\d+)?`)
	if err != nil {
		logger.Error("[%s %s] regexp compile failed %v", userId, username, err)
		return update.Answer("您输入的格式错误，请输入正确的格式<code>grant username 3</code>").DoVoid(ctx)
	}
	var user models.User
	groups := utils.FindGroups(compile, text)
	err = global.App.DB.First(&user, "username = ?", fmt.Sprintf("%s", groups["username"])).Error
	if err != nil {
		logger.Error("[%s %s] not found user, %v", userId, username, err)
		return update.Answer("您输入的用户名不存在~").DoVoid(ctx)
	}
	var day int
	if groups["days"] == "" {
		day = 3
	} else {
		day, _ = strconv.Atoi(groups["days"])
	}
	expire := time.Now().AddDate(0, 0, day)
	user.ExpiredAt = expire
	err = global.App.DB.Save(&user).Error
	if err != nil {
		logger.Error("[%s %s] update user expiredTime failed %v", userId, username, err)
		return update.Answer("服务器错误~").DoVoid(ctx)
	}

	return update.Answer("授权成功！").DoVoid(ctx)
}
