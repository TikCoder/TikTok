package db

import (
	"TikTok/model"
	"github.com/jinzhu/gorm"
)

func IsFollow(db *gorm.DB, followId, followerId int64) (bool, error) {
	err := db.Table("relation").
		Where("follow_id = ? and follower_id = ?", followId, followerId).
		First(&model.Relation{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			// todo log err
			return false, err
		}
	}
	return true, nil
}
