package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	logmw "tiktok/middlewear"
	"tiktok/pkg/log"
)

var once sync.Once

func init() {
	once.Do(func() {
		_ = log.Init()
	})
}

func main() {
	logger := log.NewLogrusLogger()
	logMiddleware := logmw.LogMiddleware(logger)

	router := gin.Default()
	router.Use(logMiddleware)

	router.GET("/", func(c *gin.Context) {
		logger.Info("Received a request")
		c.String(http.StatusOK, "Hello Gin!")
	})

	_ = router.Run(":8080")
}
