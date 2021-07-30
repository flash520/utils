/**
 * @Author: koulei
 * @Description:
 * @File: object_test
 * @Version: 1.0.0
 * @Date: 2021/7/27 01:18
 */

package test

import (
	"testing"
)

func BenchmarkByte(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]byte, 4096)
		s = append(s, 1)
	}
}

func BenchmarkPool(b *testing.B) {
	type poolChan struct {
		ch chan []byte
	}
	pool := &poolChan{
		ch: make(chan []byte, 100000000),
	}
	for i := 0; i < 100000; i++ {
		pool.ch <- make([]byte, 4096)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case B := <-pool.ch:
			B = append(B, 1)
			pool.ch <- B
		}
	}
}
