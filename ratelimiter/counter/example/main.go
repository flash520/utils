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
	"time"

	"gitee.com/flash520/utils/ratelimiter/counter"
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
	rl := counter.NewRateLimiter(2, time.Second*10)
	return func(c *gin.Context) {
		if rl.Grant() {
			c.Next()
			return
		}
		c.String(http.StatusOK, "-_-!! %s\n", "超过请求限速啦！")
		c.Abort()
	}
}
