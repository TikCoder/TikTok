package cache

import (
	"TikTok/config"
	"TikTok/model"
	"context"
	"encoding/json"
	"strconv"
	"time"
)

// GetUserCacheInfo 根据 用户ID 查询缓存
func (r *Redis) GetUserCacheInfo(userID int64) (*model.User, error) {
	redisKey := UserPrefix + strconv.FormatInt(userID, 10)
	val, err := r.RedisClient.Get(context.Background(), redisKey).Result()
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	err = json.Unmarshal([]byte(val), user)
	return user, err
}

// SetUserCacheInfo 创建缓存
func (r *Redis) SetUserCacheInfo(user *model.User) error {
	redisKey := UserPrefix + strconv.FormatInt(user.Id, 10)
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	expired := time.Second * time.Duration(config.Conf.Cache.UserExpired)
	_, err = r.RedisClient.Set(context.Background(), redisKey, val, expired*time.Second).Result()
	return err
}

// UpdateUserCacheInfo 为保证缓存命中率，采用更新缓存方式
func (r *Redis) UpdateUserCacheInfo(user *model.User) error {
	err := r.SetUserCacheInfo(user)
	if err != nil { // 如果更新失败，删除
		redisKey := UserPrefix + strconv.FormatInt(user.Id, 10)
		RedisConn.RedisClient.Del(context.Background(), redisKey).Result()
	}
	return err
}
