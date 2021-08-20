/**
 * @Author: koulei
 * @Description:
 * @File: requestid
 * @Version: 1.0.0
 * @Date: 2021/8/17 21:23
 */

package requestid

import (
	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
)

const XRequestId = "X-Request-Id"

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqId := c.GetHeader(XRequestId)
		if reqId == "" {
			reqId = uuid.New()
			c.Request.Header.Set(XRequestId, reqId)
		}

		c.Header("X-Request-Id", reqId)
		c.Next()
	}
}
