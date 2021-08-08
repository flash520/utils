/**
 * @Author: koulei
 * @Description:
 * @File: tokenBucket
 * @Version: 1.0.0
 * @Date: 2021/8/7 00:16
 */

package tokenBucket

import "time"

type RateLimiter struct {
	startAt   time.Time // 刷新时间
	size      int64     // 桶大小
	tokens    int64     // 桶中令牌数量
	tokenRate int64     // 生成令牌的速率
}

func NewRateLimiter(tokenRate, bucketSize int64) *RateLimiter {
	return &RateLimiter{
		startAt:   time.Now(),
		tokenRate: tokenRate,
		size:      bucketSize,
	}
}

func (r *RateLimiter) Grant() bool {
	now := time.Now()

	// 从上次成功处理的时间到现在生成的令牌数量
	in := now.Sub(r.startAt).Milliseconds() / 1000 * r.tokenRate

	r.tokens = min(r.size, r.tokens+in)

	if r.tokens > 0 {
		r.startAt = now
		r.tokens--
		return true
	}
	return false
}

func min(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}
