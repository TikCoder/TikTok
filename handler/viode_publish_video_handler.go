package handler

import (
	"TikTok/config"
	"TikTok/constant"
	"TikTok/model"
	"TikTok/utils"
	"TikTok/utils/jwtUtils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

type PublishVideoHandler struct {
	Req       model.PublishVideoReq
	Resp      model.PublishVideoResp
	fileName  string
	videoPath string
	ctx       *gin.Context
}

var insertVideoInfoToDBLimit = 100

func PublishVideo(c *gin.Context) {
	var handle PublishVideoHandler

	defer func() {
		handle.Resp.StatusMsg = constant.GetErrMsg(handle.Resp.StatusCode)
		c.JSON(http.StatusOK, handle.Resp)
	}()
	var err error
	handle.Req.Token = c.PostForm("token")
	handle.Req.Title = c.PostForm("title")
	handle.Req.Data, err = c.FormFile("data")
	if err != nil {
		// TODO log err
		handle.Resp.StatusCode = constant.ERR_INPUT_INVALID
		return
	}

	handle.fileName = filepath.Base(handle.Req.Data.Filename)
	handle.fileName = fmt.Sprintf("%s_%s", utils.RandomString(), handle.fileName)

	handle.videoPath = config.Conf.Common.VideoPath
	handle.ctx = c

	Run(&handle)
}

// HandleInput 输入检查
func (r *PublishVideoHandler) HandleInput() error {
	if r.Req.Token == "" || r.Req.Title == "" || r.Req.Data == nil {
		// TODO log err
		r.Resp.StatusCode = constant.ERR_INPUT_INVALID
		return constant.ERR_HANDLE_INPUT
	}
	return nil
}

// HandleProcess 处理逻辑
func (r *PublishVideoHandler) HandleProcess() error {
	// 1. 鉴权
	userId, err := jwtUtils.VerifyToken(r.Req.Token)
	if err != nil {
		r.Resp.StatusCode = constant.ERR_TOKEN
		return err
	}

	// todo 鉴黄

	// 2. 将视频存储到本地
	saveFilePath := filepath.Join(r.videoPath, r.fileName)
	// todo log info path
	if err := r.ctx.SaveUploadedFile(r.Req.Data, saveFilePath); err != nil {
		r.Resp.StatusCode = constant.ERR_SAVE_UPLOAD
		return err
	}

	// 3.1 视频上传到MinIO
	videoUrl, err := videoService.PublishVideoToMinio(saveFilePath, userId)
	if err != nil {
		r.Resp.StatusCode = constant.ERR_UPLOAD_MINIO
		return err
	}

	// 3.2 截取视频封面第一帧 并上传
	pictureUrl, err := videoService.PublishPictureToMinio(saveFilePath, userId)
	if err != nil {
		pictureUrl = "https://p6-passport.byteacctimg.com/img/user-avatar/de432cd6200bc3d3f7d633a3ccd528d8~180x180.awebp?"
		//r.Resp.StatusCode = constant.ERR_UPLOAD_MINIO
		//return err
	}

	// 4 视频信息插入数据库和缓存
	for i, err := 0, videoService.PublishVideoToDB(userId, videoUrl, pictureUrl, r.Req.Title); i != insertVideoInfoToDBLimit && err != nil; i++ {
		r.Resp.StatusCode = constant.ERR_CREATE_VIDEO
		return err
	}

	// todo 缓存

	return nil
}
