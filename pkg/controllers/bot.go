package controllers

import "github.com/mr-linch/go-tg"

type Bot struct {
	ReplayMarkup *tg.ReplyKeyboardMarkup
	Cmd          []tg.BotCommand
}

func NewBot() *Bot {
	layout := tg.NewReplyKeyboardMarkup(
		tg.NewButtonRow(
			tg.NewKeyboardButton(Menu.Overview),
			tg.NewKeyboardButton(Menu.Price),
		),
		tg.NewButtonRow(
			tg.NewKeyboardButton(Menu.Address),
			tg.NewKeyboardButton(Menu.Help),
		),
	)
	layout.ResizeKeyboard = true

	botCmd := []tg.BotCommand{
		{Command: "start", Description: Menu.Start},
		{Command: "help", Description: Menu.Help},
	}

	return &Bot{
		ReplayMarkup: layout,
		Cmd:          botCmd,
	}
}

var Menu = struct {
	Start    string
	Address  string
	Overview string
	Price    string
	Help     string
}{
	Start:    "🎉 开始",
	Help:     "💡 帮助",
	Overview: "💰钱包概览",
	Address:  "⚙️钱包管理",
	Price:    "🌏实时U价",
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

var Labels = []string{"today", "yesterday", "week", "last_week", "month", "last_month"}
