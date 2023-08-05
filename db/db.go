package db

import (
	"fmt"
	"tiktok2023/config"

	"go.uber.org/zap"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.MysqlConfig) error {
	//user := config.Conf.MySQL.User
	//url := config.Conf.MySQL.Url
	//pwd := config.Conf.MySQL.Pwd
	//dbName := config.Conf.MySQL.Dbname

	var err error
	//DB, err = Factory.CreateGorm(user, pwd, url, dbName)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
		cfg.User, cfg.Pwd, cfg.Url, cfg.Dbname,
	)
	DB, err = CreateGorm(dsn, cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.MaxIdleTimeout)
	if err != nil {
		//seelog.Errorf("CreateGorm err %s", err.Error())
		zap.L().Error("CreateGorm err", zap.Error(err))
		return err
	}
	return nil
}
