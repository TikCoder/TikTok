package service

import "tiktok2023/db"

type FavoriteService struct {
}

// IsFavorite 根据 userId, videoId 查询是否点赞
// 只要 favorite 表中存在就表示点赞
func (favorite *FavoriteService) IsFavorite(userId, videoId int64) (bool, error) {
	return db.IsFavorite(db.DB, userId, videoId)
}
