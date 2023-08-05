package db

import (
	"fmt"
	"testing"
	"tiktok2023/config"
)

func TestInitDB(t *testing.T) {
	config.TestFilePath = "../config/config-test.toml"
	// 加载配置
	if err := config.InitViper(); err != nil {
		fmt.Printf("init Viper failed, err:%v\n", err)
		return
	}
	config.Init()
	err := InitDB(config.Conf.MySQL)
	if err != nil {
		fmt.Println(err)
	}
}
