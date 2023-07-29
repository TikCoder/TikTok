package db

import (
	"fmt"
	"testing"
	"tiktok2023/config"
)

func TestInitDB(t *testing.T) {
	config.TestFilePath = "../config/config-test.toml"
	config.Init()
	err := InitDB()
	if err != nil {
		fmt.Println(err)
	}
}
