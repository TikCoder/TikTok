package initialize

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"tiktok2023/db"
	"tiktok2023/handler"
	"tiktok2023/utils/minioStore"
)

func RegisterRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")
	apiRouter.POST("/user/register/", handler.Register)
	apiRouter.POST("/user/login/", handler.Login)
	apiRouter.GET("/user/", handler.GetUserInfo)
	apiRouter.GET("/feed/", handler.Feed)
	apiRouter.POST("/publish/action/", handler.PublishVideo)
	apiRouter.GET("/publish/list/", handler.PublishList)
	apiRouter.GET("/favorite/list/", handler.FavoriteVideoList)
}

// InitInfra 初始化基础设置 Infrastructure
func InitInfra() error {
	err := db.InitDB()
	if err != nil {
		return err
	}
	err = minioStore.InitMinio()
	if err != nil {
		return err
	}
	return nil
}

// InitResource 初始化服务资源
func InitResource() error {
	err := InitInfra()
	if err != nil {
		seelog.Errorf("InitInfra err %s", err.Error())
		return err
	}
	return nil
}
