package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"go.uber.org/zap/buffer"
	"golang.org/x/exp/slices"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"telegram-monitor/pkg/core/cst"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/core/logger"
	"telegram-monitor/pkg/middleware"
	"telegram-monitor/pkg/models"
)

func SetAvator(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [set_avator] controller", userId, username)
	sess := middleware.SessionManager.Get(ctx)
	address := sess.AvatorAddress
	logger.Info("[%s %s] avator sess: %+v", userId, username, sess)
	sizes := update.Message.Photo
	slices.SortFunc(sizes, func(a, b tg.PhotoSize) bool {
		return a.Width*a.Height > b.Width*b.Height
	})
	photo := sizes[0]
	logger.Info("[%s %s] file width: %d, height: %d, size: %d", userId, username, photo.Width, photo.Height, photo.FileSize)
	_ = update.Update.Reply(ctx, tg.NewSendChatActionCall(update.Message.Chat, tg.ChatActionUploadPhoto))
	info, err := update.Client.GetFile(photo.FileID).Do(ctx)
	if err != nil {
		logger.Error("[%s %s] get file failed %v", userId, username, err)
		return err
	}
	ext := strings.ToLower(filepath.Ext(info.FilePath))
	logger.Info("[%s %s] file ext: %s", userId, username, ext)

	file, err := update.Client.Download(ctx, info.FilePath)
	if err != nil {
		logger.Error("[%s %s] download file failed %v", userId, username, err)
		return err
	}
	defer func() { _ = file.Close() }()
	uploadDir := "upload"
	_, err = os.Stat(uploadDir)
	if errors.Is(err, os.ErrNotExist) {
		_ = os.Mkdir(uploadDir, os.ModePerm)
	}
	newAvator := fmt.Sprintf("%s_%s%s", userId, address, filepath.Ext(info.FilePath))
	newFile, err := os.Create(fmt.Sprintf("upload/%s", newAvator))
	if err != nil {
		logger.Error("[%s %s] create new avator failed %v", userId, username, err)
		return err
	}
	defer func() { _ = newFile.Close() }()
	_, err = io.Copy(newFile, file)
	if err != nil {
		logger.Error("[%s %s] write new file failed %v", userId, username, err)
		return err
	}
	logger.Info("[%s %s] save new avator %s successfully", userId, username, newAvator)
	var addr models.Address
	err = global.App.DB.First(&addr, "user_id = ? AND address = ?", userId, address).Error
	if err != nil {
		logger.Error("[%s %s] not found address %s %s", userId, username)
		return err
	}
	addr.Avator = newAvator
	global.App.DB.Save(&addr)
	middleware.SessionManager.Reset(sess)
	tmpl, _ := template.ParseFiles(cst.SetAvatorTemplateFile)
	buf := new(buffer.Buffer)
	_ = tmpl.Execute(buf, address)
	return update.Answer(buf.String()).ParseMode(tg.HTML).DoVoid(ctx)
}
