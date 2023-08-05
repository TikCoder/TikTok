package db

import (
	"github.com/jinzhu/gorm"
	"tiktok2023/model"
)

// CountFavoriteByUserId 根据 userId 查询点赞数量
func CountFavoriteByUserId(db *gorm.DB, userId int64) (int64, error) {
	var res int64
	err := db.Table("favorite").Where("user_id = ?", userId).Count(&res).Error
	if err != nil {
		return int64(0), err
	}
	return res, nil
}

// IsFavorite 查询该字段记录是否存在
func IsFavorite(db *gorm.DB, userId, videoId int64) (bool, error) {
	err := db.Table("favorite").
		Where("user_id = ? and video_id = ?", userId, videoId).
		First(&model.Favorite{}).Error
	if err != nil {
		// 错误为查找记录不存在
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			// todo log err
			return false, err
		}
	}
	return true, nil
}
