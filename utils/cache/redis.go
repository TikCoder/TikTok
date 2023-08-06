package cache

import (
	"TikTok/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	RedisClient *redis.Client
}

var RedisConn Redis

func InitRedis() error {
	conf := config.Conf.Redis
	addr := fmt.Sprintf("%s:%s", conf.Url, conf.Port)
	redisConn := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Pwd,
		DB:       conf.DB,
		PoolSize: conf.PoolSize,
	})
	if redisConn == nil {
		panic("failed to call redis.NewClient")
	}
	RedisConn.RedisClient = redisConn
	_, err := RedisConn.RedisClient.Set(context.Background(), "abc", 100, 1000).Result()
	// todo log info
	_, err = RedisConn.RedisClient.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to ping redis, err:%s")
	}
	return nil
}
