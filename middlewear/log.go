package middlewear

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"tiktok/pkg/log"
	"time"
)

func LogMiddleware(logger *log.LogrusLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 让请求继续处理
		c.Next()

		// 记录请求和响应信息
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()

		fmt.Println("statusCode:", statusCode)
		log.HttpStatusCode = statusCode

		logrus.Infof("[TikTok-GIN]  %s | %3d | %13v | %15s | %s | router ==> %5s | %10s ",
			endTime.Format("2006-01-02 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			userAgent)
	}
}
