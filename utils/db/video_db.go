package db

import (
	"TikTok/model"
	"github.com/jinzhu/gorm"
)

// CountVideoByUserId 根据 userId 查询发布的视频数量
func CountVideoByUserId(db *gorm.DB, userId int64) (int64, error) {
	var res int64
	err := db.Table("video").Where("author_id = ?", userId).Count(&res).Error
	if err != nil {
		return int64(0), err
	}
	return res, nil
}

// GetVideoListByFeed 根据流信息获取视频
func GetVideoListByFeed(db *gorm.DB, currentTime int64, limit int) (*[]model.Video, error) {
	var videos []model.Video
	err := db.Table("video").Where("publish_time < ?", currentTime).Limit(limit).Order("id DESC").
		Find(&videos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	// TODO: 加 log
	return &videos, nil
}

// CreateVideoInfo 创建视频信息并插入
func CreateVideoInfo(db *gorm.DB, videoInfo *model.Video) error {
	err := db.Table("video").Create(videoInfo).Error
	if err != nil {
		// todo log err
		return err
	}
	return nil
}

func GetVideoInfoByVideoId(db *gorm.DB, videoId int64) (*model.Video, error) {
	videoInfo := &model.Video{}
	err := db.Table("video").Where("id = ?", videoId).Take(videoInfo).Error
	return videoInfo, err
}
