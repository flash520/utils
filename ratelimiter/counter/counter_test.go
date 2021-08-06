/**
 * @Author: koulei
 * @Description:
 * @File: counter_test
 * @Version: 1.0.0
 * @Date: 2021/8/6 19:44
 */

package counter

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestNewRateLimiter(t *testing.T) {
	limitCount := int64(6)
	interval := time.Second

	limiter := NewRateLimiter(limitCount, interval)
	time.Sleep(time.Millisecond * 1000)
	for i := 0; i < 10; i++ {

		if limiter.Grant() {
			log.Info("放行")
			continue
		}
		log.Warn("阻断")
	}
}
