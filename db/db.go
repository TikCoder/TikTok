package db

import (
	"github.com/cihub/seelog"
	"github.com/jinzhu/gorm"
	"tiktok2023/config"
)

var DB *gorm.DB

func InitDB() error {
	//user := config.Conf.MySQL.User
	//url := config.Conf.MySQL.Url
	//pwd := config.Conf.MySQL.Pwd
	//dbName := config.Conf.MySQL.Dbname

	var err error
	//DB, err = Factory.CreateGorm(user, pwd, url, dbName)
	DB, err = Factory.CreateGorm(config.Conf.MySQL.User,
		config.Conf.MySQL.Pwd, config.Conf.MySQL.Url, config.Conf.MySQL.Dbname)
	if err != nil {
		seelog.Errorf("CreateGorm err %s", err.Error())
		return err
	}
	return nil
}
