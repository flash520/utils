/**
 * @Author: koulei
 * @Description:
 * @File: leakyBucket
 * @Version: 1.0.0
 * @Date: 2021/8/8 12:59
 */

package uberleakybucket

import "go.uber.org/ratelimit"

type RateLimiter struct {
	limiter ratelimit.Limiter
}

func NewRateLimiter(rate int, opts ...ratelimit.Option) *RateLimiter {
	return &RateLimiter{limiter: ratelimit.New(rate)}
}

func (r *RateLimiter) Grant() bool {
	take := r.limiter.Take()
	if take.IsZero() {
		return false
	}
	return true
}
