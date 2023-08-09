package db

import (
	"TikTok/model"
	"fmt"
	"github.com/jinzhu/gorm"
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

// UpdateFavorite 点赞/取消点赞
func UpdateFavorite(db *gorm.DB, actionType int64, favoriteInfo, authorInfo *model.User, videoInfo *model.Video) error {
	favorite := &model.Favorite{
		UserId:  favoriteInfo.Id,
		VideoId: videoInfo.Id,
	}
	var incr int64 // 对应点赞数量 +/- 1

	transaction := db.Begin() // 开启事务
	if actionType == 1 {      // 点赞
		incr = 1
		// 1 favorite 表创建关系
		err := transaction.Table("favorite").Create(favorite).Error
		if err != nil {
			// 遇到错误回滚事务
			transaction.Rollback()
			// todo log err
			fmt.Println("UpdateFavorite err/n")
			return err
		}
	} else {
		incr = -1
		// 删除favorite表数据
		err := transaction.Table("favorite").
			Where("user_id = ? and video_id = ?", favoriteInfo.Id, videoInfo.Id).
			Take(favorite).Error
		if err != nil {
			transaction.Rollback()
			return err
		}
		// 删除
		err = db.Table("favorite").Delete(favorite).Error
		if err != nil {
			transaction.Rollback()
			return err
		}
	}

	// 2. 点赞用户喜欢数 +/- 1
	err := transaction.Table("user").Where("id = ?", favoriteInfo.Id).
		Update("favorite_count", gorm.Expr("favorite_count + ?", incr)).Take(favoriteInfo).Error
	if err != nil {
		// 遇到错误回滚事务
		transaction.Rollback()
		return err
	}
	// 3. 作者用户被喜欢数 +/- 1
	err = transaction.Table("user").Where("id = ?", authorInfo.Id).
		Update("total_favorited", gorm.Expr("total_favorited + ?", incr)).Take(authorInfo).Error
	if err != nil {
		// 遇到错误回滚事务
		transaction.Rollback()
		return err
	}

	// 4. 视频被喜欢数 +/- 1
	err = transaction.Table("video").Where("id = ?", videoInfo.Id).
		Update("favorite_count", gorm.Expr("favorite_count + ?", incr)).Take(videoInfo).Error
	if err != nil {
		// 遇到错误回滚事务
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil
}
