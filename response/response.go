/**
 * @Author: koulei
 * @Description: TODO
 * @File:  response
 * @Version: 1.0.0
 * @Date: 2021/5/12 00:20
 */

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	RequestStatusOkCode   int    = 1
	RequestStatusOkMsg    string = "Success"
	RequestStatusFailCode int    = 0
	RequestStatusFailMsg  string = "Failed"

	ErrorsNoAuthorization    string = "token鉴权未通过，请通过token授权接口重新获取token,"
	ValidatorParamsCheckFail string = "请求参数错误"
)

func ReturnJson(c *gin.Context, httpCode int, dataCode int, msg interface{}, data interface{}, success bool) {
	c.JSON(httpCode, gin.H{
		"code":    dataCode,
		"msg":     msg,
		"data":    data,
		"success": success,
	})
}

// Success 直接返回成功
func Success(c *gin.Context, dataCode int, msg interface{}, data interface{}) {
	ReturnJson(c, http.StatusOK, dataCode, msg, data, true)
}

// Fail 失败的业务逻辑
func Fail(c *gin.Context, dataCode int, msg interface{}, data interface{}) {
	ReturnJson(c, http.StatusOK, dataCode, msg, data, false)
	c.Abort()
}

// ErrorAuthFail 权限校验失败
func ErrorAuthFail(c *gin.Context) {
	ReturnJson(c, http.StatusUnauthorized, http.StatusUnauthorized, ErrorsNoAuthorization, nil, false)
	c.Abort()
}

// ErrorParam 参数校验失败
func ErrorParam(c *gin.Context, wrongParam interface{}) {
	ReturnJson(c, http.StatusOK, RequestStatusOkCode, ValidatorParamsCheckFail, wrongParam, false)
	c.Abort()
}
