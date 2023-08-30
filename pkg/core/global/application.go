package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/mr-linch/go-tg"
	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	mq "telegram-monitor/pkg/common/rabbit"
	"telegram-monitor/pkg/config"
)

type Application struct {
	Config *config.Configuration
	DB     *gorm.DB
	Redis  *redis.Client
	Cron   *cron.Cron
	Client *tg.Client
	MQ     mq.MessageQueue
	Cache  *cache.Cache
}

var App = new(Application)
