package mysql

import (
	"github.com/ahaostudy/calendar_reminder/model"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ahaostudy/calendar_reminder/conf"
)

var DB *gorm.DB

func InitMySQL() {
	var err error
	DB, err = gorm.Open(mysql.Open(conf.GetConf().MySQL.DSN),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
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
	)
	if err != nil {
		logrus.Fatal(err)
	}
}
