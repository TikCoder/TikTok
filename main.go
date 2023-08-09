package main

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"tiktok2023/config"
	"tiktok2023/initialize"
)

func main() {
	// 初始化配置
	config.Init()
	// 初始化资源，主要是 MySQL 连接
	err := initialize.InitResource()
	if err != nil {
		fmt.Printf("initialize.InitResource err %s", err.Error())
		seelog.Errorf("initialize.InitResource err %s", err.Error())
		return
	}
	// 创建一个 web 服务
	//gin.DefaultWriter = io.Discard // 不打印消息
	r := gin.Default()
	//debug.SetGCPercent(10)	// GC
	initialize.RegisterRouter(r)
	r.Run(":41555")
}
