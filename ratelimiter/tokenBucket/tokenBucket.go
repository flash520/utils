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
	startAt   time.Time
	size      int64
	tokens    int64
	tokenRate int64
}

func NewRateLimiter(tokenRate, size int64) *RateLimiter {
	return &RateLimiter{
		startAt:   time.Now(),
		tokenRate: tokenRate,
		size:      size,
	}
}

func (rl *RateLimiter) Grant() bool {
	now := time.Now()
	in := now.Sub(rl.startAt).Milliseconds() * rl.tokenRate

	rl.tokens = min(rl.size, rl.tokens+in)
	rl.startAt = now

	if rl.tokens > 0 {
		rl.tokens--
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
