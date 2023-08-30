package routes

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mr-linch/go-tg"
	"go.uber.org/zap/buffer"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"telegram-monitor/pkg/core/cst"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
	"telegram-monitor/pkg/middleware"
	"telegram-monitor/pkg/models"
	"telegram-monitor/pkg/services/tron"
	"telegram-monitor/pkg/utils"
	"time"
)

type USDT struct {
	TransactionId string
	FromAddress   string
	ToAddress     string
	CreateTime    string
	Balance       string //入帐或出帐金额
	Description   string
	Amount        string //钱包余额
	Type          string
	NotiferDesc   NotifierDesc
}

type NotifierDesc struct {
	Label         string
	Address       string
	ListenAddress string
	Mark          string
}

func RegisterRoutes() *gin.Engine {
	if global.App.Config.App.Env == "ttkj" ||
		global.App.Config.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	//router.Use(gin.Logger())
	// 跨域处理
	router.Use(middleware.Cors())
	// 注册 api 分组路由
	apis := router.Group("/api")
	{
		apis.POST("/callback/trc20", NotifierTRC20)
	}
	return router
}

var (
	TRC20Queue  = make(map[string][]string)
	TRC20Ticker = make(map[string]int)
)

func NotifierTRC20(ctx *gin.Context) {
	var notifier tron.Notifier
	err := ctx.ShouldBindJSON(&notifier)
	if err != nil {
		logger.Error("[callback trc20] notification bind json failed %v", err)
		ctx.JSON(http.StatusBadRequest, nil)
		ctx.Abort()
	}
	var addrs []models.Address
	global.App.DB.Find(&addrs, "id = ? AND address = ?", notifier.ID, notifier.Address)
	if len(addrs) == 0 {
		logger.Warn("[callback trc20] not found by notifier: %+v", notifier)
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	wallet, err := tron.GetAddressAccountInfo(map[string]string{"address": notifier.Address})
	if err != nil {
		logger.Warn("%s wallet can't search, please waiting for next time.")
	}
	var balance = "0.000"
	if wallet != nil {
		balance = fmt.Sprintf("%.3f", wallet.Balance)
	}
	mark := "暂无备注"
	if addrs[0].Remark != "" {
		mark = addrs[0].Remark
	}
	var (
		nType string
		desc  NotifierDesc
	)
	key := fmt.Sprintf("%s_%s", notifier.UserId, notifier.Address)
	TRC20Ticker[key]++
	for _, transfer := range notifier.Transfers {
		exist := utils.ListContains(transfer.TransactionId, TRC20Queue[key])
		if !exist {
			logger.Info("[callback trc20] found %s latest txId: %s", key, transfer.TransactionId)
			TRC20Queue[key] = append(TRC20Queue[key], transfer.TransactionId)
			quant, _ := strconv.ParseFloat(transfer.Quant, 64)
			amount := fmt.Sprintf("%.3f USDT", quant/math.Pow10(6))
			if transfer.FromAddress == notifier.Address {
				nType = "🔴支出"
				desc = NotifierDesc{
					Label:         "目标地址",
					Address:       transfer.ToAddress,
					ListenAddress: notifier.Address,
					Mark:          mark,
				}
			} else {
				nType = "🟢收入"
				desc = NotifierDesc{
					Label:         "来源地址",
					Address:       transfer.FromAddress,
					ListenAddress: notifier.Address,
					Mark:          mark,
				}
			}
			usdt := USDT{
				TransactionId: transfer.TransactionId,
				FromAddress:   transfer.FromAddress,
				ToAddress:     transfer.ToAddress,
				CreateTime:    time.UnixMilli(transfer.BlockTs).Format(cst.DateTimeFormatter),
				Balance:       balance,
				Type:          nType,
				Amount:        amount,
				NotiferDesc:   desc,
			}
			err = LatestUSDTNotifer(addrs[0], usdt)
			if err != nil {
				logger.Error("latest USDT notification failed %v", err)
				ctx.JSON(http.StatusInternalServerError, nil)
				return
			}

			if len(TRC20Queue[key]) >= 500 {
				logger.Warn("[GIN] clean oldest txnIds 0 ~ 400")
				TRC20Queue[key] = TRC20Queue[key][400:]
			}
		}
	}

	ctx.JSON(http.StatusOK, nil)
	return
}

func LatestUSDTNotifer(addr models.Address, usdt USDT) error {
	var fileArg tg.FileArg
	var flag bool
	if addr.Avator != "" {
		flag = true
		inputFile, err := tg.NewInputFileLocal(fmt.Sprintf("upload/%s", addr.Avator))
		if err != nil {
			logger.Error("[callback trc20] not found avator file err %v", err)
			flag = false
		}
		fileArg = tg.NewFileArgUpload(inputFile)
	}

	buf := new(buffer.Buffer)
	tmpl, err := template.ParseFiles(cst.LatestUsdtTemplateFile)
	if err != nil {
		logger.Error("[callback trc20] template parse file %s, failed %v", cst.LatestUsdtTemplateFile, err)
		return err
	}
	err = tmpl.Execute(buf, &usdt)
	if err != nil {
		logger.Error("[callback trc20] template execute file %s, failed %v", cst.LatestUsdtTemplateFile, err)
		return err
	}
	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](1).Row()
	layout.Insert(
		tg.NewInlineKeyboardButtonURL("查看交易详情", fmt.Sprintf("https://tronscan.org/#/transaction/%s", usdt.TransactionId)),
	)
	ikb := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)
	chatId, _ := strconv.ParseInt(addr.UserID, 10, 64)
	if flag {
		_ = global.App.Client.SendPhoto(tg.ChatID(chatId), fileArg).
			Caption(buf.String()).
			ParseMode(tg.HTML).
			ReplyMarkup(ikb).
			DoVoid(context.Background())
	} else {
		_ = global.App.Client.SendMessage(tg.ChatID(chatId), buf.String()).
			ParseMode(tg.HTML).
			DisableWebPagePreview(true).
			ReplyMarkup(ikb).
			DoVoid(context.Background())
	}

	if addr.Group.ChatID != "" {
		chatId, _ = strconv.ParseInt(addr.Group.ChatID, 0, 64)
		if flag {
			_ = global.App.Client.SendPhoto(tg.ChatID(chatId), fileArg).
				Caption(buf.String()).
				ParseMode(tg.HTML).
				ReplyMarkup(ikb).
				DoVoid(context.Background())
		} else {
			_ = global.App.Client.SendMessage(tg.ChatID(chatId), buf.String()).
				ParseMode(tg.HTML).
				ReplyMarkup(ikb).
				DisableWebPagePreview(true).
				DoVoid(context.Background())
		}
	}
	//body := map[string]interface{}{
	//	"text":       buf.String(),
	//	"parse_mode": "HTML",
	//	"reply_markup": map[string]interface{}{
	//		"inline_keyboard": []interface{}{
	//			[]map[string]interface{}{
	//				{
	//					"text": "查看交易详情",
	//					"url":  fmt.Sprintf("https://tronscan.org/#/transaction/%s", usdt.TransactionId),
	//				},
	//			},
	//		},
	//	},
	//}
	return nil
}
