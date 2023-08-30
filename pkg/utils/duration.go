package utils

import (
	"fmt"
	"time"
)

func Duration(duration string) (start, end int64, label string) {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	switch duration {
	case "today":
		start = time.Date(year, month, day, 0, 0, 0, 0, time.Local).UnixMilli()
		end = now.UnixMilli()
		label = now.Format(time.DateOnly)
	case "yesterday":
		start = time.Date(year, month, day, 0, 0, 0, 0, time.Local).AddDate(0, 0, -1).UnixMilli()
		end = time.Date(year, month, day, 0, 0, 0, 0, time.Local).UnixMilli()
		label = fmt.Sprintf("%s ~ %s", time.UnixMilli(start).Format(time.DateOnly), time.UnixMilli(end).Format(time.DateOnly))
	case "week":
		start = ThisMonday()
		end = now.UnixMilli()
		label = fmt.Sprintf("%s ~ %s", time.UnixMilli(start).Format(time.DateOnly), time.UnixMilli(end).Format(time.DateOnly))
	case "last_week":
		thisWeekMonday := ThisMonday()
		start = time.UnixMilli(thisWeekMonday).AddDate(0, 0, -7).UnixMilli()
		end = thisWeekMonday
		label = fmt.Sprintf("%s ~ %s", time.UnixMilli(start).Format(time.DateOnly), time.UnixMilli(end).AddDate(0, 0, -1).Format(time.DateOnly))
	case "month":
		start = time.Date(year, month, 1, 0, 0, 0, 0, time.Local).UnixMilli()
		end = now.UnixMilli()
		label = fmt.Sprintf("%s ~ %s", time.UnixMilli(start).Format(time.DateOnly), time.UnixMilli(end).Format(time.DateOnly))
	case "last_month":
		thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
		start = thisMonth.AddDate(0, -1, 0).UnixMilli()
		end = thisMonth.AddDate(0, 0, 0).UnixMilli()
		label = fmt.Sprintf("%s ~ %s", time.UnixMilli(start).Format(time.DateOnly), time.UnixMilli(end).AddDate(0, 0, -1).Format(time.DateOnly))
	}
	return
}

func ThisMonday() int64 {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset).UnixMilli()
}
