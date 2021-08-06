/**
 * @Author: koulei
 * @Description:
 * @File: leakyBucket
 * @Version: 1.0.0
 * @Date: 2021/8/6 22:47
 */

package leakyBucket

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type RateLimiter struct {
	startAt time.Time // 起始时间
	size    int64     // 桶大小(桶能存放的数量)
	rate    int64     // 速率
	water   int64     // 当前桶中的数量
}

func NewRateLimiter(rate, size int64) *RateLimiter {
	return &RateLimiter{
		startAt: time.Now(),
		rate:    rate,
		size:    size,
	}
}

func (r *RateLimiter) Grant() bool {
	now := time.Now()
	out := time.Since(r.startAt).Milliseconds() * r.rate

	// 漏水后残余水
	r.water = max(0, r.water-out)
	log.Warn(r.water)
	r.startAt = now
	if r.water+1 < r.size {
		r.water++
		return true
	}
	return false
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
