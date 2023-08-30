package routes

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"regexp"
	"telegram-monitor/pkg/controllers"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/middleware"
	"telegram-monitor/pkg/models"
)

func Telegram(router *tgb.Router) {
	router.Use(middleware.SessionManager)
	router.Use(tgb.MiddlewareFunc(middleware.Hook))

	router.Message(controllers.Start, tgb.Command("start"))
	router.Message(controllers.Help, tgb.Any(tgb.Command("help"), tgb.TextEqual(Menu.Help)))
	router.Message(controllers.Address, tgb.TextEqual(Menu.Address), tgb.ChatType(tg.ChatTypePrivate))
	router.Message(controllers.Overview, tgb.TextEqual(Menu.Overview), tgb.ChatType(tg.ChatTypePrivate))
	router.Message(controllers.Price, tgb.TextEqual(Menu.Price), tgb.ChatType(tg.ChatTypePrivate))
	router.Message(controllers.AddAddress, tgb.Regexp(regexp.MustCompile(`^T\w+$`)))
	router.Message(controllers.Grant, tgb.All(isAdmin(), tgb.Regexp(regexp.MustCompile(`^grant\s+\w+(?P<space>\s+)?(?P<days>\d+)?`))), tgb.ChatType(tg.ChatTypePrivate))
	router.Message(controllers.Show, tgb.Regexp(regexp.MustCompile(`show(?P<space>\s+)?(?P<username>\w+)?`)), tgb.ChatType(tg.ChatTypePrivate))
	router.Message(controllers.List, tgb.All(isAdmin(), tgb.TextEqual("list")), tgb.ChatType(tg.ChatTypePrivate))
	router.Message(controllers.Remark, middleware.IsSessionStep(middleware.SessionStepRemark))
	router.Message(controllers.SetAvator, middleware.IsSessionStep(middleware.SessionStepAvator), tgb.MessageType(tg.MessageTypePhoto))

	router.CallbackQuery(controllers.Close, tgb.TextEqual(Callback.Close))
	router.CallbackQuery(controllers.BackCallback, tgb.TextEqual(Callback.ListOverview))
	router.CallbackQuery(controllers.AddCallback, tgb.Regexp(regexp.MustCompile(`^/add\s+T\w+$`)))
	router.CallbackQuery(controllers.RemoveCallback, tgb.Regexp(regexp.MustCompile(`^/remove\s+T\w+$`)))
	router.CallbackQuery(controllers.RemarkCallback, tgb.Regexp(regexp.MustCompile(`^/remark\s+T\w+$`)))
	router.CallbackQuery(controllers.StatisticsCallback, tgb.Regexp(regexp.MustCompile(`^/statistics\s+T\w+`)))
	router.CallbackQuery(controllers.BillDurationCallback, tgb.Regexp(regexp.MustCompile(`^/duration\s+\w+\s+\w+`)))
	router.CallbackQuery(controllers.BackCallback, tgb.Regexp(regexp.MustCompile(`^/bill\s+\w+\s+T\w+`)))
	router.CallbackQuery(controllers.AvatorCallback, tgb.Regexp(regexp.MustCompile(`avator\s+\d+\s+T\w+`)))
}

var isAdmin = func() tgb.Filter {
	return tgb.FilterFunc(func(ctx context.Context, update *tgb.Update) (bool, error) {
		message := update.Message
		callback := update.CallbackQuery
		var user models.User
		var userId string
		if message != nil {
			userId = message.From.ID.PeerID()
			global.App.DB.First(&user, "user_id = ?", userId)
			if user.IsAdmin {
				return true, nil
			}
		}
		if callback != nil {
			userId = callback.From.ID.PeerID()
			global.App.DB.First(&user, "user_id = ?", userId)
			if user.IsAdmin {
				return true, nil
			}
		}
		return false, nil
	})
}

var Menu = struct {
	Start         string
	Address       string
	Overview      string
	Price         string
	Help          string
	Administrator string
}{
	Start:         "🎉 开始",
	Help:          "💡 帮助",
	Overview:      "💰钱包概览",
	Address:       "⚙️钱包管理",
	Price:         "🌏实时U价",
	Administrator: "🧰管理员操作",
}

var Callback = struct {
	ListOverview string
	Close        string
	Today        string
	Yesterday    string
	Week         string
	LastWeek     string
	Month        string
	LastMonth    string
}{
	ListOverview: "list_overview",
	Close:        "close",
	Today:        "today",
	Yesterday:    "yesterday",
	Week:         "week",
	LastWeek:     "last_week",
	Month:        "month",
	LastMonth:    "last_month",
}
