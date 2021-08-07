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
	size    int64     // 桶大小(桶最大能够容纳的请求数量)
	rate    int64     // 从桶中请求通过漏洞流出的速率
	water   int64     // 当前桶内剩余的请求数量
}

func NewRateLimiter(rate, size int64) *RateLimiter {
	return &RateLimiter{
		startAt: time.Now(), // 起始时间
		rate:    rate,       // 处理请求的速率，请求/秒
		size:    size,
	}
}

func (r *RateLimiter) Grant() bool {
	// 计算出水量
	now := time.Now()
	out := int64(time.Since(r.startAt).Seconds()) * r.rate

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
