package mysql

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ahaostudy/calendar_reminder/model"

	"github.com/ahaostudy/calendar_reminder/conf"
)

var DB *gorm.DB

func InitMySQL() {
	var err error
	loggerConf := logger.New(
		logrus.StandardLogger(),
		logger.Config{
			SlowThreshold: time.Millisecond,
			LogLevel:      logger.Warn,
			Colorful:      false,
		},
	)
	DB, err = gorm.Open(mysql.Open(conf.GetConf().MySQL.DSN),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
			Logger:                 loggerConf,
		},
	)
	if err != nil {
		logrus.Fatal(err)
	}

	migrate()
}

func migrate() {
	err := DB.AutoMigrate(
		new(model.User),
		new(model.Task),
	)
	if err != nil {
		logrus.Fatal(err)
	}
}
