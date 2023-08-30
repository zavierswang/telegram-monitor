package controllers

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"go.uber.org/zap/buffer"
	"html/template"
	"math"
	"regexp"
	"strconv"
	"strings"
	"telegram-monitor/pkg/core/cst"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
	"telegram-monitor/pkg/middleware"
	"telegram-monitor/pkg/models"
	"telegram-monitor/pkg/services/tron"
	"telegram-monitor/pkg/utils"
)

func AddCallback(ctx context.Context, callback *tgb.CallbackQueryUpdate) error {
	userId := callback.From.ID.PeerID()
	username := callback.From.Username.PeerID()
	username = strings.ToLower(strings.ReplaceAll(username, "@", ""))
	logger.Info("[%s %s] trigger action [add_monitor] callback", userId, username)

	text := callback.Data
	chatId := callback.Message.Chat.ID
	messageId := callback.Message.ID
	compile, err := regexp.Compile(`^/add\s+(?P<address>T\w+)`)
	if err != nil {
		logger.Error("compile regexp failed %v", err)
		return err
	}
	groups := utils.FindGroups(compile, text)
	var addrs []models.Address
	global.App.DB.Find(&addrs, "user_id = ? AND address = ?", userId, groups["address"])
	if len(addrs) == 0 {
		addr := models.Address{
			UserID:    userId,
			Username:  username,
			Address:   groups["address"],
			IsMonitor: true,
		}
		err = global.App.DB.Create(&addr).Error
		if err != nil {
			logger.Error("[%s %s] insert monitor address failed %v", userId, username, err)
			return err
		}
	}

	chatType := callback.Message.Chat.Type
	if chatType == tg.ChatTypeSupergroup {
		group := models.Group{
			ChatID:   callback.Message.Chat.ID.PeerID(),
			Username: callback.Message.Chat.Username.PeerID(),
			Title:    callback.Message.Chat.Title,
		}
		err = global.App.DB.Model(&models.Address{}).Where("user_id = ? AND address = ?", userId, groups["address"]).Updates(models.Address{Group: group}).Error
		if err != nil {
			logger.Error("[%s %s] update address group failed %v", userId, username, err)
			return err
		}
	}

	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](2).Row()
	layout.Insert(tg.NewInlineKeyboardButtonCallback("🏷️备注", fmt.Sprintf("/remark %s", groups["address"])))
	layout.Insert(tg.NewInlineKeyboardButtonCallback("👤头像", fmt.Sprintf("avator %s %s", userId, groups["address"])))
	layout.Insert(tg.NewInlineKeyboardButtonCallback("关闭", Callback.Close))
	ikb := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)
	tmpl, err := template.ParseFiles(cst.AddMonitorTemplateFile)
	if err != nil {
		logger.Error("[%s %s] template parse file %s, failed %v", userId, username, cst.AddMonitorTemplateFile, err)
		return err
	}
	buf := new(buffer.Buffer)
	global.App.DB.Find(&addrs, "user_id = ?", userId)
	var tpl = struct {
		Address string
		Addrs   []models.Address
	}{
		Address: groups["address"],
		Addrs:   addrs,
	}
	err = tmpl.Execute(buf, tpl)
	if err != nil {
		logger.Error("[%s %s] template execute file %s, failed %v", userId, username, cst.AddMonitorTemplateFile, err)
		return err
	}
	return callback.Client.EditMessageText(chatId, messageId, buf.String()).ParseMode(tg.HTML).ReplyMarkup(ikb).DoVoid(ctx)
}

func RemoveCallback(ctx context.Context, callback *tgb.CallbackQueryUpdate) error {
	userId := callback.From.ID.PeerID()
	username := callback.From.Username.PeerID()
	username = strings.ToLower(strings.ReplaceAll(username, "@", ""))
	logger.Info("[%s %s] trigger action [remove_monitor] callback", userId, username)

	chatId := callback.Message.Chat.ID
	messageId := callback.Message.ID
	compile, err := regexp.Compile(`^/remove\s+(?P<address>T\w+)`)
	if err != nil {
		logger.Error("[%s %s] compile regexp failed %v", userId, username, err)
		return err
	}
	groups := utils.FindGroups(compile, callback.Data)
	err = global.App.DB.Delete(&models.Address{}, "user_id = ? AND address = ?", userId, groups["address"]).Error
	if err != nil {
		logger.Error("[%s %s] delete monitor address failed %v", userId, username, err)
		return err
	}
	return callback.Client.EditMessageText(chatId, messageId, "<b>删除成功♻️</b>").ParseMode(tg.HTML).DoVoid(ctx)
}

