package db

import (
	"github.com/jinzhu/gorm"
	"tiktok2023/model"
)

// GetUserInfoByUsername 根据 user_name 获取用户信息
func GetUserInfoByUsername(db *gorm.DB, username string) (*model.User, error) {
	var user = new(model.User)
	err := db.Table("user").Where("user_name = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserInfoByUserId 根据 id 获取用户信息
func GetUserInfoByUserId(db *gorm.DB, userId int64) (*model.User, error) {
	var user = new(model.User)
	err := db.Table("user").Where("id = ?", userId).First(user).Error
	if err != nil {
		// todo: log
		return nil, err
	}
	return user, nil
}

// CreateUser 创建用户
func CreateUser(db *gorm.DB, user *model.User) error {
	err := db.Table("user").Create(user).Error

	if err != nil {
		// todo log err
		return err
	}
	return nil
}
