package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok2023/constant"
	"tiktok2023/model"
	"tiktok2023/utils/jwtUtils"
	"time"
)

type FeedHandler struct {
	Req  model.FeedReq
	Resp model.FeedResp
}

func Feed(c *gin.Context) {
	var handle FeedHandler

	defer func() {
		handle.Resp.StatusMsg = constant.GetErrMsg(handle.Resp.StatusCode)
		c.JSON(http.StatusOK, handle.Resp)
	}()

	// 1. 如果没有时间，从当前时间开始
	currentTime, err := strconv.ParseInt(c.Query("latest_time"), 10, 64)
	if err != nil || currentTime == int64(0) {
		currentTime = time.Now().UnixNano() / 1e6
	}
	handle.Req.LatestTime = currentTime
	// 2. token
	handle.Req.Token = c.Query("token")
	Run(&handle)
}

// HandleInput 输入检查
func (r *FeedHandler) HandleInput() error {
	return nil
}

// HandleProcess 处理逻辑
func (r *FeedHandler) HandleProcess() error {
	// 1. 根据时间戳获取视频列表
	videoList, err := videoService.GetVideoListByFeed(r.Req.LatestTime)
	if err != nil {
		return err
	}

	if len(*videoList) == 0 {
		return nil
	}
	r.Resp.VideoList = make([]model.VideoInfo, len(*videoList))
	r.Resp.NextTime = (*videoList)[0].PublishTime

	// 2. 根据 author_id 获取 视频的作者信息
	for k, v := range *videoList {
		// 2.1根据用户ID查询视频用户信息
		userInfo, err := userService.GetUserInfoByUserId(v.AuthorId)
		if err != nil {
			return err
		}
		r.Resp.VideoList[k].Video = v
		r.Resp.VideoList[k].Author.User = *userInfo

		// 2.2 如果登录
		if r.Req.Token != "" {
			// 2.2.1 查看是否点赞
			userId, err := jwtUtils.VerifyToken(r.Req.Token)
			if err != nil {
				r.Resp.StatusCode = constant.ERR_TOKEN
				return err
			}
			isExit, err := favoriteService.IsFavorite(userId, v.Id)
			if err != nil {
				r.Resp.StatusCode = constant.ERR_GET_FAVORITE
				return err
			}
			if isExit {
				r.Resp.VideoList[k].IsFavorite = true
			}

			// 2.2.2 查看是否关注该视频的作者
			isExit, err = relationService.IsFollow(userId, userInfo.Id)
			if err != nil {
				r.Resp.StatusCode = constant.ERR_GET_RELATION
				return err
			}
			if isExit {
				r.Resp.VideoList[k].Author.IsFollow = true
			}
		}
	}
	// todo 缓存

	return nil
}
