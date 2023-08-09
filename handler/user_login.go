package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"tiktok2023/constant"
	"tiktok2023/model"
	"tiktok2023/utils/jwtUtils"
)

type LoginHandle struct {
	Req  model.LoginReq
	Resp model.LoginResp
}

func Login(c *gin.Context) {
	// 声明响应
	var handle LoginHandle
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
func (r *LoginHandle) HandleInput() error {
	if r.Req.Username == "" || r.Req.Password == "" {
		// TODO log err
		r.Resp.StatusCode = constant.ERR_INPUT_INVALID
		return errors.New("input invalid")
	}
	return nil
}

// HandleProcess 处理逻辑
func (r *LoginHandle) HandleProcess() error {
	// 1. 检查用户是否存在
	userInfo, isExit, err := userService.UserNameIsExit(r.Req.Username)
	if err != nil {
		r.Resp.StatusCode = constant.ERR_GET_USER
		return err
	}
	// 如果不存在
	if !isExit {
		r.Resp.StatusCode = constant.ERR_USER_IS_NOT_EXIT
		return errors.New("username is not exited")
	}

	// 2. 密码校验
	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(r.Req.Password))
	if err != nil {
		r.Resp.StatusCode = constant.ERR_GET_USER
		return err
	}

	// 3. 登录 生成 token
	token, err := jwtUtils.GenToken(userInfo.Id, userInfo.Username)
	if err != nil {
		r.Resp.StatusCode = constant.ERR_GEN_TOKEN
		return err
	}

	r.Resp.UserId = userInfo.Id
	r.Resp.Token = token

	// todo 缓存

	return nil
}
