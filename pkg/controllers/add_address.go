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
	"telegram-monitor/pkg/middleware"
	"telegram-monitor/pkg/models"
	"telegram-monitor/pkg/services/tron"
	"time"
)

func AddAddress(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [add_address] controller", userId, username)
	chatId := update.Chat.ID
	text := strings.TrimSpace(update.Message.Text)
	sess := middleware.SessionManager.Get(ctx)
	sess.Step = middleware.SessionStepAvator
	sess.RemarkAddress = text
	//resp, err := grid.ValidateAddress(text)
	//if err != nil || !resp.Result {
	//	logger.Error("[%s %s] address invalid resp: %+v, err: %+v", userId, username, resp, err)
	//	return update.Answer("非法输入，请输入正确的TRC20地址").ParseMode(tg.HTML).DoVoid(ctx)
	//}
	var user models.User
	err := global.App.DB.Find(&user, "user_id = ?", userId).Error
	if err != nil {
		logger.Error("[%s %s] not found user, failed %v", userId, username, err)
		return err
	}
	delta := time.Now().UnixMilli() - user.ExpiredAt.UnixMilli()
	if !user.IsAdmin && delta > 0 {
		logger.Warn("[%s %s] account has been expired", userId, username)
		return update.Answer("您的帐户已过期，请续费或联系客服~").ParseMode(tg.HTML).DoVoid(ctx)
	}
	msg, _ := update.Answer("正在查询，请稍后~").ParseMode(tg.HTML).Do(ctx)
	info, err := tron.GetAddressAccountInfo(map[string]string{"address": text})
	if err != nil || info == nil {
		logger.Error("[%s %s] func [service.GetAddressAccountInfo] failed %v", userId, username, err)
		return update.Client.EditMessageText(chatId, msg.ID, "网络错误，请稍后重试~").DoVoid(ctx)
	}

	var addrs []models.Address
	var remark string
	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](3).Row()
	global.App.DB.Find(&addrs, "user_id = ? AND address = ?", userId, text)

	var (
		fileArg tg.FileArg
		flag    bool
	)
	if len(addrs) == 0 {
		remark = "暂未设置"
		if user.IsAdmin || delta <= 0 {
			layout.Insert(tg.NewInlineKeyboardButtonCallback("🔔添加到监控", fmt.Sprintf("/add %s", text)))
		} else {
			layout.Insert(tg.NewInlineKeyboardButtonCallback("帐户已过期", "alert"))
		}
	} else {
		remark = addrs[0].Remark
		if remark == "" {
			remark = "暂无备注"
		}
		avator := addrs[0].Avator
		if avator != "" {
			inputFile, err := tg.NewInputFileLocal(fmt.Sprintf("upload/%s", avator))
			if err == nil {
				fileArg = tg.NewFileArgUpload(inputFile)
				flag = true
			} else {
				logger.Error("[%s %s] avator file not found", userId, username)
			}
		}
		if update.Message.Chat.Type.String() == "private" {
			layout.Insert(tg.NewInlineKeyboardButtonCallback("🔕删除", fmt.Sprintf("/remove %s", text)))
		}
		layout.Insert(tg.NewInlineKeyboardButtonCallback("🏷️备注", fmt.Sprintf("/remark %s", text)))
		layout.Insert(tg.NewInlineKeyboardButtonCallback("👤头像", fmt.Sprintf("avator %s %s", userId, text)))
	}
	layout.Insert(tg.NewInlineKeyboardButtonURL("联系客服", fmt.Sprintf("https://t.me/%s", global.App.Config.App.Support)))
	layout.Insert(tg.NewInlineKeyboardButtonCallback("关闭", Callback.Close))
	ikb := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)

	tmpl, err := template.ParseFiles(cst.AddressDetailTemplateFile)
	if err != nil {
		logger.Error("[%s %s] template parse file %s, failed %v", userId, username, cst.AddressDetailTemplateFile, err)
		return err
	}
	buf := new(buffer.Buffer)
	tpl := Statistics{
		AddressDetailInfo: *info,
		Remark:            remark,
	}
	err = tmpl.Execute(buf, tpl)
	if err != nil {
		logger.Error("[%s %s] template execute file %s, failed %v", userId, username, cst.AddressDetailTemplateFile, err)
		return err
	}
	if flag {
		_ = update.Client.DeleteMessage(chatId, msg.ID).DoVoid(ctx)
		return update.AnswerPhoto(fileArg).Caption(buf.String()).ParseMode(tg.HTML).ReplyMarkup(ikb).DoVoid(ctx)
	}
	return update.Client.EditMessageText(chatId, msg.ID, buf.String()).ParseMode(tg.HTML).ReplyMarkup(ikb).DoVoid(ctx)
}
