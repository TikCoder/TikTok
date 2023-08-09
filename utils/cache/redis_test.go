package cache

import (
	"TikTok/config"
	"fmt"
	"testing"
)

func TestInitRedis(t *testing.T) {
	config.TestFilePath = "../../config/config-test.toml"
	config.Init()
	err := InitRedis()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(RedisConn)
}
