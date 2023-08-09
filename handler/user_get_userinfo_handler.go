package handler

import (
	"TikTok/constant"
	"TikTok/model"
	"TikTok/utils/jwtUtils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type GetUserInfoHandler struct {
	Req  model.GetUserInfoReq
	Resp model.GetUserInfoResp
}

func GetUserInfo(c *gin.Context) {
	// 声明响应
	var handle GetUserInfoHandler
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
func (r *GetUserInfoHandler) HandleInput() error {
	if r.Req.UserId == 0 || r.Req.Token == "" {
		// todo: log error
		r.Resp.StatusCode = constant.ERR_INPUT_INVALID
		return constant.ERR_HANDLE_INPUT
	}
	return nil
}

// HandleProcess 处理逻辑
func (r *GetUserInfoHandler) HandleProcess() error {
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

	userInfo, err := userService.GetUserInfo(r.Req.UserId)
	if err != nil {
		r.Resp.StatusCode = constant.ERR_USER_IS_NOT_EXIT
		return err
	}
	r.Resp.User.User = *userInfo
	r.Resp.User.IsFollow = false // 是否被关注

	// todo 缓存

	return nil
}
