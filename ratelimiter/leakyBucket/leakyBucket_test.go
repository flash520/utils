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

	log "github.com/sirupsen/logrus"
)

func TestLeakyBucket(t *testing.T) {
	rate := int64(6)
	size := int64(10)
	limiter := NewRateLimiter(rate, size)
	for i := 0; i < 20; i++ {
		if limiter.Grant() {
			log.Info("放行")
			continue
		}
		log.Error("阻断")
	}
}
