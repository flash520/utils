/**
 * @Author: koulei
 * @Description:
 * @File: leakyBucket_test
 * @Version: 1.0.0
 * @Date: 2021/8/6 23:27
 */

package leakyBucket

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestLeakyBucket(t *testing.T) {
	rate := int64(4)
	size := int64(10)
	limiter := NewRateLimiter(rate, size)

	for i := 0; i < 50; i++ {
		time.Sleep(time.Millisecond * 100)
		if limiter.Grant() {
			fmt.Printf("%d 放行\n", i)
			continue
		}
		fmt.Println(time.Now().Unix(), "阻断")
	}
}

func TestRate(t *testing.T) {
	limiter := rate.NewLimiter(1, 100)
	for i := 0; i < 101; i++ {
		time.Sleep(time.Millisecond * 100)
		if limiter.Allow() {
			fmt.Println("success")
		} else {
			fmt.Println("failed")
		}
	}
}
