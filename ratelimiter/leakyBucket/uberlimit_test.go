/**
 * @Author: koulei
 * @Description:
 * @File: uberlimit_test
 * @Version: 1.0.0
 * @Date: 2021/8/7 18:37
 */

package leakyBucket

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/ratelimit"
)

func TestUberLeakyBucket(t *testing.T) {
	limiter := ratelimit.New(3000000)

	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := limiter.Take()
		fmt.Println(i, now.Sub(prev))
		prev = now
	}
}
