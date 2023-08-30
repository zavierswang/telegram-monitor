package services

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap/buffer"
	"io"
	"net/http"
	"telegram-monitor/pkg/core/logger"
)

type UserInfo struct {
	FirstName string
	LogUrl    string
	Exist     bool
}

func GetUserInfo(url string) (*UserInfo, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Error("new request with url: %s, failed %v", url, err)
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := new(buffer.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		logger.Error("io copy response body failed %v", err)
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		logger.Error("goquery new document failed %v", err)
		return nil, err
	}
	var info UserInfo
	doc.Find(`div.tgme_page`).Each(func(i int, selection *goquery.Selection) {
		firstName := selection.Find(`div.tgme_page_title > span`).Text()

		logUrl, exist := selection.Find(`div.tgme_page_photo > a > img`).Attr("src")
		//logger.Info("firstName: %s, logoUrl: %s", firstName, logUrl)
		info = UserInfo{
			FirstName: firstName,
			LogUrl:    logUrl,
			Exist:     exist,
		}
	})
	return &info, nil
}
