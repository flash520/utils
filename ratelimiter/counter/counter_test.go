/**
 * @Author: koulei
 * @Description:
 * @File: counter_test
 * @Version: 1.0.0
 * @Date: 2021/8/6 19:44
 */

package counter

import (
	"fmt"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
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

func BenchmarkRateLimiter_Grant(b *testing.B) {
	limiter := NewRateLimiter(1, time.Second)
	for i := 0; i < b.N; i++ {
		limiter.Grant()
	}
}

func TestRate(t *testing.T) {
	limiter := rate.NewLimiter(1, 2)

	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 100)
		if limiter.Allow() {
			fmt.Println("放行")
			continue
		}
		fmt.Println(time.Now(), "拒绝")
	}
}
