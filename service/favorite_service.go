package service

import (
	"TikTok/model"
	"TikTok/utils/cache"
	"TikTok/utils/db"
	"fmt"
)

type FavoriteService struct {
}

// IsFavorite 根据 userId, videoId 查询是否点赞
// 只要 favorite 表中存在就表示点赞
func (favorite *FavoriteService) IsFavorite(userId, videoId int64) (bool, error) {
	return db.IsFavorite(db.DB, userId, videoId)
}

// FavoriteAction 点赞动作
func (favorite *FavoriteService) FavoriteAction(actionType int64, videoInfo *model.Video, favoriteInfo, authorInfo *model.User) error {
	// 检查是否存在
	isExit, err := db.IsFavorite(db.DB, favoriteInfo.Id, videoInfo.Id)
	if err != nil {
		return err
	}
	if actionType == 1 {
		// 如果不存在 直接返回
		if isExit {
			return nil
		}
	} else {
		// 如果不存在 直接返回
		if !isExit {
			return nil
		}
	}
	err = db.UpdateFavorite(db.DB, actionType, favoriteInfo, authorInfo, videoInfo)
	if err != nil {
		return err
	}
	// 更新缓存
	err = cache.RedisConn.UpdateUserCacheInfo(favoriteInfo)
	if err != nil {
		// todo log err
		fmt.Println("FavoriteInfo cache err", err)
	}
	err = cache.RedisConn.UpdateUserCacheInfo(authorInfo)
	if err != nil {
		// todo log err
		fmt.Println("authorInfo cache err", err)
	}
	err = cache.RedisConn.UpdateVideoCacheInfo(videoInfo)
	if err != nil {
		// todo log err
		fmt.Println("videoInfo cache err", err)
	}
	return nil
}
