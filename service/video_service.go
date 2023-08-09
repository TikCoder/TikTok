package service

import (
	"TikTok/model"
	"TikTok/utils"
	"TikTok/utils/cache"
	db "TikTok/utils/db"
	"TikTok/utils/minioStore"
	"fmt"
	"strconv"
	"time"
)

type VideoService struct {
}

var GetVideoListLimit = 30 // 获取视频的数量

// GetVideoListByFeed 视频流获取固定数量的视频
func (v *VideoService) GetVideoListByFeed(currentTime int64) (*[]model.Video, error) {
	videoList, err := db.GetVideoListByFeed(db.DB, currentTime, GetVideoListLimit)
	if err != nil {
		return videoList, err
	}
	return videoList, nil
}

// PublishVideoToMinio 上传视频到 MinIO
func (v *VideoService) PublishVideoToMinio(localFilePath string, userId int64) (string, error) {
	videoUrl, err := minioStore.UploadFile("video", localFilePath, strconv.FormatInt(userId, 10))
	// todo log info save file
	if err != nil {
		return "", err
	}
	return videoUrl, nil
}

// PublishPictureToMinio 上传图片——封面
func (v *VideoService) PublishPictureToMinio(localFilePath string, userId int64) (string, error) {
	// 截取封面
	imageFile, err := utils.GetImageFile(localFilePath)
	if err != nil {
		return "", err
	}
	pictureUrl, err := minioStore.UploadFile("picture", imageFile, strconv.FormatInt(userId, 10))
	if err != nil {
		return "", err
	}
	return pictureUrl, nil
}

// PublishVideoToDB 视频信息上传到数据库
func (v *VideoService) PublishVideoToDB(authorId int64, videoUrl, pictureUrl, title string) error {
	videoInfo := &model.Video{
		AuthorId:    authorId,
		PlayUrl:     videoUrl,
		CoverUrl:    pictureUrl,
		PublishTime: time.Now().UnixNano() / 1e6,
		Title:       title,
	}
	err := db.CreateVideoInfo(db.DB, videoInfo)
	if err != nil {
		return err
	}
	// 更新到缓存
	err = cache.RedisConn.SetVideoCaCheInfo(videoInfo)
	if err != nil {
		// todo log cache err
		fmt.Println("video cache err")
	}
	return nil
}

// GetAuthorId 根据视频ID查视频信息
func (v *VideoService) GetAuthorId(videoId int64) (*model.Video, error) {
	// 获取视频缓存
	videoInfo, err := cache.RedisConn.GetVideoCaCheInfo(videoId)
	if err != nil {
		// 没有找到，再去找数据库
		videoInfo, err = db.GetVideoInfoByVideoId(db.DB, videoId)
		if err != nil {
			return nil, err
		}
	}
	return videoInfo, err
}
