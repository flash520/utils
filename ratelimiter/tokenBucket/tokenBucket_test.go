/**
 * @Author: koulei
 * @Description:
 * @File: tokenBucket_test
 * @Version: 1.0.0
 * @Date: 2021/8/7 00:18
 */

package tokenBucket

import (
	"fmt"
	"testing"
	"time"
)

func TestTokenBucket(t *testing.T) {
	tokenRate := int64(5)
	size := int64(10)
	rateLimiter := NewRateLimiter(tokenRate, size)
	for i := 0; i < 80; i++ {
		time.Sleep(time.Microsecond * 1000)
		if rateLimiter.Grant() {
			fmt.Println("Continue to process")
			continue
		}
		fmt.Println("Exceed rate limit")
	}
}
