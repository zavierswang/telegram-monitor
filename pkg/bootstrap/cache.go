package bootstrap

import (
	"github.com/patrickmn/go-cache"
	"telegram-monitor/pkg/core/global"
	"time"
)

func NewCache() {
	global.App.Cache = cache.New(time.Hour*5, time.Hour*2)
}
