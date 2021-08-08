/**
 * @Author: koulei
 * @Description:
 * @File: main
 * @Version: 1.0.0
 * @Date: 2021/8/6 23:59
 */

package main

import (
	"net/http"

	"gitee.com/flash520/utils/ratelimiter/uberleakybucket"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(rateLimiter())
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "正常请求:%s\n", "(^^)")
	})
	_ = r.Run(":80")
}

func rateLimiter() gin.HandlerFunc {
	rl := uberleakybucket.NewRateLimiter(500)
	// rl := counter.NewRateLimiter(8, time.Second)
	// rl := slidingWindow.NewRateLimiter(time.Second, 1, func() slidingWindow.Window {
	// 	return slidingWindow.NewLocalWindow()
	// })
	return func(c *gin.Context) {
		if !rl.Grant() {
			c.Abort()
			c.String(http.StatusOK, "-_-!! %s\n", "超过请求限速啦！")
			return
		}
		c.Next()
	}
}
