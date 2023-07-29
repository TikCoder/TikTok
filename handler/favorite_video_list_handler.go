package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok2023/constant"
	"tiktok2023/model"
	"tiktok2023/utils/jwtUtils"
)

type FavoriteVideoListHandle struct {
	Req  model.FavoriteVideoListReq
	Resp model.FavoriteVideoListResp
}

func FavoriteVideoList(c *gin.Context) {
	var handle PublishListHandle
	var err error
	defer func() {
		handle.Resp.StatusMsg = constant.GetErrMsg(handle.Resp.StatusCode)
		c.JSON(http.StatusOK, handle.Resp)
	}()

	// 解析请求
	handle.Req.UserId, err = strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		// todo: log error
		handle.Resp.StatusCode = constant.ERR_INPUT_INVALID
		return
	}

	handle.Req.Token = c.Query("token")

	Run(&handle)
}

// HandleInput 输入检查
func (r *FavoriteVideoListHandle) HandleInput() error {
	if r.Req.UserId == 0 || r.Req.Token == "" {
		// todo: log error
		r.Resp.StatusCode = constant.ERR_INPUT_INVALID
		return constant.ERR_HANDLE_INPUT
	}
	return nil
}

func (r *FavoriteVideoListHandle) HandleProcess() error {
	// jwt 鉴权
	jwtId, err := jwtUtils.VerifyToken(r.Req.Token)
	if err != nil {
		r.Resp.StatusCode = constant.ERR_TOKEN
		return err
	}

	if jwtId != r.Req.UserId {
		// TODO log err
		r.Resp.StatusCode = constant.ERR_GET_USER
		return errors.New("userId is wrong")
	}
	r.Resp.VideoList = make([]model.VideoInfo, 0)
	return nil
}
