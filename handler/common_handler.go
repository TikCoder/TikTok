package handler

import "TikTok/service"

var userService service.UserService
var videoService service.VideoService
var favoriteService service.FavoriteService
var relationService service.RelationService

func Run(handler HandleInterface) error {
	err := handler.HandleInput()
	if err != nil {
		return err
	}
	err = handler.HandleProcess()
	return err
}

// HandleInterface handler 接口
type HandleInterface interface {
	HandleInput() error
	HandleProcess() error
}