func RemarkCallback(ctx context.Context, callback *tgb.CallbackQueryUpdate) error {
	userId := callback.From.ID.PeerID()
	username := callback.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [remark] callback", userId, username)

	sess := middleware.SessionManager.Get(ctx)
	sess.Step = middleware.SessionStepRemark
	compile, err := regexp.Compile(`^/remark\s+(?P<address>T\w+)$`)
	if err != nil {
		logger.Error("compile remark failed %v", err)
		return err
	}
	groups := utils.FindGroups(compile, callback.Data)
	sess.RemarkAddress = groups["address"]
	chatId := callback.Message.Chat.ID
	return callback.Update.Reply(ctx, callback.Client.
		SendMessage(chatId, "<b>请回复备注名</b>\n请尽量避免用标点符号或过长名称").
		ParseMode(tg.HTML).
		DisableWebPagePreview(true))
}

func StatisticsCallback(ctx context.Context, callback *tgb.CallbackQueryUpdate) error {
	userId := callback.From.ID.PeerID()
	username := callback.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [statistics] callback", userId, username)

	compile, err := regexp.Compile(`^/statistics\s+(?P<address>T\w+)`)
	if err != nil {
		logger.Error("compile failed %v", err)
		return err
	}
	chatId := callback.Message.Chat.ID
	messageId := callback.Message.ID
	_ = callback.Client.EditMessageText(chatId, messageId, "正在查询...").ParseMode(tg.HTML).DoVoid(ctx)

	groups := utils.FindGroups(compile, callback.Data)
	address := groups["address"]
	var addr models.Address
	err = global.App.DB.First(&addr, "user_id = ? AND address = ?", userId, address).Error
	if err != nil {
		logger.Error("[%s %s] not found address %s", userId, username, address)
		return callback.Client.EditMessageText(chatId, messageId, "该地址不存在，可能已被其它客户端删除了！").ParseMode(tg.HTML).DoVoid(ctx)
	}
	remark := addr.Remark
	if remark == "" {
		remark = "暂无备注"
	}
	var (
		fileArg tg.FileArg
		flag    bool
	)
	avator := addr.Avator
	if avator != "" {
		flag = true
		inputFile, err := tg.NewInputFileLocal(fmt.Sprintf("upload/%s", avator))
		if err != nil {
			logger.Error("[%s %s] not found avator file %s", userId, username, err)
			flag = false
		}
		fileArg = tg.NewFileArgUpload(inputFile)
	}
	info, err := tron.GetAddressAccountInfo(map[string]string{"address": address})
	if err != nil {
		logger.Error("[%s %s] tron.GetAddressAccountInfo failed %v", userId, username, err)
		return callback.Client.EditMessageText(chatId, messageId, "网络错误，请重试~").DoVoid(ctx)
	}
	tmpl, err := template.ParseFiles(cst.StatisticsTemplateFile)
	if err != nil {
		logger.Error("[%s %s] template parse file %s, failed %v", userId, username, cst.StatisticsTemplateFile, err)
		return err
	}
	buf := new(buffer.Buffer)
	tpl := Statistics{
		AddressDetailInfo: *info,
		Remark:            remark,
	}
	err = tmpl.Execute(buf, tpl)
	if err != nil {
		logger.Error("[%s %s] template execute file %s, failed %v", userId, username, cst.StatisticsTemplateFile, err)
		return err
	}
	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](4).Row()
	layout.Insert(tg.NewInlineKeyboardButtonCallback("今日", fmt.Sprintf("/duration %s %s", Callback.Today, address)))
	layout.Insert(tg.NewInlineKeyboardButtonCallback("昨日", fmt.Sprintf("/duration %s %s", Callback.Yesterday, address)))
	layout.Insert(tg.NewInlineKeyboardButtonCallback("本周", fmt.Sprintf("/duration %s %s", Callback.Week, address)))
	layout.Insert(tg.NewInlineKeyboardButtonCallback("上周", fmt.Sprintf("/duration %s %s", Callback.LastWeek, address)))
	layout.Insert(tg.NewInlineKeyboardButtonCallback("本月", fmt.Sprintf("/duration %s %s", Callback.Month, address)))
	layout.Insert(tg.NewInlineKeyboardButtonCallback("上月", fmt.Sprintf("/duration %s %s", Callback.LastMonth, address)))
	inlineKeyboardDuration := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)

	layoutBack := tg.NewButtonLayout[tg.InlineKeyboardButton](2).Row()
	layoutBack.Insert(tg.NewInlineKeyboardButtonCallback("↩️返回", Callback.ListOverview))
	inlineKeyboardBack := tg.NewInlineKeyboardMarkup(layoutBack.Keyboard()...)

	inlineKeyboardDuration.InlineKeyboard = append(inlineKeyboardDuration.InlineKeyboard, inlineKeyboardBack.InlineKeyboard...)
	if flag {
		_ = callback.Client.DeleteMessage(chatId, messageId).DoVoid(ctx)
		return callback.Client.SendPhoto(chatId, fileArg).
			Caption(buf.String()).
			ParseMode(tg.HTML).
			ReplyMarkup(inlineKeyboardDuration).
			DoVoid(ctx)
	}
	return callback.Client.EditMessageText(chatId, messageId, buf.String()).
		ParseMode(tg.HTML).
		ReplyMarkup(inlineKeyboardDuration).
		DisableWebPagePreview(true).
		DoVoid(ctx)
}

