package controllers

import (
	"context"
	"github.com/mr-linch/go-tg"
	"go.uber.org/zap/buffer"
	"html/template"
	"strconv"
	"sync"
	"telegram-monitor/pkg/core/cst"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/models"
)

func Update(token string) {
	var users []models.User
	global.App.DB.Find(&users)
	tmpl, _ := template.ParseFiles(cst.UpdateTemplateFile)
	buf := new(buffer.Buffer)
	_ = tmpl.Execute(buf, nil)
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(user models.User) {
			defer wg.Done()
			chatId, _ := strconv.ParseInt(user.UserID, 10, 64)
			_ = global.App.Client.SendMessage(tg.ChatID(chatId), buf.String()).
				ParseMode(tg.HTML).
				DisableWebPagePreview(true).
				DoVoid(context.Background())
		}(user)
	}
	wg.Wait()

}
