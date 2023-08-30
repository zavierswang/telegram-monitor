package utils

import (
	"strings"
	"telegram-monitor/pkg/core/cst"
	"time"
)

func DateTime(t time.Time) string {
	return t.Format(cst.DateTimeFormatter)
}

func Boolen(ok bool) string {
	if ok {
		return "🟢"
	}
	return "⚫️"
}

func Replace(s, text string) string {
	return strings.Replace(text, s, "", 1)
}
