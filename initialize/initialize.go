package initialize

import (
	"tiktok2023/config"
	"tiktok2023/db"
	"tiktok2023/handler"
	"tiktok2023/logger"
	"tiktok2023/utils/minioStore"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	return r
}

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
	err := db.InitDB(config.Conf.MySQL)
	if err != nil {
		zap.L().Error("init DB failed", zap.Error(err))
		return err
	}
	err = minioStore.InitMinio()
	if err != nil {
		zap.L().Error("Init Minio failed", zap.Error(err))
		return err
	}
	return nil
}

// InitResource 初始化服务资源
func InitResource() error {
	err := InitInfra()
	if err != nil {
		//seelog.Errorf("InitInfra err %s", err.Error())
		zap.L().Error("InitInfra err", zap.Error(err))
		return err
	}
	return nil
}
