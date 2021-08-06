/**
 * @Author: koulei
 * @Description:
 * @File: conunter
 * @Version: 1.0.0
 * @Date: 2021/8/6 19:15
 */

package counter

import "time"

type RateLimiter struct {
	limitCount   int64         // 速率
	interval     time.Duration // 间隔时间
	requestCount int64         // 请求计数
	startAt      time.Time     // 起始时间
}

func NewRateLimiter(limitCount int64, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		limitCount: limitCount,
		interval:   interval,
		startAt:    time.Now(),
	}
}

func (r *RateLimiter) Grant() bool {
	now := time.Now()
	if now.Before(r.startAt.Add(r.interval)) {
		if r.requestCount < r.limitCount {
			r.requestCount++
			return true
		}
		return false
	}
	r.startAt = time.Now()
	r.requestCount = 0
	return false
}
