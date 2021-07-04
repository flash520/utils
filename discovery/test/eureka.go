/**
 * @Author: koulei
 * @Description: TODO
 * @File:  eureka_test
 * @Version: 1.0.0
 * @Date: 2021/7/4 11:05
 */

package main

import (
	"time"

	"gitee.com/flash520/utils/discovery"
	log "github.com/sirupsen/logrus"
)

func main() {
	_ = discovery.InitDiscovery("eureka", "upgrade", "admin", "admin", "localhost:7000/eureka/", 80)

	for {
		time.Sleep(time.Second * 2)
		log.Infof("服务地址: %v\n", discovery.GetInstance("upgrade"))
	}
}
