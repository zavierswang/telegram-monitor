package bootstrap

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"telegram-monitor/pkg/core/global"
	"telegram-monitor/pkg/models"
	"time"
)

func ConnectDB() {
	// 根据驱动配置进行初始化
	switch global.App.Config.DB.Driver {
	case "mysql":
		global.App.DB = initMySQLGorm()
	default:
		global.App.DB = initMySQLGorm()
	}
}

func initMySQLGorm() *gorm.DB {
	dbConfig := global.App.Config.DB
	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName: dbConfig.Driver,
		DSN:        dbConfig.Dsn,
	}))
	if err != nil {
		fmt.Printf("connect mysql failed %v\n", err)
		return nil
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	migrate(db)
	return db
}

func migrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(
		models.User{},
		models.Address{},
	)
	if err != nil {
		fmt.Printf("migrate database failed %v\n", err)
		os.Exit(0)
	}
}

func getGormLogger() logger.Interface {
	var logMode logger.LogLevel

	switch global.App.Config.DB.LogMode {
	case "silent":
		logMode = logger.Silent
	case "error":
		logMode = logger.Error
	case "warn":
		logMode = logger.Warn
	case "info":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}

	return logger.New(getGormLogWriter(), logger.Config{
		SlowThreshold:             200 * time.Millisecond,                    // 慢 SQL 阈值
		LogLevel:                  logMode,                                   // 日志级别
		IgnoreRecordNotFoundError: false,                                     // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  !global.App.Config.DB.EnableFileLogWriter, // 禁用彩色打印
	})
}

// 自定义 gorm Writer
func getGormLogWriter() logger.Writer {
	var writer io.Writer

	// 是否启用日志文件
	if global.App.Config.DB.EnableFileLogWriter {
		// 自定义 Writer
		writer = &lumberjack.Logger{
			Filename:   "./logs/" + global.App.Config.DB.LogFilename,
			MaxSize:    200,
			MaxBackups: 5,
			MaxAge:     28,
			Compress:   true,
		}
	} else {
		// 默认 Writer
		writer = os.Stdout
	}
	return log.New(writer, "\r\n", log.LstdFlags)
}
