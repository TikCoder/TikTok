package handler

import (
	"TikTok/constant"
	"TikTok/model"
	"TikTok/utils/jwtUtils"
	"TikTok/utils/lock"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FavoriteActionHandler struct {
	Req  model.FavoriteActionReq
	Resp model.FavoriteActionResp
}

func FavoriteAction(c *gin.Context) {
	var handle FavoriteActionHandler
	var err error
	defer func() {
		handle.Resp.StatusMsg = constant.GetErrMsg(handle.Resp.StatusCode)
		c.JSON(http.StatusOK, handle.Resp)
	}()

	// 解析请求
	handle.Req.VideoId, err = strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		// todo: log error
		handle.Resp.StatusCode = constant.ERR_INPUT_INVALID
		return
	}
	handle.Req.ActionType, err = strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil {
		// todo: log error
		handle.Resp.StatusCode = constant.ERR_INPUT_INVALID
		return
	}

	handle.Req.Token = c.Query("token")

	Run(&handle)
}

// HandleInput 输入检查
func (r *FavoriteActionHandler) HandleInput() error {
	if (r.Req.ActionType != int64(1) && r.Req.ActionType != int64(2)) || r.Req.Token == "" || r.Req.ActionType == 0 {
		// todo: log error
		r.Resp.StatusCode = constant.ERR_INPUT_INVALID
		return constant.ERR_HANDLE_INPUT
	}
	return nil
}

func (r *FavoriteActionHandler) HandleProcess() error {
	// jwt 鉴权
	userId, err := jwtUtils.VerifyToken(r.Req.Token)
	if err != nil {
		r.Resp.StatusCode = constant.ERR_TOKEN
		return err
	}

	// 2. 执行点赞和取消，并更新缓存
	// 2.0 加锁
	// 加锁 以用户ID和视频ID为key
	// 此时和该点赞用户、被点赞视频的相关操作都等待
	favoriteMutex := lock.MutexLock.GetLock(userId)
	//authorMutex := lock.MutexLock.GetLock(authorInfo.Id)
	VideoMutex := lock.MutexLock.GetLock(r.Req.VideoId)
	defer lock.MutexLock.Unlock(favoriteMutex)
	//defer lock.MutexLock.Unlock(authorMutex)
	defer lock.MutexLock.Unlock(VideoMutex)

	// 2.1 查询点赞的用户信息
	favoriteInfo, err := userService.GetUserInfo(userId)
	// 2.2 查询被点赞的用户信息
	videoInfo, err := videoService.GetAuthorId(r.Req.VideoId)
	if err != nil {
		r.Resp.StatusCode = constant.ERR_VIDEO_IS_NOT_EXIT
		return err
	}
	// 2.3 根据作者ID获取作者信息
	authorInfo, err := userService.GetUserInfo(videoInfo.AuthorId)
	if err != nil {
		r.Resp.StatusCode = constant.ERR_USER_IS_NOT_EXIT
		return err
	}
	// 2.4 更新数据库和缓存
	err = favoriteService.FavoriteAction(r.Req.ActionType, videoInfo, favoriteInfo, authorInfo)
	if err != nil {
		r.Resp.StatusCode = constant.ERR_FAV_ACTION
	}

	// 3. 解锁
	return nil
}
