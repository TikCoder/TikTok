package model

type PublishListReq struct {
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
}

type PublishListResp struct {
	Response
	VideoList []VideoInfo `json:"video_list"` // 视频列表
}
