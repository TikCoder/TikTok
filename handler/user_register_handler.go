package handler

import (
	"TikTok/constant"
	"TikTok/model"
	"TikTok/utils/jwtUtils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterHandler struct {
	Req  model.RegisterReq
	Resp model.RegisterResponse
}

func Register(c *gin.Context) {
	// 声明响应
	var handle RegisterHandler
	defer func() {
		handle.Resp.StatusMsg = constant.GetErrMsg(handle.Resp.StatusCode)
		c.JSON(http.StatusOK, handle.Resp)
	}()

	// 解析请求
	handle.Req.Username = c.Query("username")
	handle.Req.Password = c.Query("password")

	Run(&handle)
}

// HandleInput 输入检查
func (r *RegisterHandler) HandleInput() error {

	if r.Req.Username == "" || r.Req.Password == "" {
		// TODO log err
		r.Resp.StatusCode = constant.ERR_INPUT_INVALID
		return constant.ERR_HANDLE_INPUT
	}
	return nil
}

// HandleProcess 处理逻辑
func (r *RegisterHandler) HandleProcess() error {
	// 1. 判断用户是否存在
	_, isExit, err := userService.UserNameIsExit(r.Req.Username)
	if err != nil {
		r.Resp.StatusCode = constant.ERR_GET_USER
		return err
	}
	if isExit {
		// 如果存在, 报错
		//seelog.Error("UserNameIsExist err!")
		r.Resp.StatusCode = constant.ERR_USER_IS_EXIT
		return errors.New("username has exited")
	}

	// 2. 创建用户
	userInfo, err := userService.InsertUser(r.Req.Username, r.Req.Password)
	if err != nil {
		return err
	}
	// 3. 利用 id 和 username 生成 token
	token, err := jwtUtils.GenToken(userInfo.Id, r.Req.Username)
	if err != nil {
		r.Resp.StatusCode = constant.ERR_GEN_TOKEN
		return err
	}

	r.Resp.Token = token
	r.Resp.UserID = userInfo.Id

	return nil
}
