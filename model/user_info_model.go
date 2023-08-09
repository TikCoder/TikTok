package model

type GetUserInfoReq struct {
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
}

type GetUserInfoResp struct {
	Response
	User Object
}

// Object
type Object struct {
	User
	IsFollow bool      `json:"is_follow"` // true-已关注，false-未关注
	Username *struct{} `json:"-"`
	Password *struct{} `json:"-"`
}
