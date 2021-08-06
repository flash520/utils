/**
 * @Author: koulei
 * @Description:
 * @File: slidingWindow_test
 * @Version: 1.0.0
 * @Date: 2021/8/6 23:53
 */

package slidingWindow

import (
	"fmt"
	"testing"
	"time"
)

func TestSlidingWindow(t *testing.T) {
	// allow 10 requests per second
	rateLimiter := NewRateLimiter(time.Second, 3, func() Window {
		return NewLocalWindow()
	})

	for i := 0; i < 20; i++ {
		time.Sleep(time.Millisecond * 300)
		if rateLimiter.Grant() {
			fmt.Println("Continue to process")
		} else {
			fmt.Println("Exceed rate limit")
		}
	}
}