func BillDurationCallback(ctx context.Context, callback *tgb.CallbackQueryUpdate) error {
	userId := callback.From.ID.PeerID()
	username := callback.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")

	chatId := callback.Message.Chat.ID
	messageId := callback.Message.ID
	compile, err := regexp.Compile(`^/duration\s+(?P<date>\w+)\s+(?P<address>T\w+)`)
	if err != nil {
		logger.Error("compile failed %v", err)
		return err
	}
	groups := utils.FindGroups(compile, callback.Data)
	logger.Info("[%s %s] trigger action [duration %s] callback", userId, username, groups["date"])
	var addr models.Address
	err = global.App.DB.First(&addr, "user_id = ? AND address = ?", userId, groups["address"]).Error
	if err != nil {
		logger.Error("[%s %s] not found address", userId, username)
		_ = callback.Client.DeleteMessage(chatId, messageId)
		return callback.Client.SendMessage(chatId, "地址已经被删除，请重试其它地址~").DoVoid(ctx)
	}
	if addr.Avator != "" {
		err = callback.Client.EditMessageCaption(chatId, messageId, "正在查询，请稍后...").
			ParseMode(tg.HTML).
			DoVoid(ctx)
	} else {
		err = callback.Client.EditMessageText(chatId, messageId, "正在查询，请稍后...").
			ParseMode(tg.HTML).
			DisableWebPagePreview(true).
			DoVoid(ctx)
	}
	if err != nil {
		logger.Error("[%s %s] send queries message failed %v", err)
	}

	buf, err := durationAddress(groups, userId)
	if err != nil {
		logger.Error("[%s %s] duration address statistics failed %v", err)
		return callback.Client.EditMessageText(chatId, messageId, "网络错误，请重试~").DoVoid(ctx)
	}
	if addr.Avator != "" {
		return callback.Client.EditMessageCaption(chatId, messageId, buf.String()).
			ParseMode(tg.HTML).
			ReplyMarkup(inlineKeyboard(groups)).
			DoVoid(ctx)
	}
	return callback.Client.EditMessageText(chatId, messageId, buf.String()).
		ParseMode(tg.HTML).
		DisableWebPagePreview(true).
		ReplyMarkup(inlineKeyboard(groups)).
		DoVoid(ctx)
}

func BackCallback(ctx context.Context, callback *tgb.CallbackQueryUpdate) error {
	userId := callback.From.ID.PeerID()
	username := callback.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [back] callback", userId, username)

	chatId := callback.Message.Chat.ID
	messageId := callback.Message.ID
	var addrs []models.Address
	global.App.DB.Find(&addrs, "user_id = ?", userId)
	if len(addrs) == 0 {
		return callback.Update.Reply(ctx, callback.Client.SendMessage(chatId, "<b>您还未添加钱包地址, 请输入TRC20地址👇</b>").ParseMode(tg.HTML))
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
	ikb := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)
	_ = callback.Client.DeleteMessage(chatId, messageId).DoVoid(ctx)
	return callback.Update.Reply(ctx, callback.Client.SendMessage(chatId, "<b>您的钱包地址列表</b>").
		ParseMode(tg.HTML).
		ReplyMarkup(ikb))
}

func AvatorCallback(ctx context.Context, callback *tgb.CallbackQueryUpdate) error {
	userId := callback.From.ID.PeerID()
	username := callback.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [remark] callback", userId, username)
	chatId := callback.Message.Chat.ID

	compile, err := regexp.Compile(`avator\s+(?P<userId>\d+)\s+(?P<address>\w+)`)
	if err != nil {
		logger.Error("[%s %s] compile regexp failed %v", userId, username, err)
		return err
	}
	groups := utils.FindGroups(compile, callback.Data)

	sess := middleware.SessionManager.Get(ctx)
	sess.Step = middleware.SessionStepAvator
	sess.AvatorUserId = groups["userId"]
	sess.AvatorAddress = groups["address"]
	logger.Info("[%s %s] avator sess: %+v", sess)
	return callback.Client.SendMessage(chatId, "请上传您要设置的头像图片文件, 仅支持PNG/JPG格式").DoVoid(ctx)
}

