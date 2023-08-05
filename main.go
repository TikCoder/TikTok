package main

import (
	"fmt"
	"tiktok2023/config"
	"tiktok2023/initialize"
	"tiktok2023/logger"
	"tiktok2023/pkg/snowflake"

	"go.uber.org/zap"
)

func main() {
	// 初始化配置
	config.Init()
	// 初始化日志
	if err := logger.Init(config.Conf.Log); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync() // 延迟注册，把缓冲区的日志给追加到日志文件中
	zap.L().Debug("logger init success...")

	// 初始化资源，主要是 MySQL 连接
	err := initialize.InitResource()
	if err != nil {
		//fmt.Printf("initialize.InitResource err %s", err.Error())
		// seelog.Errorf("initialize.InitResource err %s", err.Error())
		zap.L().Error("initialize.InitResource err", zap.Error(err))
		return
	}

	// snowflake to generate ID
	if err := snowflake.Init(config.Conf.Common.StartTime, config.Conf.Common.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	// 创建一个 web 服务
	//gin.DefaultWriter = io.Discard // 不打印消息
	//r := gin.Default()
	r := initialize.Setup()

	//debug.SetGCPercent(10)	// GC
	initialize.RegisterRouter(r)
	r.Run(":41555")
}
