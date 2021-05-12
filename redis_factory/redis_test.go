/**
 * @Author: koulei
 * @Description: TODO
 * @File:  redis_test
 * @Version: 1.0.0
 * @Date: 2021/5/11 22:34
 */

package redis_factory

import (
	"context"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestRedis(t *testing.T) {
	addr := []string{"127.0.0.1:6379"}
	client := CreateRedis(addr).GetConn()
	res, _ := client.Incr(context.Background(), "test").Result()
	log.Info("redis value: ", res)
	log.Infof("namingClient type: %v\n", client)
}
