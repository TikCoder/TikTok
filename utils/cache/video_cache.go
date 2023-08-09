package cache

import (
	"TikTok/config"
	"TikTok/model"
	"context"
	"encoding/json"
	"strconv"
	"time"
)

func (r *Redis) SetVideoCaCheInfo(videoInfo *model.Video) error {
	redisKey := VideoPrefix + strconv.FormatInt(videoInfo.Id, 10)
	val, err := json.Marshal(videoInfo)
	if err != nil {
		return err
	}
	_, err = r.RedisClient.Set(context.Background(), redisKey, val,
		time.Second*time.Duration(config.Conf.Cache.UserExpired)).Result()
	return err
}

func (r *Redis) GetVideoCaCheInfo(videoID int64) (*model.Video, error) {
	redisKey := VideoPrefix + strconv.FormatInt(videoID, 10)
	val, err := r.RedisClient.Get(context.Background(), redisKey).Result()
	if err != nil {
		return nil, err
	}
	videoInfo := &model.Video{}
	err = json.Unmarshal([]byte(val), videoInfo)
	return videoInfo, err
}

// UpdateVideoCacheInfo 为保证缓存命中率，采用更新缓存方式
func (r *Redis) UpdateVideoCacheInfo(video *model.Video) error {
	err := r.SetVideoCaCheInfo(video)
	if err != nil { // 如果更新失败，删除
		redisKey := VideoPrefix + strconv.FormatInt(video.Id, 10)
		RedisConn.RedisClient.Del(context.Background(), redisKey).Result()
	}
	return err
}
