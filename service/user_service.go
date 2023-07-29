package service

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"tiktok2023/db"
	"tiktok2023/model"
)

type UserService struct{}

// UserNameIsExit 查询用户是否存在
func (user *UserService) UserNameIsExit(username string) (*model.User, bool, error) {
	userInfo, err := db.GetUserInfoByUsername(db.DB, username)
	if err != nil {
		// 数据记录不存在
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, false, nil
		}
		// 数据库的错误
		return nil, false, err
	}
	// 数据记录存在
	return userInfo, true, nil
}

// InsertUser 注册用户
func (user *UserService) InsertUser(username, password string) (*model.User, error) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	userInfo := model.User{
		Username:  username,
		Password:  string(hashPassword),
		Signature: "test sign",
	}

	err := db.CreateUser(db.DB, &userInfo)
	if err != nil {

		return nil, err
	}

	return &userInfo, nil
}

// GetUserInfo 获取用户信息
func (user *UserService) GetUserInfo(userId int64) (*model.User, error) {
	// 根据 用户id 查询信息
	userInfo, err := user.GetUserInfoByUserId(userId)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

// GetUserInfoByUserId 根据用户ID 获取用户信息
func (user *UserService) GetUserInfoByUserId(userId int64) (*model.User, error) {
	userInfo, err := db.GetUserInfoByUserId(db.DB, userId)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
