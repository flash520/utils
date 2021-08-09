/**
 * @Author: koulei
 * @Description:
 * @File: daemon
 * @Version: 1.0.0
 * @Date: 2021/8/9 18:53
 */

package main

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	args := os.Args
	daemon := false
	for k, v := range args {
		if v == "-d" {
			daemon = true
			args[k] = ""
		}
	}

	if daemon {
		Daemonize(args...)
		return
	}
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Response: %s\n", "我是后台运行的哦")
	})
	_ = r.Run(":80")
}

func Daemonize(args ...string) {
	var arg []string
	if len(args) > 1 {
		arg = args[1:]
	}
	cmd := exec.Command(args[0], arg...)
	cmd.Env = os.Environ()
	_ = cmd.Start()
}
