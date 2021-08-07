/**
 * @Author: koulei
 * @Description:
 * @File: leakyBucket_test
 * @Version: 1.0.0
 * @Date: 2021/8/6 23:27
 */

package leakyBucket

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestLeakyBucket(t *testing.T) {
	rate := int64(1)
	size := int64(2)
	limiter := NewRateLimiter(rate, size)

	for i := 0; i < 20; i++ {
		time.Sleep(time.Millisecond * 100)
		if limiter.Grant() {
			log.Info("放行")
			continue
		}
		log.Error("阻断")
	}
}
