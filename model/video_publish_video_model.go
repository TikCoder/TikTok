package model

import "mime/multipart"

type PublishVideoReq struct {
	Token string                `json:"token"`
	Title string                `json:"title"`
	Data  *multipart.FileHeader `json:"data"`
}

type PublishVideoResp struct {
	Response
}
