package model

type FavoriteVideoListReq struct {
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
}

type FavoriteVideoListResp struct {
	Response
	VideoList []VideoInfo `json:"video_list"` // 视频列表
}
