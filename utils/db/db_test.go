package db

import (
	"TikTok/config"
	"fmt"
	"testing"
)

func TestInitDB(t *testing.T) {
	config.TestFilePath = "../config/config-test.toml"
	config.Init()
	err := InitDB()
	if err != nil {
		fmt.Println(err)
	}
}
