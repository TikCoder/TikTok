package model

type FeedReq struct {
	LatestTime int64  `json:"latest_time,omitempty"` // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	Token      string `json:"token,omitempty"`       // 用户登录状态下设置
}

type FeedResp struct {
	Response
	NextTime  int64       `json:"next_time"`  // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	VideoList []VideoInfo `json:"video_list"` // 视频列表
}

// VideoInfo 视频信息
type VideoInfo struct {
	Video
	Author     UserVideoInfo `json:"author"`      // 视频作者信息
	IsFavorite bool          `json:"is_favorite"` // true-已点赞，false-未点赞
	AuthorId   *struct{}     `json:"-"`
}

// UserVideoInfo 视频作者信息
type UserVideoInfo struct {
	User
	Username *struct{} `json:"-"` // 忽略结构体 User 中的 username
	Password *struct{} `json:"-"`
	IsFollow bool      `json:"is_follow"` // true-已关注，false-未关注
}
