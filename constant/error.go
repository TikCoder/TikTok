package constant

import (
	"errors"
	"fmt"
)

var (
	ERR_HANDLE_INPUT = errors.New("handle input error")
)
var (
	SUCCESS              = int64(0)
	ERR_INPUT_INVALID    = int64(8020)
	ERR_USER_IS_EXIT     = int64(8023)
	ERR_CREATE_USER      = int64(8024)
	ERR_GET_USER         = int64(8025)
	ERR_USER_IS_NOT_EXIT = int64(8026)
	ERR_PASSWORD         = int64(8027)
	ERR_GEN_TOKEN        = int64(8028)
	ERR_TOKEN            = int64(8029)
	ERR_GET_VIDEO        = int64(8031)
	ERR_GET_FAVORITE     = int64(8032)
	ERR_GET_RELATION     = int64(8033)
	ERR_SAVE_UPLOAD      = int64(8034)
	ERR_UPLOAD_MINIO     = int64(8035)
	ERR_CREATE_VIDEO     = int64(8036)
)

// 错误码对应的错误信息
var errMsgDic = map[int64]string{
	SUCCESS:              "ok",
	ERR_INPUT_INVALID:    "input invalid",
	ERR_USER_IS_EXIT:     "username has exited",
	ERR_CREATE_USER:      "create user failed",
	ERR_GET_USER:         "get user failed",
	ERR_USER_IS_NOT_EXIT: "username is not exited",
	ERR_PASSWORD:         "password wrong",
	ERR_GEN_TOKEN:        "generate token failure",
	ERR_TOKEN:            "token input is wrong",
	ERR_GET_FAVORITE:     "get favorite failed",
	ERR_GET_VIDEO:        "get video failed",
	ERR_GET_RELATION:     "get relation failed",
	ERR_SAVE_UPLOAD:      "save upload video/image failed",
	ERR_UPLOAD_MINIO:     "upload to minio failed",
	ERR_CREATE_VIDEO:     "create video failed",
}

// GetErrMsg 获取错误描述
func GetErrMsg(code int64) string {
	if msg, ok := errMsgDic[code]; ok {
		return msg
	}
	return fmt.Sprintf("unknown error code %d", code)
}
