package model

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}