func Close(ctx context.Context, callback *tgb.CallbackQueryUpdate) error {
	userId := callback.From.ID.PeerID()
	username := callback.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [close] callback", userId, username)

	chatId := callback.Message.Chat.ID
	messageId := callback.Message.ID
	sess := middleware.SessionManager.Get(ctx)
	middleware.SessionManager.Reset(sess)
	return callback.Client.DeleteMessage(chatId, messageId).DoVoid(ctx)
}

func inlineKeyboard(groups map[string]string) tg.InlineKeyboardMarkup {
	layoutDate := tg.NewButtonLayout[tg.InlineKeyboardButton](4).Row()
	for _, label := range Labels {
		display := DurationMap[label]
		if groups["date"] == label {
			layoutDate.Insert(tg.NewInlineKeyboardButtonCallback(fmt.Sprintf("✅%s", display), fmt.Sprintf("/duration %s %s", label, groups["address"])))
		} else {
			layoutDate.Insert(tg.NewInlineKeyboardButtonCallback(display, fmt.Sprintf("/duration %s %s", label, groups["address"])))
		}
	}
	layoutBack := tg.NewButtonLayout[tg.InlineKeyboardButton](1).Row()
	layoutBack.Insert(tg.NewInlineKeyboardButtonCallback("↩️返回", Callback.ListOverview))
	inlineKeyboardDate := tg.NewInlineKeyboardMarkup(layoutDate.Keyboard()...)
	inlineKeyboardBack := tg.NewInlineKeyboardMarkup(layoutBack.Keyboard()...)
	inlineKeyboardDate.InlineKeyboard = append(inlineKeyboardDate.InlineKeyboard, inlineKeyboardBack.InlineKeyboard...)
	return inlineKeyboardDate

}

func durationAddress(groups map[string]string, userId string) (*buffer.Buffer, error) {
	start, end, label := utils.Duration(groups["date"])
	var addr models.Address
	err := global.App.DB.First(&addr, "user_id = ? AND address = ?", userId, groups["address"]).Error
	if err != nil {
		logger.Error("duration address not found")
		return nil, err
	}
	remark := addr.Remark
	if remark == "" {
		remark = "暂无备注"
	}
	tmpl, err := template.ParseFiles(cst.DurationStatisticsTemplateFile)
	if err != nil {
		return nil, err
	}
	params := map[string]string{
		"limit":            "30",
		"start":            "0",
		"contract_address": cst.ContractAddress,
		"sort":             "-timestamp",
		"count":            "true",
		"filterTokenValue": "0",
		"start_timestamp":  strconv.FormatInt(start, 10),
		"end_timestamp":    strconv.FormatInt(end, 10),
		"relatedAddress":   groups["address"],
	}
	header := map[string]string{
		"TRON-PRO-API-KEY": global.App.Config.Telegram.TronScanApiKey,
	}
	transfers, err := tron.TRC20Transfer(params, header, false, false)
	if err != nil {
		logger.Error("trc20 transfers failed %v", err)
		return nil, err
	}

	var (
		inBalance  float64
		inCount    int
		outBalance float64
		outCount   int
	)
	for _, transfer := range transfers {
		quant, _ := strconv.ParseFloat(transfer.Quant, 64)
		if transfer.FromAddress == groups["address"] {
			outBalance += quant / math.Pow10(6)
			outCount += 1
		} else {
			inBalance += quant / math.Pow10(6)
			inCount += 1
		}
	}
	buf := new(buffer.Buffer)
	tpl := DurationBill{
		Tips:       fmt.Sprintf("✅%s 收支统计", DurationMap[groups["date"]]),
		Label:      label,
		Remark:     remark,
		Address:    groups["address"],
		InCount:    inCount,
		InBalance:  fmt.Sprintf("%.3f", inBalance),
		OutCount:   outCount,
		OutBalance: fmt.Sprintf("%.3f", outBalance),
	}
	err = tmpl.Execute(buf, tpl)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

var DurationMap = map[string]string{
	"today":      "今日",
	"yesterday":  "昨日",
	"week":       "本周",
	"last_week":  "上周",
	"month":      "本月",
	"last_month": "上月",
}

type DurationBill struct {
	Tips       string
	Label      string
	Remark     string
	Address    string
	InCount    int
	InBalance  string
	OutCount   int
	OutBalance string
}

type Statistics struct {
	tron.AddressDetailInfo
	Remark string
}
