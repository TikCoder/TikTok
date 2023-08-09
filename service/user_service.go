package service

import (
	"TikTok/model"
	"TikTok/utils/cache"
	"TikTok/utils/db"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

// UserNameIsExit 查询用户是否存在
func (user *UserService) UserNameIsExit(username string) (*model.User, bool, error) {
	//userInfo, err := cache.RedisConn.GetUserCacheInfo()
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
	// 存入缓存
	err = cache.RedisConn.SetUserCacheInfo(&userInfo)
	if err != nil {
		// log err cache set
		fmt.Println("CreateUser cache set err")
	}
	return &userInfo, nil
}

// GetUserInfo 获取用户信息
func (user *UserService) GetUserInfo(userId int64) (*model.User, error) {
	// 1. 查缓存
	userInfo, err := cache.RedisConn.GetUserCacheInfo(userId)
	if err == nil && userInfo.Id == userId {
		// todo log cache succ
		return userInfo, nil
	}
	// 2. 根据 用户id 查询信息
	userInfo, err = db.GetUserInfoByUserId(db.DB, userId)
	if err != nil {
		return nil, err
	}

	// 3. 创建缓存
	err = cache.RedisConn.SetUserCacheInfo(userInfo)
	if err != nil {
		// todo cache err
	}
	return userInfo, nil
}
