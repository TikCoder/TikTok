package model

type RegisterResponse struct {
	Response
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
